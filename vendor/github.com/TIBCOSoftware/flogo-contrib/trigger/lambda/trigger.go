package lambda

import (
	"context"
	"encoding/json"
	"flag"
	syslog "log"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	// Import the aws-lambda-go. Required for dep to pull on app create
	_ "github.com/aws/aws-lambda-go/lambda"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-lambda")
var singleton *LambdaTrigger

// LambdaTrigger AWS Lambda trigger struct
type LambdaTrigger struct {
	metadata *trigger.Metadata
	config   *trigger.Config
	handlers []*trigger.Handler
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &LambdaFactory{metadata: md}
}

// LambdaFactory AWS Lambda Trigger factory
type LambdaFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *LambdaFactory) New(config *trigger.Config) trigger.Trigger {
	singleton = &LambdaTrigger{metadata: t.metadata, config: config}
	return singleton
}

// Metadata implements trigger.Trigger.Metadata
func (t *LambdaTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

func (t *LambdaTrigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	return nil
}

// Invoke starts the trigger and invokes the action registered in the handler
func Invoke() (map[string]interface{}, error) {

	log.Info("Starting AWS Lambda Trigger")
	syslog.Println("Starting AWS Lambda Trigger")

	// Parse the flags
	flag.Parse()

	// Looking up the arguments
	evtArg := flag.Lookup("evt")
	var evt interface{}
	// Unmarshall evt
	if err := json.Unmarshal([]byte(evtArg.Value.String()), &evt); err != nil {
		return nil, err
	}

	log.Debugf("Received evt: '%+v'\n", evt)
	syslog.Printf("Received evt: '%+v'\n", evt)

	// Get the context
	ctxArg := flag.Lookup("ctx")
	var lambdaCtx interface{}

	// Unmarshal ctx
	if err := json.Unmarshal([]byte(ctxArg.Value.String()), &lambdaCtx); err != nil {
		return nil, err
	}

	log.Debugf("Received ctx: '%+v'\n", lambdaCtx)
	syslog.Printf("Received ctx: '%+v'\n", lambdaCtx)

	//select handler, use 0th for now
	handler := singleton.handlers[0]

	inputData := map[string]interface{}{
		"context": lambdaCtx,
		"evt":     evt,
	}

	results, err := handler.Handle(context.Background(), inputData)

	var replyData interface{}
	var replyStatus int

	if len(results) != 0 {
		dataAttr, ok := results["data"]
		if ok {
			replyData = dataAttr.Value()
		}
		code, ok := results["status"]
		if ok {
			replyStatus, _ = data.CoerceToInteger(code.Value())
		}
	}

	if err != nil {
		log.Debugf("Lambda Trigger Error: %s", err.Error())
		return nil, err
	}

	flowResponse := map[string]interface{}{
		"data":   replyData,
		"status": replyStatus,
	}
	return flowResponse, err
}

func (t *LambdaTrigger) Start() error {
	return nil
}

// Stop implements util.Managed.Stop
func (t *LambdaTrigger) Stop() error {
	return nil
}
