package v2

import (
	"bytes"
	"encoding/json"
	"html/template"

	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v1"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

type mapping struct {
	Type  int                    `json:"type"`
	Value map[string]interface{} `json:"value"`
	MapTo string                 `json:"mapTo"`
}

// Translate translates a v2 mashling gateway JSON config to a Flogo app.
func Translate(gateway *types.Schema) ([]byte, error) {
	flogoTriggers := []*ftrigger.Config{}
	flogoActions := []*faction.Config{}

	// Triggers and handlers get mapped to appropriate actionIds.
	for _, trigger := range gateway.Gateway.Triggers {
		var inputMappings []interface{}
		var handlers []*ftrigger.HandlerConfig

		// Trigger metadata.
		triggerMD, err := v1.GetLocalTriggerMetadata(trigger.Type)
		if err != nil {
			return nil, err
		}
		complexValue := make(map[string]interface{})
		for _, outputs := range triggerMD.Output {
			complexValue[outputs.Name()] = "{{$flow." + outputs.Name() + "}}"
		}
		inputMappings = append(inputMappings, mapping{Type: 4, Value: complexValue, MapTo: "mashlingPayload"})
		for _, handler := range trigger.Handlers {
			convertedName := trigger.Name + "_" + handler.Dispatch
			newHandler := &ftrigger.HandlerConfig{ActionId: convertedName, Settings: handler.Settings}
			handlers = append(handlers, newHandler)
			flogoActionMap := map[string]*faction.Config{}
			// Dispatches become actions within a Flogo app.
			for _, dispatch := range gateway.Gateway.Dispatches {
				if dispatch.Name != handler.Dispatch {
					continue
				}
				if _, exists := flogoActionMap[convertedName]; exists {
					continue
				}
				rawMappings, err := json.Marshal(inputMappings)
				if err != nil {
					return nil, err
				}
				rawRoutes, err := json.Marshal(dispatch.Routes)
				if err != nil {
					return nil, err
				}
				rawServices, err := json.Marshal(gateway.Gateway.Services)
				if err != nil {
					return nil, err
				}
				var output bytes.Buffer
				actionTemplate.Execute(&output, struct {
					Name          string
					ID            string
					Identifier    string
					Instance      string
					Routes        template.HTML
					Services      template.HTML
					InputMappings template.HTML
				}{
					Name:          convertedName,
					ID:            convertedName,
					Identifier:    convertedName,
					Instance:      gateway.Gateway.Name,
					Routes:        template.HTML(rawRoutes),
					Services:      template.HTML(rawServices),
					InputMappings: template.HTML(rawMappings),
				})
				rawAction := json.RawMessage(output.String())
				flogoActionMap[convertedName] = &faction.Config{Id: convertedName, Data: rawAction, Ref: "github.com/TIBCOSoftware/flogo-contrib/action/flow"}
			}
			for _, action := range flogoActionMap {
				flogoActions = append(flogoActions, action)
			}
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

var actionTemplate = template.Must(template.New("").Parse(`{
    "flow": {
			"explicitReply": true,
      "type": 1,
      "attributes": [],
      "rootTask": {
        "id": 1,
        "type": 1,
        "tasks": [
          {
            "id": 2,
            "name": "Invoke Mashling Core",
            "description": "Execute Mashling Core",
            "type": 1,
            "activityType": "mashling-core",
            "activityRef": "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/activity/core",
            "attributes": [
              {
                "name": "mashlingPayload",
                "required": false,
								"value": null,
                "type": "object"
              },
              {
                "name": "identifier",
                "value": "{{ .Identifier }}",
                "required": true,
                "type": "string"
              },
              {
                "name": "instance",
                "value": "{{ .Instance }}",
                "required": true,
                "type": "string"
              },
              {
                "name": "routes",
                "value": {{ .Routes }},
                "required": true,
                "type": "array"
              },
              {
                "name": "services",
                "value": {{ .Services }},
                "required": true,
                "type": "array"
              }
            ],
						"inputMappings": {{ .InputMappings }}
          }
        ]
      }
    }
  }`))
