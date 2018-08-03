---
title: REST
weight: 4706
---
# tibco-rest
This trigger provides your flogo application the ability to start a flow via REST over HTTP

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/rest
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings": [
    {
      "name": "port",
      "type": "integer"
    }
  ],
  "output": [
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
      "name": "content",
      "type": "object"
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "method",
        "type": "string",
        "required" : true
      },
      {
        "name": "path",
        "type": "string",
        "required" : true
      }
    ]
  }
}
```
## Settings
### Trigger:
| Setting     | Description    |
|:------------|:---------------|
| port | The port to listen on |         
### Endpoint:
| Setting     | Description    |
|:------------|:---------------|
| method      | The HTTP method |         
| path        | The resource path  |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the REST Trigger.

### POST
Configure the Trigger to handle a POST on /device

```json
{
  "triggers": [
    {
      "name": "flogo-rest",
      "settings": {
        "port": "8080"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://new_device_flow",
          "settings": {
            "method": "POST",
            "path": "/device"
          }
        }
      ]
    }
  ]
}
```

### GET
Configure the Trigger to handle a GET on /device/:id

```json
{
  "triggers": [
    {
      "name": "flogo-rest",
      "settings": {
        "port": "8080"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://get_device_flow",
          "settings": {
            "method": "GET",
            "path": "/device/:id"
          }
        }
      ]
    }
  ]
}
```
