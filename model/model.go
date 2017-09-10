package model

import (
	"encoding/json"

	"github.com/TIBCOSoftware/mashling-lib/types"
)

func CreateMashlingSampleModel() (types.Microgateway, error) {

	microGateway := types.Microgateway{
		Gateway: types.Gateway{
			Name:        "Test",
			Version:     "1.0.0",
			Description: "This is the first microgateway app",
			//Configurations: []types.Config{},
			Configurations: []types.Config{
				{
					Name:        "kafkaConfig",
					Type:        "github.com/wnichols/kafkasub",
					Description: "Configuration for kafka cluster",
					Settings: json.RawMessage(`{
										"BrokerUrl": "localhost:9092"
									}`),
				},
			},
			Triggers: []types.Trigger{
				{
					Name:        "rest_trigger",
					Description: "The trigger on 'pets' endpoint",
					Type:        "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
					Settings: json.RawMessage(`{
					  "port": "9096",
					  "method": "GET",
					  "path": "/pets/:petId"
					}`),
				},
			},
			EventHandlers: []types.EventHandler{
				{
					Name:        "get_pet_success_handler",
					Description: "Handle the user access",
					Params: json.RawMessage(`{
                    				"uri": "petstore.swagger.io/v2/pet/3"
					}`),
					Reference: "github.com/TIBCOSoftware/mashling-lib/flow/flogo.json",
				},
			},
			EventLinks: []types.EventLink{
				{
					Triggers: []string{
						"rest_trigger",
					},
					Dispatches: []types.Dispatch{
						{
							Path: types.Path{
								Handler: "get_pet_success_handler",
							},
						},
					},
				},
			},
		},
	}

	return microGateway, nil
}
