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

### swagger
This command is used to generater Swagger 2.0 docs for HTTP triggers in your mashling.json file.

Usage:

	mashling swagger

**options**

- *-f* : specify the mashling json (default is mashling.json).
- *-h* : the hostname where this mashling will be deployed (default is localhost).
- *-t* : the trigger name to target (default is all).
- *-o* : the output file to write the swagger.json to (default is stdout).

Example using a mashling.json with a single HTTP trigger:

	mashling swagger -f reference-gateway.json

Output:
```json
{
    "host": "localhost",
    "info": {
        "description": "This is the first microgateway app",
        "title": "demo",
        "version": "1.0.0"
    },
    "paths": {
        "/pets/{petId}": {
            "get": {
                "description": "The trigger on 'pets' endpoint",
                "parameters": [
                    {
                        "in": "path",
                        "name": "petId",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "The trigger on 'pets' endpoint"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "rest_trigger"
                ]
            }
        }
    },
    "swagger": "2.0"
}
```

Example using a more complex conditional gateway:

	mashling swagger -f rest-conditional-gateway.json

Output:
```json
{
    "host": "localhost",
    "info": {
        "description": "This is the rest based microgateway app",
        "title": "demoRestGw",
        "version": "1.0.0"
    },
    "paths": {
        "/pets": {
            "put": {
                "description": "Animals rest trigger - PUT animal details",
                "parameters": [],
                "responses": {
                    "200": {
                        "description": "Animals rest trigger - PUT animal details"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "animals_rest_trigger"
                ]
            }
        },
        "/pets/{petId}": {
            "get": {
                "description": "Animals rest trigger - get animal details",
                "parameters": [
                    {
                        "in": "path",
                        "name": "petId",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Animals rest trigger - get animal details"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "get_animals_rest_trigger"
                ]
            }
        }
    },
    "swagger": "2.0"
}
```

Example specifying a trigger:

	mashling swagger -f rest-conditional-gateway.json -t get_animals_rest_trigger

Output:
```json
{
    "host": "localhost",
    "info": {
        "description": "This is the rest based microgateway app",
        "title": "demoRestGw",
        "version": "1.0.0"
    },
    "paths": {
        "/pets/{petId}": {
            "get": {
                "description": "Animals rest trigger - get animal details",
                "parameters": [
                    {
                        "in": "path",
                        "name": "petId",
                        "required": true,
                        "type": "string"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Animals rest trigger - get animal details"
                    },
                    "default": {
                        "description": "error"
                    }
                },
                "tags": [
                    "get_animals_rest_trigger"
                ]
            }
        }
    },
    "swagger": "2.0"
}
```

Example sending output to a file instead of STDOUT:

	mashling swagger -f rest-conditional-gateway.json -o swagger.json

For more details please use:

	mashling help swagger

### publish
This command is used to publish HTTP triggers in your mashling.json file
to Mashery.

Usage:

    mashling publish -k key -s secret_key  -u username -p password -uuid  uuid -portal mashery_portal -h petstore.swagger.io

**options**

- *-f*      : specify the mashling json (default is mashling.json).
- *-k*      : the api key (required)
- *-s*      : the api secret key (required)
- *-u*      : username (required)
- *-p*      : password (required)
- *-portal* : the portal (required)
- *-uuid*   : the proxy uuid (required)
- *-h*			: the publicly available hostname where this mashling will be deployed (required)
- *-mock*		: true to mock, where it will simply display the transformed swagger doc; false to actually publish to Mashery (default is false).



Example (display transformed swagger doc only):

    mashling publish -k 12345  -s 6789  -u foo -p bar -uuid  xxxyyy -portal "tibcobanqio.api.mashery.com" -mock true  -h petstore.swagger.io

Example (publish to Mashery):

    mashling publish -k 12345  -s 6789  -u foo -p bar -uuid  xxxyyy -portal "tibcobanqio.api.mashery.com"  -h petstore.swagger.io

For more details please use:

    mashling help publish

#####

## Gateway Project

### Structure

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
