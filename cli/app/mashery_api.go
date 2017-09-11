package app

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type ApiUser struct {
	username     string
	password     string
	apiKey       string
	apiSecretKey string
	uuid         string
}

const MASHERY_URI = "https://api.mashery.com"

func (user *ApiUser) FetchOAuthToken() (string, error) {
	fmt.Println("Trying to fetch with: ", user.username)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Add("username", user.username)
	data.Add("password", user.password)
	data.Add("scope", user.uuid)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", MASHERY_URI+"/v3/token", bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	r.SetBasicAuth(user.apiKey, user.apiSecretKey)

	resp, err := client.Do(r)
	fmt.Println(resp.Status)
	if err != nil {
		return "", err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	s := string(bodyText)

	if resp.StatusCode != 200 {
		return "", errors.New("Unable to get the OAuth token")
	}

	return s, err
}
