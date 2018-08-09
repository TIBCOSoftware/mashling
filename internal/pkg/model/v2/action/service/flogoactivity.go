package service

import (
	"errors"
	"fmt"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/imdario/mergo"
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
	flogoActivityService.Request = req
	err = flogoActivityService.setRequestValues(settings)
	return flogoActivityService, err
}

// UpdateRequest updates a request on an existing FlogoActivity service instance with new values.
func (f *FlogoActivity) UpdateRequest(values map[string]interface{}) (err error) {
	return f.setRequestValues(values)
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

func (f *FlogoActivity) setRequestValues(settings map[string]interface{}) (err error) {
	for k, v := range settings {
		switch k {
		case "ref":
			ref, ok := v.(string)
			if !ok {
				return errors.New("invalid type for ref")
			}
			f.Request.Ref = ref
		case "inputs":
			inputs, ok := v.(map[string]interface{})
			if !ok {
				return errors.New("invalid type for inputs")
			}
			if err := mergo.Merge(&f.Request.Inputs, inputs, mergo.WithOverride); err != nil {
				return errors.New("unable to merge inputs values")
			}
		default:
			// ignore and move on.
		}
	}
	return nil
}

// FlogoActivityContext is an activity context in a mashling flow.
type FlogoActivityContext struct {
	details      activity.FlowDetails
	activityHost activity.Host
	TaskNameVal  string
	Attrs        map[string]*data.Attribute

	metadata *activity.Metadata
	inputs   map[string]*data.Attribute
	outputs  map[string]*data.Attribute
}

type FlogoActivityHost struct {
	HostID  string
	HostRef string

	IoMetadata    *data.IOMetadata
	HostData      data.Scope
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

	input := map[string]*data.Attribute{"Input1": data.NewZeroAttribute("Input1", data.TypeString)}
	output := map[string]*data.Attribute{"Output1": data.NewZeroAttribute("Output1", data.TypeString)}

	tc := &FlogoActivityContext{
		details: fd,
		activityHost: &FlogoActivityHost{
			HostID:     "1",
			HostRef:    "github.com/TIBCOSoftware/flogo-contrib/action/flow",
			IoMetadata: &data.IOMetadata{Input: input, Output: output},
			HostData:   data.NewSimpleScope(nil, nil),
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

func (ac *FlogoActivityHost) Name() string {
	return ""
}

func (ac *FlogoActivityHost) IOMetadata() *data.IOMetadata {
	return nil
}

func (ac *FlogoActivityHost) ID() string {
	return ac.HostID
}

func (ac *FlogoActivityHost) Ref() string {
	return ac.HostRef
}

func (ac *FlogoActivityHost) Reply(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

func (ac *FlogoActivityHost) Return(data map[string]*data.Attribute, err error) {
	//todo log reply
	ac.ReplyDataAttr = data
	ac.ReplyErr = err
}

func (ac *FlogoActivityHost) WorkingData() data.Scope {
	return ac.HostData
}

func (ac *FlogoActivityHost) GetResolver() data.Resolver {
	return definition.GetDataResolver()
}

// FlowDetails implements activity.Context.FlowDetails
func (c *FlogoActivityContext) FlowDetails() activity.FlowDetails {
	return c.details
}

// ActionContext implements activity.Context.FlowDetails
func (c *FlogoActivityContext) ActivityHost() activity.Host {
	return c.activityHost
}

// TaskName implements activity.Context.TaskName
func (c *FlogoActivityContext) TaskName() string {
	return c.TaskNameVal
}

func (c *FlogoActivityContext) Name() string {
	return c.TaskNameVal
}

// GetSetting implements activity.Context.GetSetting
func (c *FlogoActivityContext) GetSetting(setting string) (value interface{}, exists bool) {

	return nil, false
}

// GetInitValue implements activity.Context.GetInitValue
func (c *FlogoActivityContext) GetInitValue(key string) (value interface{}, exists bool) {
	return nil, false
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
