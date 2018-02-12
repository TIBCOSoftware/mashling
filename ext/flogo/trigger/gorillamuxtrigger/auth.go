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

	ldap "github.com/jtblin/go-ldap-client"
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

// Auth allows a client to be authenticated
type Auth interface {
	authenticate(clientCred string) bool
}

type basic struct {
	ldapHost         string
	ldapBase         string
	ldapBindDN       string
	ldapBindPassword string
	ldapUserFilter   string
	ldapGroupFilter  string
}

const (
	basicAuthFile    = "basicAuthFile"
	ldapHost         = "ldapHost"
	ldapBase         = "ldapBase"
	ldapBindDN       = "ldapBindDN"
	ldapBindPassword = "ldapBindPassword"
	ldapUserFilter   = "ldapUserFilter"
	ldapGroupFilter  = "ldapGroupFilter"
)

var mu sync.Mutex
var credMap map[string]hashedCred

func basicAuth(settings map[string]interface{}) Auth {
	auth := &basic{
		ldapUserFilter:  "(uid=%s)",
		ldapGroupFilter: "(memberUid=%s)",
	}
	if value, ok := settings[ldapHost]; ok {
		auth.ldapHost = value.(string)
	}
	if value, ok := settings[ldapBase]; ok {
		auth.ldapBase = value.(string)
	}
	if value, ok := settings[ldapBindDN]; ok {
		auth.ldapBindDN = value.(string)
	}
	if value, ok := settings[ldapBindPassword]; ok {
		auth.ldapBindPassword = value.(string)
	}
	if value, ok := settings[ldapUserFilter]; ok {
		auth.ldapUserFilter = value.(string)
	}
	if value, ok := settings[ldapGroupFilter]; ok {
		auth.ldapGroupFilter = value.(string)
	}
	return auth
}

// isAuthEnabled check if authentication is enabled
func isAuthEnabled(settings map[string]interface{}) bool {
	// Check if basic auth is in use
	_, hasBasicAuthFile := settings[basicAuthFile]
	_, hasLDAPHost := settings[ldapHost]
	if !hasBasicAuthFile && !hasLDAPHost {
		return false
	}

	return true
}

// setupAuth setups up authentication.
// This should be called to load the creds into a map once.
func setupAuth(settings map[string]interface{}) {
	_, hasLDAPHost := settings[ldapHost]
	if hasLDAPHost {
		return
	}

	err := loadCreds(settings[basicAuthFile].(string))
	if err != nil {
		log.Error(err)
		panic("Unable to load creds file")
	}
}

// loadCreds loads the file with the credentials.
// The file may contain lines of the form username:password
// or username:salt:sha256(salt + password)
func loadCreds(pathToFile string) error {
	mu.Lock()
	if credMap != nil {
		mu.Unlock()
		return nil
	}

	credMap = make(map[string]hashedCred)

	csvFile, err := os.Open(pathToFile)
	if err != nil {
		mu.Unlock()
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
	mu.Unlock()

	return nil
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
	auth := basicAuth(settings)
	if len(result) == 2 {
		// Now verify the client creds
		return auth.authenticate(result[1])
	}

	return false
}

// authenticate performs basic authentication against provided clientCred
func (a *basic) authenticate(clientCred string) bool {
	username, passwd := base64decode(clientCred)

	if a.ldapHost != "" {
		client := &ldap.LDAPClient{
			Base:         a.ldapBase,
			Host:         a.ldapHost,
			Port:         389,
			BindDN:       a.ldapBindDN,
			BindPassword: a.ldapBindPassword,
			UserFilter:   a.ldapUserFilter,
			GroupFilter:  a.ldapGroupFilter,
			Attributes:   []string{"givenName", "sn", "mail", "uid"},
		}
		defer client.Close()

		ok, user, err := client.Authenticate(username, passwd)
		if err != nil {
			log.Errorf("Error authenticating user %s: %+v", username, err)
		}
		log.Infof("Authenticated User: %+v", user)

		return ok
	}

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

// sha applies sha256 to the string
func sha(str string) string {
	bytes := []byte(str)

	h := sha256.New()
	h.Write(bytes)
	code := h.Sum(nil)
	codestr := hex.EncodeToString(code)

	return codestr
}
