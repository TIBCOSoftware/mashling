/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package eftl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	condition "github.com/TIBCOSoftware/mashling/lib/conditions"
	"github.com/TIBCOSoftware/mashling/lib/eftl"
	"github.com/TIBCOSoftware/mashling/lib/util"

	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"sourcegraph.com/sourcegraph/appdash"
	appdashtracing "sourcegraph.com/sourcegraph/appdash/opentracing"
)

const (
	TracerNoOP      = "noop"
	TracerZipKin    = "zipkin"
	TracerAPPDash   = "appdash"
	TracerLightStep = "lightstep"

	settingURL            = "url"
	settingID             = "id"
	settingUser           = "user"
	settingPassword       = "password"
	settingCA             = "ca"
	settingTracer         = "tracer"
	settingTracerEndpoint = "tracerEndpoint"
	settingTracerToken    = "tracerToken"
	settingTracerDebug    = "tracerDebug"
	settingTracerSameSpan = "tracerSameSpan"
	settingTracerID128Bit = "tracerID128Bit"
	settingDest           = "dest"
)

var (
	ErrorTracerEndpointRequired = errors.New("tracer endpoint required")
	ErrorInvalidTracer          = errors.New("invalid tracer")
	ErrorTracerTokenRequired    = errors.New("tracer token required")
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-eftl")

// Span is a tracing span
type Span struct {
	opentracing.Span
}

// Error is for reporting errors
func (s *Span) Error(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	s.SetTag("error", str)
	log.Error(str)
}

//OptimizedHandler optimized handler
type OptimizedHandler struct {
	defaultActionID   string
	defaultHandlerCfg *trigger.HandlerConfig
	dispatches        []*Dispatch
}

// GetActionID gets the action id of the matched handler
func (h *OptimizedHandler) GetActionID(payload string, span Span) (string, *trigger.HandlerConfig) {
	actionID := ""
	var handlerCfg *trigger.HandlerConfig

	for _, dispatch := range h.dispatches {
		expressionStr := dispatch.condition
		//Get condtion and expression type
		conditionOperation, exprType, err := condition.GetConditionOperationAndExpressionType(expressionStr)

		if err != nil || exprType == condition.EXPR_TYPE_NOT_VALID {
			span.Error("not able parse the condition '%v' mentioned for content based handler. skipping the handler.", expressionStr)
			continue
		}

		log.Debugf("Expression type: %v", exprType)
		log.Debugf("conditionOperation.LHS %v", conditionOperation.LHS)
		log.Debugf("conditionOperation.OperatorInfo %v", conditionOperation.OperatorInfo().Names)
		log.Debugf("conditionOperation.RHS %v", conditionOperation.RHS)

		//Resolve expression's LHS based on expression type and
		//evaluate the expression
		if exprType == condition.EXPR_TYPE_CONTENT {
			exprResult, err := condition.EvaluateCondition(*conditionOperation, payload)
			if err != nil {
				span.Error("not able evaluate expression - %v with error - %v. skipping the handler.", expressionStr, err)
			}
			if exprResult {
				actionID = dispatch.actionID
				handlerCfg = dispatch.handlerCfg
			}
		} else if exprType == condition.EXPR_TYPE_HEADER {
			span.Error("header expression type is invalid for eftl trigger condition")
		} else if exprType == condition.EXPR_TYPE_ENV {
			//environment variable based condition
			envFlagValue := os.Getenv(conditionOperation.LHS)
			log.Debugf("environment flag = %v, val = %v", conditionOperation.LHS, envFlagValue)
			if envFlagValue != "" {
				conditionOperation.LHS = envFlagValue
				op := conditionOperation.Operator
				exprResult := op.Eval(conditionOperation.LHS, conditionOperation.RHS)
				if exprResult {
					actionID = dispatch.actionID
					handlerCfg = dispatch.handlerCfg
				}
			}
		}

		if actionID != "" {
			log.Debugf("dispatch resolved with the actionId - %v", actionID)
			break
		}
	}

	//If no dispatch is found, use default action
	if actionID == "" {
		actionID = h.defaultActionID
		handlerCfg = h.defaultHandlerCfg
		log.Debugf("dispatch not resolved. Continue with default action - %v", actionID)
	}

	return actionID, handlerCfg
}

//Dispatch holds dispatch actionId and condition
type Dispatch struct {
	actionID   string
	condition  string
	handlerCfg *trigger.HandlerConfig
}

// Trigger is a simple EFTL trigger
type Trigger struct {
	metadata   *trigger.Metadata
	runner     action.Runner
	config     *trigger.Config
	handlers   map[string]*OptimizedHandler
	connection *eftl.Connection
	stop       chan bool
}

// Factory MQTT Trigger factory
type Factory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (f *Factory) New(config *trigger.Config) trigger.Trigger {
	return &Trigger{metadata: f.metadata, config: config}
}

//NewFactory create a new Trigger factory
func NewFactory(metadata *trigger.Metadata) trigger.Factory {
	return &Factory{metadata: metadata}
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// CreateHandlers creates handlers mapped to thier topic
func (t *Trigger) CreateHandlers() map[string]*OptimizedHandler {
	handlers := make(map[string]*OptimizedHandler)

	for _, h := range t.config.Handlers {
		t := h.Settings[settingDest]
		if t == nil {
			continue
		}
		dest := t.(string)

		handler := handlers[dest]
		if handler == nil {
			handler = &OptimizedHandler{}
			handlers[dest] = handler
		}

		if condition := h.Settings[util.Flogo_Trigger_Handler_Setting_Condition]; condition != nil {
			dispatch := &Dispatch{
				actionID:   h.ActionId,
				condition:  condition.(string),
				handlerCfg: h,
			}
			handler.dispatches = append(handler.dispatches, dispatch)
		} else {
			handler.defaultActionID = h.ActionId
			handler.defaultHandlerCfg = h
		}
	}

	return handlers
}

// Init implements ext.Trigger.Init
func (t *Trigger) Init(runner action.Runner) {
	t.runner = runner
	t.handlers = t.CreateHandlers()
}

// getLocalIP gets the public ip address of the system
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "0.0.0.0"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "0.0.0.0"
}

// configureTracer configures the distributed tracer
func (t *Trigger) configureTracer() {
	tracer := TracerNoOP
	if setting, ok := t.config.Settings[settingTracer]; ok {
		tracer = setting.(string)
	}
	tracerEndpoint := ""
	if setting, ok := t.config.Settings[settingTracerEndpoint]; ok {
		tracerEndpoint = setting.(string)
	}
	tracerToken := ""
	if setting, ok := t.config.Settings[settingTracerToken]; ok {
		tracerToken = setting.(string)
	}
	tracerDebug := false
	if setting, ok := t.config.Settings[settingTracerDebug]; ok {
		tracerDebug = setting.(bool)
	}
	tracerSameSpan := false
	if setting, ok := t.config.Settings[settingTracerSameSpan]; ok {
		tracerSameSpan = setting.(bool)
	}
	tracerID128Bit := true
	if setting, ok := t.config.Settings[settingTracerID128Bit]; ok {
		tracerID128Bit = setting.(bool)
	}

	switch tracer {
	case TracerNoOP:
		opentracing.SetGlobalTracer(&opentracing.NoopTracer{})
	case TracerZipKin:
		if tracerEndpoint == "" {
			panic(ErrorTracerEndpointRequired)
		}

		collector, err := zipkin.NewHTTPCollector(tracerEndpoint)
		if err != nil {
			panic(fmt.Sprintf("unable to create Zipkin HTTP collector: %+v\n", err))
		}

		recorder := zipkin.NewRecorder(collector, tracerDebug,
			getLocalIP(), t.config.Name)

		tracer, err := zipkin.NewTracer(
			recorder,
			zipkin.ClientServerSameSpan(tracerSameSpan),
			zipkin.TraceID128Bit(tracerID128Bit),
		)
		if err != nil {
			panic(fmt.Sprintf("unable to create Zipkin tracer: %+v\n", err))
		}

		opentracing.SetGlobalTracer(tracer)
	case TracerAPPDash:
		if tracerEndpoint == "" {
			panic(ErrorTracerEndpointRequired)
		}

		collector := appdash.NewRemoteCollector(tracerEndpoint)
		chunkedCollector := appdash.NewChunkedCollector(collector)
		tracer := appdashtracing.NewTracer(chunkedCollector)
		opentracing.SetGlobalTracer(tracer)
	case TracerLightStep:
		if tracerToken == "" {
			panic(ErrorTracerTokenRequired)
		}

		lightstepTracer := lightstep.NewTracer(lightstep.Options{
			AccessToken: tracerToken,
		})

		opentracing.SetGlobalTracer(lightstepTracer)
	default:
		panic(ErrorInvalidTracer)
	}
}

// Start implements ext.Trigger.Start
func (t *Trigger) Start() error {
	t.configureTracer()

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	ca := t.config.GetSetting(settingCA)
	if ca != "" {
		certificate, err := ioutil.ReadFile(ca)
		if err != nil {
			log.Error("can't open certificate", err)
			return err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(certificate)
		tlsConfig = &tls.Config{
			RootCAs: pool,
		}
	}

	id := t.config.GetSetting(settingID)
	user := t.config.GetSetting(settingUser)
	password := t.config.GetSetting(settingPassword)
	options := &eftl.Options{
		ClientID:  id,
		Username:  user,
		Password:  password,
		TLSConfig: tlsConfig,
	}

	url := t.config.GetSetting(settingURL)
	errorsChannel := make(chan error, 1)
	var err error
	t.connection, err = eftl.Connect(url, options, errorsChannel)
	if err != nil {
		log.Errorf("connection failed: %s", err)
		return err
	}

	messages := make(chan eftl.Message, 1000)
	for dest := range t.handlers {
		matcher := fmt.Sprintf("{\"_dest\":\"%s\"}", dest)
		_, err = t.connection.Subscribe(matcher, "", messages)
		if err != nil {
			log.Errorf("subscription failed: %s", err)
			return err
		}
	}

	t.stop = make(chan bool, 1)
	go func() {
		for {
			select {
			case message := <-messages:
				value := message["_dest"]
				dest, ok := value.(string)
				if !ok {
					log.Errorf("dest is required for valid message")
					continue
				}
				handler := t.handlers[dest]
				if handler == nil {
					log.Errorf("no handler for dest", dest)
					continue
				}
				value = message["content"]
				content, ok := value.([]byte)
				if !ok {
					content = []byte{}
				}
				t.RunAction(handler, dest, content)
			case err := <-errorsChannel:
				log.Errorf("connection error: %s", err)
			case <-t.stop:
				return
			}
		}
	}()

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *Trigger) Stop() error {
	if t.connection != nil {
		t.connection.Disconnect()
	}
	if t.stop != nil {
		t.stop <- true
	}
	return nil
}

// RunAction starts a new Process Instance
func (t *Trigger) RunAction(handler *OptimizedHandler, dest string, content []byte) {
	span := Span{
		Span: opentracing.StartSpan(dest),
	}
	defer span.Finish()

	replyTo, data := t.constructStartRequest(content, span)

	startAttrs, err := t.metadata.OutputsToAttrs(data, false)
	if err != nil {
		span.Error("Error setting up attrs: %v", err)
	}

	actionURI, handlerCfg := handler.GetActionID(string(content), span)
	action := action.Get(actionURI)
	context := trigger.NewContextWithData(context.Background(), &trigger.ContextData{Attrs: startAttrs, HandlerCfg: handlerCfg})
	_, replyData, err := t.runner.Run(context, action, actionURI, nil)
	if err != nil {
		span.Error("Error starting action: %v", err)
	}
	log.Debugf("Ran action: [%s]", actionURI)

	if replyTo == "" {
		return
	}
	reply, err := util.Marshal(replyData)
	if err != nil {
		span.Error("failed to marshal reply data: %v", err)
		return
	}
	span.SetTag("replyTo", replyTo)
	span.SetTag("reply", string(reply))
	err = t.connection.Publish(eftl.Message{
		"_dest":   replyTo,
		"content": reply,
	})
	if err != nil {
		span.Error("failed to send reply data: %v", err)
	}
}

func (t *Trigger) constructStartRequest(message []byte, span Span) (string, map[string]interface{}) {
	span.SetTag("message", string(message))

	var content map[string]interface{}
	err := util.Unmarshal("", message, &content)
	if err != nil {
		span.Error("Error unmarshaling message ", err.Error())
	}

	replyTo := ""
	pathParams := make(map[string]string)
	queryParams := make(map[string]string)

	mime := ""
	if value, ok := content[util.MetaMIME].(string); ok {
		mime = value
	}
	if mime == util.MIMEApplicationXML {
		getRoot := func() map[string]interface{} {
			body := content[util.XMLKeyBody]
			if body == nil {
				return nil
			}
			for _, e := range body.([]interface{}) {
				element, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				name, ok := element[util.XMLKeyType].(string)
				if !ok || name != util.XMLTypeElement {
					continue
				}
				return element
			}
			return nil
		}
		root := getRoot()
		fill := func(target string, params map[string]string) {
			rootBody, ok := root[util.XMLKeyBody].([]interface{})
			if !ok {
				return
			}
			for i, e := range rootBody {
				element, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				name, ok := element[util.XMLKeyName].(string)
				if !ok || name != target {
					continue
				}
				body := element[util.XMLKeyBody]
				if body == nil {
					continue
				}
				for _, e := range body.([]interface{}) {
					element, ok := e.(map[string]interface{})
					if !ok {
						continue
					}
					typ, ok := element[util.XMLKeyType].(string)
					if !ok || typ != util.XMLTypeElement {
						continue
					}
					params[element["key"].(string)] = element["value"].(string)
				}
				root[util.XMLKeyBody] = rootBody[:i+copy(rootBody[i:], rootBody[i+1:])]
				return
			}
		}

		if value, ok := root["replyTo"].(string); ok {
			replyTo = value
			delete(root, "replyTo")
		}
		fill("pathParams", pathParams)
		fill("queryParams", queryParams)
	} else {
		if value, ok := content["replyTo"].(string); ok {
			replyTo = value
			delete(content, "replyTo")
		}

		if params, ok := content["pathParams"].(map[string]interface{}); ok {
			for k, v := range params {
				if param, ok := v.(string); ok {
					pathParams[k] = param
				}
			}
			delete(content, "pathParams")
		}

		if params, ok := content["queryParams"].(map[string]interface{}); ok {
			for k, v := range params {
				if param, ok := v.(string); ok {
					queryParams[k] = param
				}
			}
			delete(content, "queryParams")
		}
	}

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	data := map[string]interface{}{
		"params":      pathParams,
		"pathParams":  pathParams,
		"queryParams": queryParams,
		"content":     content,
		"tracing":     ctx,
	}

	return replyTo, data
}
