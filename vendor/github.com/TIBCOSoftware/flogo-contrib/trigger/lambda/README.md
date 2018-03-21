# tibco-lambda
This trigger provides your flogo application the ability to start a flow as an AWS Lambda function

## Installation

```bash
flogo install trigger github.com/TIBCOSoftware/flogo-contrib/trigger/lambda
```

## Schema
Settings, Outputs:

```json
{
  "settings": [
  ],
  "outputs": [
    {
      "name": "context",
      "type": "object"
    },
    {
      "name": "evt",
      "type": "string"
    }
  ]
}
```

A sample of the context object:

```json
{
  "awsRequestId":"",
  "functionName":"",
  "functionVersion":"",
  "logGroupName":"",
  "logStreamName":"",
  "memoryLimitInMB":0
}
```