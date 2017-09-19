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
	gatewayJSON := "{ \"gateway\": { \"name\": \"demoRestGw\", \"version\": \"1.0.0\", \"display_name\":\"Rest Conditional Gateway\", \"description\": \"This is the rest based microgateway app\", \"configurations\": [ { \"name\": \"restConfig\", \"type\": \"github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger\", \"description\": \"Configuration for rest trigger\", \"settings\": { \"port\": \"9096\" } } ], \"triggers\": [ { \"name\": \"animals_rest_trigger\", \"description\": \"Animals rest trigger - PUT animal details\", \"type\": \"github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger\", \"settings\": { \"config\": \"${configurations.restConfig}\", \"method\": \"PUT\", \"path\": \"/pets\", \"optimize\":\"true\" } }, { \"name\": \"get_animals_rest_trigger\", \"description\": \"Animals rest trigger - get animal details\", \"type\": \"github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger\", \"settings\": { \"config\": \"${configurations.restConfig}\", \"method\": \"GET\", \"path\": \"/pets/{petId}\", \"optimize\":\"true\" } } ], \"event_handlers\": [ { \"name\": \"mammals_handler\", \"description\": \"Handle mammals\", \"reference\": \"github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json\" }, { \"name\": \"birds_handler\", \"description\": \"Handle birds\", \"reference\": \"github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json\" }, { \"name\": \"animals_handler\", \"description\": \"Handle other animals\", \"reference\": \"github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json\" }, { \"name\": \"animals_get_handler\", \"description\": \"Handle other animals\", \"reference\": \"github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestGetActivity.json\" } ], \"event_links\": [ { \"triggers\": [\"animals_rest_trigger\"], \"dispatches\": [ { \"if\": \"${trigger.content.name in (ELEPHANT,CAT)}\", \"handler\": \"mammals_handler\" }, { \"if\": \"${trigger.content.name == SPARROW}\", \"handler\": \"birds_handler\" }, { \"handler\": \"animals_handler\" } ] }, { \"triggers\": [\"get_animals_rest_trigger\"], \"dispatches\": [ { \"handler\": \"animals_get_handler\" } ] } ] } } "
	//expectedFlogoJSON := "{\"name\": \"testGw\",\"type\": \"flogo:app\",\"version\": \"1.0.0\",\"description\": \"This is the first microgateway app\",\"properties\": null,\"triggers\": [{\"name\": \"rest_trigger\",\"id\": \"rest_trigger\",\"ref\": \"github.com/rameshpolishetti/triggerhttpnew\",\"settings\": {\"port\": \"9096\"},\"handlers\": [{\"actionId\": \"get_pet_success_handler\",\"settings\": {\"autoIdReply\": \"false\",\"method\": \"GET\",\"path\": \"/pets/:petId\",\"useReplyHandler\": \"false\"}}],\"endpoints\": null}],\"actions\": [{\"id\": \"get_pet_success_handler\",\"ref\": \"github.com/TIBCOSoftware/flogo-contrib/action/flow\",\"data\": {\"flow\": {\"explicitReply\": true,\"type\": 1,\"attributes\": [],\"rootTask\": {\"id\": 1,\"type\": 1,\"tasks\": [{\"id\": 2,\"name\": \"Log Message\",\"description\": \"Simple Log Activity\",\"type\": 1,\"activityType\": \"github-com-tibco-software-flogo-contrib-activity-log\",\"activityRef\": \"github.com/TIBCOSoftware/flogo-contrib/activity/log\",\"attributes\": [{\"name\": \"message\",\"value\": null,\"required\": false,\"type\": \"string\"},{\"name\": \"flowInfo\",\"value\": \"false\",\"required\": false,\"type\": \"boolean\"},{\"name\": \"addToFlow\",\"value\": \"false\",\"required\": false,\"type\": \"boolean\"}],\"inputMappings\": [{\"type\": 1,\"value\": \"{T.pathParams}\",\"mapTo\": \"message\"}]},{\"id\": 3,\"name\": \"Invoke REST Service\",\"description\": \"Simple REST Activity\",\"type\": 1,\"activityType\": \"tibco-rest\",\"activityRef\": \"github.com/TIBCOSoftware/flogo-contrib/activity/rest\",\"attributes\": [{\"name\": \"method\",\"value\": \"GET\",\"required\": true,\"type\": \"string\"},{\"name\": \"uri\",\"value\": \"http://petstore.swagger.io/v2/pet/:petId\",\"required\": true,\"type\": \"string\"},{\"name\": \"pathParams\",\"value\": null,\"required\": false,\"type\": \"params\"},{\"name\": \"queryParams\",\"value\": null,\"required\": false,\"type\": \"params\"},{\"name\": \"content\",\"value\": null,\"required\": false,\"type\": \"any\"}],\"inputMappings\": [{\"type\": 1,\"value\": \"{T.pathParams}\",\"mapTo\": \"pathParams\"}]},{\"id\": 4,\"name\": \"Reply To Trigger\",\"description\": \"Simple Reply Activity\",\"type\": 1,\"activityType\": \"tibco-reply\",\"activityRef\": \"github.com/TIBCOSoftware/flogo-contrib/activity/reply\",\"attributes\": [{\"name\": \"code\",\"value\": 0,\"required\": true,\"type\": \"integer\"},{\"name\": \"data\",\"value\": null,\"required\": false,\"type\": \"any\"}],\"inputMappings\": [{\"type\": 1,\"value\": \"{A3.result}\",\"mapTo\": \"data\"}]}],\"links\": [{\"id\": 1,\"from\": 2,\"to\": 3,\"type\": 0},{\"id\": 2,\"from\": 3,\"to\": 4,\"type\": 0}],\"attributes\": []}}}}]}"
	expectedFlogoJSON := `{
		"name": "taceingGw2",
	"type": "flogo:app",
	"version": "1.0.0",
	"description": "This is the rest based microgateway app",
	"properties": null,
	"triggers": [
		{
			"name": "animals_rest_trigger",
			"id": "animals_rest_trigger",
			"ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
			"settings": {
				"port": "9096"
			},
			"outputs": null,
			"handlers": [
				{
					"actionId": "mammals_handler",
					"settings": {
						"Condition": "${trigger.content.name in (ELEPHANT,CAT)}",
						"autoIdReply": "false",
						"method": "PUT",
						"path": "/pets",
						"useReplyHandler": "false"
					},
					"outputs": null
				},
				{
					"actionId": "birds_handler",
					"settings": {
						"Condition": "${trigger.content.name == SPARROW}",
						"autoIdReply": "false",
						"method": "PUT",
						"path": "/pets",
						"useReplyHandler": "false"
					},
					"outputs": null
				},
				{
					"actionId": "animals_handler",
					"settings": {
						"autoIdReply": "false",
						"method": "PUT",
						"path": "/pets",
						"useReplyHandler": "false"
					},
					"outputs": null
				},
				{
					"actionId": "animals_get_handler",
					"settings": {
						"autoIdReply": "false",
						"method": "GET",
						"path": "/pets/{petId}",
						"useReplyHandler": "false"
					},
					"outputs": null
				}
			],
			"endpoints": null
		}
	],
	"actions": [
		{
			"id": "mammals_handler",
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
								"name": "Invoke REST Service",
								"description": "Simple REST Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-rest",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
								"attributes": [
									{
										"name": "method",
										"value": "PUT",
										"required": true,
										"type": "string"
									},
									{
										"name": "uri",
										"value": "http://petstore.swagger.io/v2/pet",
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
										"value": "{T.content}",
										"mapTo": "content"
									}
								]
							},
							{
								"id": 3,
								"name": "Log Message",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Success",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 4,
								"name": "Log Message (2)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}.id",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 12,
								"name": "Log Message (7)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "tibco-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
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
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 13,
								"name": "Reply To Trigger (3)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "tibco-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 200,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "data"
									}
								]
							},
							{
								"id": 6,
								"name": "Log Message (3)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Failed",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 7,
								"name": "Log Message (4)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 8,
								"name": "Reply To Trigger (2)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 400,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
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
								"type": 1,
								"value": "${A2.result}.id==${T.content}.id"
							},
							{
								"id": 2,
								"from": 3,
								"to": 4,
								"type": 0
							},
							{
								"id": 3,
								"from": 4,
								"to": 12,
								"type": 0
							},
							{
								"id": 4,
								"from": 12,
								"to": 13,
								"type": 0
							},
							{
								"id": 5,
								"from": 2,
								"to": 6,
								"type": 1,
								"value": "${A2.result}.code==1"
							},
							{
								"id": 6,
								"from": 6,
								"to": 7,
								"type": 0
							},
							{
								"id": 7,
								"from": 7,
								"to": 8,
								"type": 0
							}
						],
						"attributes": []
					},
					"errorHandlerTask": {
						"id": 9,
						"type": 1,
						"tasks": [
							{
								"id": 10,
								"name": "Log Message (5)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Error processing request in gateway",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 11,
								"name": "Log Message (6)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{T.content}",
										"mapTo": "message"
									}
								]
							}
						],
						"links": [
							{
								"id": 8,
								"from": 10,
								"to": 11,
								"type": 0
							}
						],
						"attributes": []
					}
				}
			}
		},
		{
			"id": "birds_handler",
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
								"name": "Invoke REST Service",
								"description": "Simple REST Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-rest",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
								"attributes": [
									{
										"name": "method",
										"value": "PUT",
										"required": true,
										"type": "string"
									},
									{
										"name": "uri",
										"value": "http://petstore.swagger.io/v2/pet",
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
										"value": "{T.content}",
										"mapTo": "content"
									}
								]
							},
							{
								"id": 3,
								"name": "Log Message",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Success",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 4,
								"name": "Log Message (2)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}.id",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 12,
								"name": "Log Message (7)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "tibco-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
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
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 13,
								"name": "Reply To Trigger (3)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "tibco-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 200,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "data"
									}
								]
							},
							{
								"id": 6,
								"name": "Log Message (3)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Failed",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 7,
								"name": "Log Message (4)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 8,
								"name": "Reply To Trigger (2)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 400,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
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
								"type": 1,
								"value": "${A2.result}.id==${T.content}.id"
							},
							{
								"id": 2,
								"from": 3,
								"to": 4,
								"type": 0
							},
							{
								"id": 3,
								"from": 4,
								"to": 12,
								"type": 0
							},
							{
								"id": 4,
								"from": 12,
								"to": 13,
								"type": 0
							},
							{
								"id": 5,
								"from": 2,
								"to": 6,
								"type": 1,
								"value": "${A2.result}.code==1"
							},
							{
								"id": 6,
								"from": 6,
								"to": 7,
								"type": 0
							},
							{
								"id": 7,
								"from": 7,
								"to": 8,
								"type": 0
							}
						],
						"attributes": []
					},
					"errorHandlerTask": {
						"id": 9,
						"type": 1,
						"tasks": [
							{
								"id": 10,
								"name": "Log Message (5)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Error processing request in gateway",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 11,
								"name": "Log Message (6)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{T.content}",
										"mapTo": "message"
									}
								]
							}
						],
						"links": [
							{
								"id": 8,
								"from": 10,
								"to": 11,
								"type": 0
							}
						],
						"attributes": []
					}
				}
			}
		},
		{
			"id": "animals_handler",
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
								"name": "Invoke REST Service",
								"description": "Simple REST Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-rest",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
								"attributes": [
									{
										"name": "method",
										"value": "PUT",
										"required": true,
										"type": "string"
									},
									{
										"name": "uri",
										"value": "http://petstore.swagger.io/v2/pet",
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
										"value": "{T.content}",
										"mapTo": "content"
									}
								]
							},
							{
								"id": 3,
								"name": "Log Message",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Success",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 4,
								"name": "Log Message (2)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}.id",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 12,
								"name": "Log Message (7)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "tibco-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
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
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 13,
								"name": "Reply To Trigger (3)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "tibco-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 200,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "data"
									}
								]
							},
							{
								"id": 6,
								"name": "Log Message (3)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Failed",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 7,
								"name": "Log Message (4)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 8,
								"name": "Reply To Trigger (2)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 400,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
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
								"type": 1,
								"value": "${A2.result}.id==${T.content}.id"
							},
							{
								"id": 2,
								"from": 3,
								"to": 4,
								"type": 0
							},
							{
								"id": 3,
								"from": 4,
								"to": 12,
								"type": 0
							},
							{
								"id": 4,
								"from": 12,
								"to": 13,
								"type": 0
							},
							{
								"id": 5,
								"from": 2,
								"to": 6,
								"type": 1,
								"value": "${A2.result}.code==1"
							},
							{
								"id": 6,
								"from": 6,
								"to": 7,
								"type": 0
							},
							{
								"id": 7,
								"from": 7,
								"to": 8,
								"type": 0
							}
						],
						"attributes": []
					},
					"errorHandlerTask": {
						"id": 9,
						"type": 1,
						"tasks": [
							{
								"id": 10,
								"name": "Log Message (5)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Error processing request in gateway",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 11,
								"name": "Log Message (6)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{T.content}",
										"mapTo": "message"
									}
								]
							}
						],
						"links": [
							{
								"id": 8,
								"from": 10,
								"to": 11,
								"type": 0
							}
						],
						"attributes": []
					}
				}
			}
		},
		{
			"id": "animals_get_handler",
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
								"name": "Invoke REST Service",
								"description": "Simple REST Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-rest",
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
										"value": "http://petstore.swagger.io/v2/pet/:id",
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
										"value": "{T.pathParams}.petId",
										"mapTo": "pathParams.id"
									}
								]
							},
							{
								"id": 3,
								"name": "Log Message",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Success",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 4,
								"name": "Log Message (2)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}.id",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 12,
								"name": "Log Message (7)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "tibco-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
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
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 13,
								"name": "Reply To Trigger (3)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "tibco-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 200,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "data"
									}
								]
							},
							{
								"id": 6,
								"name": "Log Message (3)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Failed",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 7,
								"name": "Log Message (4)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
										"mapTo": "message"
									}
								]
							},
							{
								"id": 8,
								"name": "Reply To Trigger (2)",
								"description": "Simple Reply Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-reply",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/reply",
								"attributes": [
									{
										"name": "code",
										"value": 400,
										"required": true,
										"type": "integer"
									},
									{
										"name": "data",
										"value": null,
										"required": true,
										"type": "any"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{A2.result}",
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
								"type": 1,
								"value": "${A2.result}.id==${T.pathParams}.petId"
							},
							{
								"id": 2,
								"from": 3,
								"to": 4,
								"type": 0
							},
							{
								"id": 3,
								"from": 4,
								"to": 12,
								"type": 0
							},
							{
								"id": 4,
								"from": 12,
								"to": 13,
								"type": 0
							},
							{
								"id": 5,
								"from": 2,
								"to": 6,
								"type": 1,
								"value": "${A2.result}.id!=${T.pathParams}.petId"
							},
							{
								"id": 6,
								"from": 6,
								"to": 7,
								"type": 0
							},
							{
								"id": 7,
								"from": 7,
								"to": 8,
								"type": 0
							}
						],
						"attributes": []
					},
					"errorHandlerTask": {
						"id": 9,
						"type": 1,
						"tasks": [
							{
								"id": 10,
								"name": "Log Message (5)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "Error processing request in gateway",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								]
							},
							{
								"id": 11,
								"name": "Log Message (6)",
								"description": "Simple Log Activity",
								"type": 1,
								"activityType": "github-com-tibco-software-flogo-contrib-activity-log",
								"activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
								"attributes": [
									{
										"name": "message",
										"value": "",
										"required": false,
										"type": "string"
									},
									{
										"name": "flowInfo",
										"value": "true",
										"required": false,
										"type": "boolean"
									},
									{
										"name": "addToFlow",
										"value": "true",
										"required": false,
										"type": "boolean"
									}
								],
								"inputMappings": [
									{
										"type": 1,
										"value": "{T.pathParams}.petId",
										"mapTo": "message"
									}
								]
							}
						],
						"links": [
							{
								"id": 8,
								"from": 10,
								"to": 11,
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
