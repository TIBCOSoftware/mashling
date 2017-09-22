Inline gateway recipe

Prerequisites:
Run the kafka producer with topic "orders".

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f inline-gateway.json inline-gateway

2)Go to the path: cd inline-gateway/bin

3)Run the app using command: ./inline-gateway

4)Enter any message in "orders" producer, the message gets logged in the gateway.


