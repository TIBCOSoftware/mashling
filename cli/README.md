# mashling/cli
> Command line tool for building **Mashling**-based gateways.

**Mashling** is a Micro-gateway framework written in Go. It was designed from the ground up to be robust enough for microservices.


## Installation
### Prerequisites
* The Go programming language 1.7 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system
* In order to simplify development and building in Go, we are using the **gb** build tool.  It can be downloaded from [here](https://getgb.io).  
* For a private repo, set up ssh interaction with github. Follow the instruction [here](https://help.github.com/articles/adding-a-new-ssh-key-to-your-github-account) and run the following command  
    git config --global url."git@github.com:".insteadOf "https://github.com/"

### Install Mashling
    go get github.com/TIBCOSoftware/mashling/...

### Update Mashling
    go get -u github.com/TIBCOSoftware/mashling/...

## Getting Started
A mashling gateway is created using the **mashling** CLI tool.  The tool can be used to create a gateway from an existing *mashling.json* or to create a simple base gateway to get you started.  In this example we will walk you through creating the base gateway.

To create the base gateway, which consists of a REST trigger and a simple event handler flow with a log activity, you use the following commands.


```bash
mashling create myApp

```

Start base gateway by

```bash
cd myApp/bin
./myapp
```

- Mashling will start a REST server
- Send GET request to run the flow. eg: http://localhost:9096/pets/2

The base gateway built with the following mashling.json. It can be edited to add additional triggers and event handlers. A variation of base gateway is called a mashling recipe. Mashling recipes can be download from mashling.io.


```json
{
	"mashling_schema": "0.2",
	"gateway": {
		"name": "demo",
		"version": "1.0.0",
		"description": "This is the first microgateway app",
		"configurations": [
			{
				"name": "kafkaConfig",
				"type": "github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub",
				"description": "Configuration for kafka cluster",
				"settings": {
					"BrokerUrl": "localhost:9092"
				}
			}
		],
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


For more details about the REST Trigger configuration go [here](https://github.com/TIBCOSoftware/flogo-contrib/tree/master/trigger/rest#example-configurations)

## Documentation
Additional documentation on mashling and the CLI tool

### mashling tool ###
  - details about mashling CLI commands are [here](docs/gateway.md)

### mashling triggers ###

For more details about the mashling REST Trigger, go [here](https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/gorillamuxtrigger)

For more details about the mashling KAFKA Trigger, go [here](https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/kafkasubrouter)

## Steps to create and run a mashling app using mashling.json: ##

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
### Dispatch Conditions

In the above example the condition is content based. The below formats can be used for content and header based routing.

| Condition Prefix | Description | Example |
|:----------|:-----------|:-------|
| trigger.content | Trigger content / payload based condition | trigger.content.name == CAT |
| trigger.header | HTTP trigger's header based condition | trigger.header.Accept == text/plain |

#### Preconditions:

For content based routing the content of the trigger should be a valid json.

#### Example conditions:

When the json is {"name": "CAT"} the following condition can be used trigger.content.name == CAT.

When the json is {"name": "CAT","details":{"color":"white"}} the following condition can be used trigger.content.details.color == white.

When the json is {"names":[{"nickname":"blackie"},{"nickname":"doggie"}]} the following condition can be used trigger.content.names[1].nickname == doggie

For Header based routing the condition will always be trigger.header.headername == headervalue

Also the following operators are supported and can be used in conditions:
==(equals),>(greater than),in,<(less than),!=(notequals) and notin.

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

### Build mashling CLI from source
```
$go get github.com/TIBCOSoftware/mashling/...

$cd $GOPATH/src/github.com/TIBCOSoftware/mashling

[optional, only if building from branch] 
$git checkout my_branch

[need to manually go get all dependencies for example:] 
$go get github.com/xeipuuv/gojsonschema

$go install ./... 
```
mashling CLI is built and installed in $GOPATH/bin

##License
mashling/cli is licensed under a BSD-type license. See TIBCO LICENSE.txt for license text.


### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/mashling/issues)
