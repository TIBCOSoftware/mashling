package logger

import (
	"crypto/tls"
	"errors"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/Sirupsen/logrus"
)

// KafkaLogger is a configuration for the Kafka logging hook.
type KafkaLogger struct {
	ID             string   `json:"id"`
	Levels         []string `json:"levels"`
	Brokers        []string `json:"brokers"`
	DefaultTopic   string   `json:"defaultTopic"`
	InjectHostname bool     `json:"injectHostname"`
}

// InitializeKafkaHook initializes a Kafka logger with the provided settings.
func InitializeKafkaHook(settings map[string]interface{}) (kafkaLogger *KafkaLogger, err error) {
	kafkaLogger = &KafkaLogger{}
	for k, v := range settings {
		switch k {
		case "id":
			id, ok := v.(string)
			if !ok {
				return nil, errors.New("invalid type for id")
			}
			kafkaLogger.ID = id
		case "levels":
			levels, ok := v.([]string)
			if !ok {
				return nil, errors.New("invalid type for levels")
			}
			kafkaLogger.Levels = levels
		case "brokers":
			brokers, ok := v.([]string)
			if !ok {
				return nil, errors.New("invalid type for brokers")
			}
			kafkaLogger.Brokers = brokers
		case "defaultTopic":
			defaultTopic, ok := v.(string)
			if !ok {
				return nil, errors.New("invalid type for defaultTopic")
			}
			kafkaLogger.DefaultTopic = defaultTopic
		case "injectHostname":
			injectHostname, ok := v.(bool)
			if !ok {
				return nil, errors.New("invalid type for injectHostname")
			}
			kafkaLogger.InjectHostname = injectHostname
		default:
			// ignore and move on.
		}
	}
	return kafkaLogger, err
}

func (k *KafkaLogger) Id() (string, error) {
	return k.ID, nil
}

func (k *KafkaLogger) GetHook() (logrus.Hook, error) {
	var levels []logrus.Level
	for _, level := range k.Levels {
		parsedLevel, err := logrus.ParseLevel(level)
		if err != nil {
			return nil, err
		}
		levels = append(levels, parsedLevel)
	}
	return NewKafkaLogrusHook(k.ID, levels, &logrus.JSONFormatter{}, k.Brokers, k.DefaultTopic, k.InjectHostname, nil)
}

// BELOW CODE IS FROM https://github.com/tracer0tong/kafkalogrus, which has an Apache 2.0 license
// TODO: replace with import statement once Flogo-lib import of Sirupsen/logrus case is corrected.

// KafkaLogrusHook is the primary struct
type KafkaLogrusHook struct {
	id             string
	defaultTopic   string
	injectHostname bool
	hostname       string
	levels         []logrus.Level
	formatter      logrus.Formatter
	producer       sarama.AsyncProducer
}

// NewKafkaLogrusHook creates a new KafkaHook
func NewKafkaLogrusHook(id string,
	levels []logrus.Level,
	formatter logrus.Formatter,
	brokers []string,
	defaultTopic string,
	injectHostname bool,
	tls *tls.Config) (*KafkaLogrusHook, error) {
	var err error
	var producer sarama.AsyncProducer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	kafkaConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	// check here if provided *tls.Config is not nil and assign to the sarama config
	// NOTE: we automatically enabled the TLS config because sarama would error out if our
	//       config were non-nil but disabled. To avoid issue father down the stack, we enable.
	if tls != nil {
		kafkaConfig.Net.TLS.Enable = true
		kafkaConfig.Net.TLS.Config = tls
	}

	if producer, err = sarama.NewAsyncProducer(brokers, kafkaConfig); err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send log entry to Kafka: %v\n", err)
		}
	}()

	var hostname string
	if hostname, err = os.Hostname(); err != nil {
		hostname = "localhost"
	}

	hook := &KafkaLogrusHook{
		id,
		defaultTopic,
		injectHostname,
		hostname,
		levels,
		formatter,
		producer,
	}

	return hook, nil
}

// Id returns the internal ID of the hook
func (hook *KafkaLogrusHook) Id() string {
	return hook.id
}

// Levels is required to implement the hook interface from logrus
func (hook *KafkaLogrusHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire is required to implement the hook interface from logrus
func (hook *KafkaLogrusHook) Fire(entry *logrus.Entry) error {
	var partitionKey sarama.ByteEncoder
	var b []byte
	var err error

	t, _ := entry.Data["time"].(time.Time)
	if b, err = t.MarshalBinary(); err != nil {
		return err
	}
	partitionKey = sarama.ByteEncoder(b)

	if hook.injectHostname {
		if _, ok := entry.Data["host"]; !ok {
			entry.Data["host"] = hook.hostname
		}
	}

	if b, err = hook.formatter.Format(entry); err != nil {
		return err
	}
	value := sarama.ByteEncoder(b)

	topic := hook.defaultTopic
	if tsRaw, ok := entry.Data["topic"]; ok {
		if ts, ok := tsRaw.(string); !ok {
			return errors.New("Incorrect topic filed type (should be string)")
		} else {
			topic = ts
		}
	}
	hook.producer.Input() <- &sarama.ProducerMessage{
		Key:   partitionKey,
		Topic: topic,
		Value: value,
	}
	return nil
}
