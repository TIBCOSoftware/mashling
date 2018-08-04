# eFTL
This trigger provides your mashling application the ability to start a flow via EFTL.
It supports content based conditional routing and distributed tracing.

## Schema
Settings, Outputs and Dest:

```json
{
  "name": "tibco-eftl",
  "type": "flogo:trigger",
  "ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/eftl",
  "version": "0.0.1",
  "title": "Receive EFTL Message",
  "author": "Andrew Snodgrass <asnodgra@tibco.com>",
  "description": "EFTL Trigger",
  "homepage": "https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/eftl",
  "settings":[
    {
      "name": "url",
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
      "name": "ca",
      "type": "string"
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
      "name": "tracing",
      "type": "any"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "dest",
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
| url | The URL for the EFTL server |
| id | The id sent to the EFTL server |
| user | The user for the EFTL server |
| password | The password for the EFTL server |
| ca | The certificate authority for the EFTL server |
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
| content | The EFTL request payload |
| tracing | Tracing context |

### Handler settings
| Key    | Description   |
|:-----------|:--------------|
| dest | The EFTL destination to listen to |
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

See [this recipe](https://github.com/TIBCOSoftware/mashling-recipes/tree/master/recipes/eftl-gateway)
