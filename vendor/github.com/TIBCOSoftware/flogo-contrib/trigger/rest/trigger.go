package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/rest/cors"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/julienschmidt/httprouter"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

const (
	REST_CORS_PREFIX = "REST_TRIGGER"
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-rest")

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

// RestTrigger REST trigger struct
type RestTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	server   *Server
	config   *trigger.Config
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

func (t *RestTrigger) Init(runner action.Runner) {

	router := httprouter.New()

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	if _, ok := t.config.Settings["port"]; !ok {
		panic(fmt.Sprintf("No Port found for trigger '%s' in settings", t.config.Id))
	}

	addr := ":" + t.config.GetSetting("port")
	t.runner = runner

	// Init handlers
	for _, handlerCfg := range t.config.Handlers {

		if handlerIsValid(handlerCfg) {
			method := strings.ToUpper(handlerCfg.GetSetting("method"))
			path := handlerCfg.GetSetting("path")

			log.Debugf("REST Trigger: Registering handler [%s: %s] for Action Id: [%s]", method, path, handlerCfg.ActionId)

			router.OPTIONS(path, handleCorsPreflight) // for CORS
			router.Handle(method, path, newActionHandler(t, handlerCfg.ActionId, handlerCfg))

		} else {
			panic(fmt.Sprintf("Invalid handler: %v", handlerCfg))
		}
	}

	log.Debugf("REST Trigger: Configured on port %s", t.config.Settings["port"])
	t.server = NewServer(addr, router)
}

func (t *RestTrigger) Start() error {
	return t.server.Start()
}

// Stop implements util.Managed.Stop
func (t *RestTrigger) Stop() error {
	return t.server.Stop()
}

// Handles the cors preflight request
func handleCorsPreflight(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	log.Infof("Received [OPTIONS] request to CorsPreFlight: %+v", r)

	c := cors.New(REST_CORS_PREFIX, log)
	c.HandlePreflight(w, r)
}

// IDResponse id response object
type IDResponse struct {
	ID string `json:"id"`
}

func newActionHandler(rt *RestTrigger, actionId string, handlerCfg *trigger.HandlerConfig) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Infof("REST Trigger: Received request for id '%s'", rt.config.Id)

		c := cors.New(REST_CORS_PREFIX, log)
		c.WriteCorsActualRequestHeaders(w)

		pathParams := make(map[string]string)
		for _, param := range ps {
			pathParams[param.Key] = param.Value
		}

		var content interface{}
		err := json.NewDecoder(r.Body).Decode(&content)
		if err != nil {
			switch {
			case err == io.EOF:
				// empty body
				//todo should handler say if content is expected?
			case err != nil:
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		queryValues := r.URL.Query()
		queryParams := make(map[string]string, len(queryValues))

		for key, value := range queryValues {
			queryParams[key] = strings.Join(value, ",")
		}

		outData := map[string]interface{}{
			"params":      pathParams,
			"pathParams":  pathParams,
			"queryParams": queryParams,
			"content":     content,
		}

		//todo handle error
		startAttrs, _ := rt.metadata.OutputsToAttrs(outData, false)

		act := action.Get(actionId)
		ctx := trigger.NewInitialContext(startAttrs, handlerCfg)
		results, err := rt.runner.RunAction(ctx, act, nil)

		var replyData interface{}
		var replyCode int

		if len(results) != 0 {
			dataAttr, ok := results["data"]
			if ok {
				replyData = dataAttr.Value()
			}
			codeAttr, ok := results["code"]
			if ok {
				replyCode, _ = data.CoerceToInteger(codeAttr.Value())
			}
		}

		if err != nil {
			log.Debugf("REST Trigger Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if replyData != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if replyCode == 0 {
				replyCode = 200
			}
			w.WriteHeader(replyCode)
			if err := json.NewEncoder(w).Encode(replyData); err != nil {
				log.Error(err)
			}
			return
		}

		if replyCode > 0 {
			w.WriteHeader(replyCode)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
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
