---
title: Developing and Building from Source
weight: 4200

pre: "<i class=\"fa fa-code-fork\" aria-hidden=\"true\"></i> "
---

## Development and Building from Source
This section only applies if you are actively making changes to the Mashling source code. There is no need to build from source if you just want to get started using the **mashling-gateway** or **mashling-cli** tools. Those can be downloaded for your desired platform and run immediately.

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

If Go is installed correctly this command will also build the mashling binary targets and install them into your **$GOPATH**. You should be able to run **mashling-gateway** and **mashling-cli** immediately if your **PATH** environment variable includes **$GOPATH/bin**.

#### Using Git
If you decide to download the source code from Github using Git, please make sure Git is installed.

```bash
git clone -b feature-v2-model --single-branch https://github.com/TIBCOSoftware/mashling.git $GOPATH/src/github.com/TIBCOSoftware/mashling
cd $GOPATH/src/github.com/TIBCOSoftware/mashling
go install ./...
```

If Go is installed correctly the **go install ./...** command will also build the mashling binary targets and install them into your **$GOPATH**. You should be able to run **mashling-gateway** and **mashling-cli** immediately if your **PATH** environment variable includes **$GOPATH/bin**.

#### Building

You can build the default target, assuming all dependencies are satisfied, with the default **go** commands, for instance:

```
go install ./...
```

This will compile and install the binaries into your **$GOPATH/bin**.

If you are making significant changes to the source code and have added new dependencies or file assets, please make sure to run the **setup** command from the root of your **mashling** directory:

```
go run build.go setup
```

You can then build your binaries using the automated build targets we have provided using:

```
go run build.go all
```

This will regenerate any **go** generated files, search for Flogo activities and triggers as binary assets, rebundle the CLI assets, format the generated code, vet the code, and then build the **mashling-gateway** and **mashling-cli** binaries. These will be available in your **$GOPATH/bin**.

You can build your binaries for release by doing the following:

```
go run build.go releaseall
```

This will build your binaries for all supported platforms and then compress them. The result of the **releaseall** process will be available in the **release/** folder of your **mashling** directory.

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