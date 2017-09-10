package app

import (
	"testing"
)

func TestGetGatewayDetails(t *testing.T) {
	//Need to create gateway project under temp folder and pass the same for testing.
	// _, err := GetGatewayDetails(SetupNewProjectEnv(), ALL)
	// if err != nil {
	// 	t.Error("Error while getting gateway details")
	// }
}

func TestTranslateGatewayJSON2FlogoJSON(t *testing.T) {
	gatewayJSON := "{\"gateway\": {\"name\": \"testGw\",\"version\": \"1.0.0\",\"description\": \"This is the first microgateway app\",\"configurations\": [{\"name\": \"kafkaConfig\",\"type\": \"github.com/wnichols/kafkasub\",\"description\": \"Configuration for kafka cluster\",\"settings\": {\"BrokerUrl\": \"localhost:9092\"}}],\"triggers\": [{\"name\": \"rest_trigger\",\"description\": \"The trigger on 'pets' endpoint\",\"type\": \"github.com/rameshpolishetti/triggerhttpnew\",\"settings\": {\"port\": \"9096\",\"method\": \"GET\",\"path\": \"/pets/:petId\"}}],\"event_handlers\": [{\"name\": \"get_pet_success_handler\",\"description\": \"Handle the user access\",\"reference\": \"github.com/TIBCOSoftware/mashling-lib/flow/flogo.json\",\"params\": {\"uri\": \"petstore.swagger.io/v2/pet/3\"}}],\"event_links\": [{\"triggers\": [\"rest_trigger\"],\"dispatches\": [{\"if\": \"trigger.content != undefined\",\"handler\": \"get_pet_success_handler\"}]}]}}"
	//expectedFlogoJSON := "{\"name\": \"testGw\",\"type\": \"flogo:app\",\"version\": \"1.0.0\",\"description\": \"This is the first microgateway app\",\"properties\": null,\"triggers\": [{\"name\": \"rest_trigger\",\"id\": \"rest_trigger\",\"ref\": \"github.com/rameshpolishetti/triggerhttpnew\",\"settings\": {\"port\": \"9096\"},\"handlers\": [{\"actionId\": \"get_pet_success_handler\",\"settings\": {\"autoIdReply\": \"false\",\"method\": \"GET\",\"path\": \"/pets/:petId\",\"useReplyHandler\": \"false\"}}],\"endpoints\": null}],\"actions\": [{\"id\": \"get_pet_success_handler\",\"ref\": \"github.com/TIBCOSoftware/flogo-contrib/action/flow\",\"data\": {\"flow\": {\"explicitReply\": true,\"type\": 1,\"attributes\": [],\"rootTask\": {\"id\": 1,\"type\": 1,\"tasks\": [{\"id\": 2,\"name\": \"Log Message\",\"description\": \"Simple Log Activity\",\"type\": 1,\"activityType\": \"github-com-tibco-software-flogo-contrib-activity-log\",\"activityRef\": \"github.com/TIBCOSoftware/flogo-contrib/activity/log\",\"attributes\": [{\"name\": \"message\",\"value\": null,\"required\": false,\"type\": \"string\"},{\"name\": \"flowInfo\",\"value\": \"false\",\"required\": false,\"type\": \"boolean\"},{\"name\": \"addToFlow\",\"value\": \"false\",\"required\": false,\"type\": \"boolean\"}],\"inputMappings\": [{\"type\": 1,\"value\": \"{T.pathParams}\",\"mapTo\": \"message\"}]},{\"id\": 3,\"name\": \"Invoke REST Service\",\"description\": \"Simple REST Activity\",\"type\": 1,\"activityType\": \"tibco-rest\",\"activityRef\": \"github.com/TIBCOSoftware/flogo-contrib/activity/rest\",\"attributes\": [{\"name\": \"method\",\"value\": \"GET\",\"required\": true,\"type\": \"string\"},{\"name\": \"uri\",\"value\": \"http://petstore.swagger.io/v2/pet/:petId\",\"required\": true,\"type\": \"string\"},{\"name\": \"pathParams\",\"value\": null,\"required\": false,\"type\": \"params\"},{\"name\": \"queryParams\",\"value\": null,\"required\": false,\"type\": \"params\"},{\"name\": \"content\",\"value\": null,\"required\": false,\"type\": \"any\"}],\"inputMappings\": [{\"type\": 1,\"value\": \"{T.pathParams}\",\"mapTo\": \"pathParams\"}]},{\"id\": 4,\"name\": \"Reply To Trigger\",\"description\": \"Simple Reply Activity\",\"type\": 1,\"activityType\": \"tibco-reply\",\"activityRef\": \"github.com/TIBCOSoftware/flogo-contrib/activity/reply\",\"attributes\": [{\"name\": \"code\",\"value\": 0,\"required\": true,\"type\": \"integer\"},{\"name\": \"data\",\"value\": null,\"required\": false,\"type\": \"any\"}],\"inputMappings\": [{\"type\": 1,\"value\": \"{A3.result}\",\"mapTo\": \"data\"}]}],\"links\": [{\"id\": 1,\"from\": 2,\"to\": 3,\"type\": 0},{\"id\": 2,\"from\": 3,\"to\": 4,\"type\": 0}],\"attributes\": []}}}}]}"
	expectedFlogoJSON := `{
	"name": "testGw",
	"type": "flogo:app",
	"version": "1.0.0",
	"description": "This is the first microgateway app",
	"properties": null,
	"triggers": [
		{
			"name": "rest_trigger",
			"id": "rest_trigger",
			"ref": "github.com/rameshpolishetti/triggerhttpnew",
			"settings": {
				"port": "9096"
			},
			"handlers": [
				{
					"actionId": "get_pet_success_handler",
					"settings": {
						"autoIdReply": "false",
						"method": "GET",
						"path": "/pets/:petId",
						"useReplyHandler": "false"
					}
				}
			],
			"endpoints": null
		}
	],
	"actions": [
		{
			"id": "get_pet_success_handler",
			"ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
			"data": {
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
								"name": "Log Message",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": null,
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "false",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "false",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{T.pathParams}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 3,
								"name": "Invoke REST Service",
								"description": "Simple REST Activity",
								"type": 1,
								"activityType": "tibco-rest",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
								"attributes": [
									{
										"name": "method",
										"value": "GET",
										"required": true,
										"type": "string"
									},
									{
										"name": "uri",
										"value": "http://petstore.swagger.io/v2/pet/:petId",
										"required": true,
										"type": "string"
									},
									{
										"name": "pathParams",
										"value": null,
										"required": false,
										"type": "params"
									},
									{
										"name": "queryParams",
										"value": null,
										"required": false,
										"type": "params"
									},
									{
										"name": "content",
										"value": null,
										"required": false,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{T.pathParams}",
										"mapTo": "pathParams"
									}
								]
							},
							{
								"id": 4,
								"name": "Reply To Trigger",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "tibco-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 0,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": false,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A3.result}",
										"mapTo": "data"
									}
								]
							}
						],
						"links": [
							{
								"id": 1,
								"from": 2,
								"to": 3,
								"type": 0
							},
							{
								"id": 2,
								"from": 3,
								"to": 4,
								"type": 0
							}
						],
						"attributes": []
					}
				}
			}
		}
	]
}`
	flogoJSON, err := TranslateGatewayJSON2FlogoJSON(gatewayJSON)

	if err != nil {
		t.Error("Error in TranslateGatewayJSON2FlogoJSON function. err: ", err)
	}

	if expectedFlogoJSON != flogoJSON {
		t.Error("Generated flogoJSON and expected flgoJSON are not same")
	}
}

func TestBuildMashling(t *testing.T) {
	gatewayJSON := ""

	err := BuildMashling("", gatewayJSON)
	if err == nil {
		t.Error("BuildMashling not handled empty gateway JSON")
	}
}
