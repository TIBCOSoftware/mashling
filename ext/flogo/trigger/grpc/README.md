# gRPC Trigger
gRPC trigger for mashling gateway supports method name based routing.

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
      "name": "protoname",
      "type": "string",
      "required": true
    },
    {
      "name": "servicename",
      "type": "string",
      "required": true
    }
  ],
  "outputs": [
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "grpcData",
      "type": "object"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "autoIdReply",
        "type": "boolean"
      },
      {
        "name": "useReplyHandler",
        "type": "boolean"
      },
      {
        "name": "methodName",
        "type": "string"
      }
    ]
  }
```
### Settings
| Key    | Description   |
|:-----------|:--------------|
| port | The port to listen on |
| protoname | The name of the proto file|
| servicename | The name of the service mentioned in proto file|
| enableTLS | true - To enable TLS (Transport Layer Security), false - No TLS security  |
| serverCert | Server certificate file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |
| serverKey | Server private key file in PEM format. Need to provide file name along with path. Path can be relative to gateway binary location. |

### Outputs
| Key    | Description   |
|:-----------|:--------------|
| params | Request params |
| grpcData | gRPC Method parameters |

### Handler settings
| Key    | Description   |
|:-----------|:--------------|
| methodName | Name of the method |
| autoIdReply | boolean flag to enable or disable auto reply |
| useReplyHandler | boolean flag to use reply handler |

### Sample Mashling Gateway Recipie

Following is the example mashling gateway descriptor uses a grpc trigger.

```json
{
    "mashling_schema": "1.0",
    "gateway": {
        "name": "grpc sample trigger",
        "version": "1.0.0",
        "description": "This is a simple grpc server.",
        "triggers": [
            {
                "name": "grpc sample trigger",
                "description": "This is a simple grpc server.",
                "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/grpc",
                "settings": {
                    "port": 9096,
                    "protoname":"messages",
                    "servicename":"PetService",
                    "enableTLS": "true",
                    "serverCert": "${env.SERVER_CERT}",
                    "serverKey": "${env.SERVER_KEY}"
                },
                "handlers": [
                    {
                        "dispatch": "petByIdDispatch",
                        "settings": {
                            "autoIdReply": "false",
                            "useReplyHandler": "false",
                            "methodName": "PetById"
                        }
                    },
                    {
                        "dispatch": "userByNameDispatch",
                        "settings": {
                            "autoIdReply": "false",
                            "useReplyHandler": "false",
                            "methodName": "UserByName"
                        }
                    }
                ]
            }
        ],
        "dispatches": [
            {
                "name": "petByIdDispatch",
                "routes": [
                    {
                        "steps": [
                            {
                                "service": "PetStorePets",
                                "input": {
                                    "method": "GET",
                                    "pathParams.id": "${payload.pathParams.Id}"
                                }
                            }
                        ],
                        "responses": [
                            {
                                "error": false,
                                "output": {
                                    "code": 200,
                                    "data": {
                                        "pet": "${PetStorePets.response.body}"
                                    }
                                }
                            }
                        ]
                    }
                ]
            },
            {
                "name": "userByNameDispatch",
                "routes": [
                    {
                        "steps": [
                            {
                                "service": "PetStoreUsersByName",
                                "input": {
                                    "method": "GET",
                                    "pathParams.username": "${payload.pathParams.Username}"
                                }
                            }
                        ],
                        "responses": [
                            {
                                "error": false,
                                "output": {
                                    "code": 200,
                                    "data": {
                                        "user": "${PetStoreUsersByName.response.body}"
                                    }
                                }
                            }
                        ]
                    }
                ]
            }
        ],
        "services": [
            {
                "name": "PetStorePets",
                "description": "Make calls to find pets",
                "type": "http",
                "settings": {
                    "url": "http://petstore.swagger.io/v2/pet/:id"
                }
            },
            {
                "name": "PetStoreUsersByName",
                "description": "Make calls to find users",
                "type": "http",
                "settings": {
                    "url": "http://petstore.swagger.io/v2/user/:username"
                }
            }
        ]
    }
}
```
### Sample Usage
This trigger depends on support files which can be generated with cli tool by passing proto file. Usage of tool can be found [here](https://github.com/TIBCOSoftware/mashling/tree/master/docs/cli#grpc).<br>

Sample demonstration of this trigger can be found in gRPC [recipe](https://github.com/TIBCOSoftware/mashling-recipes/tree/master/recipes).

#### Note
Currently This Trigger handles.<br>
1. Unary methods propagation.
2. GET method is supported in REST end point.
3. REST path params can be mapped through params output key.
4. Routing can be done based on method names.