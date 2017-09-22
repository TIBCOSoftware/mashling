RestTrigger to KafkaPublisher recipe

Prerequisites:
Run the kafka consumer with topic "syslog".

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f RestTrigger-To-KafkaPublisher.json RestTrigger-To-KafkaPublisher

2)Go to the path: cd RestTrigger-To-KafkaPublisher/bin

3)Run the app using command: ./RestTrigger-To-KafkaPublisher

4)Use "PUT" operation and hit the url "http://localhost:9096/petEvent" with the below sample payload:
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

5)The payload is published to the "syslog" Topic.