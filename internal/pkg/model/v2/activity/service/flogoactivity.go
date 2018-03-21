package service

import (
	"errors"
	"fmt"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// FlogoActivity is a Flogo activity service.
type FlogoActivity struct {
	Request  FlogoActivityRequest  `json:"request"`
	Response FlogoActivityResponse `json:"response"`
}

// FlogoActivityRequest is a Flogo activity service request.
type FlogoActivityRequest struct {
	Ref    string                 `json:"ref"`
	Inputs map[string]interface{} `json:"inputs"`
}

// FlogoActivityResponse is a Flogo activity service response.
type FlogoActivityResponse struct {
	Done    bool                   `json:"done"`
	Error   string                 `json:"error"`
	Outputs map[string]interface{} `json:"outputs"`
}

// InitializeFlogoActivity initializes a FlogoActivity service with provided settings.
func InitializeFlogoActivity(settings map[string]interface{}) (flogoActivityService *FlogoActivity, err error) {
	flogoActivityService = &FlogoActivity{}
	req := FlogoActivityRequest{}
	req.Inputs = make(map[string]interface{})
	for k, v := range settings {
		switch k {
		case "ref":
			ref, ok := v.(string)
			if !ok {
				return flogoActivityService, errors.New("invalid type for ref")
			}
			req.Ref = ref
		case "inputs":
			inputs, ok := v.(map[string]interface{})
			if !ok {
				return flogoActivityService, errors.New("invalid type for inputs")
			}
			req.Inputs = inputs
		default:
			// ignore and move on.
		}
		flogoActivityService.Request = req
	}
	return flogoActivityService, err
}

// Execute invokes this FlogoActivity service.
func (f *FlogoActivity) Execute() (err error) {
	fa := activity.Get(f.Request.Ref)
	if fa == nil {
		return fmt.Errorf("unable to find Flogo activity: %s", f.Request.Ref)
	}
	var done bool
	actContext := NewFlogoActivityContext(fa.Metadata())
	actContext.TaskNameVal = f.Request.Ref
	for name, value := range f.Request.Inputs {
		actContext.SetInput(name, value)
	}
	done, err = fa.Eval(actContext)
	f.Response = FlogoActivityResponse{}
	f.Response.Done = done
	if err != nil {
		f.Response.Error = err.Error()
	}
	f.Response.Outputs = actContext.GetOutputs()
	return err
}

// FlogoActivityContext is an activity context in a mashling flow.
type FlogoActivityContext struct {
	details     activity.FlowDetails
	ACtx        *MashlingActionContext
	TaskNameVal string
	Attrs       map[string]*data.Attribute

	metadata *activity.Metadata
	inputs   map[string]*data.Attribute
	outputs  map[string]*data.Attribute
}

// MashlingActionContext is an action context in a mashling flow.
type MashlingActionContext struct {
	ActionID  string
	ActionRef string
	ActionMd  *action.ConfigMetadata

	ActionData    data.Scope
	ReplyData     map[string]interface{}
	ReplyDataAttr map[string]*data.Attribute
	ReplyErr      error
}

// MashlingFlowDetails simple FlowDetails for use in a mashling flow.
type MashlingFlowDetails struct {
	FlowIDVal   string
	FlowNameVal string
}

// NewFlogoActivityContext creates a new FlogoActivityContext
func NewFlogoActivityContext(metadata *activity.Metadata) *FlogoActivityContext {

	fd := &MashlingFlowDetails{
		FlowIDVal:   "1",
		FlowNameVal: "Mashling Core",
	}
	input := []*data.Attribute{data.NewZeroAttribute("Input1", data.STRING)}
	output := []*data.Attribute{data.NewZeroAttribute("Output1", data.STRING)}

	tc := &FlogoActivityContext{
		details: fd,
		ACtx: &MashlingActionContext{
			ActionID:   "1",
			ActionRef:  "github.com/TIBCOSoftware/flogo-contrib/action/flow",
			ActionMd:   &action.ConfigMetadata{Input: input, Output: output},
			ActionData: data.NewSimpleScope(nil, nil),
		},
		TaskNameVal: "Mashling Core Flogo Activity",
		Attrs:       make(map[string]*data.Attribute),
		inputs:      make(map[string]*data.Attribute, len(metadata.Input)),
		outputs:     make(map[string]*data.Attribute, len(metadata.Output)),
	}

	for _, element := range metadata.Input {
		tc.inputs[element.Name()], _ = data.NewAttribute(element.Name(), element.Type(), nil)
	}
	for _, element := range metadata.Output {
		tc.outputs[element.Name()], _ = data.NewAttribute(element.Name(), element.Type(), nil)
	}

	return tc
}

// ID implements activity.FlowDetails.ID
func (fd *MashlingFlowDetails) ID() string {
	return fd.FlowIDVal
}

// Name implements activity.FlowDetails.Name
func (fd *MashlingFlowDetails) Name() string {
	return fd.FlowNameVal
}

// ReplyHandler implements activity.FlowDetails.ReplyHandler
func (fd *MashlingFlowDetails) ReplyHandler() activity.ReplyHandler {
	return nil
}

// FlowDetails implements activity.Context.FlowDetails
func (c *FlogoActivityContext) FlowDetails() activity.FlowDetails {
	return c.details
}

// ID implements action.Context.ID
func (ac *MashlingActionContext) ID() string {
	return ac.ActionID
}

// Ref implements action.Context.Ref
func (ac *MashlingActionContext) Ref() string {
	return ac.ActionRef
}

// InstanceMetadata implements action.Context.InstanceMetadata
func (ac *MashlingActionContext) InstanceMetadata() *action.ConfigMetadata {
	return ac.ActionMd
}

// Reply implements action.Context.Reply
func (ac *MashlingActionContext) Reply(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

// Return implements action.Context.Return
func (ac *MashlingActionContext) Return(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

// WorkingData implements action.Context.WorkingData
func (ac *MashlingActionContext) WorkingData() data.Scope {
	return ac.ActionData
}

// GetResolver implements action.Context.GetResolver
func (ac *MashlingActionContext) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

// ActionContext implements activity.Context.FlowDetails
func (c *FlogoActivityContext) ActionContext() action.Context {
	return c.ACtx
}

// TaskName implements activity.Context.TaskName
func (c *FlogoActivityContext) TaskName() string {
	return c.TaskNameVal
}

// GetAttrType implements data.Scope.GetAttrType
func (c *FlogoActivityContext) GetAttrType(attrName string) (attrType data.Type, exists bool) {
	attr, found := c.Attrs[attrName]
	if found {
		return attr.Type(), true
	}
	return 0, false
}

// GetAttrValue implements data.Scope.GetAttrValue
func (c *FlogoActivityContext) GetAttrValue(attrName string) (value string, exists bool) {
	attr, found := c.Attrs[attrName]
	if found {
		return attr.Value().(string), true
	}
	return "", false
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *FlogoActivityContext) SetAttrValue(attrName string, value string) {
	attr, found := c.Attrs[attrName]
	if found {
		attr.SetValue(value)
	}
}

// SetInput implements activity.Context.SetInput
func (c *FlogoActivityContext) SetInput(name string, value interface{}) {
	attr, found := c.inputs[name]
	if found {
		attr.SetValue(value)
	}
}

// GetInput implements activity.Context.GetInput
func (c *FlogoActivityContext) GetInput(name string) interface{} {
	attr, found := c.inputs[name]
	if found {
		return attr.Value()
	}
	return nil
}

// SetOutput implements activity.Context.SetOutput
func (c *FlogoActivityContext) SetOutput(name string, value interface{}) {
	attr, found := c.outputs[name]
	if found {
		attr.SetValue(value)
	}
}

// GetOutput implements activity.Context.GetOutput
func (c *FlogoActivityContext) GetOutput(name string) interface{} {
	attr, found := c.outputs[name]
	if found {
		return attr.Value()
	}
	return nil
}

// GetOutputs implements activity.Context.GetOutput
func (c *FlogoActivityContext) GetOutputs() map[string]interface{} {
	rawOutput := make(map[string]interface{})
	for name, attr := range c.outputs {
		rawOutput[name] = attr.Value()
	}
	return rawOutput
}
