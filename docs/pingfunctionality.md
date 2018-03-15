# Mashling CLI Pingfunctionality
Ping functionality is to know whether gateway is alive and healthy over the network.

Ping functionality can be enabled by setting environment variable MASHLING_PING_ENABLE value to TRUE. By default this feature is disabled.

To create the base gateway with ping functionality, use the following command to enable ping functionality:

```bash
export MASHLING_PING_ENABLE=TRUE
```

Ping port can be set in two ways i.e by setting environment variable "MASHLING_PING_PORT" OR by using mashling CLI command flag - "pingport". CLI flag gets more priority than environment variable. If user doesn't provide ping port while creating / building the gateway, default value 9090 will be used.

Creation of app can be done as below:

```bash
mashling create -f mashling.json -pingport 9091 SampleGatewayApp
```
Or

```bash
export MASHLING_PING_PORT=9091 && mashling create -f mashling.json SampleGatewayApp
```

Building App:

```bash
mashling build -pingport 9091 SampleGatewayApp
```

Testing:

Gateway is build and up using the provided configuration.

Use below command to check gateway service:

	curl http://<GATEWAY IP>:<PING-PORT>/ping/


Expected Result: 

	{"response":"success"}

Use below command to check gateway service with additional details:

	curl  http://<GATEWAY IP>:<PING-PORT>/ping/details/


Expected Result: 

	{"AppDescrption":"This is the first microgateway app","AppVersion":"1.0.0","FlogolibRev":"XXX......","MashlingCliLocalRev":"","MashlingCliVersion":"0.3.0","MashlingRev":"XXX......","SchemaVersion":"0.2","mashlingCliRevision":"XXX......","response":"success"}