package coap

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/dustin/go-coap"
)

const (
	methodGET    = "GET"
	methodPOST   = "POST"
	methodPUT    = "PUT"
	methodDELETE = "DELETE"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-coap")

var validMethods = []string{methodGET, methodPOST, methodPUT, methodDELETE}

type StartFunc func(payload string) (string, bool)

// CoapTrigger CoAP trigger struct
type CoapTrigger struct {
	metadata  *trigger.Metadata
	resources map[string]*CoapResource
	server    *Server
	config    *trigger.Config
}

type CoapResource struct {
	path     string
	attrs    map[string]string
	handlers map[string]*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &CoapFactory{metadata: md}
}

//CoapFactory Coap Trigger factory
type CoapFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (f *CoapFactory) New(config *trigger.Config) trigger.Trigger {
	return &CoapTrigger{metadata: f.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *CoapTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *CoapTrigger) Initialize(ctx trigger.InitContext) error {

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	port := t.config.Settings["port"]

	if port == "" {
		panic(fmt.Sprintf("No Port found for trigger '%s' in settings", t.config.Id))
	}

	mux := coap.NewServeMux()
	mux.Handle("/.well-known/core", coap.FuncHandler(t.handleDiscovery))

	t.resources = make(map[string]*CoapResource)

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		if handlerIsValid(handler) {
			method := strings.ToUpper(handler.GetStringSetting("method"))
			path := handler.GetStringSetting("path")

			log.Debugf("Registering handler for [%s: %s] - %s", method, path, handler)

			resource, exists := t.resources[path]

			if !exists {
				resource = &CoapResource{path: path, attrs: make(map[string]string), handlers: make(map[string]*trigger.Handler)}
				t.resources[path] = resource
			}

			resource.handlers[method] = handler

			mux.Handle(path, newActionHandler(resource))

		} else {
			return fmt.Errorf("invalid endpoint: %v", handler)
		}
	}

	log.Debugf("CoAP Trigger: Configured on port %s", port)

	t.server = NewServer("udp", fmt.Sprintf(":%s", port), mux)

	return nil
}

// Start implements trigger.Trigger.Start
func (t *CoapTrigger) Start() error {
	return t.server.Start()
}

// Stop implements trigger.Trigger.Start
func (t *CoapTrigger) Stop() error {
	return t.server.Stop()
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func (t *CoapTrigger) handleDiscovery(conn *net.UDPConn, addr *net.UDPAddr, msg *coap.Message) *coap.Message {

	//path := msg.PathString() //handle queries

	//todo add filter support

	var buffer bytes.Buffer

	numResources := len(t.resources)

	i := 0
	for _, resource := range t.resources {

		i++

		buffer.WriteString("<")
		buffer.WriteString(resource.path)
		buffer.WriteString(">")

		if len(resource.attrs) > 0 {
			for k, v := range resource.attrs {
				buffer.WriteString(";")
				buffer.WriteString(k)
				buffer.WriteString("=")
				buffer.WriteString(v)
			}
		}

		if i < numResources {
			buffer.WriteString(",\n")
		} else {
			buffer.WriteString("\n")
		}
	}

	payloadStr := buffer.String()

	res := &coap.Message{
		Type:      msg.Type,
		Code:      coap.Content,
		MessageID: msg.MessageID,
		Token:     msg.Token,
		Payload:   []byte(payloadStr),
	}
	res.SetOption(coap.ContentFormat, coap.AppLinkFormat)

	log.Debugf("Transmitting %#v", res)

	return res
}

func newActionHandler(resource *CoapResource) coap.Handler {

	return coap.FuncHandler(func(conn *net.UDPConn, addr *net.UDPAddr, msg *coap.Message) *coap.Message {

		log.Debugf("CoAP Trigger: Recieved request")

		method := toMethod(msg.Code)
		uriQuery := msg.Option(coap.URIQuery)
		var data map[string]interface{}

		if uriQuery != nil {
			//todo handle error
			queryValues, _ := url.ParseQuery(uriQuery.(string))

			queryParams := make(map[string]string, len(queryValues))

			for key, value := range queryValues {
				queryParams[key] = strings.Join(value, ",")
			}

			data = map[string]interface{}{
				"queryParams": queryParams,
				"payload":     string(msg.Payload),
			}
		} else {
			data = map[string]interface{}{
				"payload": string(msg.Payload),
			}
		}

		handler, exists := resource.handlers[method]

		if !exists {
			res := &coap.Message{
				Type:      coap.Reset,
				Code:      coap.MethodNotAllowed,
				MessageID: msg.MessageID,
				Token:     msg.Token,
			}

			return res
		}

		_, err := handler.Handle(context.Background(), data)

		if err != nil {
			//todo determining if 404 or 500
			res := &coap.Message{
				Type:      coap.Reset,
				Code:      coap.NotFound,
				MessageID: msg.MessageID,
				Token:     msg.Token,
				Payload:   []byte(fmt.Sprintf("Unable to execute handler '%s'", handler)),
			}

			return res
		}

		log.Debugf("Ran: %s", handler)

		if msg.IsConfirmable() {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      0,
				MessageID: msg.MessageID,
				Token:     msg.Token,
			}
			//res.SetOption(coap.ContentFormat, coap.TextPlain)

			log.Debugf("Transmitting %#v", res)
			return res
		}

		return nil
	})
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils

func handlerIsValid(handler *trigger.Handler) bool {

	method := handler.GetStringSetting("method")

	if method == "" {
		return false
	}

	if !stringInList(strings.ToUpper(method), validMethods) {
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

func toMethod(code coap.COAPCode) string {

	var method string

	switch code {
	case coap.GET:
		method = methodGET
	case coap.POST:
		method = methodPOST
	case coap.PUT:
		method = methodPUT
	case coap.DELETE:
		method = methodDELETE
	}

	return method
}
