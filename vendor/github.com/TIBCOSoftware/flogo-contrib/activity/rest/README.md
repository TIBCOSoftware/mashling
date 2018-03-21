# tibco-rest
This activity provides your flogo application the ability to invoke a REST service.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/rest
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "method",
      "type": "string",
      "required": true
    },
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "pathParams",
      "type": "params"
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "object"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| method      | The HTTP method to invoke |         
| uri         | The uri of the resource   |
| pathParams  | The path parameters |
| queryParams | The query parameters |
| content     | The message content |
| params      | The path parameters (Deprecated) |
Note: 

* **pathParams**: Is only required if you have params in your URI ( i.e. http://.../pet/:id )
* **content**: Is only used in POST, PUT, PATCH

## Configuration Examples
### Simple
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io):

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for pet 1234",
  "attributes": [
    { "name": "method", "value": "GET" },
    { "name": "uri", "value": "http://petstore.swagger.io/v2/pet/1234" }
  ]
}
```
### Using Path Params
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io) via parameters.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for Pet",
  "attributes": [
    { "name": "method", "value": "GET" },
    { "name": "uri", "value": "http://petstore.swagger.io/v2/pet/:id" },
    { "name": "params", "value": { "id": "1234"} }
  ]
}
```
### Advanced
Configure a task in flow to get pet from the [swagger petstore](http://petstore.swagger.io) using a flow attribute to specify the id.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for Pet",
  "attributes": [
    { "name": "method", "value": "GET" },
    { "name": "uri", "value": "http://petstore.swagger.io/v2/pet/:id" },
  ],
  "inputMappings": [
    { "type": 1, "value": "petId", "mapTo": "params.id" }
  ]
}
```
