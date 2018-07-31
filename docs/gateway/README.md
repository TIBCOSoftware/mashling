## Table of Contents

- [Overview](#overview)
- [Usage](#usage)
  * [Health Check](#healthcheck)
- [Configuration](#configuration)
  * [Triggers](#triggers)
  * [Dispatches](#dispatches)
  * [Routes](#routes)
  * [Steps](#steps)
  * [Services](#services)
    * [HTTP](#services-http)
    * [JS](#services-js)
    * [Flogo Activity](#services-flogo-activity)
    * [Flogo Flow](#services-flogo-flow)
    * [Anomaly](#services-anomaly)
    * [SQL Detector](#services-sqld)
    * [gRPC](#services-grpc)
  * [Responses](#responses)
  * [Policies Proposal](#policies)
    * [Simple Policy](#simple-policy)
    * [Complex Policy](#complex-policy)

## <a name="overview"></a>Overview

The mashling-gateway powers the core event driven routing engine of the Mashling project. This core binary can run all versions of the mashling schema to date, however for the purposes of this document, we will focus on the `1.0` configuration schema.

## <a name="usage"></a>Usage

The gateway binary has the following command line arguments available to setup and specify how you would like the binary to operate.

They can be found by running:

```bash
./mashling-gateway -h
```

The output and flags are:

```bash
A static binary that executes Mashling gateway logic defined in a mashling.json configuration file. Complete documentation is available at https://github.com/TIBCOSoftware/mashling

Version: v0.3.3-internal-29-gf6c81fd-dirty
Build Date: 2018-04-03T10:11:33-0400

Usage:
  mashling-gateway [flags]
  mashling-gateway [command]

Available Commands:
  help        Help about any command
  version     Prints the mashling-gateway version

Flags:
  -c, --config string          mashling gateway configuration (default "mashling.json")
  -C, --config-cache string    location of the configuration artifacts cache (default ".cache")
  -E, --config-cache-enabled   cache post-processed configuration artifacts locally (default true)
  -d, --dev                    run mashling in dev mode
  -e, --env-var-name string    name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -h, --help                   help for mashling-gateway
  -l, --load-from-env          load the mashling gateway configuration from an environment variable
  -p, --ping-enabled           enable gateway ping service (default true)
  -P, --ping-port string       configure mashling gateway ping service port (default "9090")

Use "mashling-gateway [command] --help" for more information about a command.
```

Right now, `dev` mode just reloads the running gateway instance when a change is detected in the `mashling.json` file but the behavior is inconsistent between triggers.

### <a name="healthcheck"></a>Health Check

An integrated ping service is used to know if a gateway instance is alive and healthy.

The health check ping service is enabled by default and configured to run on port `9090`. You can specify a different port at startup time via:

```bash
./mashling-gateway -c <path to mashling json> -P <ping port value>
```

You can also disable the ping service via:
```bash
./mashling-gateway -c <path to mashling json> -p=false
```

The health check endpoint is available at `http://<GATEWAY IP>:<PING-PORT>/ping` with an expected result of:
```json
{"response":"Ping successful"}
```

A more detailed health check response is available at `http://<GATEWAY IP>:<PING-PORT>/ping/details` with an example result of:
```json
{"Version":"0.2","Appversion":"1.0.0","Appdescription":"This is the first microgateway app"}
```

## <a name="configuration"></a>Configuration

The `mashling.json` configuration file is what contains all details related to the runtime behavior of a mashling-gateway instance. The file can be named anything and pointed to via the `-c` or `--config` flag.

A mashling configuration file specifies the appropriate schema version to load and validate against via the `mashling_schema` key. This is located at the top level of the configuration JSON schema. All other components specifying runtime behavior are contained within a `gateway` key and will be explained in detail below.

Example configuration files for the `1.0` schema version can be found in the [V2 example recipes folder](../../examples/recipes/v2). The corresponding schema can be found [here](../../internal/pkg/model/v2/schema/schema.json).

### <a name="triggers"></a>Triggers

Triggers in Mashling are, currently, just Flogo triggers. Any Flogo trigger *should* work with the `1.0` schema specification. For the purposes of most of our examples, Mashling implemented triggers that conform to Flogo's specification are used.

An example trigger that listens for and dispatches HTTP requests looks like:

```json
{
  "name": "MyProxy",
  "description": "Animals rest trigger - PUT animal details",
  "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
  "settings": {
    "port": "9096"
  },
  "handlers": [
    {
      "dispatch": "Pets",
      "settings": {
        "autoIdReply": "false",
        "method": "GET",
        "path": "/pets/{petId}",
        "useReplyHandler": "false"
      }
    }
  ]
}
```

You can map the execution of a trigger to a specific dispatch via the `handlers` array in the trigger configuration. This allows you to send the execution context to a different process flow based off of some specific settings. The `dispatch` value must map to a name in the `dispatches` array.

### <a name="dispatches"></a>Dispatches

Dispatches are used to map trigger invocation with a set of possible execution routes. A dispatch has a name and receives the execution context from a trigger when that name is mapped via the trigger's handler. A dispatch is simple a name and an array of routes. A simple dispatch looks like:

```json
{
  "name": "Pets",
  "routes": ["..."]
}
```

### <a name="routes"></a>Routes

Routes define the actual execution logic of a dispatch. Each route in a dispatch comes with a condition value in the `if` key. The mashling engine will evaluate this condition within the trigger context. The first route with a condition to evaluate to `true` will then be the route executed. Only **one** route is executed per triggered flow. Once a route is selected by the mashling engine the steps defined therein will be evaluated and executed in the order in which they are defined. If a route is marked as `"async": true` then the execution will be asynchronous and the trigger will immediately be returned a response.

A simple route looks like:

```json
{
  "if": "payload.pathParams.petId >= 8 && payload.pathParams.petId <= 15",
  "async": false,
  "steps": ["..."]
}
```

### <a name="steps"></a>Steps

Each route is composed of a number of steps. Each step is evaluated in the order in which it is defined via an optional `if` condition. If the condition is `true`, that step is executed. If that condition is `false` the execution context moves onto the next step in the process and evaluates that one. A blank or omitted `if` condition always evaluates to `true`.

A simple step looks like:

```json
{
  "if": "payload.pathParams.petId == 9",
  "service": "PetStorePets",
  "input": {
    "method": "GET",
    "pathParams.id": "${payload.pathParams.petId}"
  }
}
```

As you can see above, a step consists of a simple condition, a service reference, input parameters, and (not shown) output parameters. The `service` must map to a service defined in the `services` array that is defined outside of a dispatch. Input key and value pairs are translated and handed off to the service execution. Output key value pairs are translated and retained after the service has executed. Values wrapped with `${}` are evaluated as variables within the context of the execution.

### <a name="services"></a>Services

A service defines a function or activity of some sort that will be utilized in a step within an execution flow. Services have names, types, and settings. Currently supported types are `http`, `js`, `flogoActivity`, `flogoFlow`, `anomaly`, and `sqld`. Services may call external endpoints like HTTP servers or may stay within the context of the mashling gateway, like the `js` service. Once a service is defined it can be used as many times as needed within your routes and steps.

#### <a name="services-http"></a>HTTP

The `http` service type executes an HTTP request against a specified target `url` and returns a response.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| path | string | The path of the URL |
| pathParams | JSON object | Key/value pairs representing parameters to interpolate in the URL and path  |
| method | string | The method to use when invoking the HTTP request (GET, PUT, POST, PATCH, DELETE)|
| url | string | The target URL of the HTTP request |
| body | string | Body of the HTTP request |
| headers | JSON object | Key/value pairs representing headers to send to the HTTP target|
| query | JSON object | Key/value pairs representing query parameters that are appended to the URL |
| timeout | integer | Timeout in seconds for this HTTP request (default is 5 seconds) |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| statusCode | integer | The HTTP status code of the response |
| body | JSON object | The response body |
| headers | JSON object | The key/value pairs representing the headers returned from the HTTP target |

A sample `service` definition is:

```json
{
  "name": "PetStorePets",
  "description": "Make calls to find pets",
  "type": "http",
  "settings": {
    "url": "http://petstore.swagger.io/v2/pet/:id"
  }
}
```

An example `step` that invokes the above `PetStorePets` service using `pathParams` is:

```json
{
  "service": "PetStorePets",
  "input": {
    "method": "GET",
    "pathParams.id": "${payload.pathParams.petId}"
  }
}
```

Utilizing and extracting the response values can be seen in both a conditional evaluation:

```json
{"if": "PetStorePets.response.body.status == 'available'"}
```

and a response handler:

```json
{
  "if": "PetStorePets.response.body.status == 'available'",
  "error": false,
  "output": {
    "code": 200,
    "format": "json",
    "body.pet": "${PetStorePets.response.body}",
    "body.inventory": "${PetStoreInventory.response.body}"
  }
}
```

#### <a name="services-js"></a>JS

The `js` service type evaluates a javascript `script` along with provided `parameters` and returns the result as the response.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| script | string | The javascript code to evaluate |
| parameters | JSON object | Key/value pairs representing parameters to evaluate within the context of the script  |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| error | bool | The HTTP status code of the response |
| errorMessage | string | The error message |
| result | JSON object | The result object from the javascript code  |

A sample `service` definition is:

```json
{
  "name": "JSCalc",
  "description": "Make calls to a JS calculator",
  "type": "js",
  "settings": {
    "script": "result.total = parameters.num * 2;"
  }
}
```

An example `step` that invokes the above `JSCalc` service using `parameters` is:

```json
{
  "if": "PetStorePets.response.body.status == 'available'",
  "service": "JSCalc",
  "input": {
    "parameters.num": "${PetStoreInventory.response.body.available}"
  }
}
```

Utilizing the response values can be seen in a response handler:

```json
{
  "if": "PetStorePets.response.body.status == 'available'",
  "error": false,
  "output": {
    "code": 200,
    "format": "json",
    "body.pet": "${PetStorePets.response.body}",
    "body.inventory": "${PetStoreInventory.response.body}",
    "body.availableTimesTwo": "${JSCalc.response.result.total}"
  }
}
```

#### <a name="services-flogo-activity"></a>Flogo Activity

The `flogoActivity` service type executes a Flogo Activity defined by `ref` with the provided `inputs` values.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| ref | string | The URI representing the Flogo Activity |
| inputs | JSON object | Key/value pairs representing inputs to pass to the Flogo Activity execution context  |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| done | bool | If the activity is done |
| error | string | An error message if a message occurred |
| outputs | JSON object | The output of this activity execution |

Input keys nested under the `inputs` key are specific to the type of Flogo Activity that is referenced.

A sample `service` definition is:

```json
{
  "name": "PetStorePets",
  "description": "Get pets by ID from the petsore.",
  "type": "flogoActivity",
  "settings": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
    "inputs": {
      "uri": "http://petstore.swagger.io/v2/pet/:petId",
      "method": "GET"
    }
  }
}
```

An example `step` that invokes the above `PetStorePets` service using `inputs` is:

```json
{
  "service": "PetStorePets",
  "input": {
    "inputs.pathParams": "${payload.pathParams}"
  }
}
```

Utilizing the response values can be seen in a response handler:

```json
{
  "error": false,
  "output": {
    "code": 200,
    "format": "json",
    "body": "${PetStorePets.response.outputs.result}"
  }
}
```

#### <a name="services-flogo-flow"></a>Flogo Flow

The `flogoFlow` service type executes a complete Flogo Flow defined by `reference` or `definition` with the provided `inputs` values.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| reference | string | The URI representing the location of a Flogo Flow (Github, URL, local file) |
| definition | JSON object | A complete Flogo Flow defined inline |
| inputs | JSON object | Key/value pairs representing inputs to pass to the Flogo Flow execution context  |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| done | bool | If the flow is done |
| error | string | An error message if a message occurred |
| outputs | JSON object | The output of this flow execution |

Input keys nested under the `inputs` key are specific to the expectations of the specific Flogo Flow.

A sample `service` definition is:

```json
{
  "name": "FlogoRestGetFlow",
  "description": "Make GET calls against a remote HTTP service using a Flogo flow.",
  "type": "flogoFlow",
  "settings": {
    "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestGetActivity.json"
  }
}
```

An example `step` that invokes the above `FlogoRestGetFlow` service using `inputs` is:

```json
{
  "service": "FlogoRestGetFlow",
  "input": {
    "inputs.pathParams": "${payload.pathParams}"
  }
}
```

Utilizing the response values can be seen in a response handler:

```json
{
  "error": false,
  "output": {
    "code": 200,
    "format": "json",
    "body": "${FlogoRestGetFlow.response.outputs.data}"
  }
}
```

#### <a name="services-anomaly"></a>Anomaly

The `anomaly` service type implements anomaly detection for payloads. The anomaly detection algorithm is based on a [statistical model](https://fgiesen.wordpress.com/2015/05/26/models-for-adaptive-arithmetic-coding/) for compression. The anomaly detection algorithm computes the relative [complexity](https://en.wikipedia.org/wiki/Kolmogorov_complexity), K(payload | previous payloads), of a payload and then updates the statistical model. A running mean and standard deviation of the complexity is then computed using [this](https://dev.to/nestedsoftware/calculating-standard-deviation-on-streaming-data-253l) algorithm. If the complexity of a payload is some number of deviations from the mean then it is an anomaly. An anomaly is a payload that is statistically significant relative to previous payloads. The anomaly detection algorithm uses real time learning, so what is considered an anomaly can change over time.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| payload | JSON object | A payload to do anomaly detection on |
| context | string | Allows a different statistical model to be used for payloads with different sources |
| depth | number |  The size of the statistical model. Defaults to 2 |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| complexity | number | How unusual the payload is in terms of standard deviations from the mean |
| count | number | The number of payloads that have been processed |

A sample `service` definition is:

```json
{
  "name": "Anomaly",
  "description": "Look for anomalies",
  "type": "anomaly",
  "settings": {
    "context": "test",
    "depth": 3
  }
}
```

An example `step` that invokes the above `Anomaly` service using `payload` is:

```json
{
  "service": "Anomaly",
  "input": {
    "payload": "${payload}"
  }
}
```

Utilizing the response values can be seen in a response handler:

```json
{
  "if": "(Anomaly.count < 100) || (Anomaly.complexity < 3)",
  "error": false,
  "output": {
    "code": 200,
    "data": "${Update.response.body}"
  }
}
```

#### <a name="services-sqld"></a>SQL Detector

The `sqld` service type implements SQL injection attack detection. Regular expressions and a [GRU](https://en.wikipedia.org/wiki/Gated_recurrent_unit) recurrent neural network are used to detect SQL injection attacks.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| payload | JSON object | A payload to do SQL injection attack detection on |
| file | string | An optional file name for custom neural network weights |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| attack | number | The probability that the payload is a SQL injection attack |
| attackValues | JSON object | The SQL injection attack probability for each string in the payload |

A sample `service` definition is:

```json
{
  "name": "SQLSecurity",
  "description": "Look for sql injection attacks",
  "type": "sqld"
}
```

An example `step` that invokes the above `SQLSecurity` service using `payload` is:

```json
{
  "service": "SQLSecurity",
  "input": {
    "payload": "${payload}"
  }
}
```

Utilizing the response values can be seen in a response handler:

```json
{
  "if": "SQLSecurity.attack > 80",
  "error": true,
  "output": {
    "code": 403,
    "data": {
      "error": "hack attack!",
      "attackValues": "${SQLSecurity.attackValues}"
    }
  }
}
```
#### <a name="services-grpc"></a>gRPC

The `grpc` service type works as gRPC client and will communicate to the given address with given method parameters.

The service `settings` and available `input` for the request are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
| grpcMthdParamtrs | JSON object | A grpcMthdParamtrs payload which holds full information like method parameters, service name, proto name, method name etc.|
| hosturl | string | A gRPC end point url with port |

The available response outputs are as follows:

| Name   |  Type   | Description   |
|:-----------|:--------|:--------------|
|body | JSON object | The response object from gRPC end server |

A sample `service` definition is:

```json
{
    "name": "PetStoreUsers",
    "description": "Make calls to grpc end point",
    "type": "grpc",
    "settings": {
        "hosturl": "localhost:9000"
    }
}
```

An example `step` that invokes the above `PetStoreUsers` service using `grpcMthdParamtrs` is:

```json
{
 "service": "PetStoreUsers",
 "input": {
 "grpcMthdParamtrs": "${payload.grpcData}"
 }
}
```

Response handler:

```json
{
  "error": false,
  "output": {
      "code": 200,
      "data": "${PetStoreUsers.response.body}"
  }
}
```

### <a name="responses"></a>Responses

Each route has an optional set of responses that can be evaluated and returned to the invoking trigger. Much like routes, the first response with an `if` condition evaluating to true is the response that gets executed and returned. A response contains an `if` condition, an `error` boolean, a `complex` boolean, and an `output` object. The `error` boolean dictates whether or not an error should be returned to the engine. The `complex` boolean dictates whether to use the `Reply` or `ReplyWithData` function. A value of `true` causes the `ReplyWithData` function to be used when sending the response back to the trigger. The `output` is evaluated within the context of the execution and then sent back to the trigger as well.

A simple response looks like:

```json
{
  "if": "PetStorePets.response.body.status == 'available'",
  "error": false,
  "complex": false,
  "output": {
    "code": 200,
    "format": "json",
    "body.pet": "${PetStorePets.response.body}",
    "body.inventory": "${PetStoreInventory.response.body}"
  }
}
```

### <a name="policies"></a>Policies (Proposed Solution, Take 4: Updated 3-07-18)

Policies are called out in the JSON Schema and the types for the V2 package, however, they are not yet implemented. This section of the document outlines the third iteration of a proposed policy design. This has been reworked following feedback from two previous sessions with the team.

The new proposed implementation is to treat policies as distinct entities from services and to make each policy invocation atomic. The notion of hooks for policies are also introduced in this design. As with most entities in the model, a conditional expression is optional and is mostly useful for `after` policy hooks and for feedback into a policy that is invoked in the corresponding `before` hook. Lifecycle hook specification is optional. If it is omitted the behavior for all policies specified under that `policies` key is the same as if the `before` hook was used.

This iteration of the policy design adds a policy block to `dispatches` and also expands the schema definition of the policy invocation blocks to introduce the notion of hooks. These hooks look like `beforeRoute`, `afterRoute`, `beforeStep`, `afterStep`, etc... and dictate the invocation order for the included policies. The ability to add a one off lower level invocation can be achieved by adding the policy to the `policies` key in that lower level entity.

Unlike the previous proposals, an interrupt is not required to achieve any of the example policies outlined below. An interrupt is left in the example below simply because it is a useful flow construct, but it is not required for policies to function.

Providing these entry points to polices allows one to support something simple like a rate limiter that returns a simple yes or no before executing the steps. It also provides the ability to wrap a call in a circuit breaker via the `after[Dispatch|Route|Step]` policy hook.

#### <a name="simple-policy"></a>Simple Policy Example
A simple HTTP proxy example with two policies (rate limiter and API key validation) added before the HTTP backend invocation happens is below. This example also demonstrates a simplified way of declaring a policy invocation: the `before` and `after` lifecycle hooks are excluded resulting in a default of `before` invocation behavior.

```json
{
  "mashling_schema": "1.0",
  "gateway": {
    "name": "MyProxy",
    "version": "1.0.0",
    "description": "This is a simple proxy.",
    "triggers": [
      {
        "name": "MyProxy",
        "description": "Animals rest trigger - PUT animal details",
        "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
        "settings": {
          "port": "9096"
        },
        "handlers": [
          {
            "dispatch": "Pets"
          }
        ]
      }
    ],
    "dispatches": [
      {
        "name": "Pets",
        "routes": [
          {
            "policies": [
              {
                "policy": "GlobalRateLimiter",
                "input": {
                  "key": "${payload.ipAddress}"
                }
              },
              {
                "policy": "APIKeyAuth",
                "input": {
                  "key": "${payload.headers.APIKey}"
                }
              }
            ],
            "steps": [
              {
                "service": "MySpecialBackend",
                "input": {
                  "pathParams.id": "${payload.pathParams.petId}"
                }
              }
            ],
            "responses": [
              {
                "output": {
                  "code": "${MySpecialBackend.response.code}",
                  "format": "json",
                  "body.pet": "${MySpecialBackend.response.body}",
                  "body.inventory": "${MySpecialBackend.response.body}"
                }
              }
            ]
          }
        ]
      }
    ],
    "services": [
      {
        "name": "MySpecialBackend",
        "description": "Make calls to do stuff",
        "type": "http",
        "settings": {
          "url": "http://petstore.swagger.io/v2/pet/:id"
        }
      }
    ],
    "policies": [
      {
        "name": "GlobalRateLimiter",
        "description": "Rate limit all requests",
        "type": "rateLimiter",
        "settings": {
          "perSecond": 100
        }
      },
      {
        "name": "APIKeyAuth",
        "description": "Test API key.",
        "type": "apiKeyAuth",
        "settings": {
          "url": "https://www.somewherespecial.com"
        }
      }
    ]
  }
}

```

#### <a name="complex-policy"></a>Complex Policy Example
A complex configuration file that has a contrived example using all of the hooks is as follows:

```json
{
  "mashling_schema": "1.0",
  "gateway": {
    "name": "PolicyExample",
    "version": "1.0.0",
    "description": "This is a simple proxy.",
    "triggers": [
      {
        "name": "MyProxy",
        "description": "Animals rest trigger - PUT animal details",
        "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
        "settings": {
          "port": "9096"
        },
        "handlers": [
          {
            "dispatch": "Pets",
            "settings": {
              "autoIdReply": "false",
              "method": "GET",
              "path": "/pets/{petId}",
              "useReplyHandler": "false"
            }
          }
        ]
      }
    ],
    "dispatches": [
      {
        "name": "Pets",
        "policies": {
          "beforeDispatch": [
            {
              "policy": "Splunk"
            }
          ],
          "afterDispatch": [
            {
              "policy": "Splunk"
            }
          ]
        },
        "routes": [
          {
            "if": "payload.pathParams.petId >= 8 && payload.pathParams.petId <= 15",
            "policies": {
              "beforeRoute": [
                {
                  "policy": "GlobalRateLimiter",
                  "input": {
                    "key": "${payload.ipAddress}"
                  }
                },
                {
                  "policy": "CircuitBreaker"
                }
              ],
              "beforeStep": [
                {
                  "policy": "Splunk"
                }
              ],
              "afterStep": [
                {
                  "policy": "Splunk"
                }
              ],
              "beforeResponse": [
                {
                  "policy": "Splunk"
                }
              ],
              "afterResponse": [
                {
                  "policy": "Splunk"
                }
              ],
              "beforeInterrupt": [
                {
                  "policy": "Splunk"
                }
              ],
              "afterInterrupt": [
                {
                  "policy": "Splunk"
                }
              ],
              "afterRoute": [
                {
                  "if": "PetStorePets.response.error == true",
                  "policy": "CircuitBreaker",
                  "input": {
                    "failed": true
                  }
                }
              ]
            },
            "steps": [
              {
                "service": "PetStorePets",
                "input": {
                  "method": "GET",
                  "pathParams.id": "${payload.pathParams.petId}"
                },
                "interrupt": "PetStorePets.response.error == true"
              },
              {
                "if": "PetStorePets.response.body.status == 'available'",
                "policies": {
                  "beforeStep": [
                    {
                      "policy": "OneOffPolicyInvocationForJustThisStep"
                    }
                  ]
                },
                "service": "PetStoreInventory",
                "input": {
                  "method": "GET"
                }
              }
            ],
            "interrupts": [
              {
                "if": "PetStorePets.response.error == true",
                "service": "RemoteErrorNotification",
                "input": {
                  "body.message": "${PetStorePets.response.errorMessage}"
                }
              }
            ],
            "responses": [
              {
                "if": "payload.pathParams.petId == 13",
                "error": true,
                "output": {
                  "code": 404,
                  "format": "json",
                  "body": "petId is invalid"
                }
              },
              {
                "if": "PetStorePets.response.body.status != 'available'",
                "error": true,
                "output": {
                  "code": 403,
                  "format": "json",
                  "body": "Pet is unavailable."
                }
              },
              {
                "if": "PetStorePets.response.body.status == 'available'",
                "error": false,
                "output": {
                  "code": 200,
                  "format": "json",
                  "body.pet": "${PetStorePets.response.body}",
                  "body.inventory": "${PetStoreInventory.response.body}"
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
        "name": "PetStoreInventory",
        "description": "Get pet store inventory.",
        "type": "http",
        "settings": {
          "url": "http://petstore.swagger.io/v2/store/inventory"
        }
      },
      {
        "name": "RemoteErrorNotification",
        "description": "Send error details somewhere custom.",
        "type": "http",
        "settings": {
          "method": "POST",
          "url": "http://www.errorsarebad.io/report_error"
        }
      }
    ],
    "policies": [
      {
        "name": "GlobalRateLimiter",
        "description": "Rate limit all requests",
        "type": "rateLimiter",
        "settings": {
          "perSecond": 100
        }
      },
      {
        "name": "CircuitBreaker",
        "description": "Stop hitting broken routes.",
        "type": "circuitBreaker",
        "settings": {
          "maxFails": 5
        }
      },
      {
        "name": "Splunk",
        "description": "Send my information to Splunk.",
        "type": "splunk",
        "settings": {
          "format": "${time} - ${error} - ${message}"
        }
      }
    ]
  }
}
```
