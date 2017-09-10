// Copyright (c) 2015 TIBCO Software Inc.
// All Rights Reserved.

//Cors package to validate CORS requests and provide the correct headers in the response
package cors

import (
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ORIGIN_HEADER                           string = "Origin"
	ACCESS_CONTROL_REQUEST_METHOD_HEADER    string = "Access-Control-Request-Method"
	ACCESS_CONTROL_REQUEST_HEADER_HEADER    string = "Access-Control-Request-Headers"
	ACCESS_CONTROL_ALLOW_ORIGIN_HEADER      string = "Access-Control-Allow-Origin"
	ACCESS_CONTROL_ALLOW_METHODS_HEADER     string = "Access-Control-Allow-Methods"
	ACCESS_CONTROL_ALLOW_HEADERS_HEADER     string = "Access-Control-Allow-Headers"
	ACCESS_CONTROL_EXPOSE_HEADERS_HEADER    string = "Access-Control-Expose-Headers"
	ACCESS_CONTROL_ALLOW_CREDENTIALS_HEADER string = "Access-Control-Allow-Credentials"
	ACCESS_CONTROL_MAX_AGE_HEADER           string = "Access-Control-Max-Age"
)

// CORS interface
type Cors interface {
	// HandlePreflight Handles the preflight OPTIONS request
	HandlePreflight(w http.ResponseWriter, r *http.Request)
	// WriteCorsActualRequestHeaders writes the needed request headers for the CORS support
	WriteCorsActualRequestHeaders(w http.ResponseWriter)
}

type cors struct {
	// Prefix used for the CORS environment variables
	Prefix string

	log logger.Logger
}

// make sure that the cors implements the Cors interface
var _ Cors = (*cors)(nil)

//Cors constructor
func New(prefix string, log logger.Logger) Cors {
	return cors{Prefix: prefix, log: log}
}

// HandlePreflight Handles the cors preflight request setting the right headers and responding to the request
func (c cors) HandlePreflight(w http.ResponseWriter, r *http.Request) {
	// Check if it has Origin Header
	hasOrigin := HasOriginHeader(r)
	if hasOrigin == false {
		c.log.Info("Invalid CORS preflight request, no Origin header found")
		writeInvalidPreflightResponse(w)
		return
	}

	// Check Access-Control-Request-Method header
	requestMethodHeader := r.Header.Get(ACCESS_CONTROL_REQUEST_METHOD_HEADER)
	if isValidAccessControlMethod(requestMethodHeader, c.Prefix, c.log) != true {
		// Invalid Access Control Method
		writeInvalidPreflightResponse(w)
		return
	}

	// Check Access-Control-Allow-Headers header
	requestHeadersHeader := r.Header.Get(ACCESS_CONTROL_REQUEST_HEADER_HEADER)
	if isValidAccessControlHeaders(requestHeadersHeader, c.Prefix, c.log) != true {
		// Invalid Access Control Header
		writeInvalidPreflightResponse(w)
		return
	}

	writeValidPreflightResponse(w, c)
}

// HasOriginHeader returns true if the request has Origin header, false otherwise
func HasOriginHeader(r *http.Request) bool {
	h := r.Header.Get(ORIGIN_HEADER)
	if h == "" {
		return false
	}
	return true
}

// Check if the method name is valid and allowed by the environment variable
func isValidAccessControlMethod(methodName string, prefix string, log logger.Logger) bool {
	if methodName == "" {
		log.Infof("Invalid Access Control Method for preflight request: '%s'", methodName)
		return false
	}
	allowedMethodsEnv := GetCorsAllowMethods(prefix)
	allowedMethods := strings.Split(allowedMethodsEnv, ",")
	log.Debugf("Allowed Methods '%s'", allowedMethods)
	for i := range allowedMethods {
		if strings.ToLower(strings.TrimSpace(allowedMethods[i])) == strings.ToLower(strings.TrimSpace(methodName)) {
			return true
		}
	}
	log.Infof("Invalid Access Control Method for preflight request: '%s'", methodName)
	return false
}

// Check if the headers are valid and allowed by the environment variable
func isValidAccessControlHeaders(headersStr string, prefix string, log logger.Logger) bool {
	if headersStr == "" {
		return true
	}
	allowedHeadersEnv := GetCorsAllowHeaders(prefix)
	allowedHeaders := strings.Split(allowedHeadersEnv, ",")

	// Create a map for faster lookup
	allowedHeadersMap := make(map[string]struct{}, len(allowedHeaders))
	for _, s := range allowedHeaders {
		allowedHeadersMap[strings.ToLower(strings.TrimSpace(s))] = struct{}{}
	}

	headers := strings.Split(headersStr, ",")

	for i := range headers {
		_, ok := allowedHeadersMap[strings.ToLower(strings.TrimSpace(headers[i]))]
		if ok == false {
			log.Infof("Invalid Access Control Header for preflight request: '%s'", strings.TrimSpace(headers[i]))
			return false
		}
	}
	return true
}

// Writes invalid preflight response
func writeInvalidPreflightResponse(w http.ResponseWriter) {
	// Write 200 but no CORS header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Writes valid preflight response
func writeValidPreflightResponse(w http.ResponseWriter, c cors) {
	// Write 200 with CORS headers
	writeCorsPreflightHeaders(w, c)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Writes the CORS preflight request headers (origin and credential)
func writeCorsPreflightHeaders(w http.ResponseWriter, c cors) {
	c.WriteCorsActualRequestHeaders(w)
	w.Header().Set(ACCESS_CONTROL_ALLOW_METHODS_HEADER, GetCorsAllowMethods(c.Prefix))
	w.Header().Set(ACCESS_CONTROL_ALLOW_HEADERS_HEADER, GetCorsAllowHeaders(c.Prefix))
	w.Header().Set(ACCESS_CONTROL_EXPOSE_HEADERS_HEADER, GetCorsExposeHeaders(c.Prefix))
	maxAge := GetCorsMaxAge(c.Prefix)
	if maxAge != "" {
		w.Header().Set(ACCESS_CONTROL_MAX_AGE_HEADER, maxAge)
	}
}

// Writes the CORS actual request headers (origin and credential)
func (c cors) WriteCorsActualRequestHeaders(w http.ResponseWriter) {
	w.Header().Set(ACCESS_CONTROL_ALLOW_ORIGIN_HEADER, GetCorsAllowOrigin(c.Prefix))
	allowCredentials := GetCorsAllowCredentials(c.Prefix)
	if strings.TrimSpace(allowCredentials) == "true" {
		w.Header().Set(ACCESS_CONTROL_ALLOW_CREDENTIALS_HEADER, strings.TrimSpace(allowCredentials))
	}
}
