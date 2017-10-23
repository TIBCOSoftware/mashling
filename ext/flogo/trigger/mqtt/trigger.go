package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	condition "github.com/TIBCOSoftware/mashling/lib/conditions"
	"github.com/TIBCOSoftware/mashling/lib/util"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
)

var (
	ErrorTracerEndpointRequired = errors.New("tracer endpoint required")
	ErrorInvalidTracer          = errors.New("invalid tracer")
	ErrorTracerTokenRequired    = errors.New("tracer token required")
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-mqtt")

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
	defaultActionId string
	dispatches      []*Dispatch
}

// GetActionID gets the action id of the matched handler
func (h *OptimizedHandler) GetActionID(payload string, span Span) string {
	actionId := ""

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
				actionId = dispatch.actionId
			}
		} else if exprType == condition.EXPR_TYPE_HEADER {
			span.Error("header expression type is invalid for mqtt trigger condition")
		} else if exprType == condition.EXPR_TYPE_ENV {
			//environment variable based condition
			envFlagValue := os.Getenv(conditionOperation.LHS)
			log.Debugf("environment flag = %v, val = %v", conditionOperation.LHS, envFlagValue)
			if envFlagValue != "" {
				conditionOperation.LHS = envFlagValue
				op := conditionOperation.Operator
				exprResult := op.Eval(conditionOperation.LHS, conditionOperation.RHS)
				if exprResult {
					actionId = dispatch.actionId
				}
			}
		}

		if actionId != "" {
			log.Debugf("dispatch resolved with the actionId - %v", actionId)
			break
		}
	}

	//If no dispatch is found, use default action
	if actionId == "" {
		actionId = h.defaultActionId
		log.Debugf("dispatch not resolved. Continue with default action - %v", actionId)
	}

	return actionId
}

//Dispatch holds dispatch actionId and condition
type Dispatch struct {
	actionId  string
	condition string
}

// MqttTrigger is simple MQTT trigger
type MqttTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	client   mqtt.Client
	config   *trigger.Config
	handlers map[string]*OptimizedHandler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MQTTFactory{metadata: md}
}

// MQTTFactory MQTT Trigger factory
type MQTTFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *MQTTFactory) New(config *trigger.Config) trigger.Trigger {
	return &MqttTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *MqttTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *MqttTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// CreateHandlers creates handlers mapped to thier topic
func (t *MqttTrigger) CreateHandlers() map[string]*OptimizedHandler {
	handlers := make(map[string]*OptimizedHandler)

	for _, h := range t.config.Handlers {
		t := h.Settings["topic"]
		if t == nil {
			continue
		}
		topic := t.(string)

		handler := handlers[topic]
		if handler == nil {
			handler = &OptimizedHandler{}
			handlers[topic] = handler
		}

		if condition := h.Settings[util.Flogo_Trigger_Handler_Setting_Condition]; condition != nil {
			dispatch := &Dispatch{
				actionId:  h.ActionId,
				condition: condition.(string),
			}
			handler.dispatches = append(handler.dispatches, dispatch)
		} else {
			handler.defaultActionId = h.ActionId
		}
	}

	return handlers
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
func (t *MqttTrigger) configureTracer() {
	tracer := TracerNoOP
	if setting, ok := t.config.Settings["tracer"]; ok {
		tracer = setting.(string)
	}
	tracerEndpoint := ""
	if setting, ok := t.config.Settings["tracerEndpoint"]; ok {
		tracerEndpoint = setting.(string)
	}
	tracerToken := ""
	if setting, ok := t.config.Settings["tracerToken"]; ok {
		tracerToken = setting.(string)
	}
	tracerDebug := false
	if setting, ok := t.config.Settings["tracerDebug"]; ok {
		tracerDebug = setting.(bool)
	}
	tracerSameSpan := false
	if setting, ok := t.config.Settings["tracerSameSpan"]; ok {
		tracerSameSpan = setting.(bool)
	}
	tracerID128Bit := true
	if setting, ok := t.config.Settings["tracerID128Bit"]; ok {
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
func (t *MqttTrigger) Start() error {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(t.config.GetSetting("broker"))
	opts.SetClientID(t.config.GetSetting("id"))
	opts.SetUsername(t.config.GetSetting("user"))
	opts.SetPassword(t.config.GetSetting("password"))
	b, err := strconv.ParseBool(t.config.GetSetting("cleansess"))
	if err != nil {
		log.Error("Error converting \"cleansess\" to a boolean ", err.Error())
		return err
	}
	opts.SetCleanSession(b)
	if storeType := t.config.Settings["store"]; storeType != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(t.config.GetSetting("store")))
	}

	t.configureTracer()

	t.handlers = t.CreateHandlers()
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()

		span := Span{
			Span: opentracing.StartSpan(topic),
		}
		defer span.Finish()

		//TODO we should handle other types, since mqtt message format are data-agnostic
		payload := string(msg.Payload())
		log.Debug("Received msg:", payload)
		handler, found := t.handlers[topic]
		if found {
			t.RunAction(handler.GetActionID(payload, span), payload, span)
		} else {
			span.Error("Topic %s not found", topic)
		}
	})

	client := mqtt.NewClient(opts)
	t.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	i, err := strconv.Atoi(t.config.GetSetting("qos"))
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return err
	}

	for topic := range t.handlers {
		if token := t.client.Subscribe(topic, byte(i), nil); token.Wait() && token.Error() != nil {
			log.Errorf("Error subscribing to topic %s: %s", topic, token.Error())
			panic(token.Error())
		} else {
			log.Debugf("Suscribed to topic: %s", topic)
		}
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *MqttTrigger) Stop() error {
	//unsubscribe from topic

	for topic := range t.handlers {
		log.Debug("Unsubcribing from topic: ", topic)
		if token := t.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
			log.Errorf("Error unsubscribing from topic %s: %s", topic, token.Error())
		}
	}

	t.client.Disconnect(250)

	return nil
}

