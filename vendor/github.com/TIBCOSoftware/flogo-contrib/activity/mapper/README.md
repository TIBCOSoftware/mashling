---
title: Mapper
weight: 4616
---

# Mapper
This activity allows you to map values to the working attribute set of a flow.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/mapper
```

## Schema
Input and Output:

```json
{
  "input":[
    {
      "name": "mappings",
      "type": "array",
      "required": true,
      "display": {
        "name": "Mapper",
        "type": "mapper",
        "mapperOutputScope" : "action"
      }
    }
  ],
  "output": [
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| mappings    | True     | An array of mappings that are executed when the activity runs |

## Example
The below example allows you to configure the activity to reply and set the output values to literals "name" (a string) and 2 (an integer).

```json
{
  "id": "mapper_6",
  "name": "Mapper",
  "description": "Simple Mapper Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/mapper",
    "input": {
      "mappings": [
        {
          "mapTo": "FlowAttr1",
          "type": "assign",
          "value": "$activity[log_3].message"
        }
      ]
    }
  }
}
```
