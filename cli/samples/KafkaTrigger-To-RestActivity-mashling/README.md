KafkaTrigger to RestInvoker recipe

Prerequisites:
Run the kafka producer with topic "syslog".

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f KafkaTrigger-To-RestInvoker.json KafkaTrigger-To-RestInvoker

2)Go to the path: cd KafkaTrigger-To-RestInvoker/bin

3)Run the app using command: ./KafkaTrigger-To-RestInvoker

4)Enter the below sample payload without any spaces in the "syslog" producer, the payload gets added to the swagger petstore.
{"category":{"id":10,"name":"string"},"id":10,"name":"doggie","photoUrls":["string"],"status":"available","tags":[{"id":0,"name":"string"}]}

5)Use "GET" operation and hit the swagger url "http://petstore.swagger.io/v2/pet/10" to check the above added pet details.