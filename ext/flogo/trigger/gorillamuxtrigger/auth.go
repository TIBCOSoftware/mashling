/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package gorillamuxtrigger

import (
	"encoding/base64"
	"net/http"
	"regexp"
)

type Auth interface {
	authenticate(clientCred string) bool
}

type basic struct {
	username string
	password string
}

const (
	basicAuthUsername = "basic_auth_username"
	basicAuthPassword = "basic_auth_password"
)

//IsAuthEnabled check if authentication is enabled
func isAuthEnabled(settings map[string]interface{}) bool {
	// Check if basic auth is in use
	if settings[basicAuthUsername] != "" && settings[basicAuthPassword] != "" {
		return true
	}

	return false
}

//Authenticate check if the request is allowed
func authenticate(r *http.Request, settings map[string]interface{}) bool {
	if settings[basicAuthUsername] != "" && settings[basicAuthPassword] != "" {
		// Authenticate using basic auth
		auth := basicAuth(settings[basicAuthUsername].(string), settings[basicAuthPassword].(string))

		clientCred := r.Header.Get("Authorization")
		re, err := regexp.Compile(`(?i)Basic (.*)`)
		if err != nil {
			log.Error(err)
			return false
		}
		result := re.FindStringSubmatch(clientCred)
		if len(result) == 2 {
			return auth.authenticate(result[1])
		}
	}

	return false
}

func basicAuth(username string, password string) Auth {
	return &basic{username, password}
}

//Authenticate performs basic authentication against provided clientCred
func (a *basic) authenticate(clientCred string) bool {
	selfEncoded := base64encode(a.username, a.password)
	if selfEncoded == clientCred {
		return true
	}

	return false

}

//base64encode base64 encode the username and password
func base64encode(username string, password string) string {
	auth := username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return encoded
}
