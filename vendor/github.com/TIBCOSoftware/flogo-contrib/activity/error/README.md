# tibco-rest
This activity provides your flogo application the ability to cause an explicit error in the flow.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/error
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "data",
      "type": "object"
    }
  ],
  "outputs": [
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| message     | The error message |         
| data        | The error data |

## Configuration Examples

Configure a task in flow cause a simple error with a message:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-error",
  "name": "Throw Error",
  "attributes": [
    { "name": "message", "value": "Unexpected Threshold Value" }
  ]
}
```
