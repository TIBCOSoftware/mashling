---
title: MongoDB
weight: 4622
---

# MongoDB
This activity allows you to connect to a MongoDB server

## Installation
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/mongodb
```

## Schema
Inputs and Outputs:

```json
{
  "input": [
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "dbName",
      "type": "string",
      "required": true
    },
    {
      "name": "collection",
      "type": "string",
      "required": true
    },
    {
      "name": "method",
      "type": "string",
      "allowed": [
        "DELETE",
        "INSERT",
        "REPLACE",
        "UPDATE"
      ],
      "value": "INSERT",
      "required": true
    },
    {
      "name": "keyName",
      "type": "string"
    },
    {
      "name": "keyValue",
      "type": "string"
    },
    {
      "name": "value",
      "type": "any"
    }
  ],
  "output": [
    {
      "name": "output",
      "type": "any"
    },
    {
      "name": "count",
      "type": "integer"
    }
  ]
 }
```
## Settings
| Setting        | Required | Description |
|:---------------|:---------|:------------|
| uri            | True     | The MongoDB connection URI |         
| dbName         | True     | The name of the database
| collection     | True     | The collection to work on
| method         | True     | The method type (DELETE, INSERT, UPDATE or REPLACE). This field defaults to `INSERT` |
| keyName        | False    | The name of the key to use when looking up an object (used in DELETE, UPDATE, and REPLACE)
| keyValue       | False    | The value of the key to use when looking up an object (used in DELETE, UPDATE, and REPLACE)
| value          | False    | The value of the object (used in INSERT, UPDATE, and REPLACE)

## Example
The below example is for an insert into MongoDB

```json
{
  "id": "MongoDB_1",
  "name": "MongoDB connector",
  "description": "Insert MongoDB documents in a specified collection",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/mongodb",
    "input": {
      "uri": "mongodb://localhost:27017",
      "dbName": "mydb",
      "collection": "users",
      "method": "INSERT",
      "value": {"name":"theuser"}
    }
  }
}
```
