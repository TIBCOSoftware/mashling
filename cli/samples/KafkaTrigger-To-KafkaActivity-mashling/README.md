KafkaTrigger to KafkaPublisher recipe

Prerequisites:
Run the kafka producer with topic "publishpet" and the consumer with topic "subscribepet".

Instructions:

1)Place the json in folder and create the app using the below command:
mashling create -f KafkaTrigger-To-KafkaPublisher.json KafkaTrigger-To-KafkaPublisher

2)Go to the path: cd KafkaTrigger-To-KafkaPublisher/bin

3)Run the app using command: ./KafkaTrigger-To-KafkaPublisher

4)Enter any message in "publishpet" producer, the message gets published in "subscribepet" consumer.


