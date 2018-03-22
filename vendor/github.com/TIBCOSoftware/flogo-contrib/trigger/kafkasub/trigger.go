package kafkasub

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("trigger-flogo-kafkasub")

type _topichandler struct {
	topic      string
	offset     int64
	group      string
	partitions []int32
}

type _kafkaParms struct {
	brokers  []string
	handlers []_topichandler
}

// MqttTrigger is simple MQTT trigger
type KafkaSubTrigger struct {
	metadata           *trigger.Metadata
	config             *trigger.Config
	handlers           []*trigger.Handler
	kafkaParms         _kafkaParms
	shutdownChan       *chan struct{}
	signals            *chan os.Signal
	kafkaConfig        *sarama.Config
	kafkaConsumer      *sarama.Consumer
	partitionConsumers *map[string]sarama.PartitionConsumer
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &KafkasubFactory{metadata: md}
}

// KafkasubFactory Kafka subscriber Trigger factory
type KafkasubFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *KafkasubFactory) New(config *trigger.Config) trigger.Trigger {
	kafkaTrigger := &KafkaSubTrigger{metadata: t.metadata, config: config}
	kafkaTrigger.kafkaConfig = sarama.NewConfig()
	kafkaTrigger.kafkaConfig.Consumer.Return.Errors = true
	//err := initKafkaParams(kafkaTrigger)
	//if err != nil {
	//	panic(fmt.Errorf("Failed to initialize parms: [%s]", err))
	//}

	return kafkaTrigger
}

// Metadata implements trigger.Trigger.Metadata
func (t *KafkaSubTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Initialize implements trigger.Init
func (t *KafkaSubTrigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()

	err := initKafkaParams(t)
	return err
}

// Start implements ext.Trigger.Start
func (t *KafkaSubTrigger) Start() error {
	shutdownChan := make(chan struct{})
	t.shutdownChan = &shutdownChan
	signals := make(chan os.Signal, 1)
	t.signals = &signals
	signal.Notify(*t.signals, os.Interrupt)
	err := run(t)
	//log.Debug("KafkaSubTrigger Started")
	return err
}

// Stop implements ext.Trigger.Stop
func (t *KafkaSubTrigger) Stop() error {
	//unsubscribe from topic
	if t.partitionConsumers == nil {
		log.Debug("Closed called for a subscriber with no running consumers")
		log.Debug("Stopped")
		return nil
	}
	for id, partitionConsumer := range *t.partitionConsumers {
		partitionConsumer.Close()
		log.Debug("Closed partition consumer:", id)
	}
	if t.kafkaConsumer != nil {
		(*t.kafkaConsumer).Close()
		log.Debug("Closed kafka consumer")
	}
	log.Debug("Stopped")
	return nil
}

func run(t *KafkaSubTrigger) error {
	kafkaConsumer, err := sarama.NewConsumer(t.kafkaParms.brokers, t.kafkaConfig)
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer for reason [%s]", err)
	}
	t.kafkaConsumer = &kafkaConsumer
	consumers := make(map[string]sarama.PartitionConsumer)
	t.partitionConsumers = &consumers
	for id, handler := range t.kafkaParms.handlers {
		validPartitions, err := kafkaConsumer.Partitions(handler.topic)
		if err != nil {
			return fmt.Errorf("failed to get valid partitions for topic [%s] for reason [%s].  Aborting subscriber",
				handler.topic, err)
		}
		log.Debugf("Subscribing to topic [%s]", handler.topic)

		log.Debugf("Valid partitions for topic [%s] detected as: [%v]", handler.topic, validPartitions)
		if handler.partitions == nil { //subscribe to all valid partitions
			handler.partitions = validPartitions
			for _, part := range validPartitions {
				startConsumer(t, part, id)
			}
		} else { //subscribe to a subset of valid partitions
			configPartitions := handler.partitions
			for _, confPart := range configPartitions {
				for _, valPart := range validPartitions {
					if confPart == valPart {
						startConsumer(t, confPart, id)
						break
					}
					log.Errorf("Configured partition [%d] on topic [%s] does not exist and will not be subscribed", confPart, handler.topic)
				}
			}
		}
		if len(*t.partitionConsumers) < 1 {
			return fmt.Errorf("Kafka consumer is not configured for any valid partitions")
		}
		log.Debugf("Kafka consumers for topic [%s] started", handler.topic)
	}
	return nil
}

