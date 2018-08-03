---
title: Coachbase
weight: 4608
---

# Couchbase
This activity allows you to connect to a Couchbase server

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/couchbase
```

## Schema
Inputs and Outputs:

```json
{
  "input": [
    {
      "name": "key",
      "type": "string",
      "required": true
    },
    {
      "name": "data",
      "type": "string"
    },
    {
      "name": "method",
      "type": "string",
      "allowed": [
        "Insert",
        "Upsert",
        "Remove",
        "Get"
      ],
      "value": "Insert",
      "required": true
    },
    {
      "name": "expiry",
      "type": "integer",
      "value": 0,
      "required": true
    },
    {
      "name": "server",
      "type": "string",
      "required": true
    },
    {
      "name": "username",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
    },
    {
      "name": "bucket",
      "type": "string",
      "required": true
    },
    {
      "name": "bucketPassword",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "output",
      "type": "any"
    }
  ]
}
```
## Settings
| Setting        | Required | Description |
|:---------------|:---------|:------------|
| key            | True     | The document key identifier |         
| data           | False    | The document data (when the method is `get` this field is ignored) |
| method         | True     | The method type (Insert, Upsert, Remove or Get). This field defaults to `Insert` |
| expiry         | True     | The document expiry (default: 0) |
| server         | True     | The Couchbase server (e.g. *couchbase://127.0.0.1*) |
| username       | False    | Cluster username |
| password       | False    | Cluster password |
| bucket         | True     | The bucket name |
| bucketPassword | False    | The bucket password if any |

## Example
The below example is for an upsert into Couchbase

```json
{
  "id": "couchbase_1",
  "name": "Couchbase connector",
  "description": "Manage Couchbase documents in a specified bucket",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/couchbase",
    "input": {
      "key": "foo",
      "data": "bar",
      "method": "Upsert",
      "expiry": 0,
      "server": "couchbase://127.0.0.1",
      "username": "Administrator",
      "password": "password",
      "bucket": "test",
      "bucketPassword": ""
    }
  }
}
```