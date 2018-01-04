package app

import "testing"

const sampleJSON = `{
	"mashling_schema": "0.2",
	"gateway": {
	  "name": "gwTunables",
	  "version": "1.0.0",
	  "display_name":"Tunable HTTP Router",
	  "display_image":"displayImage.svg",
	  "description": "This gateway queries different end-points based on the context supplied as environment flag.",
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
		  "name": "get_rest_trigger",
		  "description": "Rest trigger",
		  "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
		  "settings": {
			"config": "${configurations.restConfig}",
			"method": "GET",
			"path": "/id/{id}",
			"optimize":"true"
		  }
		}
	  ],
	  "event_handlers": [
		{
		  "name": "get_handler",
		  "description": "Handle REST GET call",
		  "reference": "github.com/nareshkumarthota/sampleflows/tunable.json"
		}
	  ],
	  "event_links": [
		{
		  "triggers": ["get_rest_trigger"],
		  "dispatches": [
			{
			  "if": "${env.API_CONTEXT == PETS}",
			  "handler": "get_handler",
			  "inputParams":{
				"id":"${pathParams.id}",
				"endPoint":"http://petstore.swagger.io/v2/pet/:id"
			  }
			},
			{
			  "if": "${env.API_CONTEXT == USERS}",
			  "handler": "get_handler",
			  "inputParams":{
				"id":"${pathParams.id}",
				"endPoint":"http://petstore.swagger.io/v2/user/:id"
			  }
			},
			{
			  "handler": "get_handler",
			  "inputParams":{
				"id":"${pathParams.id}",
				"endPoint":"http://petstore.swagger.io/v2/pet/:id"
			  }
			}
		  ]
		}
	  ]
	}
  }`

func TestGenerateConsulDef(t *testing.T) {

	testArr, err := generateConsulDef(sampleJSON)

	if err != nil {
		t.Errorf("Test Failed Due to error in generating consul definition : %s", err)
	}

	if len(testArr) == 0 {
		t.Error("Consul definition array is empty")
	}

	for _, content := range testArr {
		if len(content.Name) == 0 || len(content.Port) == 0 {
			t.Error("Values parsing error")
		}
	}

}

func TestGenerateFlogoTriggers(t *testing.T) {

	testTriggers, err := generateFlogoTriggers(sampleJSON)

	if err != nil {
		t.Errorf("Test Failed Due to error in generating triggers : %s", err)
	}

	if len(testTriggers) == 0 {
		t.Error("Triggers array is empty")
	}

	for _, content := range testTriggers {
		if len(content.Name) == 0 || len(content.Settings) == 0 {
			t.Error("Values parsing error in triggers")
		}
	}
}
