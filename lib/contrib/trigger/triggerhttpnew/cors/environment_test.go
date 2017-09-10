// Copyright (c) 2015 TIBCO Software Inc.
// All Rights Reserved.
package cors

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Test GetCorsAllowOrigin method
func TestGetCorsAllowOriginOk(t *testing.T) {
	allowOrigin := GetCorsAllowOrigin(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, CORS_ALLOW_ORIGIN_DEFAULT, allowOrigin, "Allow Origin should be default value")
}

// Test GetCorsAllowOrigin method
func TestGetCorsAllowOriginOkModified(t *testing.T) {
	previous := os.Getenv(TEST_CORS_PREFIX + CORS_ALLOW_ORIGIN_KEY)
	defer os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_ORIGIN_KEY, previous)

	newValue := "fooAllowedOrigin"

	// Change value
	os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_ORIGIN_KEY, newValue)

	allowOrigin := GetCorsAllowOrigin(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, newValue, allowOrigin, "Allow Origin should be "+newValue)

}

// Test GetCorsAllowMethods method
func TestGetCorsAllowMethodsOk(t *testing.T) {
	envValue := GetCorsAllowMethods(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, CORS_ALLOW_METHODS_DEFAULT, envValue, "Allow Method should be default value")
}

// Test GetCorsAllowOrigin method
func TestGetCorsAllowMethodsOkModified(t *testing.T) {
	previous := os.Getenv(TEST_CORS_PREFIX + CORS_ALLOW_METHODS_KEY)
	defer os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_METHODS_KEY, previous)

	newValue := "fooAllowedMethods"

	// Change value
	os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_METHODS_KEY, newValue)

	envValue := GetCorsAllowMethods(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, newValue, envValue, "Allow Methods should be "+newValue)
}

// Test GetCorsAllowHeaders method
func TestGetCorsAllowHeadersOk(t *testing.T) {
	envValue := GetCorsAllowHeaders(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, CORS_ALLOW_HEADERS_DEFAULT, envValue, "Allow Headers should be default value")
}

// Test GetCorsAllowHeaders method
func TestGetCorsAllowHeadersOkModified(t *testing.T) {
	previous := os.Getenv(TEST_CORS_PREFIX + CORS_ALLOW_HEADERS_KEY)
	defer os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_HEADERS_KEY, previous)

	newValue := "fooAllowedHeaders"

	// Change value
	os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_HEADERS_KEY, newValue)

	envValue := GetCorsAllowHeaders(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, newValue, envValue, "Allow Headers should be "+newValue)
}

// Test GetCorsAllowCredentials method
func TestGetCorsAllowCredentialsOk(t *testing.T) {
	envValue := GetCorsAllowCredentials(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, CORS_ALLOW_CREDENTIALS_DEFAULT, envValue, "Allow Credentials should be default value")
}

// Test GetCorsAllowCredentials method
func TestGetCorsAllowCredentialsOkModified(t *testing.T) {
	previous := os.Getenv(TEST_CORS_PREFIX + CORS_ALLOW_CREDENTIALS_KEY)
	defer os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_CREDENTIALS_KEY, previous)

	newValue := "true"

	// Change value
	os.Setenv(TEST_CORS_PREFIX+CORS_ALLOW_CREDENTIALS_KEY, newValue)

	envValue := GetCorsAllowCredentials(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, newValue, envValue, "Allow Credentials should be "+newValue)
}

// Test GetCorsMaxAge method
func TestGetCorsMaxAgeOk(t *testing.T) {
	envValue := GetCorsMaxAge(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, CORS_MAX_AGE_DEFAULT, envValue, "Max Age should be default value")
}

// Test GetCorsAllowCredentials method
func TestGetCorsMaxAgeOkModified(t *testing.T) {
	previous := os.Getenv(TEST_CORS_PREFIX + CORS_MAX_AGE_KEY)
	defer os.Setenv(TEST_CORS_PREFIX+CORS_MAX_AGE_KEY, previous)

	newValue := "21"

	// Change value
	os.Setenv(TEST_CORS_PREFIX+CORS_MAX_AGE_KEY, newValue)

	envValue := GetCorsMaxAge(TEST_CORS_PREFIX)

	// assert Success
	assert.Equal(t, newValue, envValue, "Max Age should be "+newValue)
}
