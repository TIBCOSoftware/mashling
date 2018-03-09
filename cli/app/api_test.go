/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"os"
	"strings"
	"testing"

	"github.com/TIBCOSoftware/mashling/lib/util"

	"github.com/TIBCOSoftware/mashling/lib/model"
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

func TestGetGatewayDetails(t *testing.T) {
	// //Need to create gateway project under temp folder and pass the same for testing.
	// _, err := GetGatewayDetails(SetupNewProjectEnv(), ALL)
	// if err != nil {
	// 	t.Error("Error while getting gateway details")
	// }
}

func TestTranslateGatewayJSON2FlogoJSON(t *testing.T) {

	_, _, _, err := TranslateGatewayJSON2FlogoJSON(gatewayJSON, "9090", "")
	if err != nil {
		t.Error("Error in TranslateGatewayJSON2FlogoJSON function. err: ", err)
	}
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
