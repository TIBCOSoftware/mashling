package service

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/imdario/mergo"
)

const defaultTimeout = 5

// HTTP is an HTTP service.
type HTTP struct {
	Request HTTPRequest `json:"request"`
}

// HTTPRequest is an http service request.
type HTTPRequest struct {
	Path       string                 `json:"path"`
	PathParams map[string]interface{} `json:"pathParams"`
	Method     string                 `json:"method"`
	URL        string                 `json:"url"`
	Body       string                 `json:"body"`
	Headers    map[string]interface{} `json:"headers"`
	Query      map[string]string      `json:"query"`
	Timeout    int                    `json:"timeout"`
}

// HTTPResponse is an http service response.
type HTTPResponse struct {
	StatusCode int                    `json:"statusCode"`
	Body       interface{}            `json:"body"`
	Headers    map[string]interface{} `json:"headers"`
}

// Execute invokes this HTTP service.
func (h *HTTP) Execute(requestValues map[string]interface{}) (Response, error) {
	response := HTTPResponse{}
	request, err := h.createRequest(requestValues)
	if err != nil {
		return response, err
	}
	client := &http.Client{Timeout: time.Duration(request.Timeout) * time.Second}
	body := bytes.NewReader([]byte(request.Body))

	req, err := http.NewRequest(request.Method, request.CompleteURL(), body)
	if err != nil {
		return response, err
	}
	AddHeaders(req.Header, request.Headers)
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	response.StatusCode = resp.StatusCode
	response.Headers = DesliceValues(resp.Header)
	bodyReader := resp.Body
	if resp.ContentLength > 0 && resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, err = gzip.NewReader(bodyReader)
		if err != nil {
			return response, err
		}
	}
	defer bodyReader.Close()
	if resp.Header.Get("Content-Type") == "application/json" {
		err = json.NewDecoder(bodyReader).Decode(&response.Body)
	} else {
		respbody, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return response, err
		}
		response.Body = string(respbody)
	}
	return response, err
}

// InitializeHTTP initializes an HTTP service with provided settings.
func InitializeHTTP(settings map[string]interface{}) (httpService *HTTP, err error) {
	httpService = &HTTP{}
	req := HTTPRequest{}
	req.PathParams = make(map[string]interface{})
	req.Headers = make(map[string]interface{})
	req.Query = make(map[string]string)
	req.Timeout = defaultTimeout
	httpService.Request = req
	httpService.Request, err = httpService.createRequest(settings)
	return httpService, err
}

func (h *HTTP) createRequest(settings map[string]interface{}) (HTTPRequest, error) {
	request := HTTPRequest{}
	for k, v := range settings {
		switch k {
		case "url":
			url, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for url")
			}
			request.URL = url
		case "method":
			method, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for method")
			}
			request.Method = method
		case "path":
			path, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for path")
			}
			request.Path = path
		case "headers":
			headers, ok := v.(map[string]interface{})
			if !ok {
				return request, errors.New("invalid type for headers")
			}
			request.Headers = headers
			if err := mergo.Merge(&request.Headers, h.Request.Headers); err != nil {
				return request, errors.New("unable to merge header values")
			}
		case "query":
			query, ok := v.(map[string]string)
			if !ok {
				return request, errors.New("invalid type for query")
			}
			request.Query = query
		case "pathParams":
			pathParams, ok := v.(map[string]interface{})
			if !ok {
				return request, errors.New("invalid type for pathParams")
			}
			request.PathParams = pathParams
			if err := mergo.Merge(&request.PathParams, h.Request.PathParams); err != nil {
				return request, errors.New("unable to merge pathParams values")
			}
		case "timeout":
			timeout, ok := v.(int)
			if !ok {
				return request, errors.New("invalid type for timeout")
			}
			request.Timeout = timeout
		default:
			// ignore and move on.
		}
	}
	if err := mergo.Merge(&request, h.Request); err != nil {
		return request, errors.New("unable to merge request values")
	}
	return request, nil
}

// AddHeaders adds the headers in headers to headers.
func AddHeaders(h http.Header, headers map[string]interface{}) {
	for key, value := range headers {
		switch value := value.(type) {
		case string:
			h.Add(key, value)
		case []interface{}:
			for _, v := range value {
				AddHeaders(h, map[string]interface{}{key: v})
			}
		}
	}
}

// DesliceValues is used to collapse single value string slices from map values.
func DesliceValues(slice map[string][]string) map[string]interface{} {
	desliced := make(map[string]interface{})
	for k, v := range slice {
		if len(v) == 1 {
			desliced[k] = v[0]
		} else {
			desliced[k] = v
		}
	}
	return desliced
}

// CompleteURL returns the full URL including query params
func (h *HTTPRequest) CompleteURL() string {
	if h.Path != "" {
		if strings.HasPrefix(h.Path, "/") || strings.HasSuffix(h.URL, "/") {
			h.URL = h.URL + h.Path
		} else {
			h.URL = h.URL + "/" + h.Path
		}
	}
	if len(h.PathParams) > 0 {
		for k, v := range h.PathParams {
			h.URL = strings.Replace(h.URL, fmt.Sprintf(":%s", k), fmt.Sprintf("%v", v), -1)
		}
	}
	if len(h.Query) > 0 {
		params := url.Values{}
		for k, v := range h.Query {
			params.Add(k, v)
		}
		if strings.Contains(h.URL, "?") {
			return fmt.Sprintf("%s&%s", h.URL, params.Encode())
		}
		return fmt.Sprintf("%s?%s", h.URL, params.Encode())
	}
	return h.URL
}
