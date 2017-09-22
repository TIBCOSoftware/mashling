Inline gateway recipe

Prerequisites:
Run the kafka producer with topic "users".

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f kafka-conditional-gateway.json kafka-conditional-gateway

2)Go to the path: cd kafka-conditional-gateway/bin

3)Run the app using command: ./kafka-conditional-gateway

4)Enter {"country":"USA"} message in "users" producer, 'usa_users_topic_handler' gets picked and the message gets logged in the gateway.

5)Enter {"country":"IND"} message in "users" producer, 'asia_users_topic_handler' gets picked and the message gets logged in the gateway.
