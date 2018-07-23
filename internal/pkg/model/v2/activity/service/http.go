package service

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/imdario/mergo"

	"github.com/TIBCOSoftware/mashling/lib/util"
)

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodPATCH  = "PATCH"
	methodDELETE = "DELETE"

	contentTypeApplicationJSON = "application/json; charset=UTF-8"

	defaultTimeout = 5
)

// HTTP is an HTTP service.
type HTTP struct {
	Request  HTTPRequest  `json:"request"`
	Response HTTPResponse `json:"response"`
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
	NetError   string                 `json:"netError"`
	StatusCode int                    `json:"statusCode"`
	Body       interface{}            `json:"body"`
	Headers    map[string]interface{} `json:"headers"`
}

// Execute invokes this HTTP service.
func (h *HTTP) Execute() (err error) {
	h.Response = HTTPResponse{}
	if h.Request.Timeout == 0 {
		h.Request.Timeout = defaultTimeout
	}
	client := &http.Client{Timeout: time.Duration(h.Request.Timeout) * time.Second}
	body := bytes.NewReader([]byte(h.Request.Body))

	req, err := http.NewRequest(h.Request.Method, h.Request.CompleteURL(), body)
	if err != nil {
		return err
	}
	AddHeaders(req.Header, h.Request.Headers)

	resp, err := client.Do(req)
	if err != nil {
		if netError, ok := err.(net.Error); ok {
			h.Response.NetError = netError.Error()
			return nil
		}
		return err
	}
	h.Response.StatusCode = resp.StatusCode
	h.Response.Headers = DesliceValues(resp.Header)
	bodyReader := resp.Body
	if resp.ContentLength > 0 && resp.Header.Get("Content-Encoding") == "gzip" {
		bodyReader, err = gzip.NewReader(bodyReader)
		if err != nil {
			return err
		}
	}
	defer bodyReader.Close()
	respbody, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return err
	}
	contentType := resp.Header.Get("Content-Type")
	err = util.Unmarshal(contentType, respbody, &h.Response.Body)
	if err != nil {
		return err
	}

	return err
}

// InitializeHTTP initializes an HTTP service with provided settings.
func InitializeHTTP(settings map[string]interface{}) (httpService *HTTP, err error) {
	httpService = &HTTP{}
	req := HTTPRequest{}
	req.PathParams = make(map[string]interface{})
	req.Headers = make(map[string]interface{})
	req.Query = make(map[string]string)
	httpService.Request = req
	err = httpService.setRequestValues(settings)
	return httpService, err
}

// UpdateRequest updates a request on an existing HTTP service instance with new values.
func (h *HTTP) UpdateRequest(values map[string]interface{}) (err error) {
	return h.setRequestValues(values)
}

func (h *HTTP) setRequestValues(settings map[string]interface{}) (err error) {
	var body interface{}
	for k, v := range settings {
		switch k {
		case "url":
			url, ok := v.(string)
			if !ok {
				return errors.New("invalid type for url")
			}
			h.Request.URL = url
		case "method":
			method, ok := v.(string)
			if !ok {
				return errors.New("invalid type for method")
			}
			h.Request.Method = method
		case "path":
			path, ok := v.(string)
			if !ok {
				return errors.New("invalid type for path")
			}
			h.Request.Path = path
		case "headers":
			headers, ok := v.(map[string]interface{})
			if !ok {
				return errors.New("invalid type for headers")
			}
			if err := mergo.Merge(&h.Request.Headers, headers, mergo.WithOverride); err != nil {
				return errors.New("unable to merge header values")
			}
		case "query":
			query, ok := v.(map[string]string)
			if !ok {
				return errors.New("invalid type for query")
			}
			h.Request.Query = query
		case "pathParams":
			pathParams, ok := v.(map[string]interface{})
			if !ok {
				return errors.New("invalid type for pathParams")
			}
			if err := mergo.Merge(&h.Request.PathParams, pathParams, mergo.WithOverride); err != nil {
				return errors.New("unable to merge pathParams values")
			}
		case "body":
			body = v
		default:
			// ignore and move on.
		}
	}
	if body != nil {
		if method := h.Request.Method; method == methodPOST || method == methodPUT || method == methodPATCH {
			contentType := contentTypeApplicationJSON
			if object, ok := body.(map[string]interface{}); ok {
				if mime, ok := object[util.MetaMIME]; ok {
					if s, ok := mime.(string); ok {
						contentType = s
					}
				}
			}
			if _, ok := h.Request.Headers["Content-Type"]; !ok {
				h.Request.Headers["Content-Type"] = contentType
			}

			data, err := util.Marshal(body)
			if err != nil {
				return err
			}
			h.Request.Body = string(data)
		}
	}
	return nil
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
