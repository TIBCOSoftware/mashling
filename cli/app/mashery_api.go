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
}

const (
	masheryUri   = "https://api.mashery.com"
	restUri      = "/v3/rest/"
	transformUri = "/v3/rest/transform"
	accessToken  = "access_token"
)

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// CreateAPI sends the transformed swagger doc to the Mashery API.
func (user *ApiUser) Create(resource string, tfSwaggerDoc string, oauthToken string) (string, error) {

	client := &http.Client{}
	r, _ := http.NewRequest("POST", masheryUri+restUri+resource, bytes.NewReader([]byte(tfSwaggerDoc)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Authorization", "Bearer "+oauthToken)

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

	client := &http.Client{}
	r, _ := http.NewRequest("GET", masheryUri+restUri+resource+"?filter="+filter+"&fields="+fields, nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Authorization", "Bearer "+oauthToken)

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

// TransformSwagger sends the swagger doc to Mashery API to be
// transformed into the masheryapi format.
func (user *ApiUser) TransformSwagger(swaggerDoc string, sourceFormat string, targetFormat string, oauthToken string) (string, error) {
	v := url.Values{}
	v.Set("sourceFormat", sourceFormat)
	v.Add("targetFormat", targetFormat)
	v.Add("publicDomain", user.portal)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", masheryUri+transformUri+"?"+v.Encode(), bytes.NewReader([]byte(swaggerDoc)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Authorization", "Bearer "+oauthToken)

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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unable to transform the swagger doc: status code %v", resp.StatusCode)
	}

	return string(bodyText), err
}

// FetchOAuthToken exchanges the creds for an OAuth token
func (user *ApiUser) FetchOAuthToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", user.username)
	data.Set("password", user.password)
	data.Set("scope", user.uuid)

	client := &http.Client{}
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

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unable to get the OAuth token: status code %v", resp.StatusCode)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(string(bodyText)), &dat); err != nil {
		return "", errors.New("Unable to unmarshal JSON")
	}

	return dat[accessToken].(string), err
}
