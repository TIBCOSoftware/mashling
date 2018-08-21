package v2

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

type mashlingActionData struct {
	Dispatch types.Dispatch  `json:"dispatch"`
	Services []types.Service `json:"services"`
}

// Translate translates a v2 mashling gateway JSON config to a Flogo JSON.
func Translate(gateway *types.Schema) ([]byte, error) {
	flogoTriggers := []*ftrigger.Config{}
	flogoActions := []*faction.Config{}

	// Triggers and handlers get mapped to appropriate actionIds.
	for _, trigger := range gateway.Gateway.Triggers {
		var handlers []*ftrigger.HandlerConfig
		flogoActionMap := map[string]*faction.Config{}
		for _, dispatch := range gateway.Gateway.Dispatches {
			actionData := &mashlingActionData{Dispatch: dispatch, Services: gateway.Gateway.Services}
			rawAction, err := json.Marshal(actionData)
			if err != nil {
				return nil, err
			}
			flogoActionMap[dispatch.Name] = &faction.Config{Id: dispatch.Name, Data: rawAction, Ref: "github.com/TIBCOSoftware/mashling/pkg/flogo/action"}
		}
		for _, action := range flogoActionMap {
			flogoActions = append(flogoActions, action)
		}
		for _, handler := range trigger.Handlers {
			actConfig := &ftrigger.ActionConfig{Config: &faction.Config{Id: handler.Dispatch}}
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
	}

	flogoJSON, err := json.MarshalIndent(flogoApp, "", "    ")
	if err != nil {
		return flogoJSON, err
	}

	return flogoJSON, nil
}
