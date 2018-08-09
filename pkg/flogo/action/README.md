# Mashling Action
This action provides all of the features of Mashling for consumption within Flogo as an action. By mapping this action to a handler and a trigger you can invoke a Mashling ruleset much like a Flogo flow. This can be done from within a `flogo.json` file or programmatically via the Flogo API.

## Action Configuration Data

```json
{
  "data":[
    {
      "name": "dispatch",
      "type": "object"
    },
    {
      "name": "services",
      "type": "object"
    }]
}
```

### Data
| Key    | Description   |
|:-----------|:--------------|
| dispatch | A mashling dispatch as JSON |
| services | Mashling services used as a dispatch as JSON |

## Example Flogo JSON Usage of a Mashling Action

```json
{
    "name": "MyProxy",
    "type": "flogo:app",
    "version": "1.0.0",
    "description": "This is a simple proxy.",
    "properties": null,
    "triggers": [
        {
            "name": "MyProxy",
            "id": "MyProxy",
            "ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
            "settings": {
                "port": "9096"
            },
            "output": null,
            "handlers": [
                {
                    "settings": {
                        "autoIdReply": "false",
                        "method": "GET",
                        "path": "/pets/{petId}",
                        "useReplyHandler": "false"
                    },
                    "output": null,
                    "Action": null,
                    "actionId": "Pets",
                    "outputs": null
                }
            ],
            "outputs": null
        }
    ],
    "resources": null,
    "actions": [
        {
            "ref": "github.com/TIBCOSoftware/mashling/pkg/flogo/action",
            "data": {
                "dispatch": {
                    "name": "Pets",
                    "routes": [
                        {
                            "if": "payload.pathParams.petId \u003e= 8 \u0026\u0026 payload.pathParams.petId \u003c= 15",
                            "steps": [
                                {
                                    "service": "PetStorePets",
                                    "input": {
                                        "method": "GET",
                                        "pathParams.id": "${payload.pathParams.petId}"
                                    }
                                },
                                {
                                    "if": "PetStorePets.response.body.status == 'available'",
                                    "service": "PetStoreInventory",
                                    "input": {
                                        "method": "GET"
                                    }
                                }
                            ],
                            "responses": [
                                {
                                    "if": "payload.pathParams.petId == 13",
                                    "error": true,
                                    "output": {
                                        "code": 404,
                                        "data": {
                                            "error": "petId is invalid"
                                        }
                                    }
                                },
                                {
                                    "if": "PetStorePets.response.body.status != 'available'",
                                    "error": true,
                                    "output": {
                                        "code": 403,
                                        "data": {
                                            "error": "Pet is unavailable."
                                        }
                                    }
                                },
                                {
                                    "if": "PetStorePets.response.body.status == 'available'",
                                    "error": false,
                                    "output": {
                                        "code": 200,
                                        "data": {
                                            "inventory": "${PetStoreInventory.response.body}",
                                            "pet": "${PetStorePets.response.body}"
                                        }
                                    }
                                }
                            ]
                        }
                    ]
                },
                "services": [
                    {
                        "name": "PetStorePets",
                        "type": "http",
                        "description": "Make calls to find pets",
                        "settings": {
                            "url": "http://petstore.swagger.io/v2/pet/:id"
                        }
                    },
                    {
                        "name": "PetStoreInventory",
                        "type": "http",
                        "description": "Get pet store inventory.",
                        "settings": {
                            "url": "http://petstore.swagger.io/v2/store/inventory"
                        }
                    }
                ]
            },
            "mappings": null,
            "Act": null,
            "id": "Pets",
            "metadata": null
        }
    ]
}

```

## Example Flogo API Usage of a Mashling Action

```go
package main

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger"
	_ "github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry"
	_ "github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry/triggers"
	mAction "github.com/TIBCOSoftware/mashling/pkg/flogo/action"
)

func main() {
	app := flogo.NewApp()

	trg := app.NewTrigger(&gorillamuxtrigger.RestTrigger{}, map[string]interface{}{"port": "8080"})

	h1 := trg.NewHandler(map[string]interface{}{"method": "GET", "path": "/pets/{petId}", "autoIdReply": "false", "useReplyHandler": "true"})
	a := h1.NewAction(&mAction.MashlingAction{}, map[string]interface{}{"dispatch": DefaultDispatch, "services": DefaultServices, "instance": "MyPetsAction", "identifier": "Pets"})
	a.SetInputMappings("content=$trigger.content", "header=$trigger.header", "params=$trigger.params", "pathParams=$trigger.pathParams", "queryParams=$trigger.queryParams", "tracing=$trigger.tracing")
	e, err := flogo.NewEngine(app)
	if err != nil {
		logger.Error(err)
		return
	}

	engine.RunEngine(e)
}

var DefaultServices = json.RawMessage(`[
  {
    "name": "PetStorePets",
    "description": "Make calls to find pets",
    "type": "http",
    "settings": {
      "url": "http://petstore.swagger.io/v2/pet/:id"
    }
  },
  {
    "name": "PetStoreInventory",
    "description": "Get pet store inventory.",
    "type": "http",
    "settings": {
      "url": "http://petstore.swagger.io/v2/store/inventory"
    }
  }
]`)

var DefaultDispatch = json.RawMessage(`{
  "name": "Pets",
  "routes": [
    {
      "if": "payload.pathParams.petId >= 8 && payload.pathParams.petId <= 15",
      "steps": [
        {
          "service": "PetStorePets",
          "input": {
            "method": "GET",
            "pathParams.id": "${payload.pathParams.petId}"
          }
        },
        {
          "if": "PetStorePets.response.body.status == 'available'",
          "service": "PetStoreInventory",
          "input": {
            "method": "GET"
          }
        }
      ],
      "responses": [
        {
          "if": "payload.pathParams.petId == 13",
          "error": true,
          "output": {
            "code": 404,
            "data": {
              "error": "petId is invalid"
            }
          }
        },
        {
          "if": "PetStorePets.response.body.status != 'available'",
          "error": true,
          "output": {
            "code": 403,
            "data": {
              "error": "Pet is unavailable."
            }
          }
        },
        {
          "if": "PetStorePets.response.body.status == 'available'",
          "error": false,
          "output": {
            "code": 200,
            "data": {
              "pet": "${PetStorePets.response.body}",
              "inventory": "${PetStoreInventory.response.body}"
            }
          }
        }
      ]
    }
  ]
}`)

```
