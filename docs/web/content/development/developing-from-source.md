---
title: Developing and Building from Source
weight: 4200

pre: "<i class=\"fa fa-code-fork\" aria-hidden=\"true\"></i> "
---

## Developing and Building from Source

This section only applies if you are actively making changes to the Mashling source code or trying to build customized mashling-gateway binaries using the mashling-cli tool without using Docker.

The easiest way to use the CLI or build from source locally is to have Docker installed locally and let that do the heavily lifting related to configuration and setup of all the development dependencies. Either way, first you need to get the source.

### Downloading the Mashling Source Code
#### Using Git
If you download the source code from Github using Git, please make sure Git is installed.

```
git clone -b feature-v2-model --single-branch https://github.com/TIBCOSoftware/mashling.git $GOPATH/src/github.com/TIBCOSoftware/mashling
cd $GOPATH/src/github.com/TIBCOSoftware/mashling
```

####Using Go
If you get the Mashling source using Go commands, please make sure Go 1.10 is installed.

```
go get -u github.com/TIBCOSoftware/mashling/...
```

If Go is installed correctly, this command will also build the mashling binary targets and install them into your $GOPATH.

####Using Github's Zip file downloads
Visit the [Mashling repository on Github](https://github.com/TIBCOSoftware/mashling) and click **Clone or download** and choose **Download ZIP**. This is a useful shortcut if you are just using Docker to run your build commands and do not plan on contributing back to the repository or tracking your changes.

####Using Docker
If you have docker installed, you can get started right away after downloading the Mashling source code.

From the root of the Mashling repository, run the following command to verify everything is working as expected:

```
docker run -v "$(PWD):/mashling" --rm -t mashling/mashling-compile /bin/bash -c "make help"
```

You should see the following output:

```
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

```
docker run -v "$(PWD):/mashling" --rm -t mashling/mashling-compile /bin/bash -c "make"
```

To regenerate all the assets and then rebuild, run the following:

```
docker run -v "$(PWD):/mashling" --rm -t mashling/mashling-compile /bin/bash -c "make all"
```

These commands will build both the ```mashling-gateway``` and ```mashling-cli``` and place them under the ```/bin``` folder in the root of the mashling repository.

To build a binary for a specific operating system, like Windows, use the **two** following commands.

**First**, run:

```
docker run -v "$(PWD):/mashling" --rm -t mashling/mashling-compile /bin/bash -c "make setup assets generate fmt"
```

**Then**, run:

```
docker run -e="GOOS=windows" -v "$(PWD):/mashling" --rm -t mashling/mashling-compile /bin/bash -c "make"
```

The supported GOOS values are windows, linux, and darwin.

You can also specify a target architecture. This can be done via:

```
docker run -e="GOOS=linux" -e="GOARCH=amd64" -v "$(PWD):/mashling" --rm -t mashling/mashling-compile /bin/bash -c "make"
```

The supported GOARCH values are amd64 and arm64. arm64 is only supported for linux.

On most systems, amd64 will be the default architecture.

The make help output shown above lists all the other commands that are available.

###Using a Native Toolchain
####Prerequisites
* The Go programming language 1.10 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system.
* Mashling uses *make* for asset generation, packing, and building of both the gateway and cli targets. The *make* tool comes with Windows by default, or it can be downloaded from [here](https://sourceforge.net/projects/gnuwin32/files/make/).
####Getting Started
Start by pulling the repository down to your local machine using one of the approaches described above. This project uses *Go* and *make* by default. All dependencies are either bundled into the repository under the *vendor* folder or pulled on demand by the appropriate make target.

####Building
You can build the default target, assuming all dependencies are satisfied with (the default target is *build*):

```
make
```
This will compile the contents of the *cmd/mashling-gateway/* and *cmd/mashling-cli/* directory and put the resulting binary into the *bin/* directory.

If you have added new dependencies or file assets, please make sure to run:

```
make all
```
This will check your *vendor/* directory for any new triggers, activities, or flows and make sure the appropriate Flogo factory registrations take place on startup of the compiled binary. This command will also pull any non *.go files in the *internal/app/gateway/assets/* directory into the compiled binary.