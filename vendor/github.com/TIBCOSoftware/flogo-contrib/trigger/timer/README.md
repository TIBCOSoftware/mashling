---
title: Timer
weight: 4707
---
# tibco-timer
This trigger provides your flogo application the ability to schedule a flow via scheduling service

## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/timer
```

## Schema
Outputs and Handler:

```json
{
  "output": [
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "repeating",
        "type": "string",
        "value": "false",
      },
      {
        "name": "notImmediate",
        "type": "string",
        "value": "true",
      },
      {
        "name": "startDate",
        "type": "string",
        "value": "2018-01-01T12:00:00Z00:00"
      },
      {
        "name": "hours",
        "type": "string"
      },
      {
        "name": "minutes",
        "type": "string"
      },
      {
        "name": "seconds",
        "type": "string"
      }
    ]
  }
}
```

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Timer Trigger.

### Only once and immediate
Configure the Trigger to run a flow immediately

```json
{
  "triggers": [
    {
      "name": "flogo-timer",
      "settings": {
      },
      "handlers": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "repeating": "false"
          }
        }
      ]
    }
  ]
}
```

### Only once at schedule time
Configure the Trigger to run a flow at a certain date/time. "startDate" settings format = "mm/dd/yyyy, hours:minutes:seconds"

```json
{
  "triggers": [
    {
      "name": "flogo-rest",
      "settings": {
        "port": "8080"
      },
      "handlers": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "repeating": "false",
            "startDate" : "2018-01-01T12:00:00Z00:00"
          }
        }
      ]
    }
  ]
}
```

### Repeating
Configure the Trigger to run a flow repeating every hours|minutes|seconds. If "notImmediate" set to true, the trigger will not fire immediately.  In this case the first execution will occur in 24 hours. If set to false the first execuction will will occur immediately.

```json
{
  "triggers": [
    {
      "name": "flogo-rest",
      "settings": {
        "port": "8080"
      },
      "handlers": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "repeating": "true",
            "notImmediate": "true",
            "hours": "24"
          }
        }
      ]
    }
  ]
}
```

### Repeating with start date
Configure the Trigger to run a flow at a certain date/time and repeating every hours|minutes|seconds

```json
{
  "triggers": [
    {
      "name": "flogo-rest",
      "settings": {
        "port": "8080"
      },
      "handlers": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "repeating": "true",
            "startDate" : "2018-01-01T12:00:00Z00:00",
            "hours": "64"
          }
        }
      ]
    }
  ]
}
```