func startConsumer(t *KafkaSubTrigger, part int32, id int) error {
	log.Debugf("Creating PartitionConsumer for valid partition: [%s:%d]", t.kafkaParms.handlers[id].topic, part)
	consumer := *t.kafkaConsumer
	partitionConsumer, err := consumer.ConsumePartition(t.kafkaParms.handlers[id].topic, part, t.kafkaParms.handlers[id].offset)
	if err != nil {
		log.Errorf("Creating PartitionConsumer for valid partition: [%s:%d] failed for reason: %s", t.kafkaParms.handlers[id].topic, part, err)
		return err
	}
	consumerName := fmt.Sprintf("%d_%d", part, id)
	(*t.partitionConsumers)[consumerName] = partitionConsumer
	go consumePartition(t, partitionConsumer, part)
	return nil
}

func consumePartition(t *KafkaSubTrigger, consumer sarama.PartitionConsumer, part int32) {
	for {
		select {
		case err := <-consumer.Errors():
			if err == nil {
				//been shutdown
				return
			}
			log.Warnf("PartitionConsumer [%d] got error: [%s]", part, err)
			time.Sleep(time.Millisecond * 300)
		case msg := <-consumer.Messages():
			onMessage(t, msg)
		case <-*t.signals:
			log.Infof("Partition consumer got SIGINT; exiting")
			*t.shutdownChan <- struct{}{}
			return
		case <-*t.shutdownChan:
			log.Infof("Partition consumer got SHUTDOWN signal; exiting")
			return
		}
	}
}

func getTopics(t *KafkaSubTrigger) []string {
	return strings.Split(t.config.Settings["Topics"].(string), ",")
}

func initKafkaParams(t *KafkaSubTrigger) error {
	brokersString := t.config.Settings["BrokerUrl"]
	if brokersString == nil || brokersString.(string) == "" {
		return fmt.Errorf("BrokerUrl not provided")
	}
	brokers := strings.Split(brokersString.(string), ",")
	if len(brokers) < 1 {
		return fmt.Errorf("BrokerUrl [%s] is invalid, require at least one broker", brokersString.(string))
	}
	t.kafkaParms.brokers = make([]string, len(brokers))
	for brokerNo, broker := range brokers {
		err := validateBrokerUrl(broker)
		if err != nil {
			return fmt.Errorf("BrokerUrl [%s] format invalid for reason: [%s]", broker, err.Error())
		}
		t.kafkaParms.brokers[brokerNo] = broker
	}
	//clientKeystore
	/*
		Its worth mentioning here that when the keystore for kafka is created it must support RSA keys via
		the -keyalg RSA option.  If not then there will be ZERO overlap in supported cipher suites with java.
		see:   https://issues.apache.org/jira/browse/KAFKA-3647
		for more info
	*/
	if trustStore := t.config.Settings["truststore"]; trustStore != nil &&
		len(trustStore.(string)) > 0 {
		trustPool, err := getCerts(trustStore.(string))
		if err != nil {
			return err
		}
		config := tls.Config{
			RootCAs:            trustPool,
			InsecureSkipVerify: true}
		t.kafkaConfig.Net.TLS.Enable = true
		t.kafkaConfig.Net.TLS.Config = &config
	}
	// SASL
	if t.config.Settings["user"] != nil {
		var password string
		user := t.config.Settings["user"].(string)
		if len(user) > 0 {
			if t.config.Settings["password"] == nil ||
				len(t.config.Settings["password"].(string)) < 1 {
				return fmt.Errorf("password not provided for user: %s", user)
			}
			password = t.config.Settings["password"].(string)
			t.kafkaConfig.Net.SASL.Enable = true
			t.kafkaConfig.Net.SASL.User = user
			t.kafkaConfig.Net.SASL.Password = password
		}
	}

	// _topichandlers section
	if len(t.handlers) == 0 {
		return fmt.Errorf("Kafka trigger requires at least one handler containing a valid topic name")
	}

	t.kafkaParms.handlers = make([]_topichandler, len(t.handlers))

	for handlerNum, handler := range t.handlers {
		if handler.GetStringSetting("Topic") == "" {
			return fmt.Errorf("topic string was not provided for handler: [%s]", handler)
		}
		t.kafkaParms.handlers[handlerNum].topic = handler.GetStringSetting("Topic")

		//offset
		if handler.GetStringSetting("offset") != "" {
			i, err := strconv.Atoi(handler.GetStringSetting("offset"))
			if err != nil {
				log.Warnf("Offset [%s] specified for handler [%s] is not a valid number, using latest for offset",
					handler.GetStringSetting("offset"), handler)
				t.kafkaParms.handlers[handlerNum].offset = sarama.OffsetNewest
			} else {
				t.kafkaParms.handlers[handlerNum].offset = int64(i)
			}
		} else {
			t.kafkaParms.handlers[handlerNum].offset = sarama.OffsetNewest
		}

		//partitions
		if handler.GetStringSetting("partitions") != "" {
			partitions := handler.GetStringSetting("partitions")
			i := 0
			parts := strings.Split(partitions, ",")
			t.kafkaParms.handlers[handlerNum].partitions = make([]int32, len(parts))
			for _, p := range parts {
				n, err := strconv.Atoi(p)
				if err == nil {
					t.kafkaParms.handlers[handlerNum].partitions[i] = int32(n)
					i++
				} else {
					log.Warnf("Partition [%s] specified for handler [%s] is not a valid number and is discarded",
						p, handler)
				}
			}
		}

		//group
		if handler.GetStringSetting("group") != "" {
			t.kafkaParms.handlers[handlerNum].group = handler.GetStringSetting("group")
		}
	}
	return nil
}

