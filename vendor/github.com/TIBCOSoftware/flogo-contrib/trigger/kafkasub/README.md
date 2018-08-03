---
title: Receive Kafka Message
weight: 4703
---
# tibco-kafkasub
This trigger provides your flogo application with the ability to subscribe to messages from a kafka cluster and start a flow with the contents of the message.  It is assumed that the messages plain text.  The trigger supports TLS and SASL.  
To make a TLS connection specifiy a trust dir containing the caroots for your kafka server and a broker URL which points to an SSL port.
To use SASL simply provide the username and password in the settings config.


## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub
```

## Schema
Settings, Outputs and Endpoint:

```json
{
 "settings":[
    {
      "name": "BrokerUrl",
      "type": "string"
    },
    {
      "name": "user",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
    },
    {
      "name": "truststore",
      "type": "string"
    }

  ],
  "output": [
    {
      "name": "message",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "Topic",
        "type": "string"
      },
      {
        "name": "partitions",
        "type": "string"
      },
      {
        "name": "group",
        "type": "string"
      },
      {
        "name": "offset",
        "type": "int"
      },
    ]
  }
```

## Example Configurations
This example flow subscribes to the syslog subject of bilbo's kafka server using a plain text connection with no authentication.

```json
{
  "name": "testkafka",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "My flogo application description",
  "triggers": [
    {
      "id": "my_kafka_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub",
      "settings": {
        "BrokerUrl": "bilbo:9092"
      },
      "handlers": [
        {
          "actionId": "my_simple_flow",
          "settings": {
            "Topic": "syslog"
          }
        }
      ]
    }
  ],
  "actions": [
    {
      "id": "my_simple_flow",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
      "data": {
        "flow": {
          "attributes": [],
          "rootTask": {
            "id": 1,
            "type": 1,
            "tasks": [
              {
                "id": 2,
                "type": 1,
                "activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
                "name": "log",
                "attributes": [
                  {
                    "name": "message",
                    "value": "Simple Log",
                    "type": "string"
                  }
                ],
                "inputMappings": [
                  {
                    "type": 1,
                    "value": "{T.message}",
                    "mapTo": "message"
                  }
                ]
              }
            ],
            "links": []
          }
        }
      }
    }
  ]
}```


To connect to a TLS port on a kafka cluster member:

  "triggers": [
    {
      "id": "my_kafka_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub",
      "settings": {
        "BrokerUrl": "bilbo:9093",
        "truststore": "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust"
      },
      "handlers": [
        {
          "actionId": "my_simple_flow",
          "settings": {
            "Topic": "syslog"
          }
        }
      ]
    }
  ],
In this scenario the kafka server on bilbo is running TLS on port 9093.  The CACert used to sign the server's certificate has been copied to the truststore directory to allow clients to connect.  At this time mutual auth is not implemented.



To connect to a port on a kafka cluster where SASL authorization is enabled

  "triggers": [
    {
      "id": "my_kafka_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub",
      "settings": {
        "BrokerUrl": "bilbo:9094",
        "user": "foo",
        "password": "bar"
      },
      "handlers": [
        {
          "actionId": "my_simple_flow",
          "settings": {
            "Topic": "syslog"
          }
        }
      ]
    }
  ],
In this scenario the kafka server on bilbo is running SASL enabled port 9094. The user and password will be used to authenticate the user.

