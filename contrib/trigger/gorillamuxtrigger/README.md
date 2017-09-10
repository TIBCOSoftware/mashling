# gorillamuxtrigger
gorillamuxtrigger

## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger
```

## Schema
settings, outputs and handler:

```json
"settings": [
    {
      "name": "port",
      "type": "integer",
      "required": true
    },
    {
      "name": "enableTLS",
      "type": "boolean"
    },
    {
      "name": "serverCert",
      "type": "string"
    },
    {
      "name": "serverKey",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "pathParams",
      "type": "params"
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "method",
        "type": "string",
        "required" : true,
        "allowed" : ["GET", "POST", "PUT", "PATCH", "DELETE"]
      },
      {
        "name": "path",
        "type": "string",
        "required" : true
      },
      {
        "name": "autoIdReply",
        "type": "boolean"
      },
      {
        "name": "useReplyHandler",
        "type": "boolean"
      },
      {
        "name": "Condition",
        "type": "string"
      }
    ]
  }
```

### Settings
| Key    | Description   |
|:-----------|:--------------|
| port | The port to listen on |
| enableTLS | true - To enable TLS (Transport Layer Security), false - No TLS security  |
| serverCert | Server certificate file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |
| serverKey | Server private key file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |

### Outputs
| Key    | Description   |
|:-----------|:--------------|
| params | HTTP request params |
| pathParams | HTTP request path params |
| queryParams | HTTP request query params |
| content | HTTP request paylod |

### Handler settings
| Key    | Description   |
|:-----------|:--------------|
| method | HTTP request method. It can be  |
| path | URL path to be registered with handler |
| Condition | Handler condtion |
| autoIdReply | boolean flag to enable or disable auto reply |
| useReplyHandler | boolean flag to use reply handler |

#### Supported Handler conditions

| Condition Prefix | Description | Example |
|:----------|:-----------|:-------|
| trigger.content | Trigger content / payload based condition | trigger.content.region == APAC |
| trigger.header | HTTP trigger's header based condition | trigger.header.Accept == text/plain |
| env | Environment flag / variable based condition | env.APP_ENVIRONMENT == UAT |


### Sample Mashling Gateway Recipie

Following is the example mashling gateway descriptor uses gorillamuxtrigger as a http trigger.

```json
{
  "gateway": {
    "name": "rest",
    "version": "1.0.0",
    "description": "This is the rest based microgateway app",
    "configurations": [
      {
        "name": "restConfig",
        "type": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
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
        "type": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
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
        "type": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
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
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "birds_handler",
        "description": "Handle birds",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "content_type_multipart_handler",
        "description": "Handle reptils",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "env_prod_handler",
        "description": "Handle prod environment",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "animals_handler",
        "description": "Handle other animals",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "animals_get_handler",
        "description": "Handle other animals",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestGetActivity.json"
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
            "if": "${trigger.header.Content-Type == multipart/form-data}",
            "handler": "content_type_multipart_handler"
          },
          {
            "if": "${env.APP_ENVIRONMENT == PRODUCTION}",
            "handler": "env_prod_handler"
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
}
```
#### Sample request payload

Follwing is the sample payload. Try changing the value of name ("CAT" to some other value) to notice handler routing. 

```json
{
    "category": {
        "id": 1,
        "name": "Mammals"
    },
    "id": 38,
    "name": "CAT",
    "photoUrls": [
        "string"
    ],
    "status": "sold",
    "tags": [
        {
            "id": 0,
            "name": "string"
        }
    ]
}
```
