# tibco-log
This activity provides your flogo application with rudementary logging.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "message",
      "type": "string",
      "value": ""
    },
    {
      "name": "flowInfo",
      "type": "boolean",
      "value": "false"
    },
    {
      "name": "addToFlow",
      "type": "boolean",
      "value": "false"
    }
  ],
  "outputs": [
    {
      "name": "message",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| message   | The message to log |         
| flowInfo  | Append the flow information to the log message |
| addToFlow | Add the log message to the 'message' output of the activity |


## Configuration Examples
### Simple
Configure a task to log a message 'test message':

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-log",
  "name": "Log Message",
  "attributes": [
    { "name": "message", "value": "test message" }
  ]
}
```
### Advanced
Configure a task to log a 'petId' attribute as a message:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-log",
  "name": "Log PetId",
  "attributes": [],
  "inputMappings": [
    { "type": 1, "value": "petId", "mapTo": "message" }
  ]
}
```
