# tibco-reply
This activity provides your flogo application the ability to reply to a trigger invocation.

## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/reply
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "code",
      "type": "integer",
      "required": true
    },
    {
      "name": "data",
      "type": "any"
    }
  ],
  "outputs": [
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| code        | The response code |         
| data        | The response data |

## Configuration Examples
### Simple
Configure a activity to respond with a simple http success code.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-reply",
  "name": "Respond OK",
  "attributes": [
    { "name": "code", "value": 200 }
  ]
}
```
