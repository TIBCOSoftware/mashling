# mashling/cli
> Command line tool for building **Mashling** based gateways.

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
A Mashling gateway is created using the **mashling** CLI tool.  The tool can be used to create a gateway from an existing *mashling.json* or to create a simple base gateway to get you started.  In this example we will walk you through creating the base gateway.

To create the base gateway, which consists of a REST trigger and a simple event handler flow with a log activity, use the following command:


```bash
mashling create myApp

```

Start the base gateway by performing the following commands:

```bash
cd myApp/bin
./myapp
```

- Mashling will start a REST server
- Send GET request to run the flow. eg: http://localhost:9096/pets/2

The base gateway is built with the following mashling.json. It can be edited to add additional triggers and event handlers. A variation of base gateway is called a Mashling recipe. Mashling recipes can be download from mashling.io.


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
Additional documentation on Mashling and the CLI tool

### Mashling cli tool ###
Details about the Mashling cli commands are [here](docs/gateway.md)

### Mashling json configuration ###
Details about the Mashling json configuration are [here](docs/gateway.md)

### Mashling triggers ###
For more details about the Mashling REST Trigger, go [here](https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/gorillamuxtrigger)

For more details about the Mashling Kafka Trigger, go [here](https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/kafkasubrouter)

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

### Build Mashling CLI from source
```
$go get github.com/TIBCOSoftware/mashling/...

$cd $GOPATH/src/github.com/TIBCOSoftware/mashling

[optional, only if building from branch]
$git checkout my_branch

[need to manually go get all dependencies for example:]
$go get github.com/xeipuuv/gojsonschema

$go install ./...
```
Mashling CLI is built and installed in $GOPATH/bin

### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/mashling/issues)
