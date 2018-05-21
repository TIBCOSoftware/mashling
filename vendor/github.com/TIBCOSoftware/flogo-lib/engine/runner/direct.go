package runner

import (
	"context"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
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

// Deprecated: Use Execute() instead
func (runner *DirectRunner) Run(ctx context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {

	return 0, nil, errors.New("unsupported")
}

// Deprecated: Use Execute() instead
func (runner *DirectRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {

	return nil, errors.New("unsupported")
}

// Execute implements action.Runner.Execute
func (runner *DirectRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {

	if act == nil {
		return nil, errors.New("Action not specified")
	}

	md := action.GetMetadata(act)

	if !md.Async {
		syncAct := act.(action.SyncAction)
		return syncAct.Run(ctx, inputs)
	} else {
		asyncAct := act.(action.AsyncAction)

		handler := &SyncResultHandler{done: make(chan bool, 1)}

		err = asyncAct.Run(ctx, inputs, handler)

		if err != nil {
			return nil, err
		}

		<-handler.done

		return handler.Result()
	}

}

// SyncResultHandler simple result handler to use in synchronous case
type SyncResultHandler struct {
	done       chan bool
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
