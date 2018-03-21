package coap

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/dustin/go-coap"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-coap")

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodDELETE = "DELETE"

	typeCON = "CONFIRMABLE"
	typeNON = "NONCONFIRMABLE"
	typeACK = "ACKNOWLEDGEMENT"
	typeRST = "RESET"

	ivMethod      = "method"
	ivURI         = "uri"
	ivQueryParams = "queryParams"
	ivType        = "type"
	ivPayload     = "payload"
	ivMessageID   = "messageId"
	ivOptions     = "options"

	ovResponse = "response"
)

var validMethods = []string{methodGET, methodPOST, methodPUT, methodDELETE}
var validTypes = []string{typeCON, typeNON}

// CoAPActivity is an Activity that is used to send a CoAP message
// inputs : {method,type,payload,messageId}
// outputs: {result}
type CoAPActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new CoAP activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CoAPActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *CoAPActivity) Metadata() *activity.Metadata {
	return a.metadata
}

//todo enhance CoAP client code

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *CoAPActivity) Eval(context activity.Context) (done bool, err error) {

	method, ok := getStringValue(context, ivMethod, nil, true)

	if !ok {
		activity.NewError("Method not specified", "", nil)
	}

	uri, ok := getStringValue(context, ivURI, nil, false)

	if !ok {
		activity.NewError("URI not specified", "", nil)
	}

	msgType, _ := getStringValue(context, ivType, typeNON, true)
	payload, hasPayload := getStringValue(context, ivPayload, nil, false)
	messageID, _ := getIntValue(context, ivMessageID, 0)

	coapURI, err := url.Parse(uri)
	if err != nil {
		return false, activity.NewError(err.Error(), "", nil)
	}

	scheme := coapURI.Scheme
	if scheme != "coap" {
		return false, activity.NewError("URI scheme must be 'coap'", "", nil)
	}

	req := coap.Message{
		Type:      toCoapType(msgType),
		Code:      toCoapCode(method),
		MessageID: uint16(messageID),
	}

	if hasPayload {
		req.Payload = []byte(payload)
	}

	val := context.GetInput(ivOptions)
	if val != nil {
		options := val.(map[string]string)

		for k, v := range options {
			op, val := toOption(k, v)
			req.SetOption(op, val)
		}
	}

	if context.GetInput(ivQueryParams) != nil {
		queryParams := context.GetInput(ivQueryParams).(map[string]string)

		qp := url.Values{}

		for key, value := range queryParams {
			qp.Set(key, value)
		}

		queryStr := qp.Encode()
		req.SetOption(coap.URIQuery, queryStr)
		log.Debugf("CoAP Message: [%s] %s?%s\n", method, coapURI.Path, queryStr)

	} else {
		log.Debugf("CoAP Message: [%s] %s\n", method, coapURI.Path)
	}

	req.SetPathString(coapURI.Path)

	c, err := coap.Dial("udp", coapURI.Host)
	if err != nil {
		return false, activity.NewError(err.Error(), "", nil)
	}

	log.Debugf("conn: %v\n", c)

	rv, err := c.Send(req)
	if err != nil {
		return false, activity.NewError(err.Error(), "", nil)
	}

	if rv != nil {

		if rv.Code > 100 {
			return false, activity.NewError(fmt.Sprintf("CoAP Error: %s", rv.Code.String()), rv.Code.String(), nil)
		}

		if rv.Payload != nil {
			log.Debugf("Response payload: %s", rv.Payload)
			context.SetOutput(ovResponse, string(rv.Payload))
		}
	}

	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

func toCoapCode(method string) coap.COAPCode {

	var code coap.COAPCode

	switch method {
	case methodGET:
		code = coap.GET
	case methodPOST:
		code = coap.POST
	case methodPUT:
		code = coap.PUT
	case methodDELETE:
		code = coap.DELETE
	}

	return code
}

func toCoapType(typeStr string) coap.COAPType {

	var ctype coap.COAPType

	switch typeStr {
	case typeCON:
		ctype = coap.Confirmable
	case typeNON:
		ctype = coap.NonConfirmable
	case typeACK:
		ctype = coap.Acknowledgement
	case typeRST:
		ctype = coap.Reset
	}

	return ctype
}

func toOption(name string, value string) (coap.OptionID, interface{}) {

	var opID coap.OptionID
	var val interface{}

	val = value

	switch name {
	case "IFMATCH":
		opID = coap.IfMatch
	case "URIHOST":
		opID = coap.URIHost
	case "ETAG":
		opID = coap.ETag
	//case "IFNONEMATCH":
	//	opID = coap.IfNoneMatch
	case "OBSERVE":
		opID = coap.Observe
		val, _ = strconv.Atoi(value)
	case "URIPORT":
		opID = coap.URIPort
		val, _ = strconv.Atoi(value)
	case "LOCATIONPATH":
		opID = coap.LocationPath
	case "URIPATH":
		opID = coap.URIPath
	case "CONTENTFORMAT":
		opID = coap.ContentFormat
		val, _ = strconv.Atoi(value)
	case "MAXAGE":
		opID = coap.MaxAge
		val, _ = strconv.Atoi(value)
	case "URIQUERY":
		opID = coap.URIQuery
	case "ACCEPT":
		opID = coap.IfMatch
		val, _ = strconv.Atoi(value)
	case "LOCATIONQUERY":
		opID = coap.LocationQuery
	case "PROXYURI":
		opID = coap.ProxyURI
	case "PROXYSCHEME":
		opID = coap.ProxyScheme
	case "SIZE1":
		opID = coap.Size1
		val, _ = strconv.Atoi(value)
	default:
		opID = 0
		val = nil
	}

	return opID, val
}

func methodIsValid(method string) bool {

	if !stringInList(method, validMethods) {
		return false
	}

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

func getStringValue(context activity.Context, attrName string, defValue interface{}, uc bool) (string, bool) {

	val := context.GetInput(attrName)
	found := true

	if val == nil {
		found = false

		if defValue == nil {
			return "", false
		}
		val = defValue
	}

	if uc {
		return strings.ToUpper(val.(string)), found
	}

	return val.(string), found
}

func getIntValue(context activity.Context, attrName string, defValue interface{}) (int, bool) {

	val := context.GetInput(attrName)
	found := true

	if val == nil {
		found = false

		if defValue == nil {
			return 0, false
		}
		val = defValue
	}

	return val.(int), found
}
