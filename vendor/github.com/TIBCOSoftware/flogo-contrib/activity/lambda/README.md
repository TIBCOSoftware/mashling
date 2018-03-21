# tibco-lambda
This activity provides native Lambda invocation capabilities to your Flogo apps. You can invoke a lambda function via ARN and provide the access key and secret for authentication.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/lambda
```

## Schema
Inputs and Outputs:

```json
{
"inputs":[
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
  "outputs": [
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
| Setting     | Description                                    |
|:------------|:-----------------------------------------------|
| arn         | The ARN for the Lambda function to invoke      |
| region      | The AWS region                                 |
| accessKey   | Access key for the user to invoke the function |
| secretKey   | The users secret key                           |
| payload     | The payload. A JSON object.                    |

## Output
| Setting     | Description                                    |
|:------------|:-----------------------------------------------|
| value       | A struct containing the status and response    |
| result      | The response                                   |
| status      | The status                                     |
