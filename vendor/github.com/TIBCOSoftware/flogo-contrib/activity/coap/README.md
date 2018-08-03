---
title: CoAP
weight: 4607
---

# CoAP
This activity allows you to send a CoAP message.


## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/coap
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "method",
      "type": "string",
      "required": true,
      "allowed" : ["GET", "POST", "PUT", "DELETE"]
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
  "output": [
    {
      "name": "response",
      "type": "string"
    }
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| uri         | True     | The CoAP resource URI |
| method      | True     | The CoAP method (Accepted values are POST, GET, PUT, and DELETE) |
| queryParams | False    | The query parameters |
| type        | False    | Message Type (Confirmable, NonConfirmable, Acknowledgement, Reset) |
| messageId   | False    | ID used to detect duplicates and for optional reliability |
| options     | False    | CoAP options |
| payload     | False    | The message payload |


## Example
The below example sends a "hello world" message via CoAP:

```json
{
  "id": "coap",
  "name": "Send CoAP Message",
  "description": "Simple CoAP Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/coap",
    "input": {
      "method": "POST",
      "address": "coap://localhost:5683/device",
      "type": "Confirmable",
      "messageId": 12345,
      "payload": "hello world",
      "options": {"ETag":"tag", "MaxAge":2 }
    }
  }
}
```