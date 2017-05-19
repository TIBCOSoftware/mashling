# gateways
> Details on mashling gateway projects and associated CLI commands.

## Commands
#### create
This command is used to create a mashling gateway project.

*Create the base sample project with a specific name.*
	
	mashling create my_app
	
*Create a mashling gateway project from an existing mashling gateway descriptor.*
	
	flogo create -f myapp.json

### help
This command is used to display help on a particular command
	
	flogo help build 

##Gateway Project

###Structure

The create command creates a basic structure and files for a gateway.


	my_app/
		flogo.json
		mashling.json
		src/
			my_app/
				imports.go
				main.go
		vendor/
		
**files**

- *flogo.json* : flogo project application configuration descriptor file
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
		"name": "demoapp",
		"version": "1.0.0",
		"description": "This is the first microgateway app",
		"configurations": [
			{
				"name": "kafkaConfig",
				"type": "github.com/wnichols/kafkasub",
				"description": "Configuration for kafka cluster",
				"settings": {
					"BrokerUrl": "localhost:9092"
				}
			}
		],
		"triggers": [
			{
				"name": "rest_trigger",
				"description": "The trigger on 'users' endpoint",
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
				"name": "get_user_success_handler",
				"description": "Handle the user access",
				"reference": "github.com/aambhaik/resources/response-flow.json",
				"params": {
					"uri": "petstore.swagger.io/v2/pet/3"
				}
			}
		],
		"event_links": [
			{
				"trigger": "rest_trigger",
				"success_paths": [
					{
						"handler": "get_user_success_handler"
					}
				]
			}
		]
	}
}
```