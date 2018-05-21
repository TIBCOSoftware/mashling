# tibco-lambda
This activity provides native Lambda invocation capabilities to your Flogo apps. You can invoke a lambda function via ARN and provide the access key and secret for authentication.

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
    },
    {
      "name": "tracing",
      "type": "any",
      "required" : false
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
    },
    {
      "name": "tracing",
      "type": "any"
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
| tracing     | The tracing context to forward                 |

## Output
The output is
| Setting     | Description                                    |
|:------------|:-----------------------------------------------|
| value       | A struct containing the Status and response Payload from the function. |
| result      | The reponse payload     |
| status      | Thestatus               |
| tracing     | The tracing context to forward                 |
