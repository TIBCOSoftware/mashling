package kafkapub

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"sync"

	"github.com/Shopify/sarama"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var flogoLogger = logger.GetLogger("activity-tibco-kafkapub")

// MyActivity is a stub for your Activity implementation
type KafkaPubActivity struct {
	sync.Mutex
	metadata        *activity.Metadata
	syncProducerMap *map[string]sarama.SyncProducer
}

type KafkaParms struct {
	kafkaConfig  *sarama.Config
	brokers      []string
	topic        string
	syncProducer sarama.SyncProducer
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	flogoLogger.Debugf("Kafkapub NewActivity")
	pKafkPubActivity := &KafkaPubActivity{metadata: metadata}
	producers := make(map[string]sarama.SyncProducer)
	pKafkPubActivity.syncProducerMap = &producers
	return pKafkPubActivity
}

// Metadata implements activity.Activity.Metadata
func (a *KafkaPubActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *KafkaPubActivity) Eval(context activity.Context) (done bool, err error) {
	var parms (KafkaParms)
	flogoLogger.Debugf("Kafkapub Eval")
	err = initParms(a, context, &parms)
	if err != nil {
		flogoLogger.Errorf("Kafkapub parameters initialization got error: [%s]", err.Error())
		return false, err
	}
	if message := context.GetInput("Message"); message != nil && message.(string) != "" {
		msg := &sarama.ProducerMessage{
			Topic: parms.topic,
			Value: sarama.StringEncoder(message.(string)),
		}
		partition, offset, err := parms.syncProducer.SendMessage(msg)
		if err != nil {
			return false, fmt.Errorf("kafkapub failed to send message for reason [%s]", err.Error())
		}
		context.SetOutput("partition", partition)
		context.SetOutput("offset", offset)
		flogoLogger.Debugf("Kafkapub message [%v] sent successfully on partition [%d] and offset [%d]",
			message, partition, offset)
		return true, nil
	}
	return false, fmt.Errorf("kafkapub called without a message to publish")
}

func initParms(a *KafkaPubActivity, context activity.Context, params *KafkaParms) error {
	var producerkey (string)
	if context.GetInput("BrokerUrls") != nil && context.GetInput("BrokerUrls").(string) != "" {
		params.kafkaConfig = sarama.NewConfig()
		params.kafkaConfig.Producer.Return.Errors = true
		params.kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
		params.kafkaConfig.Producer.Retry.Max = 5
		params.kafkaConfig.Producer.Return.Successes = true
		brokerUrls := strings.Split(context.GetInput("BrokerUrls").(string), ",")
		brokers := make([]string, len(brokerUrls))
		for brokerNo, broker := range brokerUrls {
			error := validateBrokerUrl(broker)
			if error != nil {
				return fmt.Errorf("BrokerUrl [%s] format invalid for reason: [%s]", broker, error.Error())
			}
			brokers[brokerNo] = broker
			producerkey += broker
		}
		params.brokers = brokers
		flogoLogger.Debugf("Kafkapub brokers [%v]", brokers)
	} else {
		return fmt.Errorf("Kafkapub activity is not configured with at least one BrokerUrl")
	}
	if context.GetInput("Topic") != nil && context.GetInput("Topic").(string) != "" {
		params.topic = context.GetInput("Topic").(string)
		flogoLogger.Debugf("Kafkapub topic [%s]", params.topic)
	} else {
		return fmt.Errorf("Topic input parameter not provided and is required")
	}

	//clientKeystore
	/*
		Its worth mentioning here that when the keystore for kafka is created it must support RSA keys via
		the -keyalg RSA option.  If not then there will be ZERO overlap in supported cipher suites with java.
		see:   https://issues.apache.org/jira/browse/KAFKA-3647
		for more info
	*/
	if trustStore := context.GetInput("truststore"); trustStore != nil && trustStore.(string) != "" {
		if trustPool, err := getCerts(trustStore.(string)); err == nil {
			config := tls.Config{
				RootCAs:            trustPool,
				InsecureSkipVerify: true}
			params.kafkaConfig.Net.TLS.Enable = true
			params.kafkaConfig.Net.TLS.Config = &config

			flogoLogger.Debugf("Kafkapub initialized truststore from [%v]", trustStore)
		} else {
			return err
		}
		producerkey += trustStore.(string)
	}
	// SASL
	if user := context.GetInput("user"); user != nil && user.(string) != "" {
		var password (interface{})
		if password = context.GetInput("password"); password == nil {
			password = ""
		}
		params.kafkaConfig.Net.SASL.Enable = true
		params.kafkaConfig.Net.SASL.User = user.(string)
		params.kafkaConfig.Net.SASL.Password = password.(string)
		flogoLogger.Debugf("Kafkapub SASL parms initialized; user [%v]  password[########]", user)
		producerkey += user.(string)
	}
	a.Lock()
	defer func() {
		a.Unlock()
	}()

	if (*a.syncProducerMap)[producerkey] == nil {
		syncProducer, err := sarama.NewSyncProducer(params.brokers, params.kafkaConfig)
		if err != nil {
			return fmt.Errorf("Kafkapub failed to create a SyncProducer.  Check any TLS or SASL parameters carefully.  Reason given: [%s]", err)
		}
		params.syncProducer = syncProducer
		(*a.syncProducerMap)[producerkey] = syncProducer
		flogoLogger.Debugf("Kafkapub cacheing connection [%s]", producerkey)
	} else {
		params.syncProducer = (*a.syncProducerMap)[producerkey]
		flogoLogger.Debugf("Kafkapub reusing cached connection [%s]", producerkey)
	}
	flogoLogger.Debug("Kafkapub synchronous producer created")
	return nil
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
			flogoLogger.Warnf("Failed to read trusted certificate [%s] ... continuing", trustCertFile.Name())
		} else if trustCertBytes != nil {
			certPool.AppendCertsFromPEM(trustCertBytes)
		}
	}
	if len(certPool.Subjects()) < 1 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  After processing all files in the directory no valid trusted certs were found", trustStore)
	}
	return certPool, nil
}
