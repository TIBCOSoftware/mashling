---
title: Aggregate
weight: 4603
---

# Aggregate
This activity allows you to aggregate data and calculate an average or sliding average.


## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/aggregate
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "function",
      "type": "string",
      "required": true,
      "allowed" : ["block_avg", "moving_avg", "timeblockavg"]
    },
    {
      "name": "windowSize",
      "type": "integer",
      "required": true
    },
    {
      "name": "value",
      "type": "number"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "number"
    },
    {
      "name": "report",
      "type": "boolean"
    }
  ]
}
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| function    | True     | The aggregate fuction, currently only average is supported |
| windowSize  | True     | The window size of the values to aggregate |
| value       | False    | The value to aggregate |


## Example
The below example aggregates a 'temperature' attribute with a moving window of size 5:

```json
"id": "aggregate_4",
"name": "Aggregate",
"description": "Simple Aggregator Activity",
"activity": {
  "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/aggregate",
  "input": {
    "function": "average",
    "windowSize": "5"
  },
  "mappings": {
    "input": [
      {
        "type": "assign",
        "value": "temperature",
        "mapTo": "value"
      }
    ]
  }
```