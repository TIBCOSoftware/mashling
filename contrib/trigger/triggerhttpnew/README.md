# triggerhttpnew
triggerhttpnew

## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/mashling-lib/contrib/trigger/triggerhttpnew
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
| path | URL path to be registered with handler. Example: "/pets/:petId" where petId is path param |
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
			"ref": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/triggerhttpnew",
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
			"ref": "github.com/TIBCOSoftware/mashling-lib/contrib/trigger/triggerhttpnew",
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