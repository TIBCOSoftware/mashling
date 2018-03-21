// Copyright (c) 2015 TIBCO Software Inc.
// All Rights Reserved.

package cors

import (
	"os"
)

// List of constants default values that can be overriden by environment variables
const (
	CORS_ALLOW_ORIGIN_KEY          string = "CORS_ALLOW_ORIGIN"
	CORS_ALLOW_ORIGIN_DEFAULT      string = "*"
	CORS_ALLOW_METHODS_KEY         string = "CORS_ALLOW_METHODS"
	CORS_ALLOW_METHODS_DEFAULT     string = "POST, GET, OPTIONS, PUT, DELETE, PATCH"
	CORS_ALLOW_HEADERS_KEY         string = "CORS_ALLOW_HEADERS"
	CORS_EXPOSE_HEADERS_KEY        string = "CORS_EXPOSE_HEADERS"
	CORS_ALLOW_HEADERS_DEFAULT     string = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Requested-With, Accept, Accept-Language"
	CORS_EXPOSE_HEADERS_DEFAULT    string = ""
	CORS_ALLOW_CREDENTIALS_KEY     string = "CORS_ALLOW_CREDENTIALS"
	CORS_ALLOW_CREDENTIALS_DEFAULT string = "false"
	CORS_MAX_AGE_KEY               string = "CORS_MAX_AGE"
	CORS_MAX_AGE_DEFAULT           string = ""
)

//GetCorsAllowOrigin get the value for CORS 'AllowOrigin' param from environment variable and the default BS_CORS_ALLOW_ORIGIN_DEFAULT will be used if not found
func GetCorsAllowOrigin(prefix string) string {
	envalue := os.Getenv(prefix + CORS_ALLOW_ORIGIN_KEY)
	if envalue == "" {
		return CORS_ALLOW_ORIGIN_DEFAULT
	} else {
		return envalue
	}
}

//GetCorsAllowMethods get the allowed method for CORS from environment variable and the default BS_CORS_ALLOW_METHODS_DEFAULT will be used if not found
func GetCorsAllowMethods(prefix string) string {
	envalue := os.Getenv(prefix + CORS_ALLOW_METHODS_KEY)
	if envalue == "" {
		return CORS_ALLOW_METHODS_DEFAULT
	} else {
		return envalue
	}
}

//GetCorsAllowHeaders get the value for CORS 'AllowHeaders' param from environment variable and the default BS_CORS_ALLOW_HEADERS_DEFAULT will be used if not found
func GetCorsAllowHeaders(prefix string) string {
	envalue := os.Getenv(prefix + CORS_ALLOW_HEADERS_KEY)
	if envalue == "" {
		return CORS_ALLOW_HEADERS_DEFAULT
	} else {
		return envalue
	}
}

//GetCorsExposeHeaders get the value for CORS 'ExposeHeaders' param from environment variable and the default BS_CORS_EXPOSE_HEADERS_DEFAULT will be used if not found
func GetCorsExposeHeaders(prefix string) string {
	envalue := os.Getenv(prefix + CORS_EXPOSE_HEADERS_KEY)
	if envalue == "" {
		return CORS_EXPOSE_HEADERS_DEFAULT
	} else {
		return envalue
	}
}

//GetCorsAllowCredentials get the value for CORS 'AllowCredentials' param from environment variable and the default BS_CORS_ALLOW_CREDENTIALS_DEFAULT will be used if not found
func GetCorsAllowCredentials(prefix string) string {
	envalue := os.Getenv(prefix + CORS_ALLOW_CREDENTIALS_KEY)
	if envalue == "" {
		return CORS_ALLOW_CREDENTIALS_DEFAULT
	} else {
		return envalue
	}
}

//GetCorsMaxAge get the value for CORS 'Max Age' param from environment variable and the default BS_CORS_ALLOW_CREDENTIALS_DEFAULT will be used if not found
func GetCorsMaxAge(prefix string) string {
	envalue := os.Getenv(prefix + CORS_MAX_AGE_KEY)
	if envalue == "" {
		return CORS_MAX_AGE_DEFAULT
	} else {
		return envalue
	}
}
