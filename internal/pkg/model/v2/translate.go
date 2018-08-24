package v2

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

type mashlingActionData struct {
	MashlingURI   string                 `json:"mashlingURI,omitempty"`
	Dispatch      *types.Dispatch        `json:"dispatch,omitempty"`
	Services      []types.Service        `json:"services,omitempty"`
	Pattern       string                 `json:"pattern,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
}

// Translate translates a v2 mashling gateway JSON config to a Flogo JSON.
func Translate(gateway *types.Schema) ([]byte, error) {
	flogoTriggers := []*ftrigger.Config{}
	flogoActions := []*faction.Config{}
	flogoResources := []*resource.Config{}
	flogoActionMap := map[string]*faction.Config{}
	flogoResourceMap := map[string]*resource.Config{}

	// Map dispatches to Actions
	for _, dispatch := range gateway.Gateway.Dispatches {
		// Add action data with mashling resource reference
		actionData := &mashlingActionData{MashlingURI: "mashling:" + dispatch.Name}
		rawAction, err := json.Marshal(actionData)
		if err != nil {
			return nil, err
		}
		flogoActionMap[dispatch.Name] = &faction.Config{Id: "mashling:" + dispatch.Name, Data: rawAction, Ref: "github.com/TIBCOSoftware/mashling/pkg/flogo/action"}
		// Add full action data definition as resource
		var resourceActionData *mashlingActionData
		if dispatch.Pattern == "" {
			resourceActionData = &mashlingActionData{Dispatch: &dispatch, Services: gateway.Gateway.Services, Configuration: dispatch.Configuration}
		} else {
			resourceActionData = &mashlingActionData{Pattern: dispatch.Pattern, Configuration: dispatch.Configuration}
		}
		rawResourceAction, err := json.Marshal(resourceActionData)
		if err != nil {
			return nil, err
		}
		flogoResourceMap[dispatch.Name] = &resource.Config{ID: "mashling:" + dispatch.Name, Data: rawResourceAction}
	}
	for _, resource := range flogoResourceMap {
		flogoResources = append(flogoResources, resource)
	}
	for _, action := range flogoActionMap {
		flogoActions = append(flogoActions, action)
	}
	// Triggers and handlers get mapped to appropriate actions.
	for _, trigger := range gateway.Gateway.Triggers {
		var handlers []*ftrigger.HandlerConfig
		for _, handler := range trigger.Handlers {
			actConfig := &ftrigger.ActionConfig{Config: &faction.Config{Id: "mashling:" + handler.Dispatch}}
			newHandler := &ftrigger.HandlerConfig{Settings: handler.Settings, Action: actConfig}
			handlers = append(handlers, newHandler)
		}

		flogoTriggers = append(flogoTriggers, &ftrigger.Config{Id: trigger.Name, Name: trigger.Name, Settings: trigger.Settings, Ref: trigger.Type, Handlers: handlers})
	}

	flogoApp := app.Config{
		Name:        gateway.Gateway.Name,
		Type:        "flogo:app",
		Version:     gateway.Gateway.Version,
		Description: gateway.Gateway.Description,
		Triggers:    flogoTriggers,
		Actions:     flogoActions,
		Resources:   flogoResources,
	}

	flogoJSON, err := json.MarshalIndent(flogoApp, "", "    ")
	if err != nil {
		return flogoJSON, err
	}

	return flogoJSON, nil
}
