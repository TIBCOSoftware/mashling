package tester

import (
	"context"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	FLOW_REF = "github.com/TIBCOSoftware/flogo-contrib/action/flow"
)

// RequestProcessor processes request objects and invokes the corresponding
// flow Manager methods
type RequestProcessor struct {
	runner action.Runner
}

// NewRequestProcessor creates a new RequestProcessor
func NewRequestProcessor() *RequestProcessor {

	var rp RequestProcessor
	rp.runner = runner.NewDirect()

	return &rp
}

// StartFlow handles a StartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a StartRequest.
func (rp *RequestProcessor) StartFlow(startRequest *StartRequest) (results map[string]*data.Attribute, err error) {

	execOptions := &instance.ExecOptions{Interceptor: startRequest.Interceptor, Patch: startRequest.Patch}

	attrs := startRequest.Attrs

	dataLen := len(startRequest.Data)

	// attrs, not supplied so attempt to create attrs from Data
	if attrs == nil && dataLen > 0 {
		attrs = make([]*data.Attribute, 0, dataLen)

		for k, v := range startRequest.Data {

			//todo handle error
			t, _ := data.GetType(v)
			attr,_ := data.NewAttribute(k, t, v)
			attrs = append(attrs, attr)
		}
	}

	factory := action.GetFactory(FLOW_REF)
	act := factory.New(&action.Config{Id: "flow"})

	ctx := trigger.NewContext(context.Background(), attrs)

	ro := &instance.RunOptions{Op: instance.OpStart, ReturnID: true,  FlowURI: startRequest.FlowURI, ExecOptions: execOptions}
	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = ro

	return rp.runner.RunAction(ctx, act, newOptions)
}

// RestartFlow handles a RestartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a RestartRequest.
func (rp *RequestProcessor) RestartFlow(restartRequest *RestartRequest)  (results map[string]*data.Attribute, err error) {

	execOptions := &instance.ExecOptions{Interceptor: restartRequest.Interceptor, Patch: restartRequest.Patch}

	ctx := context.Background()

	if restartRequest.Data != nil {

		logger.Debugf("Updating flow attrs: %v", restartRequest.Data)
		attrs := make([]*data.Attribute, len(restartRequest.Data))

		for k, v := range restartRequest.Data {
			attr,_ := data.NewAttribute(k, data.ANY, v)
			attrs = append(attrs,attr)
		}

		ctx = trigger.NewContext(context.Background(), attrs)
	}

	factory := action.GetFactory(FLOW_REF)
	act := factory.New(&action.Config{Id: "flow"})

	ro := &instance.RunOptions{Op: instance.OpRestart, ReturnID: true, FlowURI: restartRequest.InitialState.FlowURI, InitialState: restartRequest.InitialState, ExecOptions: execOptions}
	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = ro

	return rp.runner.RunAction(ctx, act, newOptions)
}

// ResumeFlow handles a ResumeRequest for a FlowInstance.  This will
// queue a RestartRequest.
func (rp *RequestProcessor) ResumeFlow(resumeRequest *ResumeRequest)  (results map[string]*data.Attribute, err error) {

	execOptions := &instance.ExecOptions{Interceptor: resumeRequest.Interceptor, Patch: resumeRequest.Patch}

	ctx := context.Background()

	if resumeRequest.Data != nil {

		logger.Debugf("Updating flow attrs: %v", resumeRequest.Data)
		attrs := make([]*data.Attribute, len(resumeRequest.Data))

		for k, v := range resumeRequest.Data {
			attr,_ := data.NewAttribute(k, data.ANY, v)
			attrs = append(attrs,attr)
		}

		ctx = trigger.NewContext(context.Background(), attrs)
	}

	factory := action.GetFactory(FLOW_REF)
	act := factory.New(&action.Config{Id: "flow"})

	ro := &instance.RunOptions{Op: instance.OpResume, ReturnID: true, FlowURI:resumeRequest.State.FlowURI, InitialState : resumeRequest.State, ExecOptions: execOptions}
	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = ro

	return rp.runner.RunAction(ctx, act, newOptions)
}

// StartRequest describes a request for starting a FlowInstance
type StartRequest struct {
	FlowURI     string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Attrs       []*data.Attribute      `json:"attrs"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}

// RestartRequest describes a request for restarting a FlowInstance
// todo: can be merged into StartRequest
type RestartRequest struct {
	InitialState *instance.Instance     `json:"initialState"`
	Data         map[string]interface{} `json:"data"`
	Interceptor  *support.Interceptor   `json:"interceptor"`
	Patch        *support.Patch         `json:"patch"`
}

// ResumeRequest describes a request for resuming a FlowInstance
//todo: Data for resume request should be directed to waiting task
type ResumeRequest struct {
	State       *instance.Instance     `json:"state"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
}
