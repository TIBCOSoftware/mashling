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
    },
    {
      "name": "mashlingURI",
      "type": "string"
    },
    {
      "name": "pattern",
      "type": "string"
    },
    {
      "name": "configuration",
      "type": "object"
    }
  ]
}
```

### Data
| Key    | Description   |
|:-----------|:--------------|
| dispatch | A mashling dispatch as JSON |
| services | Mashling services used as a dispatch as JSON |
| mashlingURI | Mashling URI to resource data in Flogo |
| pattern | Mashling out-of-the-box gateway pattern to use |
| configuration | Mashling configuration specific to this action |

## Example Flogo JSON Usage of a Mashling Action

```json
{
    "name": "MyProxy",
    "type": "flogo:app",
    "version": "1.0.0",
    "description": "This is a simple proxy.",
    "triggers": [
        {
            "name": "MyProxy",
            "id": "MyProxy",
            "ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
            "settings": {
                "port": "9096"
            },
            "handlers": [
                {
                    "settings": {
                        "autoIdReply": "false",
                        "method": "GET",
                        "path": "/magic",
                        "useReplyHandler": "false"
                    },
                    "Action": {
                        "id": "mashling:HTTP"
                    }
                }
            ]
        }
    ],
    "resources": [
        {
            "id": "mashling:HTTP",
            "compressed": false,
            "data": {
                "pattern": "DefaultHttpPattern",
                "configuration": {
                    "backendUrl": "https://petstore.swagger.io/v2/pet/32",
                    "jwtKey": "qwertyuiopasdfghjklzxcvbnm123456",
                    "useCircuitBreaker": true,
                    "useJWT": true
                }
            }
        }
    ],
    "actions": [
        {
            "ref": "github.com/TIBCOSoftware/mashling/pkg/flogo/action",
            "data": {
                "mashlingURI": "mashling:HTTP"
            },
            "id": "mashling:HTTP"
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
	app.AddResource("mashling:HTTP", Data)

	trg := app.NewTrigger(&gorillamuxtrigger.RestTrigger{}, map[string]interface{}{"port": "8080"})
	h1 := trg.NewHandler(map[string]interface{}{"method": "GET", "path": "/magic", "autoIdReply": "false", "useReplyHandler": "true"})
	h1.NewAction(&mAction.MashlingAction{}, map[string]interface{}{"mashlingURI": "mashling:HTTP"})
	e, err := flogo.NewEngine(app)
	if err != nil {
		logger.Error(err)
		return
	}

	engine.RunEngine(e)
}

var Data = json.RawMessage(`{
    "pattern": "DefaultHttpPattern",
    "configuration": {
        "backendUrl": "https://petstore.swagger.io/v2/pet/32",
        "jwtKey": "qwertyuiopasdfghjklzxcvbnm123456",
        "useCircuitBreaker": true,
        "useJWT": true
    }
}`)
```
