/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package gorillamuxtrigger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/rest/cors"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/lib/util"

	condition "github.com/TIBCOSoftware/mashling/lib/conditions"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	REST_CORS_PREFIX = "REST_TRIGGER"
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-rest")

//OptimizedHandler optimized handler
type OptimizedHandler struct {
	defaultActionId   string
	settings          map[string]interface{}
	dispatches        []*Dispatch
	defaultHandlerCfg *trigger.HandlerConfig
}

//Dispatch holds dispatch actionId and condition
type Dispatch struct {
	actionId   string
	condition  string
	handlerCfg *trigger.HandlerConfig
}

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

// RestTrigger REST trigger struct
type RestTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	server   *Server
	config   *trigger.Config
	localIP  string
	tracer   util.Tracer
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &RestFactory{metadata: md}
}

// RestFactory REST Trigger factory
type RestFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *RestFactory) New(config *trigger.Config) trigger.Trigger {
	return &RestTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *RestTrigger) Metadata() *trigger.Metadata {
	return t.metadata
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

//Init trigger initialization
func (t *RestTrigger) Init(runner action.Runner) {

	// router := httprouter.New()
	router := mux.NewRouter()

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	if _, ok := t.config.Settings["port"]; !ok {
		panic(fmt.Sprintf("No Port found for trigger '%s' in settings", t.config.Id))
	}

	//Substitute for any environment variables referenced in the settings.
	//Expressions will be in the format ${env.SERVER_KEY} where SERVER_KEY is the env variable
	err := util.ResolveEnvironmentProperties(t.config.Settings)
	if err != nil {
		panic(fmt.Sprint(err))
	}

	addr := ":" + t.config.GetSetting("port")
	t.runner = runner
	t.localIP = getLocalIP()

	err = t.tracer.ConfigureTracer(t.config.Settings, t.localIP+":"+t.config.GetSetting("port"),
		t.config.Name)
	if err != nil {
		panic(err)
	}

	if isAuthEnabled(t.config.Settings) {
		setupAuth(t.config.Settings)
	}

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

	//Check whether client auth is enabled
	enableClientAuth := false
	trustStore := ""
	if _, ok := t.config.Settings["enableClientAuth"]; ok {
		enableClientAuthSetting, err := strconv.ParseBool(t.config.GetSetting("enableClientAuth"))
		if err == nil && enableClientAuthSetting {
			//Client auth is enabled. get client trust store (i.e. client CAs)
			enableClientAuth = true
			if _, ok := t.config.Settings["trustStore"]; !ok {
				panic(fmt.Sprintf("Client auth is enabled but client trust store is not provided for trigger '%s' in settings", t.config.Id))
			}
			trustStore = t.config.GetSetting("trustStore")
		}
	}

	//optimize flog-handlers i.e merge handlers having same settings
	optHandlers := []*OptimizedHandler{}
	for _, handler := range t.config.Handlers {
		//check if there is any handler already added with same settings
		handlerAdded := false
		for _, optHandler := range optHandlers {
			//loop through all settings
			settingsMatched := true
			for k, v := range optHandler.settings {
				if v != handler.Settings[k] {
					settingsMatched = false
					break
				}
			}
			if settingsMatched {
				//check for dispatch condition
				if dispatchCondition := handler.Settings[util.Flogo_Trigger_Handler_Setting_Condition]; dispatchCondition != nil {
					tmpDispatch := &Dispatch{
						actionId:   handler.ActionId,
						condition:  dispatchCondition.(string),
						handlerCfg: handler,
					}
					optHandler.dispatches = append(optHandler.dispatches, tmpDispatch)
				} else {
					//no dispatch condition, hence make it as default action
					optHandler.defaultActionId = handler.ActionId
					optHandler.defaultHandlerCfg = handler
				}
				handlerAdded = true
				break
			}
		}

		if !handlerAdded {
			tmpSettings := make(map[string]interface{})
			for k, v := range handler.Settings {
				if k != util.Flogo_Trigger_Handler_Setting_Condition {
					tmpSettings[k] = v
				}
			}

			var tmpDispatches []*Dispatch
			//check for dispatch condition
			if dispatchCondition := handler.Settings[util.Flogo_Trigger_Handler_Setting_Condition]; dispatchCondition != nil {
				tmpDispatch := &Dispatch{
					actionId:   handler.ActionId,
					condition:  handler.Settings[util.Flogo_Trigger_Handler_Setting_Condition].(string),
					handlerCfg: handler,
				}
				tmpDispatches = append(tmpDispatches, tmpDispatch)
			}

			optHandler := OptimizedHandler{
				defaultActionId:   handler.ActionId,
				settings:          tmpSettings,
				dispatches:        tmpDispatches,
				defaultHandlerCfg: handler,
			}

			optHandlers = append(optHandlers, &optHandler)
		}
	}

	// Init handlers
	for _, optHandler := range optHandlers {
		if handlerIsValid(optHandler) {
			method := strings.ToUpper(optHandler.settings["method"].(string))
			path := optHandler.settings["path"].(string)
			url := "http://"
			if enableTLS {
				url = "https://"
			}
			url += t.localIP + ":" + t.config.GetSetting("port") + path

			log.Debugf("REST Trigger: Registering handler [%s: %s] with default Action Id: [%s]", method, path, optHandler.defaultActionId)
			//Register Cross-Origin Resource Sharing (CORS) handler
			router.HandleFunc(path, handleCorsPreflight).
				Methods("OPTIONS")
			//register action handler
			router.HandleFunc(path, newActionHandler(t, optHandler, method, url)).
				Methods(method)
		}
	}

	log.Debugf("REST Trigger: Configured on port %s", t.config.Settings["port"])
	t.server = NewServer(addr, router, enableTLS, serverCert, serverKey, enableClientAuth, trustStore)
}

func (t *RestTrigger) Start() error {
	return t.server.Start()
}

// Stop implements util.Managed.Stop
func (t *RestTrigger) Stop() error {
	err := t.server.Stop()
	if err != nil {
		return err
	}
	return t.tracer.Close()
}

// Handles the cors preflight request
func handleCorsPreflight(w http.ResponseWriter, r *http.Request) {

	log.Infof("Received [OPTIONS] request to CorsPreFlight: %+v", r)

	c := cors.New(REST_CORS_PREFIX, log)
	c.HandlePreflight(w, r)
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func newActionHandler(rt *RestTrigger, handler *OptimizedHandler, method, url string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		log.Infof("REST Trigger: Received request for id '%s'", rt.config.Id)

		wireContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))

		var serverSpan opentracing.Span
		if err == nil {
			serverSpan = opentracing.StartSpan(
				r.URL.String(), ext.RPCServerOption(wireContext))
		} else {
			serverSpan = opentracing.StartSpan(r.URL.String())
		}
		defer serverSpan.Finish()

		serverSpan.SetTag("http.method", method)
		serverSpan.SetTag("http.url", url)

		ctx := opentracing.ContextWithSpan(context.Background(), serverSpan)

		c := cors.New(REST_CORS_PREFIX, log)
		c.WriteCorsActualRequestHeaders(w)

		// authenticate
		if isAuthEnabled(rt.config.Settings) {
			log.Debugf("Authenticating the request")
			if !authenticate(r, rt.config.Settings) {
				replyCode := http.StatusForbidden
				log.Debugf("Authentication failed, returning with status code[%d]", replyCode)
				serverSpan.SetTag("http.status_code", replyCode)
				w.WriteHeader(replyCode)
				return
			}
		}

		// get path params
		vars := mux.Vars(r)
		pathParams := make(map[string]string)
		for k, v := range vars {
			pathParams[k] = v
		}

		var content interface{}
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				serverSpan.SetTag("error", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mime := r.Header.Get("Content-Type")
			err = util.Unmarshal(mime, data, &content)
			if err != nil {
				serverSpan.SetTag("error", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		// get query params
		queryValues := r.URL.Query()
		queryParams := make(map[string]string, len(queryValues))

		for key, value := range queryValues {
			queryParams[key] = strings.Join(value, ",")
		}

		// get headers
		header := make(map[string]interface{})
		for key, value := range r.Header {
			//If header has single value, then add to map as a string
			//otherwise add as an array
			if len(value) == 1 {
				header[key] = value[0]
			} else {
				header[key] = value
			}
		}

		data := map[string]interface{}{
			"params":      pathParams,
			"pathParams":  pathParams,
			"queryParams": queryParams,
			"header":      header,
			"content":     content,
			"tracing":     ctx,
		}

		//pick action based on dispatch condition
		contentBytes, err := json.Marshal(content)
		contentStr := string(contentBytes)
		actionId := ""
		var handlerCfg *trigger.HandlerConfig

		for _, dispatch := range handler.dispatches {
			expressionStr := dispatch.condition
			//Get condtion and expression type
			conditionOperation, exprType, err := condition.GetConditionOperationAndExpressionType(expressionStr)

			if err != nil || exprType == condition.EXPR_TYPE_NOT_VALID {
				str := fmt.Sprintf("not able parse the condition '%v' mentioned for content based handler. skipping the handler.", expressionStr)
				serverSpan.SetTag("error", str)
				log.Error(str)
				continue
			}

			log.Debugf("Expression type: %v", exprType)
			log.Debugf("conditionOperation.LHS %v", conditionOperation.LHS)
			log.Debugf("conditionOperation.OperatorInfo %v", conditionOperation.OperatorInfo().Names)
			log.Debugf("conditionOperation.RHS %v", conditionOperation.RHS)

			//Resolve expression's LHS based on expression type and
			//evaluate the expression
			if exprType == condition.EXPR_TYPE_CONTENT {
				exprResult, err := condition.EvaluateCondition(*conditionOperation, contentStr)
				if err != nil {
					str := fmt.Sprintf("not able evaluate expression - %v with error - %v. skipping the handler.", expressionStr, err)
					serverSpan.SetTag("error", str)
					log.Error(str)
				}
				if exprResult {
					actionId = dispatch.actionId
					handlerCfg = dispatch.handlerCfg
				}
			} else if exprType == condition.EXPR_TYPE_HEADER {
				//resolve LHS i.e header value from http request
				headerVal := r.Header.Get(conditionOperation.LHS)
				log.Debugf("header key = %v, val = %v", conditionOperation.LHS, headerVal)
				if headerVal != "" {
					conditionOperation.LHS = headerVal
					op := conditionOperation.Operator
					exprResult := op.Eval(conditionOperation.LHS, conditionOperation.RHS)
					if exprResult {
						actionId = dispatch.actionId
						handlerCfg = dispatch.handlerCfg
					}
				}
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
						handlerCfg = dispatch.handlerCfg
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
			actionId = handler.defaultActionId
			handlerCfg = handler.defaultHandlerCfg
			log.Debugf("dispatch not resolved. Continue with default action - %v", actionId)
		}

		// upgrade to websocket protocol if wsUpgradeRequired is defined to "true" in handler settings
		wsUpgradeRequiredSetting, ok := handlerCfg.Settings["wsUpgradeRequired"]
		wsUpgradeRequired := false
		if ok && wsUpgradeRequiredSetting == "true" {
			wsUpgradeRequired = true
		}
		if wsUpgradeRequired {
			log.Debug("wsUpgradeRequired is defined to 'true'. upgrading to websocket protocol...")
			var upgrader = websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				serverSpan.SetTag("error", err.Error())
				log.Errorf("websocket protocol upgrade error: %s", err.Error())
				return
			}
			log.Info("connection upgraded to websocket protocol")
			//add connection object into context data
			data["wsconnection"] = conn
		}

		//todo handle error
		startAttrs, _ := rt.metadata.OutputsToAttrs(data, false)

		action := action.Get(actionId)
		log.Debugf("Found action' %+x'", action)

		context := trigger.NewContextWithData(context.Background(), &trigger.ContextData{Attrs: startAttrs, HandlerCfg: handlerCfg})

		replyCode, replyData, err := rt.runner.Run(context, action, actionId, nil)

		if err != nil {
			serverSpan.SetTag("error", err.Error())
			log.Debugf("REST Trigger Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if wsUpgradeRequired {
			// ignore replyCode, replyData and just return in case of websocket protocol
			return
		}

		if replyData != nil {
			if object, ok := replyData.(map[string]interface{}); ok {
				if mime, ok := object[util.MetaMIME]; ok {
					if s, ok := mime.(string); ok {
						w.Header().Set("Content-Type", s)
					}
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				}
				w.WriteHeader(replyCode)

				data, err := util.Marshal(replyData)
				if err != nil {
					serverSpan.SetTag("error", err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				_, err = io.Copy(w, bytes.NewReader(data))
				if err != nil {
					serverSpan.SetTag("error", err.Error())
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
			return
		}

		if replyCode > 0 {
			serverSpan.SetTag("http.status_code", replyCode)
			w.WriteHeader(replyCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Utils
func handlerIsValid(handler *OptimizedHandler) bool {
	if handler.settings == nil {
		return false
	}

	if handler.settings["method"] == "" {
		return false
	}

	if !stringInList(strings.ToUpper(handler.settings["method"].(string)), validMethods) {
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
