package test

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
)

// NewTestActivityContext creates a new TestActivityContext
func NewTestActivityContext(metadata *activity.Metadata) *TestActivityContext {

	input := []*data.Attribute{data.NewZeroAttribute("Input1", data.STRING)}
	output := []*data.Attribute{data.NewZeroAttribute("Output1", data.STRING)}

	ac := &TestActionCtx{
		ActionId:   "1",
		ActionRef:  "github.com/TIBCOSoftware/flogo-contrib/action/flow",
		ActionMd:   &action.ConfigMetadata{Input: input, Output: output},
		ActionData: data.NewSimpleScope(nil, nil),
	}

	return NewTestActivityContextWithAction(metadata, ac)
}

// NewTestActivityContextWithAction creates a new TestActivityContext
func NewTestActivityContextWithAction(metadata *activity.Metadata, actionCtx *TestActionCtx) *TestActivityContext {

	fd := &TestFlowDetails{
		FlowIDVal:   "1",
		FlowNameVal: "Test Flow",
	}

	tc := &TestActivityContext{
		details:     fd,
		ACtx:        actionCtx,
		TaskNameVal: "Test Task",
		Attrs:       make(map[string]*data.Attribute),
		inputs:      make(map[string]*data.Attribute, len(metadata.Input)),
		outputs:     make(map[string]*data.Attribute, len(metadata.Output)),
	}

	for _, element := range metadata.Input {
		tc.inputs[element.Name()] = data.NewZeroAttribute(element.Name(), element.Type())
	}
	for _, element := range metadata.Output {
		tc.outputs[element.Name()] = data.NewZeroAttribute(element.Name(), element.Type())
	}

	return tc
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestFlowDetails

// TestFlowDetails simple FlowDetails for use in testing
type TestFlowDetails struct {
	FlowIDVal   string
	FlowNameVal string
}

// ID implements activity.FlowDetails.ID
func (fd *TestFlowDetails) ID() string {
	return fd.FlowIDVal
}

// Name implements activity.FlowDetails.Name
func (fd *TestFlowDetails) Name() string {
	return fd.FlowNameVal
}

// ReplyHandler implements activity.FlowDetails.ReplyHandler
func (fd *TestFlowDetails) ReplyHandler() activity.ReplyHandler {
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestActionCtx

type TestActionCtx struct {
	ActionId  string
	ActionRef string
	ActionMd  *action.ConfigMetadata

	ActionData    data.Scope
	ReplyData     map[string]interface{}
	ReplyDataAttr map[string]*data.Attribute
	ReplyErr      error
}

func (ac *TestActionCtx) ID() string {
	return ac.ActionId
}

func (ac *TestActionCtx) Ref() string {
	return ac.ActionRef
}

func (ac *TestActionCtx) InstanceMetadata() *action.ConfigMetadata {
	return ac.ActionMd
}

func (ac *TestActionCtx) Reply(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

func (ac *TestActionCtx) Return(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

func (ac *TestActionCtx) WorkingData() data.Scope {
	return ac.ActionData
}

func (ac *TestActionCtx) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestActivityContext

// TestActivityContext is a dummy ActivityContext to assist in testing
type TestActivityContext struct {
	details     activity.FlowDetails
	ACtx        *TestActionCtx
	TaskNameVal string
	Attrs       map[string]*data.Attribute

	metadata *activity.Metadata
	inputs   map[string]*data.Attribute
	outputs  map[string]*data.Attribute
}

// FlowDetails implements activity.Context.FlowDetails
func (c *TestActivityContext) FlowDetails() activity.FlowDetails {
	return c.details
}

// TaskName implements activity.Context.TaskName
func (c *TestActivityContext) TaskName() string {
	return c.TaskNameVal
}

// GetAttrType implements data.Scope.GetAttrType
func (c *TestActivityContext) GetAttrType(attrName string) (attrType data.Type, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Type(), true
	}

	return 0, false
}

// GetAttrValue implements data.Scope.GetAttrValue
func (c *TestActivityContext) GetAttrValue(attrName string) (value string, exists bool) {

	attr, found := c.Attrs[attrName]

	if found {
		return attr.Value().(string), true
	}

	return "", false
}

// SetAttrValue implements data.Scope.SetAttrValue
func (c *TestActivityContext) SetAttrValue(attrName string, value string) {

	attr, found := c.Attrs[attrName]

	if found {
		attr.SetValue(value)
	}
}

// SetInput implements activity.Context.SetInput
func (c *TestActivityContext) SetInput(name string, value interface{}) {

	attr, found := c.inputs[name]

	if found {
		attr.SetValue(value)
	} else {
		//error?
	}
}

// GetInput implements activity.Context.GetInput
func (c *TestActivityContext) GetInput(name string) interface{} {

	attr, found := c.inputs[name]

	if found {
		return attr.Value()
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (c *TestActivityContext) SetOutput(name string, value interface{}) {

	attr, found := c.outputs[name]

	if found {
		attr.SetValue(value)
	} else {
		//error?
	}
}

// GetOutput implements activity.Context.GetOutput
func (c *TestActivityContext) GetOutput(name string) interface{} {

	attr, found := c.outputs[name]

	if found {
		return attr.Value()
	}

	return nil
}

func (c *TestActivityContext) ActionContext() action.Context {
	return c.ACtx
}
