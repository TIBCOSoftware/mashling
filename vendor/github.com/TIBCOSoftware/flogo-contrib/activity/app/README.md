---
title: App
weight: 4604
---

# App
This activity allows you to set and use global attributes throughout your app.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/actreply
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
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
  "output": [
    {
      "name": "value",
      "type": "any"
    }
  ]
}
```

## Settings
| Setting        | Required | Description |
|:---------------|:---------|:------------|
| attribute      | True     | The name of the attribute |         
| operation      | True     | The operation to perform |
| type           | False    | The type of the attribute, only used with NEW operation |
| value (input)  | False    | The value of the attribute, only used with ADD and UPDATE |
| value (output) |          | The returned value of the attribute, only used with GET and UPDATE |

## Configuration Examples
### New
Add a new 'myAttr' attribute of type string with the initial value of _foo_:

```json
{
  "id": "app_5",
  "name": "Use Global Attribute",
  "description": "Simple Global App Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/app",
    "input": {
      "attribute": "myAttr",
      "operation": "ADD",
      "type": "string",
      "value": "MyValue"
    }
  }
}
```

### Get
Get the value of the 'myAttr' attribute:

```json
{
  "id": "app_5",
  "name": "Use Global Attribute",
  "description": "Simple Global App Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/app",
    "input": {
      "attribute": "myAttr",
      "operation": "GET",
      "type": "string",
      "value": "MyValue"
    }
  }
}
```

### Update
Update the value of the 'myAttr' attribute to _bar_:

```json
{
  "id": "app_5",
  "name": "Use Global Attribute",
  "description": "Simple Global App Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/app",
    "input": {
      "attribute": "myAttr",
      "operation": "UPDATE",
      "value": "MyValue"
    }
  }
}
```
