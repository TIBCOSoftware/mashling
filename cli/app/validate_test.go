/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"testing"
)

const validGatewayJSON string = `{
	"mashling_schema": "0.2",
	"gateway": {
	  "name": "demoRestGw",
	  "version": "1.0.0",
	  "display_name":"Rest Conditional Gateway",
	  "display_image":"displayImage.svg",
	  "description": "This is the rest based microgateway app",
	  "configurations": [
		{
		  "name": "restConfig",
		  "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
		  "description": "Configuration for rest trigger",
		  "settings": {
			"port": "9096"
		  }
		}
	  ],
	  "triggers": [
		{
		  "name": "animals_rest_trigger",
		  "description": "Animals rest trigger - PUT animal details",
		  "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
		  "settings": {
			"config": "${configurations.restConfig}",
			"method": "PUT",
				"path": "/pets",
			"optimize":"true"
		  }
		},
		{
		  "name": "get_animals_rest_trigger",
		  "description": "Animals rest trigger - get animal details",
		  "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
		  "settings": {
			"config": "${configurations.restConfig}",
			"method": "GET",
				"path": "/pets/{petId}",
			"optimize":"true"
		  }
		}
	  ],
	  "event_handlers": [
		{
		  "name": "mammals_handler",
		  "description": "Handle mammals",
		  "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
		},
		{
		  "name": "birds_handler",
		  "description": "Handle birds",
		  "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
		},
		{
		  "name": "animals_handler",
		  "description": "Handle other animals",
		  "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
		},
		{
		  "name": "animals_get_handler",
		  "description": "Handle other animals",
		  "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestGetActivity.json"
		}
	  ],
	  "event_links": [
		{
		  "triggers": ["animals_rest_trigger"],
		  "dispatches": [
			{
			  "if": "${trigger.content.name in (ELEPHANT,CAT)}",
			  "handler": "mammals_handler"
			},
			{
			  "if": "${trigger.content.name == SPARROW}",
			  "handler": "birds_handler"
			},
			{
			  "handler": "animals_handler"
			}
		  ]
		},
		{
		  "triggers": ["get_animals_rest_trigger"],
		  "dispatches": [
			{
			  "handler": "animals_get_handler"
			}
		  ]
		}
	  ]
	}
}`

func TestValidate(t *testing.T) {
	isValid, _ := IsValidGateway(validGatewayJSON)

	if !isValid {
		t.Error("A valid gateway json failed to pass validation.")
	}
}
