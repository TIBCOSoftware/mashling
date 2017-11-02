# mqtt
This activity provides your Mashling application the ability to send MQTT messages.

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "content",
      "type": "any",
      "required": false
    },
    {
      "name": "topic",
      "type": "string",
      "required": true
    },
    {
      "name": "broker",
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
      "required": true
    },
    {
      "name": "password",
      "type": "string",
      "required": true
    },
    {
      "name": "qos",
      "type": "number",
      "required": true
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
| topic       | The MQTT topic to send the mesages to   |
| broker      | The MQTT broker URL |
| id          | The id for this MQTT client |
| user        | The user name for the MQTT broker |
| password    | The password for the MQTT broker |
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
      "name": "topic",
      "value": "test",
      "type": "string"
    },
    {
      "name": "broker",
      "value": "tcp://localhost:1883",
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
    },
    {
      "name": "qos",
      "value": "0",
      "type": "number"
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
