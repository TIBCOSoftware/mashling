# mashling-ping-activity
This activity provides your mashling application the ability to reply gateway information as a json.

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
  "activityType": "mashling-ping-activity",
  "name": "Respond OK",
  "attributes": [
    { "name": "code", "value": 200 }
  ]
}
```
