---
title: Amazon SNS
weight: 4606
---

# Amazon SNS
This activity allows you to send SMS text messages using Amazon Simple Notification Services.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/awssns
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "accessKey",
      "type": "string",
      "required": "true"
    },
    {
      "name": "secretKey",
      "type": "string",
      "required": "true"
    },
    {
      "name": "region",
      "type": "string",
      "required": "true",
      "allowed" : ["us-east-2","us-east-1","us-west-1","us-west-2","ap-south-1","ap-northeast-2","ap-southeast-1","ap-southeast-2","ap-northeast-1","cn-northwest-1","ca-central-1","eu-central-1","eu-west-1","eu-west-2","sa-east-1"]
    },
    {
      "name": "smsType",
      "type": "string",
      "allowed" : ["Promotional", "Transactional"],
	  "value": "Promotional"
    },
    {
      "name": "from",
      "type": "string",
      "required": "true"
    },
    {
      "name": "to",
      "type": "string",
      "required": "true"
    },
    {
      "name": "message",
      "type": "string",
      "required": "true"
    }
  ],
  "output": [
    {
      "name": "messageId",
      "type": "string"
    }
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| accessKey   | True     | Your Amazon access key ID |         
| secretKey   | True     | Your Amazon secret sccess Key |
| region      | True     | The default AWS region to use. See [here](http://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html) for more detailed information on supported regions |
| smsType     | True     | The type of SMS to be sent (This can be either Promotional or Transactional) |
| from        | True     | The Sender ID for the SMS |
| to          | True     | The phone number (in international format) to which to send the SMS |
| message     | True     | The message you want to send |
| messageId   | False    | The unique message ID returned by AWS SNS |


## Example
Coming soon...