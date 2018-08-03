---
title: REST
weight: 4618
---

# REST
This activity allows you to invoke a REST service.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/rest
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "method",
      "type": "string",
      "required": true,
      "allowed" : ["GET", "POST", "PUT", "PATCH", "DELETE"]
    },
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "proxy",
      "type": "string",
      "required": false
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
      "name": "header",
      "type": "params"
    },
    {
      "name": "skipSsl",
      "type": "boolean",
      "value": "false"
    },
    {
      "name": "content",
      "type": "any"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "any"
    },
    {
      "name": "status",
      "type": "integer"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| method      | True     | The HTTP method to invoke (Allowed values are GET, POST, PUT, DELETE, and PATCH) |         
| uri         | True     | The URI of the service to invoke |
| proxy       | False    | The address of the proxy server to be used |
| pathParams  | False    | The path parameters. This field is only required if you have params in your URI (for example http://.../pet/:id) |
| queryParams | False    | The query parameters |
| header      | False    | The header parameters |
| skipSsl     | False    | If set to true, skips the SSL validation (defaults to false)
| content     | False    | The message content you want to send. This field is only used in POST, PUT, and PATCH |


## Examples
### Simple
The below example retrieves a pet with number '1234' from the [swagger petstore](http://petstore.swagger.io):

```json
{
  "id": "rest_2",
  "name": "Invoke REST Service",
  "description": "Simple REST Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
    "input": {
      "method": "GET",
      "uri": "http://petstore.swagger.io/v2/pet/1234"
    }
  }
}
```

### Using Path Params
The below example is the same as above, itretrieves a pet with number '1234' from the [swagger petstore](http://petstore.swagger.io), but uses a URI parameter to configure the ID:

```json
{
  "id": "rest_2",
  "name": "Invoke REST Service",
  "description": "Simple REST Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
    "input": {
      "method": "GET",
      "uri": "http://petstore.swagger.io/v2/pet/:id",
      "params": { "id": "1234"}
    }
  }
}
```