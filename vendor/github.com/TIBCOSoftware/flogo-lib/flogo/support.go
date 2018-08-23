package flogo

import (
	"context"
	"strings"

	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"reflect"
	"strconv"
)

// toAppConfig converts an App to the core app configuration model
func toAppConfig(a *App) *app.Config {

	appCfg := &app.Config{}
	appCfg.Name = "app"
	appCfg.Version = "1.0.0"
	appCfg.Properties = a.Properties()
	appCfg.Resources = a.resources

	var triggerConfigs []*trigger.Config
	for id, trg := range a.Triggers() {

		triggerConfigs = append(triggerConfigs, toTriggerConfig(strconv.Itoa(id+1), trg))
	}

	appCfg.Triggers = triggerConfigs

	return appCfg
}

// toTriggerConfig converts Trigger to the core Trigger configuration model
func toTriggerConfig(id string, trg *Trigger) *trigger.Config {

	triggerConfig := &trigger.Config{Id:id, Ref: trg.ref, Settings: trg.Settings()}

	//todo add output
	//trigger.Config struct { Output   map[string]interface{} `json:"output"` }

	var handlerConfigs []*trigger.HandlerConfig
	for _, handler := range trg.Handlers() {
		h := &trigger.HandlerConfig{Settings: handler.Settings()}
		//todo add output
		//trigger.HandlerConfig struct { Output   map[string]interface{} `json:"output"` }

		//todo only handles one action for now
		for _, act := range handler.Actions() {
			h.Action = toActionConfig(act)
			break
		}

		handlerConfigs = append(handlerConfigs, h)
	}

	triggerConfig.Handlers = handlerConfigs
	return triggerConfig
}

// toActionConfig converts Action to the core Action configuration model
func toActionConfig(act *Action) *trigger.ActionConfig {
	actionCfg := &trigger.ActionConfig{}

	if act.act != nil {
		actionCfg.Act = act.act
		return actionCfg
	}

	actionCfg.Ref = act.ref

	//todo handle error
	jsonData, _ := json.Marshal(act.Settings())
	actionCfg.Data = jsonData

	mappings := &data.IOMappings{}

	if len(act.inputMappings) > 0 {
		mappings.Input, _ = toMappingDefs(act.inputMappings)
	}
	if len(act.outputMappings) > 0 {
		mappings.Output, _ = toMappingDefs(act.outputMappings)
	}
	actionCfg.Mappings = mappings

	return actionCfg
}

func toMappingDefs(mappings []string) ([]*data.MappingDef, error) {

	var mappingDefs []*data.MappingDef
	for _, strMapping := range mappings {

		idx := strings.Index(strMapping, "=")
		lhs := strings.TrimSpace(strMapping[:idx])
		rhs := strings.TrimSpace(strMapping[idx+1:])

		mType, mValue := getMappingValue(rhs)
		mappingDef := &data.MappingDef{Type: mType, MapTo: lhs, Value: mValue}
		mappingDefs = append(mappingDefs, mappingDef)
	}
	return mappingDefs, nil
}

func getMappingValue(strValue string) (data.MappingType, interface{}) {

	//todo add support for other mapping types
	return data.MtExpression, strValue
}

// ProxyAction

type proxyAction struct {
	handlerFunc HandlerFunc
	metadata    *action.Metadata
}

func NewProxyAction(f HandlerFunc) action.Action {
	return &proxyAction{
		handlerFunc: f,
		metadata:    &action.Metadata{Async: false},
	}
}

// Metadata get the Action's metadata
func (a *proxyAction) Metadata() *action.Metadata {
	return a.metadata
}

// IOMetadata get the Action's IO metadata
func (a *proxyAction) IOMetadata() *data.IOMetadata {
	return nil
}

// Run implementation of action.SyncAction.Run
func (a *proxyAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	return a.handlerFunc(ctx, inputs)
}

// EvalActivity evaluates the specified activity using the provided inputs
func EvalActivity(act activity.Activity, inputs map[string]interface{}) (map[string]*data.Attribute, error) {

	if act.Metadata() == nil {
		//try loading activity with metadata
		value := reflect.ValueOf(act)
		value = value.Elem()
		ref := value.Type().PkgPath()

		act = activity.Get(ref)
	}

	if act.Metadata() == nil {
		//return error
	}

	ac := &activityContext{inputScope: data.NewFixedScope(act.Metadata().Input),
		outputScope: data.NewFixedScope(act.Metadata().Output)}

	for key, value := range inputs {
		ac.inputScope.SetAttrValue(key, value)
	}

	_, evalErr := act.Eval(ac)

	if evalErr != nil {
		return nil, evalErr
	}

	return ac.outputScope.GetAttrs(), nil
}

/////////////////////////////////////////
// activity.Context Implementation

type activityContext struct {
	inputScope  *data.FixedScope
	outputScope *data.FixedScope
}

func (aCtx *activityContext) ActivityHost() activity.Host {
	return aCtx
}

func (aCtx *activityContext) Name() string {
	return ""
}

func (aCtx *activityContext) GetSetting(setting string) (value interface{}, exists bool) {
	return nil, false
}

func (aCtx *activityContext) GetInitValue(key string) (value interface{}, exists bool) {
	return nil, false
}

// GetInput implements activity.Context.GetInput
func (aCtx *activityContext) GetInput(name string) interface{} {

	val, found := aCtx.inputScope.GetAttr(name)
	if found {
		return val.Value()
	}

	return nil
}

// GetOutput implements activity.Context.GetOutput
func (aCtx *activityContext) GetOutput(name string) interface{} {

	val, found := aCtx.outputScope.GetAttr(name)
	if found {
		return val.Value()
	}

	return nil
}

// SetOutput implements activity.Context.SetOutput
func (aCtx *activityContext) SetOutput(name string, value interface{}) {
	aCtx.outputScope.SetAttrValue(name, value)
}

func (aCtx *activityContext) GetSharedData() map[string]interface{} {
	return nil
}

//Deprecated
func (aCtx *activityContext) TaskName() string {
	//ignore
	return ""
}

//Deprecated
func (aCtx *activityContext) FlowDetails() activity.FlowDetails {
	//ignore
	return nil
}

/////////////////////////////////////////
// activity.Host Implementation

func (aCtx *activityContext) ID() string {
	//ignore
	return ""
}

func (aCtx *activityContext) IOMetadata() *data.IOMetadata {
	return nil
}

func (aCtx *activityContext) Reply(replyData map[string]*data.Attribute, err error) {
	// ignore
}

func (aCtx *activityContext) Return(returnData map[string]*data.Attribute, err error) {
	//ignore
}

func (aCtx *activityContext) WorkingData() data.Scope {
	return nil
}

func (aCtx *activityContext) GetResolver() data.Resolver {
	return data.GetBasicResolver()
}
