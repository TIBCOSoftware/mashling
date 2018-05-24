---
title: Gateway
weight: 3200
pre: "<i class=\"fa fa-cog\" aria-hidden=\"true\"></i> "
---

### Overview

The [mashling-gateway](https://github.com/TIBCOSoftware/mashling/tree/master/docs/gateway) powers the core event driven routing engine of the Mashling project. This core binary can run all versions of the Mashling schema to date, however for the purposes of this document, we will focus on the 1.0 configuration schema.

### Usage
The gateway binary has the following command line arguments available to setup and specify how you would like the binary to operate.

They can be found by running:

```bash
./mashling-gateway -h
```
The output and flags are:

```
A static binary that executes Mashling gateway logic defined in a mashling.json configuration file. Complete documentation is available at https://github.com/TIBCOSoftware/mashling

Version: v0.3.3-internal-29-gf6c81fd-dirty
Build Date: 2018-04-03T10:11:33-0400

Usage:
  mashling-gateway [flags]
  mashling-gateway [command]

Available Commands:
  help        Help about any command
  version     Prints the mashling-gateway version

Flags:
  -c, --config string          mashling gateway configuration (default "mashling.json")
  -C, --config-cache string    location of the configuration artifacts cache (default ".cache")
  -E, --config-cache-enabled   cache post-processed configuration artifacts locally (default true)
  -d, --dev                    run mashling in dev mode
  -e, --env-var-name string    name of the environment variable that contains the base64 encoded mashling gateway configuration (default "MASHLING_CONFIG")
  -h, --help                   help for mashling-gateway
  -l, --load-from-env          load the mashling gateway configuration from an environment variable
  -p, --ping-enabled           enable gateway ping service (default true)
  -P, --ping-port string       configure mashling gateway ping service port (default "9090")

Use "mashling-gateway [command] --help" for more information about a command.
```

Currently, **dev** mode just reloads the running gateway instance when a change is detected in the *mashling.json* file.

#### Health Check
An integrated ping service is used to determine if a gateway instance is up and running.

The health check ping service is enabled by default and configured to run on port *9090*. You can specify a different port at startup time using the command:

```
./mashling-gateway -c <path to mashling json> -P <ping port value>
```

You can also disable the ping service using the command:

```
./mashling-gateway -c <path to mashling json> -p=false
```

The health check endpoint is available at *http://<GATEWAY IP>:<PING-PORT>/ping* with an expected result of:

```
{"response":"Ping successful"}
```

A more detailed health check response is available at *http://<GATEWAY IP>:<PING-PORT>/ping/details* with an example result of:

```
{"Version":"0.2","Appversion":"1.0.0","Appdescription":"This is the first microgateway app"}
```
