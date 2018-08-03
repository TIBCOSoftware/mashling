---
title: Error
weight: 4610
---

# Error
This activity allows you to cause an explicit error in the flow (throw an error).


## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/error
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "data",
      "type": "object"
    }
  ],
  "output": [
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| message     | False    | The error message you want to throw |         
| data        | False    | The error data you want to throw |

## Configuration Examples
The below example throws a simple error with a message:

```json
{
  "id": "error_1",
  "name": "Throw Error",
  "description": "Simple Error Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/error",
    "input": {
      "message": "Unexpected Threshold Value"
    }
  }
}
```