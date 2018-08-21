---
title: Channel
weight: 4603
---

# Channel
This activity allows you to put a value on a named channel in the flogo engine.


## Installation
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/channel
```

## Schema
Inputs and Outputs:

```json
{
  "settings": [
    {
      "name": "channel",
      "type": "string",
      "required": true
    }
  ],
  "input":[
    {
      "name": "channel",
      "type": "string"
    },
    {
      "name": "value",
      "type": "interface{}"
    }
  ],
  "output": [
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| channel    | True     | The channel to put the value on |

