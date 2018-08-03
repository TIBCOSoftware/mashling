---
title: WebSocket Message
weight: 4621
---

# Send WebSocket message
This activity allows you to send a message to a WebSocket server.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/sendwsmessage
```

## Schema
Inputs and Outputs:

```json
{
"input":[
    {
      "name": "Server",
      "type": "string",
      "value": ""
    },
    {
      "name": "Channel",
      "type": "string",
      "value": ""
    },
    {
      "name": "Destination",
      "type": "string",
      "value": ""
    },
    {
      "name": "Message",
      "type": "string",
      "value": ""
    },
    {
      "name": "Username",
      "type": "string",
      "value": ""
    },
    {
      "name": "Password",
      "type": "string",
      "value": ""
    }
  ],
  "output": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| Server      | False    | The WebSocket server to connect to (e.g. `localhost:9191`) |         
| Channel     | False    | The channel to send the message to (e.g. `/channel`)   |
| Destination | False    | The destination to send the message to (e.g. `sample`) |
| Message     | False    | The message to send |
| Username    | False    | The username to connect to the WebSocket server (e.g. `user`) |
| Password    | False    | The password to connect to the WebSocket server (e.g. `user`) |
| output      | False    | A string with the result of the action |

## Configuration Examples
The below example sends a message `Hello World`
```json
{
  "id": "wsmessage",
  "name": "Send WebSocket Message",
  "description": "This activity sends a message to a WebSocket enabled servers",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/wsmessage",
    "input": {
      "Server": "localhost:9191",
      "Channel": "/channel",
      "Destination": "sample",
      "Message": "Hello World",
      "Username": "user",
      "Password": "passwd"
    }
  }
}
```