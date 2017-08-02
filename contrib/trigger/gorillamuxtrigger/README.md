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

## Example configuration

Triggers are configured via triggers.json of your application. The follwing are some example configuration of the REST trigger.

### POST

```json
{
    "triggers": [
		{
		    "name": "rest_trigger",
			"id": "rest_trigger",
			"ref": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
			"settings": {
				"port": "9096"
			},
			"handlers": [
				{
					"actionId": "get_pet_success_handler_usa",
					"settings": {
						"Condition": "${trigger.content.country == USA}",
						"autoIdReply": "false",
						"method": "POST",
						"path": "/test",
						"useReplyHandler": "false"
					}
				}
            ]
        }
    ]
}
```

### Multiple handlers

```json
{
    "triggers": [
		{
		    "name": "rest_trigger",
			"id": "rest_trigger",
			"ref": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
			"settings": {
				"port": "9096"
			},
			"handlers": [
				{
					"actionId": "get_pet_success_handler_india",
					"settings": {
						"Condition": "${trigger.content.country == INDIA}",
						"autoIdReply": "false",
						"method": "POST",
						"path": "/test",
						"useReplyHandler": "false"
					}
				},
                {
					"actionId": "get_pet_success_handler_usa",
					"settings": {
						"Condition": "${trigger.content.country == USA}",
						"autoIdReply": "false",
						"method": "POST",
						"path": "/test",
						"useReplyHandler": "false"
					}
				}
            ]
        }
    ]
}
```
### Mashling Gateway Recipie

Following is the example mashling gateway descriptor uses gorillamuxtrigger as a http trigger.

```json
{
  "gateway": {
    "name": "muxGw",
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
        "name": "users_trigger",
        "description": "Users rest trigger",
        "type": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
        "settings": {
          "config": "${configurations.restConfig}",
          "method": "PUT",
		      "path": "/users/{petId}",
          "optimize":"true"
        }
      },
      {
        "name": "users_get_trigger",
        "description": "Users rest trigger",
        "type": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/gorillamuxtrigger",
        "settings": {
          "config": "${configurations.restConfig}",
          "method": "GET",
		      "path": "/users/{petId}",
          "optimize":"true"
        }
      }
    ],
    "event_handlers": [
      {
        "name": "usa_users_http_handler",
        "description": "Handle the user access for USA users",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "asia_users_http_handler",
        "description": "Handle the user access for asia users",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "global_users_http_handler",
        "description": "Handle the user access for asia users",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "get_users_http_handler",
        "description": "Handle the user access for asia users",
        "reference": "github.com/TIBCOSoftware/mashling-lib/flow/RestTriggerToRestGetActivity.json"
      }
    ],
    "event_links": [
      {
        "triggers": ["users_trigger"],
        "dispatches": [
          {
            "if": "${trigger.content.category.name == USA}",
            "handler": "usa_users_http_handler"
          },
          {
            "if": "${trigger.content.category.name in (IND,CHN,JPN)}",
            "handler": "asia_users_http_handler"
          },
          {
            "handler": "global_users_http_handler"
          }
        ]
      },
      {
        "triggers": ["users_get_trigger"],
        "dispatches": [
          {
            "handler": "get_users_http_handler"
          }
        ]
      }
    ]
  }
}
```
#### Sample request payload

Follwing is the sample payload. Try changing the value of category.name ("USA" to some other value) to notice handler routing. 

```json
{
    "category": {
        "id": 10,
        "name": "USA"
    },
    "id": 10,
    "name": "OSI",
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