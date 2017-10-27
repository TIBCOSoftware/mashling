# tracer
This activity provides your Mashling application the ability to trace.

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "tracing",
      "type": "any",
      "required": false
    },
    {
      "name": "span",
      "type": "any",
      "required": false
    }
  ],
  "outputs": [
    {
      "name": "tracing",
      "type": "any"
    },
    {
      "name": "span",
      "type": "any"    }
  ]
}
```
## Inputs
| Key     | Description    |
|:------------|:---------------|
| tracing | The tracing context to span off of |
| span | The current span to close |

### Outputs
| Key    | Description   |
|:-----------|:--------------|
| tracing | The tracing context to forward |
| span | The current span to forward |

## Configuration Examples
### Start span
Configure an activity to start a span:

```json
{
  "id": 2,
  "type": 1,
  "activityRef": "github.com/jpollock/mashling/ext/flogo/activity/tracer",
  "name": "name-of-span",
  "attributes": [
    {
      "name": "tracing",
      "value": null
    }
  ],
  "inputMappings": [
    {
      "type": 1,
      "value": "{T.tracing}",
      "mapTo": "tracing"
    }
  ]
}
```

### Stop span
Configure an activity to stop a span:

```json
{
  "id": 4,
  "type": 1,
  "activityRef": "github.com/jpollock/mashling/ext/flogo/activity/tracer",
  "name": "stop-span",
  "attributes": [
    {
      "name": "span",
      "value": null
    }
  ],
  "inputMappings": [
    {
      "type": 1,
      "value": "{A2.span}",
      "mapTo": "span"
    }
  ]
},
```
