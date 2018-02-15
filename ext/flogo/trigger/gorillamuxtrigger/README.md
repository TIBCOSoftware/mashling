# gorillamuxtrigger

Gorillamux based http trigger for Mashling gateway. It supports header/content based conditional routing.

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
      "name": "tracer",
      "type": "string",
      "required": false
    },
    {
      "name": "tracerEndpoint",
      "type": "string",
      "required": false
    },
    {
      "name": "tracerToken",
      "type": "string",
      "required": false
    },
    {
      "name": "tracerDebug",
      "type": "boolean",
      "required": false
    },
    {
      "name": "tracerSameSpan",
      "type": "boolean",
      "required": false
    },
    {
      "name": "tracerID128Bit",
      "type": "boolean",
      "required": false
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
    },{
      "name": "enableClientAuth",
      "type": "string"
    },
    {
      "name": "trustStore",
      "type": "string"
    },
    {
      "name": "basicAuthFile",
      "type": "string"
    },
    {
      "name": "ldapHost",
      "type": "string"
    },
    {
      "name": "ldapBase",
      "type": "string"
    },
    {
      "name": "ldapBindDN",
      "type": "string"
    },
    {
      "name": "ldapBindPassword",
      "type": "string"
    },
    {
      "name": "ldapUserFilter",
      "type": "string"
    },
    {
      "name": "ldapGroupFilter",
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
      "type": "any"
    },
    {
      "name": "tracing",
      "type": "any"
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
| tracer | The tracer to use: noop, zipkin, appdash, or lightstep |
| tracerEndpoint | The url or address of the tracer (ZipKin, APPDash)|
| tracerToken | The token for tracing access (LightStep) |
| tracerDebug | Debug mode for the tracer (ZipKin) |
| tracerSameSpan | Put client side and server side annotations in same span (ZipKin) |
| tracerID128Bit | Use 128 bit ids (ZipKin) |
| enableTLS | true - To enable TLS (Transport Layer Security), false - No TLS security  |
| serverCert | Server certificate file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |
| serverKey | Server private key file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |
| enableClientAuth | true - To enable client AUTH, false - Client AUTH is not enabled |
| trustStore | Trust dir containing clinet CAs |
| basicAuthFile | Path to a password file with username/passwords. An environment variable can be used here. |
| ldapHost | The host name of the ldap server. |
| ldapBase | The LDAP base distinguished name. |
| ldapBindDN | The bind distinguished name. |
| ldapBindPassword | The bind password for the bind distinguished name. |
| ldapUserFilter | The filter to use for authentication. |
| ldapGroupFilter | The filter to use for find user groups. |

### Outputs
| Key    | Description   |
|:-----------|:--------------|
| params | HTTP request params |
| pathParams | HTTP request path params |
| queryParams | HTTP request query params |
| content | HTTP request paylod |
| tracing | Tracing context |

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
        "name": "content_type_multipart_handler",
        "description": "Handle reptils",
        "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "env_prod_handler",
        "description": "Handle prod environment",
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

#### Basic Authentication

To use basic authentication, the necessary descriptors must be in place. For example:

```json
"configurations": [
  {
      "name": "restConfig",
      "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
      "description": "Configuration for rest trigger",
      "settings": {
        "port": "9096",
        "basicAuthFile": "${env.BASIC_AUTH_FILE}"
      }
  }
],
```

This specifies that BASIC_AUTH_FILE is an environment variable whose value will be read into "basicAuthFile" when the gateway starts. This value should be a path to a password file.

Plain username/password file: /home/test/password.txt
```
foo:bar
moo:poo
```

Alternatively, you can also use a salted password file where the format is: username:salt:sha256(salt + password)
```
foo:5VvmQnTXZ10wGZu_Gkjb8umfUPIOQTQ3p1YFadAWTl8=:6267beb3f851b7fee14011f6aa236556f35b186a6791b80b48341e990c367643
```

Start the gateway:
```
BASIC_AUTH_FILE=/home/test/password.txt myApp
```

**NOTE**: It is important to limit access to the password.txt on your environment.

#### Basic Authentication with LDAP

```json
"configurations": [
  {
    "name": "restConfig",
    "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
    "description": "Configuration for rest trigger",
    "settings": {
      "port": "9096",
      "ldapHost": "localhost",
      "ldapBase": "dc=example,dc=com"
    }
  }
],
```
