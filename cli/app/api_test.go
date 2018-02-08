/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/TIBCOSoftware/mashling/lib/util"

	"github.com/TIBCOSoftware/mashling/lib/model"
	"github.com/stretchr/testify/assert"
)

const gatewayJSON string = `{
	"mashling_schema": "0.2",
	"gateway": {
		"name": "mashlingApp",
		"version": "1.0.0",
		"display_name": "Gateway Application",
		"display_image": "GatewayIcon.jpg",
		"description": "This is the first microgateway ping app",
		"configurations": [],
		"triggers": [
			{
				"name": "rest_trigger",
				"description": "The trigger on 'pets' endpoint",
				"type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
				"settings": {
					"port": "9096",
					"method": "GET",
					"path": "/pets/{petId}"
				}
			}
		],
		"event_handlers": [
			{
				"name": "get_pet_success_handler",
				"description": "Handle the user access",
				"reference": "github.com/TIBCOSoftware/mashling/lib/flow/flogo.json",
				"params": {
					"uri": "petstore.swagger.io/v2/pet/3"
				}
			}
		],
		"event_links": [
			{
				"triggers": [
					"rest_trigger"
				],
				"dispatches": [
					{
						"handler": "get_pet_success_handler"
					}
				]
			}
		]
	}
}`
const expectedFlogoJSON string = `{
	"name": "mashlingApp",
	"type": "flogo:app",
	"version": "1.0.0",
	"description": "This is the first microgateway ping app",
	"properties": null,
	"triggers": [
		{
			"name": "rest_trigger",
			"id": "rest_trigger",
			"ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
			"settings": {
				"port": "9096"
			},
			"output": null,
			"handlers": [
				{
					"actionId": "get_pet_success_handler",
					"settings": {
						"autoIdReply": "false",
						"method": "GET",
						"path": "/pets/{petId}",
						"useReplyHandler": "false"
					},
					"output": null,
					"actionMappings": {},
					"outputs": null
				}
			],
			"outputs": null
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
										"value": "{T.pathParams}.petId",
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
										"value": 200,
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
			},
			"metadata": null
		}
	]
}`

func TestGetGatewayDetails(t *testing.T) {
	// //Need to create gateway project under temp folder and pass the same for testing.
	// _, err := GetGatewayDetails(SetupNewProjectEnv(), ALL)
	// if err != nil {
	// 	t.Error("Error while getting gateway details")
	// }
}

func TestTranslateGatewayJSON2FlogoJSON(t *testing.T) {

	flogoJSON, err := TranslateGatewayJSON2FlogoJSON(gatewayJSON, "9090")
	if err != nil {
		t.Error("Error in TranslateGatewayJSON2FlogoJSON function. err: ", err)
	}
	isEqualJSON, err := AreEqualJSON(expectedFlogoJSON, flogoJSON)
	assert.NoError(t, err, "Error: Error comparing expected and actual Flogo JSON %v", err)
	assert.True(t, isEqualJSON, "Error: Expected and actual Flogo JSON contents are not equal.")
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

func TestAppendPingDescriptor(t *testing.T) {

	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		t.Fail()
	}

	os.Setenv(util.Mashling_Ping_Embed_Config_Property, "TRUE")
	descriptor, err = appendPingDescriptor("9090", descriptor)
	if err != nil {
		t.Fail()
	}

	pingDescrptr, err := CreateMashlingPingModel("9090")
	if err != nil {
		t.Fail()
	}

	var apendPingFunctionality bool
	apendPingFunctionality = false
	for _, trigger := range pingDescrptr.Gateway.Triggers {
		for _, descTrigger := range descriptor.Gateway.Triggers {
			if strings.Compare(trigger.Name, descTrigger.Name) == 0 {
				apendPingFunctionality = true
				break
			}
		}
	}
	if !apendPingFunctionality {
		t.Fail()
	}
	os.Setenv(util.Mashling_Ping_Embed_Config_Property, "FALSE")
}
