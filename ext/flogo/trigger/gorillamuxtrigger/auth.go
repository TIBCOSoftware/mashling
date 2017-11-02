/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package gorillamuxtrigger

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

type hashedCred struct {
	salt     string
	username string
	password string
}

type plainCred struct {
	username string
	password string
}

type Auth interface {
	authenticate(clientCred string) bool
}

type basic struct {
}

const (
	basicAuthFile = "basicAuthFile"
)

var mu sync.Mutex
var credMap map[string]hashedCred

func basicAuth() Auth {
	return &basic{}
}

// isAuthEnabled check if authentication is enabled
func isAuthEnabled(settings map[string]interface{}) bool {
	// Check if basic auth is in use
	if _, ok := settings[basicAuthFile]; !ok {
		return false
	}

	return true
}

func setupAuth(settings map[string]interface{}) {
	err := loadCreds(settings[basicAuthFile].(string))
	if err != nil {
		log.Error(err)
		panic("Unable to load creds file")
	}
}

func loadCreds(pathToFile string) error {
	mu.Lock()
	if credMap != nil {
		mu.Unlock()
		return nil
	}

	credMap = make(map[string]hashedCred)
	mu.Unlock()

	csvFile, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ':'

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Error(error)
		}

		if len(line) == 3 {
			// hashed password
			credMap[line[0]] = hashedCred{
				salt:     line[1],
				password: line[2],
			}
		} else if len(line) == 2 {
			// plan text password
			credMap[line[0]] = hashedCred{
				password: line[1],
			}
		}
	}

	return nil
}

// sha applies sha256 to the string
func sha(str string) string {
	bytes := []byte(str)

	h := sha256.New()
	h.Write(bytes)
	code := h.Sum(nil)
	codestr := hex.EncodeToString(code)

	return codestr
}

// Authenticate check if the request is allowed
func authenticate(r *http.Request, settings map[string]interface{}) bool {
	// Authenticate using basic auth
	clientCred := r.Header.Get("Authorization")
	re, err := regexp.Compile(`(?i)Basic (.*)`)
	if err != nil {
		log.Error(err)
		return false
	}
	result := re.FindStringSubmatch(clientCred)

	auth := basicAuth()
	if len(result) == 2 {
		// Now verify the client creds
		return auth.authenticate(result[1])
	}

	return false
}

// authenticate performs basic authentication against provided clientCred
func (a *basic) authenticate(clientCred string) bool {
	username, passwd := base64decode(clientCred)

	credStruct, ok := credMap[username]
	if !ok {
		return false
	}

	if credStruct.salt == "" {
		if credStruct.password == passwd {
			return true
		}
	} else {
		if sha(credStruct.salt+passwd) == credStruct.password {
			return true
		}
	}

	return false
}

// base64encode base64 encode the username and password
func base64encode(username string, password string) string {
	auth := username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return encoded
}

// base64decode base64 decodes the creds
// returns the username, password
func base64decode(creds string) (string, string) {
	decoded, _ := base64.StdEncoding.DecodeString(creds)

	s := strings.Split(string(decoded), ":")
	return s[0], s[1]
}
