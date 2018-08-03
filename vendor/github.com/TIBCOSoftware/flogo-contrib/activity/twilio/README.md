---
title: Twilio
weight: 4620
---

# Twilio
This activity allows you to send a SMS via Twilio.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/twilio
```

## Schema
Inputs and Outputs:
```json
{
  "input":[
    {
      "name": "accountSID",
      "type": "string"
    },
    {
      "name": "authToken",
      "type": "string"
    },
    {
      "name": "from",
      "type": "string"
    },
    {
      "name": "to",
      "type": "string"
    },
    {
      "name": "message",
      "type": "string"
    }
  ],
  "output": []
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| accountSID  | False    | The Twilio account SID |         
| authToken   | False    | The Twilio auth token  |
| from        | False    | The Twilio number you are sending the SMS from |
| to          | False    | The number you are sending the SMS to. This field should be in the format '+15555555555' |
| message     | False    | The SMS message |

## Examples
The below example sends 'my text message' to '617-555-5555' via Twilio:
```json
{
  "id": "twilio",
  "name": "Send SMS Via Twilio",
  "description": "Simple Twilio Activity",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/twilio",
    "input": {
      "accountSID": "A...9",
      "authToken": "A...9",
      "from": "+12016901385",
      "to": "+16175555555",
      "message": "my text message"
    }
  }
}
```