/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ApiUser struct {
	username     string
	password     string
	apiKey       string
	apiSecretKey string
	uuid         string
	portal       string
	noop         bool
	skipVerify   bool
}

const (
	masheryUri   = "https://api.mashery.com"
	restUri      = "/v3/rest/"
	transformUri = "transform"
	accessToken  = "access_token"
)

type Responder func(*http.Request) (*http.Response, error)

type NopTransport struct {
	responders map[string]Responder
}

var DefaultNopTransport = &NopTransport{}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

func init() {
	DefaultNopTransport.responders = make(map[string]Responder)
}

func (n *NopTransport) RegisterResponder(method, url string, responder Responder) {
	n.responders[method+" "+url] = responder
}

func (n *NopTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.String()

	// Scan through the responders
	for k, r := range n.responders {
		if k != key {
			continue
		}
		return r(req)
	}

	return nil, errors.New("No responder found")
}

func RegisterResponder(method, url string, responder Responder) {
	DefaultNopTransport.RegisterResponder(method, url, responder)
}

func newHttp(nop bool, skipVerify bool) *http.Client {
	client := &http.Client{}
	if nop {
		client.Transport = DefaultNopTransport
	}

	if skipVerify {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	return client
}

func setContentType(r *http.Request) {
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
}

func setOauthToken(r *http.Request, oauthToken string) {
	r.Header.Add("Authorization", "Bearer "+oauthToken)
}

func readBody(body io.Reader) ([]byte, error) {
	bodyText, err := ioutil.ReadAll(body)
	if err != nil {
		return bodyText, err
	}
	return bodyText, nil
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) CreateAPI(tfSwaggerDoc string, oauthToken string) (string, error) {
	return user.CreateUpdate("POST", "services", "", tfSwaggerDoc, oauthToken)
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) Create(resource string, fields string, content string, oauthToken string) (string, error) {
	return user.CreateUpdate("POST", resource, fields, content, oauthToken)
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) CreateUpdate(method string, resource string, fields string, content string, oauthToken string) (string, error) {
	fullUri := masheryUri + restUri + resource
	if fields != "" {
		fullUri = fullUri + "?fields=" + fields
	}
	client := newHttp(user.noop, user.skipVerify)
	r, _ := http.NewRequest(method, fullUri, bytes.NewReader([]byte(content)))
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	s := string(bodyText)
	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("Unable to create the api: status code %v", resp.StatusCode)
	}

	return s, err
}

// Read fetch data
func (user *ApiUser) Read(resource string, filter string, fields string, oauthToken string) (string, error) {

	fullUri := masheryUri + restUri + resource
	if fields != "" && filter == "" {
		fullUri = fullUri + "?fields=" + fields
	} else if fields == "" && filter != "" {
		fullUri = fullUri + "?filter=" + filter
	} else {
		fullUri = fullUri + "?fields=" + fields + "&filter=" + filter
	}

	client := newHttp(user.noop, user.skipVerify)

	r, _ := http.NewRequest("GET", masheryUri+restUri+resource+"?filter="+filter+"&fields="+fields, nil)
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	s := string(bodyText)
	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("Unable to create the api: status code %v", resp.StatusCode)
	}

	return s, err
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) Update(resource string, fields string, content string, oauthToken string) (string, error) {
	return user.CreateUpdate(http.MethodPut, resource, fields, content, oauthToken)
}

// TransformSwagger sends the swagger doc to Mashery API to be
// transformed into the target format.
func (user *ApiUser) TransformSwagger(swaggerDoc string, sourceFormat string, targetFormat string, oauthToken string) (string, error) {
	// New client
	client := newHttp(user.noop, user.skipVerify)

	v := url.Values{}
	v.Set("sourceFormat", sourceFormat)
	v.Add("targetFormat", targetFormat)
	v.Add("publicDomain", user.portal)

	r, _ := http.NewRequest("POST", masheryUri+restUri+transformUri+"?"+v.Encode(), bytes.NewReader([]byte(swaggerDoc)))
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	if bodyText, err := readBody(resp.Body); err == nil {
		if resp.StatusCode != http.StatusOK {
			return string(bodyText), fmt.Errorf("Unable to transform the swagger doc: status code %v", resp.StatusCode)
		}
		return string(bodyText), nil
	} else {
		return string(bodyText), err
	}
}

// FetchOAuthToken exchanges the creds for an OAuth token
func (user *ApiUser) FetchOAuthToken() (string, error) {
	// New client
	client := newHttp(user.noop, user.skipVerify)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", user.username)
	data.Set("password", user.password)
	data.Set("scope", user.uuid)

	r, _ := http.NewRequest("POST", masheryUri+"/v3/token", strings.NewReader(data.Encode()))
	r.SetBasicAuth(user.apiKey, user.apiSecretKey)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "*/*")

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	if bodyText, err := readBody(resp.Body); err == nil {
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("Unable to get the OAuth token: status code (%v), message (%v)", resp.StatusCode, string(bodyText))
		}

		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(string(bodyText)), &dat); err != nil {
			return "", errors.New("Unable to unmarshal JSON")
		}

		accessToken, ok := dat[accessToken].(string)
		if !ok {
			return "", errors.New("Invalid json. Expected a field with access_token")
		}

		return accessToken, nil
	} else {
		return string(bodyText), err
	}
}
