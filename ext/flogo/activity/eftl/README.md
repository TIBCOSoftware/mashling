# eftl
This activity provides your Mashling application the ability to send EFTL messages.

## Schema
Inputs and Outputs:

```json
{
  "name": "eftl",
  "version": "0.0.1",
  "type": "flogo:activity",
  "ref": "github.com/TIBCOSoftware/mashling/ext/flogo/activity/eftl",
  "title": "A EFTL message producer",
  "description": "A EFTL message producer",
  "author": "Andrew Snodgrass <asnodgra@tibco.com>",
  "homepage": "https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/activity/eftl",
  "inputs":[
    {
      "name": "content",
      "type": "any",
      "required": false
    },
    {
      "name": "dest",
      "type": "string",
      "required": true
    },
    {
      "name": "url",
      "type": "string",
      "required": true
    },
    {
      "name": "id",
      "type": "string",
      "required": true
    },
    {
      "name": "user",
      "type": "string",
      "required": false
    },
    {
      "name": "password",
      "type": "string",
      "required": false
    },
    {
      "name": "ca",
      "type": "string",
      "required": false
    },
    {
      "name": "tracing",
      "type": "any",
      "required": false
    }
  ],
  "outputs": [
    {
      "name": "tracing",
      "type": "any"
    }
  ]
}
```

## Inputs
| Setting     | Description    |
|:------------|:---------------|
| content     | The message to send |
| dest        | The EFTL dest to send the messages to |
| url         | The EFTL server URL |
| id          | The id for this EFTL client |
| user        | The user name for the EFTL server |
| password    | The password for the EFTL server |
| ca          | The certificate authority for the EFTL client |
| tracing     | The tracing context |

## Outputs
| Setting     | Description    |
|:------------|:---------------|
| tracing     | The output tracing context |

## Configuration Example
```json
{
  "id": 2,
  "type": 1,
  "activityRef": "github.com/TIBCOSoftware/mashling/ext/flogo/activity/mqtt",
  "name": "mqtt",
  "attributes": [
    {
      "name": "content",
      "value": "test",
      "type": "string"
    },
    {
      "name": "dest",
      "value": "test",
      "type": "string"
    },
    {
      "name": "url",
      "value": "ws://localhost:9191/channel",
      "type": "string"
    },
    {
      "name": "id",
      "value": "mashling",
      "type": "string"
    },
    {
      "name": "user",
      "value": "",
      "type": "string"
    },
    {
      "name": "password",
      "value": "",
      "type": "string"
    }
  ],
  "inputMappings": [
    {
      "type": 1,
      "value": "${trigger.content}",
      "mapTo": "content"
    }
  ]
}
```
