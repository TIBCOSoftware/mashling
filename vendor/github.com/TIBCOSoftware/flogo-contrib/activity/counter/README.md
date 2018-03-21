# tibco-counter
This activity provides your flogo application the ability to use a global counter.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/counter
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
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
  "outputs": [
    {
      "name": "value",
      "type": "integer"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| counterName | The name of the counter |         
| increment   | Increment the counter |
| reset       | Reset the counter |
Note: if reset is set to true, increment is ignored
## Configuration Examples
### Increment
Configure a task to increment a 'messages' counter:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-counter",
  "name": "Increment Message Count",
  "attributes": [
    { "name": "counterName", "value": "messages" },
    { "name": "increment", "value": true }
  ]
}
```
### Get
Configure a task to get the 'messages' counter:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-counter",
  "name": "Get Message Count",
  "attributes": [
    { "name": "counterName", "value": "messages" }
  ]
}
```
### Reset
Configure a task to reset the 'messages' counter:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-counter",
  "name": "Reset Message Count",
  "attributes": [
    { "name": "counterName", "value": "messages" }
    { "name": "reset", "value": true }
  ]
}
```
