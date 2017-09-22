RestTrigger to RestInvoker with conditional dispatch recipe

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f rest-conditional-gateway.json rest-conditional-gateway

2)Go to the path: cd rest-conditional-gateway/bin

3)Run the app using command: ./rest-conditional-gateway

4)Use "PUT" operation and hit the url "http://localhost:9096/pets" with the below sample payload:
{
	"category": {
		"id": 16,
		"name": "Animals"
	},
	"id": 16,
	"name": "SPARROW",
	"photoUrls": [
		"string"
	],
	"status": "sold",
	"tags": [{
		"id": 0,
		"name": "string"
	}]
}

5)Use "GET" operation and hit the url "http://localhost:9096/pets/16" to check the above added pet details.