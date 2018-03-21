# tibco-twilio
This activity provides your flogo application the ability to send a SMS via Twilio.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/twilio
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
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
  "outputs": []
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| accountSID | The Twilio account SID |         
| authToken  | The Twilio auth token  |
| from       | The Twilio number you are sending the SMS from |
| to         | The number you are sending the SMS to |
| message    | The SMS message |
Note: 
Phone numbers should be in the format '+15555555555'

## Configuration Examples
### Simple
Configure a task in flow to send 'my text message' to '617-555-5555' via Twilio:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-twilio",
  "name": "Send Text Message",
  "attributes": [
    { "name": "accountSID", "value": "A...9" },
    { "name": "authToken", "value": "A...9" },
    { "name": "from", "value": "+12016901385" },
    { "name": "to", "value": "+16175555555" },
    { "name": "message", "value": "my text message" }
  ]
}
```

### Advanced
Configure a task in flow to send 'my text message' to a number from a REST trigger's query parameter:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-twilio",
  "name": "Send Text Message",
  "attributes": [
    { "name": "accountSID", "value": "A...9" },
    { "name": "authToken", "value": "A...9" },
    { "name": "from", "value": "+12016901385" },
    { "name": "message", "value": "my text message" }
  ],
  "inputMappings": [
    { "type": 1, "value": "[T.queryParams].From", "mapTo": "to" }
  ]
}
```
