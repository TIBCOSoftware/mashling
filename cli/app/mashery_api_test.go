package app

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFetchOAuthToken(t *testing.T) {
	fn := func(r *http.Request) (*http.Response, error) {
		buf := `{"access_token": "123"}`
		resp := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(buf)),
		}
		resp.StatusCode = http.StatusOK
		return &resp, nil
	}

	DefaultNopTransport.RegisterResponder("POST", masheryUri+"/v3/token", fn)

	u := ApiUser{"foo", "pass", "key", "secretKey", "uuid", "api.example.com", true, false}
	token, err := u.FetchOAuthToken()

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, "123", token)
}

func TestFetchOAuthTokenToReturnError(t *testing.T) {
	fn := func(r *http.Request) (*http.Response, error) {
		buf := `{"access_tokenXXX": "123"}`
		resp := http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(buf)),
		}
		resp.StatusCode = http.StatusOK
		return &resp, nil
	}

	DefaultNopTransport.RegisterResponder("POST", masheryUri+"/v3/token", fn)

	u := ApiUser{"foo", "pass", "key", "secretKey", "uuid", "api.example.com", true, false}
	_, err := u.FetchOAuthToken()
	assert.Error(t, err)
}