// RunAction starts a new Process Instance
func (t *MqttTrigger) RunAction(actionURI string, payload string, span Span) {
	span.SetTag("broker", t.config.GetSetting("broker"))

	req := t.constructStartRequest(payload, span)

	startAttrs, err := t.metadata.OutputsToAttrs(req.Data, false)
	if err != nil {
		span.Error("Error setting up attrs: %v", err)
	}

	action := action.Get(actionURI)
	context := trigger.NewContext(context.Background(), startAttrs)
	_, replyData, err := t.runner.Run(context, action, actionURI, nil)
	if err != nil {
		span.Error("Error starting action: %v", err)
	}
	log.Debugf("Ran action: [%s]", actionURI)

	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			span.Error(err.Error())
		} else if req.ReplyTo != "" {
			t.publishMessage(req.ReplyTo, string(data), span)
		}
	}
}

func (t *MqttTrigger) publishMessage(topic string, message string, span Span) {
	span.SetTag("replyTo", topic)
	span.SetTag("reply", message)

	log.Debug("ReplyTo topic: ", topic)
	log.Debug("Publishing message: ", message)

	qos, err := strconv.Atoi(t.config.GetSetting("qos"))
	if err != nil {
		span.Error("Error converting \"qos\" to an integer %v", err)
		return
	}
	if len(topic) == 0 {
		log.Warn("Invalid empty topic to publish to")
		return
	}
	token := t.client.Publish(topic, byte(qos), false, message)
	sent := token.WaitTimeout(5000 * time.Millisecond)
	if !sent {
		// Timeout occurred
		span.Error("Timeout occurred while trying to publish to topic '%s'", topic)
		return
	}
}

func (t *MqttTrigger) constructStartRequest(message string, span Span) *StartRequest {
	span.SetTag("message", message)

	req := &StartRequest{}

	var content map[string]interface{}
	err := json.Unmarshal([]byte(message), &content)
	if err != nil {
		span.Error("Error unmarshaling message ", err.Error())
	}

	if replyTo := content["replyTo"]; replyTo != nil {
		req.ReplyTo = replyTo.(string)
		delete(content, "replyTo")
	}

	pathParams := make(map[string]string)
	if params, ok := content["pathParams"].(map[string]interface{}); ok {
		for k, v := range params {
			if param, ok := v.(string); ok {
				pathParams[k] = param
			}
		}
		delete(content, "pathParams")
	}

	queryParams := make(map[string]string)
	if params, ok := content["queryParams"].(map[string]interface{}); ok {
		for k, v := range params {
			if param, ok := v.(string); ok {
				queryParams[k] = param
			}
		}
		delete(content, "queryParams")
	}

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	data := map[string]interface{}{
		"params":      pathParams,
		"pathParams":  pathParams,
		"queryParams": queryParams,
		"content":     content,
		"message":     message,
		"tracing":     ctx,
	}
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}
