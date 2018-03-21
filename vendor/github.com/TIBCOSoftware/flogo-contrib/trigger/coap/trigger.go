package coap

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
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
var log = logger.GetLogger("trigger-tibco-coap")

var validMethods = []string{methodGET, methodPOST, methodPUT, methodDELETE}

type StartFunc func(payload string) (string, bool)

// CoapTrigger CoAP trigger struct
type CoapTrigger struct {
	metadata  *trigger.Metadata
	runner    action.Runner
	resources map[string]*CoapResource
	server    *Server
	config    *trigger.Config
}

type CoapResource struct {
	path     string
	attrs    map[string]string
	handlers map[string]string
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

func (t *CoapTrigger) Init(runner action.Runner) {

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	port := t.config.Settings["port"]

	if port == "" {
		panic(fmt.Sprintf("No Port found for trigger '%s' in settings", t.config.Id))
	}

	t.runner = runner
	mux := coap.NewServeMux()
	mux.Handle("/.well-known/core", coap.FuncHandler(t.handleDiscovery))

	t.resources = make(map[string]*CoapResource)

	// Init handlers
	for _, handler := range t.config.Handlers {

		if handlerIsValid(handler) {
			method := strings.ToUpper(handler.GetSetting("method"))
			path := handler.GetSetting("path")

			log.Debugf("COAP Trigger: Registering handler [%s: %s] for Action Id: [%s]", method, path, handler.ActionId)

			resource, exists := t.resources[path]

			if !exists {
				resource = &CoapResource{path: path, attrs: make(map[string]string), handlers: make(map[string]string)}
				t.resources[path] = resource
			}

			resource.handlers[method] = handler.ActionId

			mux.Handle(path, newActionHandler(t, resource))

		} else {
			panic(fmt.Sprintf("Invalid endpoint: %v", handler))
		}
	}

	log.Debugf("CoAP Trigger: Configured on port %s", port)

	t.server = NewServer("udp", fmt.Sprintf(":%s", port), mux)

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

func newActionHandler(rt *CoapTrigger, resource *CoapResource) coap.Handler {

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

		actionId, exists := resource.handlers[method]

		if !exists {
			res := &coap.Message{
				Type:      coap.Reset,
				Code:      coap.MethodNotAllowed,
				MessageID: msg.MessageID,
				Token:     msg.Token,
			}

			return res
		}

		//todo handle error
		startAttrs, _ := rt.metadata.OutputsToAttrs(data, false)

		//rh := &AsyncReplyHandler{addr: addr.String(), msg: msg}
		//rh.addr2 = addr
		//rh.conn = conn

		act := action.Get(actionId)
		ctx := trigger.NewInitialContext(startAttrs, nil)
		_, err := rt.runner.RunAction(ctx, act, nil)

		if err != nil {
			//todo determining if 404 or 500
			res := &coap.Message{
				Type:      coap.Reset,
				Code:      coap.NotFound,
				MessageID: msg.MessageID,
				Token:     msg.Token,
				Payload:   []byte(fmt.Sprintf("Flow '%s' not found", actionId)),
			}

			return res
		}

		log.Debugf("Ran Action: %s", actionId)

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

func handlerIsValid(handler *trigger.HandlerConfig) bool {
	if handler.Settings == nil {
		return false
	}

	if handler.Settings["method"] == "" {
		return false
	}

	if !stringInList(strings.ToUpper(handler.GetSetting("method")), validMethods) {
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
