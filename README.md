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

Because the `mashling-cli` is building custom binaries there are a few extra dependencies that need to be installed for it to work. The easiest is to just have [Docker](https://www.docker.com) installed on your local machine and let the `mashling-cli` binary use the local Docker install to run all the build commands through a pre-built Docker image. The other option is to [satisfy the prerequisite dependencies listed below in the development section](#prerequisites).

Detailed usage information, documentation, and examples can be found in the [mashling-cli documentation](docs/cli/README.md).

A simple usage example is:

```
./mashling-cli create -c <path-to-mashling-config-with-custom-dependencies>
```

By default this will use `Docker` on your local machine, if detected, to perform all of the custom asset identification, packaging, and compilation.

A simple custom Flogo trigger example that works with the above command can be found [here](examples/recipes/v2/customized-simple-synchronous-pattern.json).

Again, in depth configuration and usage documentation for the CLI can be found [here](docs/cli/README.md).

## Development and Building from Source
With the new v2 model this section only applies if you are actively making changes to the Mashling source code or trying to build customized `mashling-gateway` binaries using the `mashling-cli` tool **without** using Docker.

The easiest way to use the CLI or build from source locally is to have Docker installed locally and let that do the heavily lifting related to configuration and setup of all the development dependencies. Either way, first you need to get the source.

### Downloading the Mashling Source Code

#### Using Git
If you download the source code from Github using Git, please make sure Git is installed.

```bash
git clone -b feature-v2-model --single-branch https://github.com/TIBCOSoftware/mashling.git $GOPATH/src/github.com/TIBCOSoftware/mashling
cd $GOPATH/src/github.com/TIBCOSoftware/mashling
```

#### Using Go
If you get the Mashling source using Go commands, please make sure Go 1.10 is installed.

```bash
go get -u github.com/TIBCOSoftware/mashling/...
```

If go is installed correctly this command will also build the mashling binary targets and install them into your `$GOPATH`.

#### Using Github's Zip file downloads
Visit the [Mashling repository on Github](https://github.com/TIBCOSoftware/mashling) and click `Clone or download` and choose `Downlaod ZIP`. This is a useful shortcut if you are just using Docker to run your build commands and do not plan on contributing back to the repository or tracking your changes.

### Using Docker
If you have docker installed you can get started right away after downloading the Mashling source code.

From the root of the Mashling repository, run the following command to verify everything is working as expected:

```bash
docker run -v "$(PWD):/mashling" --rm -t jeffreybozek/mashling:compile /bin/bash -c "make help"
```

You should see the following output:

```bash
build           Build gateway and cli binaries
all             Create assets, run go generate, fmt, vet, and then build then gateway and cli binaries
buildgateway    Build gateway binary
allgateway      Satisfy pre-build requirements and then build the gateway binary
buildcli        Build CLI binary
allcli          Satisfy pre-build requirements and then build the cli binary
buildlegacy     Build legacy CLI binary
release         Build all executables for release against all targets
docker          Build a minimal docker image containing the gateway binary
setup           Setup the dev environment
hooks           Setup the git commit hooks
lint            Run golint
metalinter      Run gometalinter
fmt             Run gofmt on all source files
vet             Run go vet on all source files
generate        Run go generate on source
cligenerate     Run go generate on CLI source
assets          Run asset generation
cliassets       Run asset generation for CLI
list            List packages
dep             Make sure dependencies are vendored
depadd          Add new dependencies
clean           Cleanup everything
```

Now building the project from source is as easy as:

```bash
docker run -v "$(PWD):/mashling" --rm -t jeffreybozek/mashling:compile /bin/bash -c "make"
```

To regenerate all the assets and then rebuild, run the following:

```bash
docker run -v "$(PWD):/mashling" --rm -t jeffreybozek/mashling:compile /bin/bash -c "make all"
```
These commands will build both the `mashling-gateway` and `mashling-cli` and place them under the `/bin` folder in the root of the mashling repository.

**To build a binary for a specific operating system**, like Windows, use the **two** following commands.

**First**, run:

```bash
docker run -v "$(PWD):/mashling" --rm -t jeffreybozek/mashling:compile /bin/bash -c "make setup assets generate fmt"
```

**Then**, run:

```bash
docker run -e="GOOS=windows" -v "$(PWD):/mashling" --rm -t jeffreybozek/mashling:compile /bin/bash -c "make"
```

The supported `GOOS` values are `windows`, `linux`, and `darwin`.

The `make help` output shown above lists all the other commands that are available.

### Using a Native Toolchain

#### <a name="prerequisites"></a>Prerequisites
* The Go programming language 1.10 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system.
* Mashling uses `make` for asset generation, packing, and building of both the gateway and cli targets. The `make` tool come by default with Windows, it can be downloaded from [here](https://sourceforge.net/projects/gnuwin32/files/make/).

#### Getting Started

Start by pulling the repository down to your local machine using one of the approaches described above. This project uses Go and `make` by default. All dependencies are either bundled into the repository under the `vendor` folder or pulled on demand by the appropriate make target.

#### Building

You can build the default target, assuming all dependencies are satisfied with (the default target is `build`):

```
make
```

This will compile the contents of the `cmd/mashling-gateway/` and `cmd/mashling-cli/` directory and put the resulting binary into the `bin/` directory.

If you have added new dependencies or file assets, please make sure to run:

```
make all
```

This will check your `vendor/` directory for any new triggers, activities, or flows and make sure the appropriate Flogo factory registrations take place on startup of the compiled binary. This command will also pull any non `*.go` files in the `internal/app/gateway/assets/` directory into the compiled binary.

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
