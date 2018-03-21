![gofmt status](https://img.shields.io/badge/gofmt-compliant-green.svg?style=flat-square) ![golint status](https://img.shields.io/badge/golint-compliant-green.svg?style=flat-square) ![automated test coverage](https://img.shields.io/badge/test%20coverage-1%20testcase-orange.svg?style=flat-square)

# sendWSMessage
This activity sends a message to a WebSocket enabled server like TIBCO eFTL.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/sendwsmessage
```

## Schema
Inputs and Outputs:

```json
{
"inputs":[
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
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| Server      | The WebSocket server to connect to (e.g. `localhost:9191`) |         
| Channel     | The channel to send the message to (e.g. `/channel`)   |
| Destination | The destination to send the message to (e.g. `sample`) |
| Message     | The actual message to send |
| Username    | The username to connect to the WebSocket server (e.g. `user`) |
| Password    | The password to connect to the WebSocket server (e.g. `user`) |

## Configuration Examples
The below configuration would connect to a WebSocket server based on TIBCO eFTL and send a message saying `Hello World`
```json
      {
        "id": 2,
        "name": "Send a message to a WebSocket server",
        "type": 1,
        "activityType": "sendWSMessage",
        "attributes": [
          {
            "name": "Server",
            "value": "localhost:9191",
            "type": "string"
          },
          {
            "name": "Channel",
            "value": "/channel",
            "type": "string"
          },
          {
            "name": "Destination",
            "value": "sample",
            "type": "string"
          },
          {
            "name": "Message",
            "value": "Hello World",
            "type": "string"
          },
          {
            "name": "Username",
            "value": "user",
            "type": "string"
          },
          {
            "name": "Password",
            "value": "password",
            "type": "string"
          }
        ]
      }
```

## Contributors
[Leon Stigter](https://github.com/retgits)