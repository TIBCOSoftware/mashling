# Mashling Pingfunctionality
Ping functionality is to know whether gateway is alive and healthy over the network.

Ping port can be set in two ways i.e by setting environment variable "MASHLING_PING_PORT" OR by using mashling-gateway command flag - "-p". CLI flag gets more priority than environment variable. If user doesn't provide ping port, default value 9090 will be used.

Testing:

Run below command:
	./mashling-gateway -c <path-to-mashling-config> -p <PING-PORT>

Use below command to check gateway service:

	curl http://<GATEWAY IP>:<PING-PORT>/ping


Expected Result: 

	{"response":"Ping successful"}

Use below command to check gateway service with additional details:

	curl  http://<GATEWAY IP>:<PING-PORT>/ping/details


Expected Result: 

	{"Version":"0.2","Appversion":"1.0.0","Appdescription":"This is the first microgateway app"}