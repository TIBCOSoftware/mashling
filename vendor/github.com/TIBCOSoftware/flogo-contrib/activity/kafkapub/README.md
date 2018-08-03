---
title: Publish Kafka Message
weight: 4613
---

# Publish Kafka Message
This activity allows you to send a Kafka message

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "BrokerUrls",
      "type": "string",
      "required": true
    },
    {
      "name": "Topic",
      "type": "string",
      "required": true
    },
    {
      "name": "Message",
      "type": "string",
      "required": true
    },
    {
      "name": "user",
      "type": "string",
      "required": false
    },
    {
      "name": "password",
      "type": "string",
      "required": false
    },
    {
      "name": "truststore",
      "type": "string",
      "required": false
    }
  ],
  "output": [
    {
      "name": "partition",
      "type": "int"
    },
    {
      "name": "offset",
      "type": "long"
    }
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| BrokerUrls  | True     | The Kafka cluster to connect to |
| Token       | True     | The Kafka topic on which to place the message |
| Message     | True     | The text message to send |
| user        | False    | If connectiong to a SASL enabled port, the userid to use for authentication |
| password    | False    | If connectiong to a SASL enabled port, the password to use for authentication |
| truststore  | False    | If connectiong to a TLS secured port, the directory containing the certificates representing the trust chain for the connection.  This is usually just the CACert used to sign the server's certificate |
| partition   | False    | Documents the partition that the message was placed on |
| offset      | False    | Documents the offset for the message                   |

## Examples
The below example sends a message to the 'syslog' topic.
```json
{
  "id": "kafkapub_1",
  "name": "Publish Kafka message",
  "description": "Publish a message to a kafka topic",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub",
    "input": {
      "BrokerUrls": "bilbo:9092",
      "Topic": "syslog",
      "Message": "mary had a little lamb"
    }
  }
}
```