package tester

import (
	"context"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/instance"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
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

	logger.Debugf("Tester starting flow")

	factory := action.GetFactory(FLOW_REF)
	act, _ := factory.New(&action.Config{})

	var inputs map[string]*data.Attribute

	if len(startRequest.Attrs) > 0 {

		logger.Debugf("Starting with flow attrs: %#v", startRequest.Attrs)

		inputs = make(map[string]*data.Attribute, len(startRequest.Attrs)+1)
		for _, attr := range startRequest.Attrs {
			inputs[attr.Name()] = attr
		}
	} else if len(startRequest.Data) > 0 {

		logger.Debugf("Starting with flow attrs: %#v", startRequest.Data)

		inputs = make(map[string]*data.Attribute, len(startRequest.Data)+1)

		for k, v := range startRequest.Data {
			t, err := data.GetType(v)
			if err != nil {
				t = data.TypeAny
			}
			attr, _ := data.NewAttribute(k, t, v)
			inputs[k] = attr
		}
	} else {
		inputs = make(map[string]*data.Attribute, 1)
	}

	execOptions := &instance.ExecOptions{Interceptor: startRequest.Interceptor, Patch: startRequest.Patch}
	ro := &instance.RunOptions{Op: instance.OpStart, ReturnID: true, FlowURI: startRequest.FlowURI, ExecOptions: execOptions}
	attr, _ := data.NewAttribute("_run_options", data.TypeAny, ro)
	inputs[attr.Name()] = attr

	return rp.runner.Execute(context.Background(), act, inputs)
}

// RestartFlow handles a RestartRequest for a FlowInstance.  This will
// generate an ID for the new FlowInstance and queue a RestartRequest.
func (rp *RequestProcessor) RestartFlow(restartRequest *RestartRequest) (results map[string]*data.Attribute, err error) {

	logger.Debugf("Tester restarting flow")

	factory := action.GetFactory(FLOW_REF)
	act, _ := factory.New(&action.Config{})

	inputs := make(map[string]*data.Attribute, len(restartRequest.Data)+1)

	if restartRequest.Data != nil {

		logger.Debugf("Updating flow attrs: %v", restartRequest.Data)

		for k, v := range restartRequest.Data {
			attr, _ := data.NewAttribute(k, data.TypeAny, v)
			inputs[k] = attr
		}
	}

	execOptions := &instance.ExecOptions{Interceptor: restartRequest.Interceptor, Patch: restartRequest.Patch}
	ro := &instance.RunOptions{Op: instance.OpRestart, ReturnID: true, FlowURI: restartRequest.InitialState.FlowURI(), InitialState: restartRequest.InitialState, ExecOptions: execOptions}
	attr, _ := data.NewAttribute("_run_options", data.TypeAny, ro)
	inputs[attr.Name()] = attr

	return rp.runner.Execute(context.Background(), act, inputs)
}

// ResumeFlow handles a ResumeRequest for a FlowInstance.  This will
// queue a RestartRequest.
func (rp *RequestProcessor) ResumeFlow(resumeRequest *ResumeRequest) (results map[string]*data.Attribute, err error) {

	logger.Debugf("Tester resuming flow")

	factory := action.GetFactory(FLOW_REF)
	act, _ := factory.New(&action.Config{})

	inputs := make(map[string]*data.Attribute, len(resumeRequest.Data)+1)

	if resumeRequest.Data != nil {

		logger.Debugf("Updating flow attrs: %v", resumeRequest.Data)

		for k, v := range resumeRequest.Data {
			attr, _ := data.NewAttribute(k, data.TypeAny, v)
			inputs[k] = attr
		}
	}

	execOptions := &instance.ExecOptions{Interceptor: resumeRequest.Interceptor, Patch: resumeRequest.Patch}
	ro := &instance.RunOptions{Op: instance.OpResume, ReturnID: true, FlowURI: resumeRequest.State.FlowURI(), InitialState: resumeRequest.State, ExecOptions: execOptions}
	attr, _ := data.NewAttribute("_run_options", data.TypeAny, ro)
	inputs[attr.Name()] = attr

	return rp.runner.Execute(context.Background(), act, inputs)
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
	InitialState *instance.IndependentInstance `json:"initialState"`
	Data         map[string]interface{}        `json:"data"`
	Interceptor  *support.Interceptor          `json:"interceptor"`
	Patch        *support.Patch                `json:"patch"`
}

// ResumeRequest describes a request for resuming a FlowInstance
//todo: Data for resume request should be directed to waiting task
type ResumeRequest struct {
	State       *instance.IndependentInstance `json:"state"`
	Data        map[string]interface{}        `json:"data"`
	Interceptor *support.Interceptor          `json:"interceptor"`
	Patch       *support.Patch                `json:"patch"`
}