func getCerts(trustStore string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	fileInfo, err := os.Stat(trustStore)
	if err != nil {
		return certPool, fmt.Errorf("Truststore [%s] does not exist", trustStore)
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		break
	case mode.IsRegular():
		return certPool, fmt.Errorf("Truststore [%s] is not a directory.  Must be a directory containing trusted certificates in PEM format",
			trustStore)
	}
	trustedCertFiles, err := ioutil.ReadDir(trustStore)
	if err != nil || len(trustedCertFiles) == 0 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  Must be a directory containing trusted certificates in PEM format", trustStore)
	}
	for _, trustCertFile := range trustedCertFiles {
		fqfName := fmt.Sprintf("%s%c%s", trustStore, os.PathSeparator, trustCertFile.Name())
		trustCertBytes, err := ioutil.ReadFile(fqfName)
		if err != nil {
			log.Warnf("Failed to read trusted certificate [%s] ... continueing", trustCertFile.Name())
		}
		certPool.AppendCertsFromPEM(trustCertBytes)
	}
	if len(certPool.Subjects()) < 1 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  After processing all files in the directory no valid trusted certs were found", trustStore)
	}
	return certPool, nil
}

//Ensure that this string meets the host:port definition of a kafka hostspec
//Kafka calls it a url but its really just host:port, which for numeric ip addresses is not a valid URI
//technically speaking.
func validateBrokerUrl(broker string) error {
	hostport := strings.Split(broker, ":")
	if len(hostport) != 2 {
		return fmt.Errorf("BrokerUrl must be composed of sections like \"host:port\"")
	}
	i, err := strconv.Atoi(hostport[1])
	if err != nil || i < 0 || i > 32767 {
		return fmt.Errorf("Port specification [%s] is not numeric and between 0 and 32767", hostport[1])
	}
	return nil
}

func onMessage(t *KafkaSubTrigger, msg *sarama.ConsumerMessage) {
	if msg == nil {
		return
	}

	log.Debugf("Kafka subscriber triggering job from topic [%s] on partition [%d] with key [%s] at offset [%d]",
		msg.Topic, msg.Partition, msg.Key, msg.Offset)

	for _, handler := range t.handlers {

		//actionID := action.Get(handler.ActionId)
		//log.Debugf("Found action: '%+x' for ActionID: %s", actionID, handler.ActionId)
		if t.metadata == nil {
			log.Infof("Kafka subscriber running in test mode received message on topic [%s] on partition [%d] with key [%s] at offset [%d]",
				msg.Topic, msg.Partition, msg.Key, msg.Offset)
			log.Infof("Content: [%s]", string(msg.Value))

			return
		}

		data := make(map[string]interface{})
		data["message"] = string(msg.Value)

		//if(t.metadata.Metadata.OutPuts
		startAttrs, errorAttrs := t.metadata.OutputsToAttrs(data, true)
		if errorAttrs != nil || startAttrs == nil {
			log.Errorf("Failed to create output attributes for kafka message for handler [%s] for reason [%s] message lost", handler, errorAttrs)
		}

		_, err := handler.Handle(context.Background(), data)

		if err != nil {
			log.Errorf("Run action for handler [%s] failed for reason [%s] message lost", handler, err)
		}
	}

}
