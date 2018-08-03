---
title: Trigger Lambda Function
weight: 4614
---

# Trigger Lambda function
This activity allows you to invoke an AWS Lambda function via ARN and provide the access key and secret for authentication.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/lambda
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "arn",
      "type": "string",
      "required": true
    },
    {
      "name": "region",
      "type": "string",
      "required": true
    },
    {
      "name": "accessKey",
      "type": "string",
      "required" : false
    },
    {
      "name": "secretKey",
      "type": "string",
      "required" : false
    },
    {
      "name": "payload",
      "type": "any",
      "required" : true
    }
  ],
  "output": [
    {
      "name": "value",
      "type": "any"
    },
    {
      "name": "result",
      "type": "any"
    },
    {
      "name": "status",
      "type": "integer"
    }
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| arn         | True     | The ARN of the Lambda function to invoke      |
| region      | True     | The AWS region in which you want to invoke the function |
| accessKey   | False    | AWS access key for the user to invoke the function |
| secretKey   | False    | AWS secret key for the user to invoke te function |
| payload     | True     | The payload to send to the function. This must be a valid JSON object. |
| value       | False    | A struct containing the status and response    |
| result      | False    | The response from the invocation |
| status      | False    | The status of the invocation |

## Examples
Coming soon...