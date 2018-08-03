---
title: Counter
weight: 4609
---

# Counter
This activity allows you to use a global counter.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/counter
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "counterName",
      "type": "string",
      "required": true
    },
    {
      "name": "increment",
      "type": "boolean"
    },
    {
      "name": "reset",
      "type": "boolean"
    }
  ],
  "output": [
    {
      "name": "value",
      "type": "integer"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| counterName | True     | The name of the counter |         
| increment   | False    | If this field is set to true, increment the counter by one |
| reset       | False    | Reset the counter. _If reset is set to true, increment is ignored_|
| value       | False    | The value of the counter after executing the increment or reset |

## Examples
### Increment
The below example increments a 'messages' counter:

```json
{
  "id": "counter_1",
  "name": "Increment Counter",
  "description": "Simple Global Counter Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/counter",
    "input": {
      "counterName": "messages",
      "increment": true
    }
  }
}
```

### Get
The below example retrieves the last value of the 'messages' counter:

```json
{
  "id": "counter_1",
  "name": "Increment Counter",
  "description": "Simple Global Counter Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/counter",
    "input": {
      "counterName": "messages"
    }
  }
}
```

### Reset
The below example resets the 'messages' counter:

```json
{
  "id": "counter_1",
  "name": "Increment Counter",
  "description": "Simple Global Counter Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/counter",
    "input": {
      "counterName": "messages",
      "reset": true
    }
  }
}
```