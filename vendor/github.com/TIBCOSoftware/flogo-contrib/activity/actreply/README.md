---
title: Reply
weight: 4601
---

# Reply
This activity allows you to reply to a trigger invocation and map output values. After replying to the trigger, this activity will allow the flow to continue further.

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
| mappings    | True     | An array of mappings that are executed when the activity runs |

## Example
The below example allows you to configure the activity to reply and set the output values to literals "name" (a string) and 2 (an integer).

```json
{
  "id": "reply",
  "name": "Reply",
  "description": "Simple Reply Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/actreply",
    "input": {
  	"mappings":[
      { "type": "literal", "value": "name", "mapTo": "Output1" },
      { "type": "literal", "value": 2, "mapTo": "Output2" }
    ]
  }
}
```