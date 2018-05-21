package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/rest/cors"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/julienschmidt/httprouter"
)

const (
	REST_CORS_PREFIX = "REST_TRIGGER"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-rest")

var validMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}

// RestTrigger REST trigger struct
type RestTrigger struct {
	metadata *trigger.Metadata
	//runner   action.Runner
	server *Server
	config *trigger.Config
	//handlers []*handler.Handler
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

func (t *RestTrigger) Initialize(ctx trigger.InitContext) error {
	router := httprouter.New()

	if t.config.Settings == nil {
		return fmt.Errorf("no Settings found for trigger '%s'", t.config.Id)
	}

	if _, ok := t.config.Settings["port"]; !ok {
		return fmt.Errorf("no Port found for trigger '%s' in settings", t.config.Id)
	}

	addr := ":" + t.config.GetSetting("port")

	pathMap := make(map[string]string)

	// Init handlers
	for _, handler := range ctx.GetHandlers() {

		err := validateHandler(t.config.Id, handler)
		if err != nil {
			return err
		}

		method := strings.ToUpper(handler.GetStringSetting("method"))
		path := handler.GetStringSetting("path")

		log.Debugf("Registering handler [%s: %s]", method, path)

		if _, ok := pathMap[path]; !ok {
			pathMap[path] = path
			router.OPTIONS(path, handleCorsPreflight) // for CORS
		}

		//router.OPTIONS(path, handleCorsPreflight) // for CORS
		router.Handle(method, path, newActionHandler(t, handler))
	}

	log.Debugf("Configured on port %s", t.config.Settings["port"])
	t.server = NewServer(addr, router)

	return nil
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

func newActionHandler(rt *RestTrigger, handler *trigger.Handler) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		log.Infof("Received request for id '%s'", rt.config.Id)

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
		header := make(map[string]string, len(r.Header))

		for key, value := range r.Header {
			header[key] = strings.Join(value, ",")
		}

		for key, value := range queryValues {
			queryParams[key] = strings.Join(value, ",")
		}

		triggerData := map[string]interface{}{
			"params":      pathParams,
			"pathParams":  pathParams,
			"queryParams": queryParams,
			"header":      header,
			"content":     content,
		}

		results, err := handler.Handle(context.Background(), triggerData)

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

func validateHandler(triggerId string, handler *trigger.Handler) error {

	method := handler.GetStringSetting("method")
	if method == "" {
		return fmt.Errorf("no method specified in handler for trigger '%s'", triggerId)
	}

	if !stringInList(strings.ToUpper(method), validMethods) {
		return fmt.Errorf("invalid method '%s' specified in handler for trigger '%s'", method, triggerId)
	}

	//validate path

	return nil
}

func stringInList(str string, list []string) bool {
	for _, value := range list {
		if value == str {
			return true
		}
	}
	return false
}
