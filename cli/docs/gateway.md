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
				"type": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
				"settings": {
					"port": "9096",
					"method": "GET",
					"path": "/pets/:petId"
				}
			}
		],
		"event_handlers": [
			{
				"name": "get_pet_success_handler",
				"description": "Handle the user access",
				"reference": "github.com/TIBCOSoftware/mashling-lib/flow/flogo.json",
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
						"if": "trigger.content != undefined",
						"handler": "get_pet_success_handler"
					}
				]
			}
		]
	}
}
```