# gateways
> Details on mashling gateway projects and associated CLI commands.

## Commands
#### create
This command is used to create a mashling gateway project.

*Create the base sample project with a specific name.*
	
	mashling create my_app
	
*Create a mashling gateway project from an existing mashling gateway descriptor.*
	
	mashling create -f myapp.json

### help
This command is used to display help on a particular command
	
	mashling help create

### list
This command is used to display components of a mashling gateway

	mashling help list

##Gateway Project

###Structure

The create command creates a basic structure and files for a gateway.


	my_app/
		mashling.json
		src/
			my_app/
				imports.go
				main.go
		vendor/
		
**files**

- *mashling.json* : mashling gateway configuration descriptor file
- *imports.go* : contains go imports for contributions (activities, triggers and models) used by the gateway
- *main.go* : main file for the engine.

**directories**	
	
- *vendor* : go libraries


## Gateway Configuration

### Gateway


The *mashling.json* file is the metadata describing the gateway project.

```json
{
	"gateway": {
		"name": "demo",
		"version": "1.0.0",
		"description": "This is the first microgateway app",
		"configurations": [],
		"triggers": [
			{
				"name": "rest_trigger",
				"description": "The trigger on 'pets' endpoint",
				"type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
				"settings": {
					"port": "9096",
					"method": "GET",
					"path": "/pets/{petId}"
				}
			}
		],
		"event_handlers": [
			{
				"name": "get_pet_success_handler",
				"description": "Handle the user access",
				"reference": "github.com/TIBCOSoftware/mashling/lib/flow/flogo.json",
				"params": {
					"uri": "petstore.swagger.io/v2/pet/3"
				}
			}
		],
		"event_links": [
			{
				"triggers": [
					"rest_trigger"
				],
				"dispatches": [
					{
						"handler": "get_pet_success_handler"
					}
				]
			}
		]
	}
}
```
### Steps to create and run a mashling app using mashling.json: ###

The mashling.json can be modified accordingly and new app can be created using the below command.

mashling create -f mashlingname.json gatewayname

Using command : "mashling create -f mashling.json mygateway" , mygateway will be created.

cd mygateway/bin

Run the App mygateway.exe

The below is the sample mashling.json:

```
{
  "mashling_schema": "0.2",
  "gateway": {
    "name": "demoRestGw",
    "version": "1.0.0",
    "display_name":"Rest Conditional Gateway",
    "description": "This is the rest based microgateway app",
    "configurations": [
      {
        "name": "restConfig",
        "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
        "description": "Configuration for rest trigger",
        "settings": {
          "port": "9096"
        }
      }
    ],
    "triggers": [
      {
        "name": "animals_rest_trigger",
        "description": "Animals rest trigger - PUT animal details",
        "type": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
        "settings": {
          "config": "${configurations.restConfig}",
          "method": "PUT",
		      "path": "/pets",
          "optimize":"true"
        }
      }
    ],
    "event_handlers": [
      {
        "name": "mammals_handler",
        "description": "Handle mammals",
        "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "birds_handler",
        "description": "Handle birds",
        "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
      },
      {
        "name": "animals_handler",
        "description": "Handle other animals",
        "reference": "github.com/TIBCOSoftware/mashling/lib/flow/RestTriggerToRestPutActivity.json"
      }
    ],
    "event_links": [
      {
        "triggers": ["animals_rest_trigger"],
        "dispatches": [
          {
            "if": "${trigger.content.name in (ELEPHANT,CAT)}",
            "handler": "mammals_handler"
          },
          {
            "if": "${trigger.content.name == SPARROW}",
            "handler": "birds_handler"
          },
          {
            "handler": "animals_handler"
          }
        ]
      }
    ]
  }
}
```
#### Dispatch Conditions

In the above example the condition is content based. The below formats can be used for content and header based routing.

| Condition Prefix | Description | Example |
|:----------|:-----------|:-------|
| trigger.content | Trigger content / payload based condition | trigger.content.name == CAT |
| trigger.header | HTTP trigger's header based condition | trigger.header.Accept == text/plain |

##### Preconditions:

For content based routing the content of the trigger should be a valid json.

##### Example conditions:

When the json is {"name": "CAT"} the following condition can be used trigger.content.name == CAT.

When the json is {"name": "CAT","details":{"color":"white"}} the following condition can be used trigger.content.details.color == white.

When the json is {"names":[{"nickname":"blackie"},{"nickname":"doggie"}]} the following condition can be used trigger.content.names[1].nickname == doggie

For Header based routing the condition will always be trigger.header.headername == headervalue

Also the following operators are supported and can be used in conditions:
==(equals),>(greater than),in,<(less than),!=(notequals) and notin.

