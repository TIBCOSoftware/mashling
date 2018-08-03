---
title: CoAP
weight: 4702
---
# tibco-coap
This trigger provides your flogo application the ability to start a flow via CoAP

## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/coap
```

## Schema
Settings, Outputs and Endpoint:

```json
"settings": [
  {
    "name": "port",
    "type": "integer",
  }
],
"output": [
  {
    "name": "payload",
    "type": "string"
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
```
## Settings
### Trigger:
| Setting     | Description    |
|:------------|:---------------|
| port | Used to override the standard CoAP server port of 5683 |         
### Endpoint:
| Setting     | Description    |
|:------------|:---------------|
| method      | The CoAP method |         
| path        | The resource path  |


## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the CoAP Trigger.

### POST
Configure the Trigger to handle a CoAP POST message with path /device/refresh

```json
{
  "triggers": [
    {
      "name": "flogo-coap",
      "settings": {},
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "method": "POST",
            "path": "/device/refresh"
          }
        }
      ]
    }
  ]
}
```

## Testing

Do to some simple testing of the CoAP trigger, you can use the [Copper (Cu)](https://addons.mozilla.org/en-US/firefox/addon/copper-270430) plugin for Firefox.<br><br>
Once you have the plugin installed, you can interact with the trigger by going to: [coap://localhost:5683](coap://localhost:5683) in Firefox.
