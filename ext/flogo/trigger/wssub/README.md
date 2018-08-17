# tibco-wssub
This trigger provides your Mashling application with the ability to subscribe to websocket message events and invokes `dispatch` with the contents of the message.

## Schema
Settings, Outputs and Handlers:

```json
{
 "settings":[
    {
      "name": "url",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "content",
      "type": "any"
    }
  ],
  "handler": {
    "settings": []
  }
```

## Example Configurations
This example flow subscribes to the syslog subject of bilbo's kafka server using a plain text connection with no authentication.

```json
{
  "triggers": [
    {
      "name": "WSMessageTrigger",
      "id": "WSMessageTrigger",
      "ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/wssub",
      "settings": {
        "url": "ws://localhost:8080/ws"
      },
      "output": null,
      "handlers": [
        {
          "settings": {
            "format": "json"
          },
          "output": null,
          "Action": null,
          "actionId": "flow1",
          "outputs": null
        }
      ],
      "outputs": null
    }
  ]
}
```
