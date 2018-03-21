# tibco-coap
This activity provides your flogo application the ability to send a CoAP message.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/coap
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "method",
      "type": "string",
      "required": true
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "type",
      "type": "string"
    },
    {
      "name": "messageId",
      "type": "integer"
    },
    {
      "name": "options",
      "type": "params"
    },
    {
      "name": "payload",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "response",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| method    | The CoAP method (POST,GET,PUT,DELETE)|
| uri   | The CoAP resource URI |         
| queryParams | The query parameters |
| type      | Message Type (Confirmable, NonConfirmable, Acknowledgement, Reset) |
| messageId | ID used to detect duplicates and for optional reliability |
| options   | CoAP options |
| payload   | The message payload |


## Configuration Examples
### Simple
Configure a task in flow to send a "hello world" message via CoAP:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-coap",
  "name": "Send CoAP Message",
  "attributes": [
    { "name": "method", "value": "POST" },
    { "name": "address", "value": "coap://localhost:5683/device" },
    { "name": "type", "value": "Confirmable" },
    { "name": "messageId", "value": 12345 },
    { "name": "payload", "value": "hello world" },
    { "name": "options", "value": {"ETag":"tag", "MaxAge":2 }
  ]
}
```
