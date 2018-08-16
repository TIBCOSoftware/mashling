package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"google.golang.org/grpc/credentials"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling/lib/util"
	"google.golang.org/grpc"
)

var addr string

const settingDest = "dest"

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-grpc")

// GRPCTriggerFactory is a gRPC Trigger factory
type GRPCTriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory creates a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &GRPCTriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *GRPCTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &GRPCTrigger{metadata: t.metadata, config: config}
}

// GRPCTrigger is a stub for your Trigger implementation
type GRPCTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	handlers map[string]*OptimizedHandler
	server   *grpc.Server
	TLSConfig
}

// TLSConfig is to hold tls support data
type TLSConfig struct {
	enableTLS bool
	serveKey  string
	serveCert string
}

// Init implements trigger.Trigger.Init
func (t *GRPCTrigger) Init(runner action.Runner) {

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	if _, ok := t.config.Settings["port"]; !ok {
		panic(fmt.Sprintf("No Port found for trigger '%s' in settings", t.config.Id))
	}

	addr = ":" + t.config.GetSetting("port")

	t.runner = runner

	t.handlers = t.CreateHandlers()

	//Check whether TLS (Transport Layer Security) is enabled for the trigger
	enableTLS := false
	serverCert := ""
	serverKey := ""
	if _, ok := t.config.Settings["enableTLS"]; ok {
		enableTLSSetting, err := strconv.ParseBool(t.config.GetSetting("enableTLS"))

		if err == nil && enableTLSSetting {
			//TLS is enabled, get server certificate & key
			enableTLS = true
			if _, ok := t.config.Settings["serverCert"]; !ok {
				panic(fmt.Sprintf("No serverCert found for trigger '%s' in settings", t.config.Id))
			}
			serverCert = t.config.GetSetting("serverCert")

			if _, ok := t.config.Settings["serverKey"]; !ok {
				panic(fmt.Sprintf("No serverKey found for trigger '%s' in settings", t.config.Id))
			}
			serverKey = t.config.GetSetting("serverKey")
		}
	}

	log.Debug("enableTLS: ", enableTLS)
	if enableTLS {
		log.Debug("serverCert: ", serverCert)
		log.Debug("serverKey: ", serverKey)
	}
	t.enableTLS = enableTLS
	t.serveCert = serverCert
	t.serveKey = serverKey
}

// Metadata implements trigger.Trigger.Metadata
func (t *GRPCTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Stop implements trigger.Trigger.Start
func (t *GRPCTrigger) Stop() error {
	// stop the trigger
	t.server.GracefulStop()
	return nil
}

// Start implements trigger.Trigger.Start
func (t *GRPCTrigger) Start() error {
	// start the trigger
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(err)
	}

	opts := []grpc.ServerOption{}

	if t.enableTLS {
		creds, err := credentials.NewServerTLSFromFile(t.serveCert, t.serveKey)
		if err != nil {
			log.Error(err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	t.server = grpc.NewServer(opts...)

	serviceName := t.config.GetSetting("serviceName")
	protoName := t.config.GetSetting("protoName")
	protoName = strings.Split(protoName, ".")[0]

	servRegFlag := false
	if len(ServiceRegistery.ServerServices) != 0 {
		for k, service := range ServiceRegistery.ServerServices {
			if strings.Compare(k, protoName+serviceName) == 0 {
				log.Infof("Registered Proto [%v] and Service [%v]", protoName, serviceName)
				service.RunRegisterServerService(t.server, t)
				servRegFlag = true
			}
		}
		if !servRegFlag {
			log.Errorf("Proto [%s] and Service [%s] not registered", protoName, serviceName)
		}
	} else {
		log.Error("gRPC server services not registered")
	}

	log.Debug("Starting server on port", addr)

	go func() {
		t.server.Serve(lis)
	}()

	log.Info("Server started")
	return nil
}

// Dispatch holds dispatch actionId and condition
type Dispatch struct {
	actionID   string
	condition  string
	handlerCfg *trigger.HandlerConfig
}

// OptimizedHandler optimized handler
type OptimizedHandler struct {
	defaultActionID   string
	defaultHandlerCfg *trigger.HandlerConfig
	dispatches        []*Dispatch
}

// CreateHandlers creates handlers mapped to their topic
func (t *GRPCTrigger) CreateHandlers() map[string]*OptimizedHandler {
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

// CallHandler is to call a particular handler based on method name
func (t *GRPCTrigger) CallHandler(grpcData map[string]interface{}) (int, interface{}, error) {
	log.Info("CallHandler method invoked")
	// getting values from inputrequestdata and mapping it to params which can be used in different services like HTTP pathparams etc.
	s := reflect.ValueOf(grpcData["reqdata"]).Elem()
	typeOfT := s.Type()
	params := make(map[string]interface{})
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		params[typeOfT.Field(i).Name] = f.Interface()
	}

	// assign req data content to trigger content
	var content interface{}
	dataBytes, err := util.Marshal(grpcData["reqdata"])
	if err != nil {
		log.Error("Marshal failed on grpc request data")
	}

	err = util.Unmarshal("application/json", dataBytes, &content)
	if err != nil {
		log.Error("Unmarshal failed on grpc request data")
	}

	grpcData["serviceName"] = t.config.GetSetting("serviceName")
	grpcData["protoName"] = t.config.GetSetting("protoName")

	data := map[string]interface{}{
		"params":   params,
		"grpcData": grpcData,
		"content":  content,
	}

	//todo handle error
	startAttrs, _ := t.metadata.OutputsToAttrs(data, false)

	handlers := t.config.Handlers

	//calling particular handler based on method name specification in gateway json file
	for _, hand := range handlers {
		if strings.Compare(hand.GetSetting("methodName"), grpcData["methodname"].(string)) == 0 {
			log.Debug("Dispatch Found for ", hand.GetSetting("methodName"), " Handler Invoked: ", hand.ActionId)
			actID := action.Get(hand.ActionId)
			context := trigger.NewContextWithData(context.Background(), &trigger.ContextData{Attrs: startAttrs, HandlerCfg: hand})
			replyCode, replyData, err := t.runner.Run(context, actID, hand.ActionId, nil)
			return replyCode, replyData, err
		}
	}

	//calling default handler if method name not specified
	for _, hand := range handlers {
		if len(hand.GetSetting("methodName")) == 0 {
			log.Debug("Default Dispatch Invoked: ", hand.ActionId)
			actID := action.Get(hand.ActionId)
			context := trigger.NewContextWithData(context.Background(), &trigger.ContextData{Attrs: startAttrs, HandlerCfg: hand})
			replyCode, replyData, err := t.runner.Run(context, actID, hand.ActionId, nil)
			return replyCode, replyData, err
		}
	}

	log.Error("Dispatch not found")
	return 0, nil, errors.New("Dispatch not found")
}
