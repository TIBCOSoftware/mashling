# Project Mashling

[![Build Status](https://travis-ci.org/TIBCOSoftware/mashling.svg?branch=master)](https://travis-ci.org/TIBCOSoftware/mashling)

Project Mashling<sup>TM</sup> is an open source event-driven microgateway.

Project Mashling highlights include:
* Ultra lightweight: 10-50x times less compute resource intensive
* Event-driven by design
* Complements Service Meshes
* Co-exists with API management platforms in a federated API Gateway model

Project Mashling consists of the following open source repos:
* [mashling](http://github.com/TIBCOSoftware/mashling): This is the main repo that includes the below components
	- CLI to build Mashling apps
	- Mashling triggers and activities
	- Library to build Mashling extensions
* [mashling-recipes](http://github.com/TIBCOSoftware/mashling-recipes): This is the repo that includes recipes that illustrate configuration of common microgateway patterns. These recipes are curated and searchable via [mashling.io](http://mashling.io)

Additional developer tooling is included in below open source repo that contains the VSCode plugin for Mashling configuration:
* [VSCode Plugin for Mashling](https://github.com/TIBCOSoftware/vscode-extension-mashling)

## Installation

### Prerequisites
* The Go programming language 1.7 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system
* In order to simplify development and building in Go, we are using the **gb** build tool.  It can be downloaded from [here](https://getgb.io).

### Install Mashling
    go get github.com/TIBCOSoftware/mashling/...

### Update Mashling
    go get -u github.com/TIBCOSoftware/mashling/...

## Getting Started
A Mashling microgateway is created using the **Mashling** CLI tool.  The tool can be used to create a gateway from an existing *mashling.json* or to create a simple base gateway to get you started.  In this example we will walk you through creating the base/sample gateway.

To create the base gateway, which consists of a REST trigger and a simple event handler flow with a log activity, you use the following commands.


```bash
mashling create myApp

```

```bash
cd myApp/bin folder
./myapp
```

- Mashling will start a REST server
- Test it by sending sample HTTP events eg: http://localhost:9096/pets/2

The built in sample microgateway is based off the following mashling.json.  This file can be modified to add additional triggers and event handlers.

```json
{
	"mashling_schema": "0.2",
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


For more details about the REST Trigger go [here](https://github.com/TIBCOSoftware/mashling/tree/master/ext/flogo/trigger/gorillamuxtrigger)

## Documentation
For additional documentation on **Mashling** CLI tool, go [here](https://github.com/TIBCOSoftware/mashling/blob/master/cli/README.md)



## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

### Build Mashling from source
```
$go get github.com/TIBCOSoftware/mashling/cli/...

$cd $GOPATH/src/github.com/TIBCOSoftware/mashling/cli

[optional, only if building from branch]
$git checkout my_branch

[need to manually go get all dependencies for example:]
$go get github.com/xeipuuv/gojsonschema

$go install ./...
```

## License
Mashling is licensed under a BSD-type license. See license text [here](https://github.com/TIBCOSoftware/mashling/blob/master/TIBCO%20LICENSE.txt).


### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/mashling/issues)
