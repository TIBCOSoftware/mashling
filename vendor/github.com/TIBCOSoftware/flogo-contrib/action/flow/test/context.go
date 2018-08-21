package test

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

//todo needs to move to lib
// NewTestActivityContext creates a new TestActivityContext
func NewTestActivityContext(metadata *activity.Metadata) *TestActivityContext {

	input := map[string]*data.Attribute{"Input1": data.NewZeroAttribute("Input1", data.TypeString)}
	output := map[string]*data.Attribute{"Output1": data.NewZeroAttribute("Output1", data.TypeString)}

	ac := &TestActivityHost{
		HostId:     "1",
		HostRef:    "github.com/TIBCOSoftware/flogo-contrib/action/flow",
		IoMetadata: &data.IOMetadata{Input: input, Output: output},
		HostData:   data.NewSimpleScope(nil, nil),
	}

	return NewTestActivityContextWithAction(metadata, ac)
}

// NewTestActivityContextWithAction creates a new TestActivityContext
func NewTestActivityContextWithAction(metadata *activity.Metadata, activityHost *TestActivityHost) *TestActivityContext {

	fd := &TestFlowDetails{
		FlowIDVal:   "1",
		FlowNameVal: "Test Flow",
	}

	tc := &TestActivityContext{
		metadata:     metadata,
		details:      fd,
		activityHost: activityHost,
		TaskNameVal:  "Test TaskOld",
		Attrs:        make(map[string]*data.Attribute),
		inputs:       make(map[string]*data.Attribute, len(metadata.Input)),
		outputs:      make(map[string]*data.Attribute, len(metadata.Output)),
		settings:     make(map[string]*data.Attribute, len(metadata.Settings)),
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
// TestActivityHost

type TestActivityHost struct {
	HostId  string
	HostRef string

	IoMetadata    *data.IOMetadata
	HostData      data.Scope
	ReplyData     map[string]interface{}
	ReplyDataAttr map[string]*data.Attribute
	ReplyErr      error
}

func (ac *TestActivityHost) Name() string {
	return ""
}

func (ac *TestActivityHost) IOMetadata() *data.IOMetadata {
	return nil
}

func (ac *TestActivityHost) ID() string {
	return ac.HostId
}

func (ac *TestActivityHost) Ref() string {
	return ac.HostRef
}

func (ac *TestActivityHost) Reply(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

func (ac *TestActivityHost) Return(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

func (ac *TestActivityHost) WorkingData() data.Scope {
	return ac.HostData
}

func (ac *TestActivityHost) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TestActivityContext

// TestActivityContext is a dummy ActivityContext to assist in testing
type TestActivityContext struct {
	details      activity.FlowDetails
	TaskNameVal  string
	Attrs        map[string]*data.Attribute
	activityHost activity.Host

	metadata *activity.Metadata
	settings map[string]*data.Attribute
	inputs   map[string]*data.Attribute
	outputs  map[string]*data.Attribute

	shared map[string]interface{}
}

func (c *TestActivityContext) FlowDetails() activity.FlowDetails {
	return c.details
}

func (c *TestActivityContext) ActivityHost() activity.Host {
	return c.activityHost
}

func (c *TestActivityContext) Name() string {
	return c.TaskNameVal
}

// GetSetting implements activity.Context.GetSetting
func (c *TestActivityContext) GetSetting(setting string) (value interface{}, exists bool) {

	attr, found := c.settings[setting]

	if found {
		return attr.Value(), true
	}

	return nil, false
}

// GetInitValue implements activity.Context.GetInitValue
func (c *TestActivityContext) GetInitValue(key string) (value interface{}, exists bool) {
	return nil, false
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

// SetInput implements activity.Context.SetInput
func (c *TestActivityContext) SetSetting(name string, value interface{}) {

	attr, found := c.metadata.Settings[name]
	if found {
		s, _ := data.NewAttribute(name, attr.Type(), value)
		c.settings[name] = s
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

func (c *TestActivityContext) GetSharedTempData() map[string]interface{} {

	if c.shared == nil {
		c.shared = make(map[string]interface{})
	}
	return c.shared
}
