package flogo

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine"
)

// App is the structure that defines an application
type App struct {
	properties []*data.Attribute
	triggers   []*Trigger
	resources  []*resource.Config
}

// Trigger is the structure that defines a Trigger for the application
type Trigger struct {
	ref      string
	settings map[string]interface{}
	handlers []*Handler
}

// Handler is the structure that defines the handler for a Trigger
type Handler struct {
	settings map[string]interface{}
	actions  []*Action
}

// HandlerFunc is the signature for a function to use as a handler for a Trigger
type HandlerFunc func(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error)

// Action is the structure that defines the Action for a Handler
type Action struct {
	ref            string
	act            action.Action
	settings       map[string]interface{}
	inputMappings  []string
	outputMappings []string
}

// NewApp creates a new Flogo application
func NewApp() *App {
	return &App{}
}

// NewTrigger adds a new trigger to the application
func (a *App) NewTrigger(trg trigger.Trigger, settings map[string]interface{}) *Trigger {

	value := reflect.ValueOf(trg)
	value = value.Elem()
	ref := value.Type().PkgPath()

	newTrg := &Trigger{ref: ref, settings: settings}
	a.triggers = append(a.triggers, newTrg)

	return newTrg
}

// AddProperty adds a shared property to the application
func (a *App) AddProperty(name string, dataType data.Type, value interface{}) error {
	property, err := data.NewAttribute(name, dataType, value)
	if err != nil {
		return err
	}
	a.properties = append(a.properties, property)
	return nil
}

// AddResource adds a Flogo resource to the application
func (a *App) AddResource(id string, data json.RawMessage) {

	res := &resource.Config{ID: id, Data: data}
	a.resources = append(a.resources, res)
}

// Properties gets the shared properties of the application
func (a *App) Properties() []*data.Attribute {

	return a.properties
}

// Triggers gets the Triggers of the application
func (a *App) Triggers() []*Trigger {

	return a.triggers
}

// Settings gets the Trigger's settings
func (t *Trigger) Settings() map[string]interface{} {

	return t.settings
}

// NewHandler adds a new Handler to the Trigger
func (t *Trigger) NewHandler(settings map[string]interface{}) *Handler {

	newHandler := &Handler{settings: settings}
	t.handlers = append(t.handlers, newHandler)

	return newHandler
}

// NewFuncHandler adds a new Function Handler to the Trigger
func (t *Trigger) NewFuncHandler(settings map[string]interface{}, handlerFunc HandlerFunc) *Handler {

	newHandler := &Handler{settings: settings}
	newAct := &Action{act: NewProxyAction(handlerFunc)}
	newHandler.actions = append(newHandler.actions, newAct)

	t.handlers = append(t.handlers, newHandler)

	return newHandler
}

// Handlers gets the Trigger's Handlers
func (t *Trigger) Handlers() []*Handler {
	return t.handlers
}

// Settings gets the Handler's settings
func (h *Handler) Settings() map[string]interface{} {
	return h.settings
}

// NewAction creates a new Action for the Handler
// note: Currently only the first Action is executed for the Handler
func (h *Handler) NewAction(act action.Action, settings map[string]interface{}) *Action {

	value := reflect.ValueOf(act)
	value = value.Elem()
	ref := value.Type().PkgPath()

	newAct := &Action{ref: ref, settings: settings}
	h.actions = append(h.actions, newAct)

	return newAct
}

// Actions gets the Actions of the Handler
func (h *Handler) Actions() []*Action {
	return h.actions
}

// Settings gets the settings of the Action
func (a *Action) Settings() map[string]interface{} {
	return a.settings
}

// SetInputMappings sets the input mappings for the Action, which maps
// the outputs of the Trigger to the inputs of the Action
func (a *Action) SetInputMappings(mappings ...string) {
	a.inputMappings = mappings
}

// SetOutputMappings sets the output mappings for the Action, which maps
// the outputs of the Action to the return of the Trigger
func (a *Action) SetOutputMappings(mappings ...string) {
	a.outputMappings = mappings
}

// InputMappings gets the Action's input mappings
func (a *Action) InputMappings() []string {
	return a.inputMappings
}

// OutputMappings gets the Action's output mappings
func (a *Action) OutputMappings() []string {
	return a.outputMappings
}

// NewEngine creates a new flogo Engine from the specified App
func NewEngine(a *App) (engine.Engine, error) {
	appConfig := toAppConfig(a)
	return engine.New(appConfig)
}
