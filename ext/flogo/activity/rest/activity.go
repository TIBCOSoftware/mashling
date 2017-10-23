/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package rest

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	opentracing "github.com/opentracing/opentracing-go"
	ctx "golang.org/x/net/context"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-rest")

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodPATCH  = "PATCH"
	methodDELETE = "DELETE"

	ivMethod      = "method"
	ivURI         = "uri"
	ivPathParams  = "pathParams"
	ivQueryParams = "queryParams"
	ivContent     = "content"
	ivParams      = "params"
	ivTracing     = "tracing"
	ivServerCert  = "serverCert"
	ivServerKey   = "serverKey"
	ivTrustStore  = "trustStore"

	ovResult  = "result"
	ovTracing = "tracing"
	ovStatus  = "status"
)

var validMethods = []string{methodGET, methodPOST, methodPUT, methodPATCH, methodDELETE}

// RESTActivity is an Activity that is used to invoke a REST Operation
// inputs : {method,uri,params}
// outputs: {result}
type RESTActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new RESTActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &RESTActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *RESTActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *RESTActivity) Eval(context activity.Context) (done bool, err error) {
	var span opentracing.Span
	if tracing := context.GetInput(ivTracing); tracing != nil {
		span = opentracing.SpanFromContext(tracing.(ctx.Context))
	}

	if span != nil {
		span = opentracing.StartSpan(
			context.TaskName(),
			opentracing.ChildOf(span.Context()))
		context.SetOutput(ovTracing, opentracing.ContextWithSpan(ctx.Background(), span))
		defer span.Finish()
	}

	setTag := func(key string, value interface{}) {
		if span != nil {
			span.SetTag(key, value)
		}
	}

	method := strings.ToUpper(context.GetInput(ivMethod).(string))
	uri := context.GetInput(ivURI).(string)

	containsParam := strings.Index(uri, "/:") > -1

	if containsParam {

		val := context.GetInput(ivPathParams)

		if val == nil {
			val = context.GetInput(ivParams)

			if val == nil {
				err := activity.NewError("Path Params not specified, required for URI: "+uri, "", nil)
				setTag("error", err.Error())
				return false, err
			}
		}

		pathParams := val.(map[string]string)
		uri = BuildURI(uri, pathParams)
	}

	if context.GetInput(ivQueryParams) != nil {
		queryParams := context.GetInput(ivQueryParams).(map[string]string)

		qp := url.Values{}

		for key, value := range queryParams {
			qp.Set(key, value)
		}

		uri = uri + "?" + qp.Encode()
	}

	setTag("method", method)
	setTag("uri", uri)

	log.Debugf("REST Call: [%s] %s\n", method, uri)

	var reqBody io.Reader

	contentType := "application/json; charset=UTF-8"

	if method == methodPOST || method == methodPUT || method == methodPATCH {

		content := context.GetInput(ivContent)

		contentType = getContentType(content)

		if content != nil {
			if str, ok := content.(string); ok {
				reqBody = bytes.NewBuffer([]byte(str))
				setTag("payload", str)
			} else {
				b, _ := json.Marshal(content) //todo handle error
				reqBody = bytes.NewBuffer([]byte(b))
				setTag("payload", str)
			}
		}
	} else {
		reqBody = nil
	}

	setTag("contentType", contentType)

	req, err := http.NewRequest(method, uri, reqBody)
	if reqBody != nil {
		req.Header.Set("Content-Type", contentType)
	}

	if span != nil {
		opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
	}

	//Enable transport layaer security
	serverCert, _ := context.GetInput(ivServerCert).(string)
	serverKey, _ := context.GetInput(ivServerKey).(string)
	trustStore, _ := context.GetInput(ivTrustStore).(string)

	tlsConfig := &tls.Config{}

	if serverCert != "" && serverKey != "" {
		//gateway certificates are available
		//load gateway certificate-key pair
		log.Debug("Loading gateway certificate - key pair...")
		cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
		if err != nil {
			log.Errorf("Unable to load cert - %v", err)
		} else {
			tlsConfig.Certificates = []tls.Certificate{cert}
			log.Debug("Loading gateway certificate - key pair DONE")
		}
	}

	if trustStore != "" {
		//trust store directory is available
		//load trusted certificates present inside the dir
		log.Debug("Loading truststore...")
		trustRootCAPool, err := getCerts(trustStore)
		if err != nil {
			log.Errorf("Error while loading trust store - %v", err)
		} else {
			tlsConfig.RootCAs = trustRootCAPool
			log.Debug("Loading truststore DONE")
		}
	}

	tlsConfig.BuildNameToCertificate()

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	resp, err := client.Do(req)
	if err != nil {
		setTag("error", err.Error())
		return false, err
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)

	var result interface{}

	d := json.NewDecoder(bytes.NewReader(respBody))
	d.UseNumber()
	err = d.Decode(&result)

	//json.Unmarshal(respBody, &result)

	log.Debug("response Body:", result)
	setTag("response", string(respBody))
	setTag("responseStatus", resp.Status)
	context.SetOutput(ovResult, result)
	context.SetOutput(ovStatus, resp.Status)

	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

//todo just make contentType a setting
func getContentType(replyData interface{}) string {

	contentType := "application/json; charset=UTF-8"

	switch v := replyData.(type) {
	case string:
		if !strings.HasPrefix(v, "{") && !strings.HasPrefix(v, "[") {
			contentType = "text/plain; charset=UTF-8"
		}
	case int, int64, float64, bool, json.Number:
		contentType = "text/plain; charset=UTF-8"
	default:
		contentType = "application/json; charset=UTF-8"
	}

	return contentType
}

func methodIsValid(method string) bool {

	if !stringInList(method, validMethods) {
		return false
	}

	//validate path

	return true
}

func stringInList(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}

// BuildURI is a temporary crude URI builder
func BuildURI(uri string, values map[string]string) string {

	var buffer bytes.Buffer
	buffer.Grow(len(uri))

	addrStart := strings.Index(uri, "://")

	i := addrStart + 3

	for i < len(uri) {
		if uri[i] == '/' {
			break
		}
		i++
	}

	buffer.WriteString(uri[0:i])

	for i < len(uri) {
		if uri[i] == ':' {
			j := i + 1
			for j < len(uri) && uri[j] != '/' {
				j++
			}

			if i+1 == j {

				buffer.WriteByte(uri[i])
				i++
			} else {

				param := uri[i+1 : j]
				value := values[param]
				buffer.WriteString(value)
				if j < len(uri) {
					buffer.WriteString("/")
				}
				i = j + 1
			}

		} else {
			buffer.WriteByte(uri[i])
			i++
		}
	}

	return buffer.String()
}

func getCerts(trustStore string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	fileInfo, err := os.Stat(trustStore)
	if err != nil {
		return certPool, fmt.Errorf("Truststore [%s] does not exist", trustStore)
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		break
	case mode.IsRegular():
		return certPool, fmt.Errorf("Truststore [%s] is not a directory.  Must be a directory containing trusted certificates in PEM format",
			trustStore)
	}
	trustedCertFiles, err := ioutil.ReadDir(trustStore)
	if err != nil || len(trustedCertFiles) == 0 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  Must be a directory containing trusted certificates in PEM format", trustStore)
	}
	for _, trustCertFile := range trustedCertFiles {
		fqfName := fmt.Sprintf("%s%c%s", trustStore, os.PathSeparator, trustCertFile.Name())
		trustCertBytes, err := ioutil.ReadFile(fqfName)
		if err != nil {
			log.Warnf("Failed to read trusted certificate [%s] ... continueing", trustCertFile.Name())
		}
		log.Debugf("Loading cert file - %v", fqfName)
		certPool.AppendCertsFromPEM(trustCertBytes)
	}
	if len(certPool.Subjects()) < 1 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  After processing all files in the directory no valid trusted certs were found", trustStore)
	}
	return certPool, nil
}
