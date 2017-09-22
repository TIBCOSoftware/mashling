Kafka Reference recipe

Prerequisites:
Run the kafka producer with topic "users".

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f kafka-reference-gateway.json kafka-reference-gateway

2)Go to the path: cd kafka-reference-gateway/bin

3)Run the app using command: ./kafka-reference-gateway

4)Enter any message in "users" producer, 'user_topic_handler' gets picked and the message gets logged in the gateway.