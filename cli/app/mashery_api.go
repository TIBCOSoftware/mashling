/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bytes"
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
}

const (
	masheryUri   = "https://api.mashery.com"
	servicesUri  = "/v3/rest/services"
	transformUri = "/v3/rest/transform"
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

func newHttp(nop bool) *http.Client {
	client := &http.Client{}
	if nop {
		client.Transport = DefaultNopTransport
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
	// New client
	client := newHttp(user.noop)

	r, _ := http.NewRequest("POST", masheryUri+servicesUri, bytes.NewReader([]byte(tfSwaggerDoc)))
	setContentType(r)
	setOauthToken(r, oauthToken)

	resp, err := client.Do(r)
	if err != nil {
		return "", err
	} else {
		defer resp.Body.Close()
	}

	if bodyText, err := readBody(resp.Body); err == nil {
		s := string(bodyText)
		if resp.StatusCode != http.StatusOK {
			return s, fmt.Errorf("Unable to create the api: status code %v", resp.StatusCode)
		}

		return string(bodyText), nil
	} else {
		return string(bodyText), err
	}
}

// TransformSwagger sends the swagger doc to Mashery API to be
// transformed into the masheryapi format.
func (user *ApiUser) TransformSwagger(swaggerDoc string, oauthToken string) (string, error) {
	// New client
	client := newHttp(user.noop)

	v := url.Values{}
	v.Set("sourceFormat", "swagger2")
	v.Add("targetFormat", "masheryapi")
	v.Add("publicDomain", user.portal)

	r, _ := http.NewRequest("POST", masheryUri+transformUri+"?"+v.Encode(), bytes.NewReader([]byte(swaggerDoc)))
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
	client := newHttp(user.noop)

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
