# tibco-mqtt
This trigger provides your mashling application the ability to start a flow via MQTT.
It supports content based conditional routing and distributed tracing.

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
    {
      "name": "broker",
      "type": "string"
    },
    {
      "name": "id",
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
      "name": "store",
      "type": "string"
    },
    {
      "name": "qos",
      "type": "number"
    },
    {
      "name": "cleansess",
      "type": "boolean"
    },
    {
      "name": "tracer",
      "type": "string",
      "required": false
    },
    {
      "name": "tracerEndpoint",
      "type": "string",
      "required": false
    },
    {
      "name": "tracerToken",
      "type": "string",
      "required": false
    },
    {
      "name": "tracerDebug",
      "type": "boolean",
      "required": false
    },
    {
      "name": "tracerSameSpan",
      "type": "boolean",
      "required": false
    },
    {
      "name": "tracerID128Bit",
      "type": "boolean",
      "required": false
    }
  ],
  "outputs": [
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "pathParams",
      "type": "params"
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
    },
    {
      "name": "message",
      "type": "string"
    },
    {
      "name": "tracing",
      "type": "any"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "topic",
        "type": "string"
      },
      {
        "name": "Condition",
        "type": "string"
      }
    ]
  }
}
```

### Settings
| Key    | Description   |
|:-----------|:--------------|
| broker | The MQTT broker |
| id | The id sent to the broker |
| user | The user for the broker |
| password | The password for the broker |
| store | The store for the broker |
| qos | The quality of service for the broker |
| cleansess | Use clean session with the broker |
| tracer | The tracer to use: noop, zipkin, appdash, or lightstep |
| tracerEndpoint | The url or address of the tracer (ZipKin, APPDash)|
| tracerToken | The token for tracing access (LightStep) |
| tracerDebug | Debug mode for the tracer (ZipKin) |
| tracerSameSpan | Put client side and server side annotations in same span (ZipKin) |
| tracerID128Bit | Use 128 bit ids (ZipKin) |

### Outputs
| Key    | Description   |
|:-----------|:--------------|
| params | Emulated HTTP request params |
| pathParams | Emulated HTTP request path params |
| queryParams | Emulated HTTP request query params |
| content | Emulated HTTP request payload |
| message | The MQTT payload |
| tracing | Tracing context |

### Handler settings
| Key    | Description   |
|:-----------|:--------------|
| topic | The MQTT topic to listen to |
| Condition | Handler condtion |

#### Supported Handler conditions

| Condition Prefix | Description | Example |
|:----------|:-----------|:-------|
| trigger.content | Trigger content / payload based condition | trigger.content.region == APAC |
| env | Environment flag / variable based condition | env.APP_ENVIRONMENT == UAT |

### MQTT Payload Options

The below keys are reserved in the MQTT payload.

| Key | Description | Example |
|:----------|:-----------|:-------|
| replyTo | Which topic to send the reply to | {"replyTo": "atopic"} |
| pathParams | The parameters that would be found in the HTTP URL path | {"pathParams": {"id": "1"}} |
| queryParams | The query that would be found in the HTTP URL query | {"queryParams": {"names": "NameA,NameB,NameC"}} |

### Sample Mashling Gateway Recipe

See [this recipe](https://github.com/TIBCOSoftware/mashling-recipes/tree/master/recipes/mqtt-gateway)
