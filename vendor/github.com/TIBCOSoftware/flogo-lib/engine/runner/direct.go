package runner

import (
	"context"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// DirectRunner runs an action synchronously
type DirectRunner struct {
}

// NewDirectRunner create a new DirectRunner
func NewDirect() *DirectRunner {
	return &DirectRunner{}
}

// Start will start the engine, by starting all of its workers
func (runner *DirectRunner) Start() error {
	//op-op
	return nil
}

// Stop will stop the engine, by stopping all of its workers
func (runner *DirectRunner) Stop() error {
	//no-op
	return nil
}

//Run
//Deprecated
func (runner *DirectRunner) Run(ctx context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {

	if act == nil {
		return 0, nil, errors.New("Action not specified")
	}

	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = options

	handler := &SyncResultHandler{done: make(chan bool, 1)}

	var ctxData *trigger.ContextData

	if ctx != nil {
		var exists bool
		ctxData, exists = trigger.ExtractContextData(ctx)

		if !exists {
			logger.Warn("Trigger data not applied to context")
		}
	}

	inputs := generateInputs(act, ctxData)

	err = act.Run(ctx, inputs, newOptions, handler)

	if err != nil {
		return 0, nil, err
	}

	<-handler.done

	ndata, err := handler.Result()

	results := generateOutputs(act, ctxData, ndata)

	if len(ndata) != 0 {
		defData, ok := results["data"]
		if ok {
			data = defData.Value()
		}
		defCode, ok := results["code"]
		if ok {
			code = defCode.Value().(int)
		}
	}

	return code, data, err
}

// Run the specified action
func (runner *DirectRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {

	if act == nil {
		return nil, errors.New("Action not specified")
	}

	handler := &SyncResultHandler{done: make(chan bool, 1)}

	ctxData, exists := trigger.ExtractContextData(ctx)

	if !exists {
		logger.Warn("Trigger data not applied to context")
	}

	inputs := generateInputs(act, ctxData)

	err = act.Run(ctx, inputs, options, handler)

	if err != nil {
		return nil, err
	}

	<-handler.done

	actionOutput, err := handler.Result()

	if err != nil {
		return nil, err
	}

	return generateOutputs(act, ctxData, actionOutput), nil
}

// SyncResultHandler simple result handler to use in synchronous case
type SyncResultHandler struct {
	done       chan (bool)
	resultData map[string]*data.Attribute
	err        error
	set        bool
}

// HandleResult implements action.ResultHandler.HandleResult
func (rh *SyncResultHandler) HandleResult(resultData map[string]*data.Attribute, err error) {

	if !rh.set {
		rh.set = true
		rh.resultData = resultData
		rh.err = err
	}
}

// Done implements action.ResultHandler.Done
func (rh *SyncResultHandler) Done() {
	rh.done <- true
}

// Result returns the latest Result set on the handler
func (rh *SyncResultHandler) Result() (resultData map[string]*data.Attribute, err error) {
	return rh.resultData, rh.err
}
