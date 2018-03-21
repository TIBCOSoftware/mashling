# tibco-app
This activity provides your Flogo application the ability to use a global attributes.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/app
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "attribute",
      "type": "string",
      "required": true
    },
    {
      "name": "operation",
      "type": "string",
      "required" : true,
      "allowed" : ["ADD","GET","UPDATE"]
    },
    {
      "name": "type",
      "type": "string",
      "allowed" : [	"string", "integer", "number", "boolean", "object", "array", "params"]
    },
    {
      "name": "value",
      "type": "any"
    }
  ],
  "outputs": [
    {
      "name": "value",
      "type": "any"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| attribute   | The name of the attribute |         
| operation   | The operation to perform |
| type        | The type of the attribute, only used with NEW operation |
| value       | The value of the attribute, only used with ADD and UPDATE |
## Configuration Examples
### New
Configure a task to add a new 'myAttr' attribute:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-app",
  "name": "Add myAttr to application",
  "attributes": [
    { "name": "attribute", "value": "myAttr" },
    { "name": "operation", "value": "NEW" },
    { "name": "type", "value": "string" },
    { "name": "value", "value": "test" },    
  ]
}
```
### Get
Configure a task to get the 'myAttr' attribute:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-app",
  "name": "Get myAttr from Application",
  "attributes": [
      { "name": "attribute", "value": "myAttr" },
      { "name": "operation", "value": "GET" },
  ]
}
```
### Update
Configure a task to update the 'myAttr' attribute:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-app",
  "name": "Update myAttr Application attribute",
  "attributes": [
    { "name": "attribute", "value": "myAttr" },
    { "name": "operation", "value": "UDPATE" },
    { "name": "value", "value": "test" },    
  ]
}
```
