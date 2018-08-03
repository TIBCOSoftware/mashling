---
title: CLI
weight: 4701
---
# tibco-cli
This trigger provides your flogo application the ability to run as a CLI app, that is, accept input via the CLI & run once till completion and return the results to stdout.

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/cli
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "output": [
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
        "type": "boolean"
      }
    ]
  }
}
```
## Settings
### Trigger
| Output      | Description                        |
|:------------|:-----------------------------------|
| args        | An array of command line arguments |  


The array contains the command-line flags from `os.Args[2:]`, `os.Args[1]` is used to determine which flow is called. So a Flogo app with a CLI trigger that is started like:
```
$ ./cliapp -myflow -param1 foo -param2 bar 
```
will result in the engine executing the flow where the `handler.settings.command` field is `myflow` and pass on the other four arguments in the array `args`.

### Handler:
| Setting      | Description                          |
|:-------------|:-------------------------------------|
| command      | The command invoked                  |         
| default      | Indicates if its the default command |


## Example Configurations

Triggers are configured via the triggers section of your application. The following are some example configuration of the CLI Trigger.

### No command
Configure the Trigger to execute one flow

```json
{
"triggers": [
    {
      "id": "cli_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/cli",
      "name": "CLI Trigger",
      "description": "Simple CLI Trigger",
      "settings": {},
      "handlers": [
        {
          "action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:version"
            }
          },
          "settings": {
            "command": "version",
            "default": true
          }
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
      "id": "cli_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/cli",
      "name": "CLI Trigger",
      "description": "Simple CLI Trigger",
      "settings": {},
      "handlers": [
        {
          "action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:version"
            }
          },
          "settings": {
            "command": "version",
            "default": false
          }
        },
        {
          "action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:search"
            },
            "mappings": {
              "input": [
                {
                  "mapTo": "commandline",
                  "type": "assign",
                  "value": "$.args"
                }
              ]
            }
          },
          "settings": {
            "command": "search",
            "default": false
          }
        }
      ]
    }
  ]
}
```

### Map arguments
To use the _args_ array from the trigger you'll need to specify an input parameter of type array and map the args from the trigger to it.
```json
"mappings": {
  "input": [
    {
      "mapTo": "commandline",
      "type": "assign",
      "value": "$.args"
    }
  ]
}
```