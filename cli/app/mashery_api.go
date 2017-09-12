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
	servicesUri  = "/v3/rest/services"
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

func (user *ApiUser) CreateAPI(tfSwaggerDoc string, oauthToken string) (string, error) {

	client := &http.Client{}
	r, _ := http.NewRequest("POST", masheryUri+servicesUri, bytes.NewReader([]byte(tfSwaggerDoc)))
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

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Unable to create the api: status code %g", resp.StatusCode)
	}

	return string(bodyText), err
}

func (user *ApiUser) TransformSwagger(swaggerDoc string, oauthToken string) (string, error) {
	v := url.Values{}
	v.Set("sourceFormat", "swagger2")
	v.Add("targetFormat", "masheryapi")
	v.Add("publicDomain", user.portal)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", masheryUri + transformUri + "?" + v.Encode(), bytes.NewReader([]byte(swaggerDoc)))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Authorization", "Bearer " + oauthToken)

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

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Unable to transform the swagger doc: status code %g", resp.StatusCode)
	}

	return string(bodyText), err
}

func (user *ApiUser) FetchOAuthToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", user.username)
	data.Set("password", user.password)
	data.Set("scope", user.uuid)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", masheryUri + "/v3/token", strings.NewReader(data.Encode()))
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

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Unable to get the OAuth token: status code %g", resp.StatusCode)
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(string(bodyText)), &dat); err != nil {
		return "", errors.New("Unable to unmarshal JSON")
	}

	return dat[accessToken].(string), err
}
