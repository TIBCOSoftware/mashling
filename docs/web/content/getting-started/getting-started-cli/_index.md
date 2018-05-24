---
title: CLI
weight: 3300
pre: "<i class=\"fa fa-terminal\" aria-hidden=\"true\"></i> "
---

Create a Mashling application using [CLI](https://github.com/TIBCOSoftware/mashling/blob/master/docs/cli/README.md) (the Command line tool for building Mashling-based gateways).

### Overview

The mashling-cli powers all non-runtime functionality associated with a Mashling configuration file. These are actions like constructing a customized Mashling binary, generating Swagger docs, publishing to Mashery, and more. This core binary can run all versions of the mashling schema.

### Usage
The cli binary has the following command line arguments available to specify commands and operation.

They can be found by running:

```bash
./mashling-cli -h
```
The output and flags are:

```
A CLI to build custom mashling-gateway instances, publish configurations to Mashery, and more. Complete documentation is available at https://github.com/TIBCOSoftware/mashling

Usage:
  mashling-cli [command]

Available Commands:
  create      Creates a customized mashling-gateway
  help        Help about any command
  publish     Publishes to supported platforms
  swagger     Creates a swagger 2.0 doc
  validate    Validates a mashling.json configuration file
  version     Prints the mashling-cli version

Flags:
  -c, --config string         mashling gateway configuration (default "mashling.json")
  -e, --env-var-name string   name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -h, --help                  help for mashling-cli
  -l, --load-from-env         load the mashling gateway configuration from an environment variable

Use "mashling-cli [command] --help" for more information about a command.
```

### Commands
Below is a list of the currently support commands via the mashling-cli binary.

#### Create
Create allows you to build customized mashling-gateway binaries that are re-usable and contain all of your custom dependencies.

The command details are as follows:

```
./mashling-cli create -h
```

```
Create a reusable customized mashling-gateway binary based off of the dependencies listed in your mashling.json configuration file

Usage:
  mashling-cli create [flags]

Flags:
  -A, --arch string   target architecture to build for (default is amd64, arm64 is only compatible with Linux)
  -h, --help          help for create
  -n, --name string   customized mashling-gateway name (default "mashling-custom")
  -N, --native        build the customized binary natively instead of using Docker
  -O, --os string     target OS to build for (default is the host OS, valid values are windows, darwin, and linux)

Global Flags:
  -c, --config string         mashling gateway configuration (default "mashling.json")
  -e, --env-var-name string   name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -l, --load-from-env         load the mashling gateway configuration from an environment variable
```

A simple example usage is:

```
./mashling-cli create -c examples/recipes/v2/customized-simple-synchronous-patter.json
```

By default, all of the build commands will run through Docker so as to simplify the setup required on your development machine. You can run these commands natively (assuming your development environment is setup correclty), by passing the *-N* flag to the *create* command.

You can also specify which target OS to build the customized binary for via the *-O* flag. Supported values are **windows**, **darwin** (for macOS), and **linux**. The default value is whatever the host operating system is at the time the create command is executed.

A target architecture to build the customized binary for can be specified via the *-A* flag. Supported values are **amd64**, and **arm64**. The default value is **amd64** and will suffice for the vast majority of use cases. Linux is the only compatible target OS for arm64 architectures.

#### Swagger
Swagger allows you to generate a Swagger 2.0 document based off of the provided mashling.json configuration file. Currently, it only works with HTTP-based triggers.

The command details are as follows:
```
./mashling-cli swagger -h
```
```
Creates a swagger 2.0 doc based off of the HTTP triggers in the mashling.json configuration file

Usage:
  mashling-cli swagger [flags]

Flags:
  -h, --help             help for swagger
  -H, --host string      the hostname where this mashling will be deployed (default "localhost")
  -o, --output string    the output file to write the swagger.json to (default is stdout)
  -t, --trigger string   the trigger name to target (default is all))

Global Flags:
  -c, --config string         mashling gateway configuration (default "mashling.json")
  -e, --env-var-name string   name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -l, --load-from-env         load the mashling gateway configuration from an environment variable
```
A simple example usage is:

```
./mashling-cli swagger -c examples/recipes/v1/rest-conditional-gateway.json
```
The resulting output is:
```
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
You can override the published hostname via the *-H* flag.

#### Publish
This command is used to publish HTTP triggers in your mashling.json file to the currently supported publish targets, namely Mashery and Consul.

Details are:
```
./mashling-cli publish -h
```
```
Publishes details of the mashling.json configuration file to various support platforms (currently Mashery and Consul)

Usage:
  mashling-cli publish [command]

Available Commands:
  consul      Publishes to Consul
  mashery     Publishes to Mashery

Flags:
  -h, --help   help for publish

Global Flags:
  -c, --config string         mashling gateway configuration (default "mashling.json")
  -e, --env-var-name string   name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -l, --load-from-env         load the mashling gateway configuration from an environment variable

Use "mashling-cli publish [command] --help" for more information about a command.
```
##### Mashery
Publishing to Mashery will take the HTTP triggers defined in your mashling.json configuration file and push them to your Mashery account. These account details are provided via command line arguments.

Details are:
```
./mashling-cli publish mashery -h
```
```
Publishes the details of the mashling.json configuration file Mashery

Usage:
  mashling-cli publish mashery [flags]

Flags:
  -k, --apiKey string        the API key
  -T, --apiTemplate string   json file that contains defaults for api/endpoint settings in mashery
  -d, --areaDomain string    the public domain of the Mashery gateway
  -i, --areaID string        the Mashery area id
  -h, --help                 help for mashery
  -H, --host string          the publicly available hostname where this mashling will be deployed (e.g. hostip:port)
  -I, --iodocs               true to create iodocs
  -m, --mock                 true to mock, where it will simply display the transformed swagger doc; false to actually publish to Mashery
  -p, --password string      password
  -s, --secretKey string     the secret key
  -t, --testplan             true to create package, plan and test app/key
  -u, --username string      username

Global Flags:
  -c, --config string         mashling gateway configuration (default "mashling.json")
  -e, --env-var-name string   name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -l, --load-from-env         load the mashling gateway configuration from an environment variable
```

Example mock usage that displays transformed swagger doc only:
```
./mashling-cli publish mashery -k 12345 -s 6789 -u foo -p bar -i xxxyyy -d "tibcobanqio.api.mashery.com" -m true -H petstore.swagger.io
```
Example usage that actually publishes to Mashery:
```
./mashling-cli publish mashery -k 12345 -s 6789 -u foo -p bar -i xxxyyy -d "tibcobanqio.api.mashery.com" -H petstore.swagger.io
```
##### Consul
Publishing to Consul will register and de-register services with the Consul server. This command will take the HTTP triggers defined in your mashling.json configuration file and push them to the Consul server specified by the command line arguments.

Details are:
```
./mashling-cli publish consul -h
```
```
Publishes the details of the mashling.json configuration file Consul

Usage:
  mashling-cli publish consul [flags]

Flags:
  -d, --consulDeRegister      de-register services with consul (required -d & -r mutually exclusive)
  -D, --consulDefDir string   service definition folder
  -r, --consulRegister        register services with consul (required -d & -r mutually exclusive) (default true)
  -t, --consulToken string    consul agent security token
  -h, --help                  help for consul
  -H, --host string           the hostname where consul is running (e.g. hostip:port)

Global Flags:
  -c, --config string         mashling gateway configuration (default "mashling.json")
  -e, --env-var-name string   name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -l, --load-from-env         load the mashling gateway configuration from an environment variable
```
Example registering a service with Consul:
```
./mashling-cli publish consul -r -c mashling-gateway-consul.json -t abcd1234 -H 192.45.32.31:8500
```
Example registering a service with Consul using the service definition folder:
```
./mashling-cli publish consul -r -c mashling-gateway-consul.json -t abcd1234 -D /etc/consul/configfiles/
```
