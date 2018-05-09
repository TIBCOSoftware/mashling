# Project Mashling
[![Build Status](https://travis-ci.org/TIBCOSoftware/mashling.svg?branch=master)](https://travis-ci.org/TIBCOSoftware/mashling)

Project Mashling<sup>TM</sup> is an open source event-driven microgateway.

Project Mashling highlights include:
* Ultra lightweight: 10-50x times less compute resource intensive
* Event-driven by design
* Complements Service Meshes
* Co-exists with API management platforms in a federated API Gateway model

Project Mashling consists of the following components:

* [mashling](http://github.com/TIBCOSoftware/mashling): This is the main repo that includes the below components:
  - A mashling-cli to build customized Mashling apps
  - A mashling-gateway to run supported features out of the box
  - Mashling triggers and activities
  - Library to build Mashling extensions

* [mashling-recipes](http://github.com/TIBCOSoftware/mashling-recipes): This is the repo that includes recipes that illustrate configuration of common microgateway patterns. These recipes are curated and searchable via [mashling.io](http://mashling.io).

* [mashling.io](http://mashling.io): Project Mashling also comes with a searchable collection of curated recipes for common microgateway patterns. To get started:

  - Find a recipe that is of interest to you
  - Browse the details on the recipe
  - Use the "Try it now" button to download the corresponding pre-created Mashling application
  - Each recipe comes with detailed usage instruction. A recipe README file, as an example [here](https://github.com/TIBCOSoftware/mashling-recipes/blob/master/recipes/event-dispatcher-router-mashling/README.md)

Additional developer tooling is included in below open source repo that contains the VSCode plugin for Mashling configuration:
* [VSCode Plugin for Mashling](https://github.com/TIBCOSoftware/vscode-extension-mashling)

## Installation and Usage

Starting with the v0.4.0 release both the `mashling-cli` and `mashling-gateway` binaries can be used by downloading them from the release page. Be sure to select the appropriate binary for your operating system.

For now you can download pre-built binaries from these links:

#### Linux
  - [mashling-gateway](https://s3.amazonaws.com/downloads.mashling.io/mashling-gateway-linux-amd64)
  - [mashling-cli](https://s3.amazonaws.com/downloads.mashling.io/mashling-cli-linux-amd64)

#### macOS
  - [mashling-gateway](https://s3.amazonaws.com/downloads.mashling.io/mashling-gateway-darwin-amd64)
  - [mashling-cli](https://s3.amazonaws.com/downloads.mashling.io/mashling-cli-darwin-amd64)

#### Windows
  - [mashling-gateway](https://s3.amazonaws.com/downloads.mashling.io/mashling-gateway-windows-amd64)
  - [mashling-cli](https://s3.amazonaws.com/downloads.mashling.io/mashling-cli-windows-amd64)

### mashling-gateway
The `mashling-gateway` is a static runtime for Mashling instances that provides the ability to load the mashling v2.0 modeling while also being backwards compatible for the original schema. It contains all the necessary dependencies to run existing recipes.

The `mashling-gateway` binary is what does the actual processing of requests and events according to the rules you've outlined in your `mashling.json` file.

Detailed usage information, documentation, and examples can be found in the [mashling-gateway documentation](docs/gateway/README.md).

A simple usage example is:

```
./mashling-gateway -c <path-to-mashling-config>
```

The default value of the `config` argument is `mashling.json`.

Any of the bundled configurations in the `examples/recipes/` folder will work. Realistically any of the recipes available on [mashling.io](https://www.mashling.io) should also work with the compiled binary created from this project.

The intent of this binary is to be used with *any* of these configuration files without re-compiling the source of this project.

Again, in depth configuration and usage documentation for the gateway can be found [here](docs/gateway/README.md).

### mashling-cli
The `mashling-cli` binary is used to create customized `mashling-gateway` binaries that contain triggers, actions, and activities not included in the default `mashling-gateway`. Much like the default binary, once your customized binary is built it can be reused with any `mashling.json` configuration file that has its dependencies satisfied by this new customized binary.

#### <a name="prerequisites"></a>Prerequisites

Because the `mashling-cli` is building custom binaries there are a few extra dependencies that need to be installed for it to work. You have two options:

1 - Use Go natively, which requires the following:
* The Go programming language 1.10 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system.

2 - Install [Docker](https://www.docker.com).

If Docker is installed locally the `mashling-cli` binary will use the local Docker install to run all the build commands through a pre-built Docker image.

If Docker is not installed but Go is then the CLI will attempt to use your native Go installation.

Detailed usage information, documentation, and examples can be found in the [mashling-cli documentation](docs/cli/README.md).

A simple usage example is:

```
./mashling-cli create -c <path-to-mashling-config-with-custom-dependencies>
```

By default this will use `Docker` on your local machine, if detected, to perform all of the custom asset identification, packaging, and compilation.

A simple custom Flogo trigger example that works with the above command can be found [here](examples/recipes/v2/customized-simple-synchronous-pattern.json).

Again, in depth configuration and usage documentation for the CLI can be found [here](docs/cli/README.md).

## Development and Building from Source
With the new v2 model this section only applies if you are actively making changes to the Mashling source code. There is no need to build from source if you just want to get started using the `mashling-gateway` or `mashling-cli` tools. Those can be downloaded for your desired platform and run immediately.

### <a name="prerequisites"></a>Prerequisites
* The Go programming language 1.10 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system.

#### Getting Started

Start by pulling the repository down to your local machine using one of the approaches following approaches.

#### Using Go
You can pull the Mashling source code using default Go commands.

```bash
go get -u github.com/TIBCOSoftware/mashling/...
```

If Go is installed correctly this command will also build the mashling binary targets and install them into your `$GOPATH`. You should be able to run `mashling-gateway` and `mashling-cli` immediately if your `PATH` environment variable includes `$GOPATH/bin`.

#### Using Git
If you decide to download the source code from Github using Git, please make sure Git is installed.

```bash
git clone -b feature-v2-model --single-branch https://github.com/TIBCOSoftware/mashling.git $GOPATH/src/github.com/TIBCOSoftware/mashling
cd $GOPATH/src/github.com/TIBCOSoftware/mashling
go install ./...
```

If Go is installed correctly the `go install ./...` command will also build the mashling binary targets and install them into your `$GOPATH`. You should be able to run `mashling-gateway` and `mashling-cli` immediately if your `PATH` environment variable includes `$GOPATH/bin`.

#### Building

You can build the default target, assuming all dependencies are satisfied, with the default `go` commands, for instance:

```
go install ./...
```

This will compile and install the binaries into your `$GOPATH/bin`.

If you are making significant changes to the source code and have added new dependencies or file assets, please make sure to run the `setup` command from the root of your `mashling` directory:

```
go run build.go setup
```

You can then build your binaries using the automated build targets we have provided using:

```
go run build.go all
```

This will regenerate any `go` generated files, search for Flogo activities and triggers as binary assets, rebundle the CLI assets, format the generated code, vet the code, and then build the `mashling-gateway` and `mashling-cli` binaries. These will be available in your `$GOPATH/bin`.

You can build your binaries for release by doing the following:

```
go run build.go releaseall
```

This will build your binaries for all supported platforms and then compress them. The result of the `releaseall` process will be available in the `release/` folder of your `mashling` directory.

To build for a specific target platform, you can run the following command to build the gateway:

```
go run build.go releasegateway -os=windows -arch=amd64
```

Similarly, to build the CLI, run the following:

```
go run build.go releasecli -os=windows -arch=amd64
```

Supported platforms are:

- darwin/amd64
- linux/amd64
- linux/arm64
- windows/amd64

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.

Please submit a github issue if you would like to propose a significant change or request a new feature.

## License
Mashling is licensed under a BSD-type license. See license text [here](https://github.com/TIBCOSoftware/mashling/blob/master/TIBCO%20LICENSE.txt).

### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/mashling/issues)
