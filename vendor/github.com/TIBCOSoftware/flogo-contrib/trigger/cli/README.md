# tibco-cli
This trigger provides your flogo application the ability to start a flow via REST over HTTP

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/cli
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "outputs": [
    {
      "name": "args",
      "type": "array"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "command",
        "type": "string"
      },
      {
        "name": "default",
        "type": "bool"
      }
    ]
  }
}
```
## Settings
### Trigger:
| Output     | Description    |
|:------------|:---------------|
| args | The array of arguments |         
### Handler:
| Setting     | Description    |
|:------------|:---------------|
| command      | The command invoked |         
| default      | Indicates if its the default command  |


## Example Configurations

Triggers are configured via the triggers section of your application. The following are some example configuration of the CLI Trigger.

### No command
Configure the Trigger to execute one flow

```json
{
    "triggers": [
      {
        "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/cli",
        "description": "Simple CLI trigger",
        "settings": {},
        "id": "main",
        "handlers": [
          {
            "settings": {
              "default": true
            },
            "actionId": "log_cli"
          }
        ]
      }
    ]
}
```

### Multiple Commands
Configure the Trigger to handle multiple commands

```json
{
    "triggers": [
      {
        "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/cli",
        "description": "Simple CLI trigger",
        "settings": {},
        "id": "main",
        "handlers": [
          {
            "settings": {
              "command": "list"
            },
            "actionId": "list_flow"
          },
          {
            "settings": {
              "command": "run"
            },
            "actionId": "run_flow"
          }
        ]
      }
    ]
}
```
