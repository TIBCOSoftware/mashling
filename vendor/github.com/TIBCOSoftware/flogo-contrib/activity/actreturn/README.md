---
title: Return
weight: 4602
---

# Return
This activity allows you to reply to a trigger invocation and map output values. After replying to the trigger, the flow ends (this will be the last actvity in your flow).

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/actreply
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
        "mapperOutputScope" : "action.output"
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
| mappings    | True     | An array of mapping that are executed when the activity runs |


## Example
The below example allows you to configure the activity to reply and set the output values to literals "name" (a string) and 2 (an integer).

```json
{
  "id": "actreturn_5",
  "name": "Return",
  "description": "Simple Return Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/actreturn",
    "input": {
  	"mappings":[
      { "type": "literal", "value": "name", "mapTo": "Output1" },
      { "type": "literal", "value": 2, "mapTo": "Output2" }
    ]
  }
}
```