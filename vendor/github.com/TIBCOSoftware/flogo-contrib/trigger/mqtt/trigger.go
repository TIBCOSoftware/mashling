package mqtt

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/support"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/eclipse/paho.mqtt.golang"
)

// log is the default package logger
var log = logger.GetLogger("trigger-tibco-mqtt")

// MqttTrigger is simple MQTT trigger
type MqttTrigger struct {
	metadata         *trigger.Metadata
	runner           action.Runner
	client           mqtt.Client
	config           *trigger.Config
	topicToActionURI map[string]string
	topicToHandlerCfg map[string]*trigger.HandlerConfig
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MQTTFactory{metadata: md}
}

// MQTTFactory MQTT Trigger factory
type MQTTFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *MQTTFactory) New(config *trigger.Config) trigger.Trigger {
	return &MqttTrigger{metadata: t.metadata, config: config}
}

// Metadata implements trigger.Trigger.Metadata
func (t *MqttTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *MqttTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *MqttTrigger) Start() error {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(t.config.GetSetting("broker"))
	opts.SetClientID(t.config.GetSetting("id"))
	opts.SetUsername(t.config.GetSetting("user"))
	opts.SetPassword(t.config.GetSetting("password"))
	b, err := data.CoerceToBoolean(t.config.Settings["cleansess"])
	if err != nil {
		log.Error("Error converting \"cleansess\" to a boolean ", err.Error())
		return err
	}
	opts.SetCleanSession(b)
	if storeType := t.config.Settings["store"]; storeType != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(t.config.GetSetting("store")))
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		//TODO we should handle other types, since mqtt message format are data-agnostic
		payload := string(msg.Payload())
		log.Debug("Received msg:", payload)
		handlerCfg, found := t.topicToHandlerCfg[topic]
		if found {
			t.RunAction(handlerCfg, payload)
		} else {
			log.Errorf("Topic %s not found", t.topicToActionURI[topic])
		}
	})

	client := mqtt.NewClient(opts)
	t.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	i, err := data.CoerceToNumber(t.config.Settings["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return err
	}

	t.topicToActionURI = make(map[string]string)
	t.topicToHandlerCfg = make(map[string]*trigger.HandlerConfig)

	for _, handlerCfg := range t.config.Handlers {

		topic := handlerCfg.GetSetting("topic")

		if token := t.client.Subscribe(topic, byte(i), nil); token.Wait() && token.Error() != nil {
			log.Errorf("Error subscribing to topic %s: %s", topic, token.Error())
			panic(token.Error())
		} else {
			log.Debugf("Suscribed to topic: %s, will trigger actionId: %s", topic, handlerCfg.ActionId)
			t.topicToActionURI[topic] = handlerCfg.ActionId
			t.topicToHandlerCfg[topic] = handlerCfg
		}
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *MqttTrigger) Stop() error {
	//unsubscribe from topic
	log.Debug("Unsubcribing from topic: ", t.config.Settings["topic"])
	for _, handlerCfg := range t.config.Handlers {
		if token := t.client.Unsubscribe(handlerCfg.GetSetting("topic")); token.Wait() && token.Error() != nil {
			log.Errorf("Error unsubscribing from topic %s: %s", handlerCfg.Settings["topic"], token.Error())
		}
	}

	t.client.Disconnect(250)

	return nil
}

// RunAction starts a new Process Instance
func (t *MqttTrigger) RunAction(handlerCfg *trigger.HandlerConfig, payload string) {

	req := t.constructStartRequest(payload)
	//err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	//if err != nil {
	//	//http.Error(w, err.Error(), http.StatusBadRequest)
	//	log.Error("Error Starting action ", err.Error())
	//	return
	//}
	req.ReplyTo = t.config.GetSetting("topic")

	//todo handle error
	startAttrs, _ := t.metadata.OutputsToAttrs(req.Data, false)

	act := action.Get(handlerCfg.ActionId)

	ctx := trigger.NewInitialContext(startAttrs, handlerCfg)
	results, err := t.runner.RunAction(ctx, act, nil)

	if err != nil {
		log.Error("Error starting action: ", err.Error())
	}
	log.Debugf("Ran action: [%s]", handlerCfg.ActionId)

	var replyData interface{}

	if len(results) != 0 {
		dataAttr, ok := results["data"]
		if ok {
			replyData = dataAttr.Value()
		}
	}

	if replyData != nil {
		data, err := json.Marshal(replyData)
		if err != nil {
			log.Error(err)
		} else {
			t.publishMessage(req.ReplyTo, string(data))
		}
	}
}

func (t *MqttTrigger) publishMessage(topic string, message string) {

	log.Debug("ReplyTo topic: ", topic)
	log.Debug("Publishing message: ", message)

	qos, err := strconv.Atoi(t.config.GetSetting("qos"))
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}
	if len(topic) == 0 {
		log.Warn("Invalid empty topic to publish to")
		return
	}
	token := t.client.Publish(topic, byte(qos), false, message)
	sent := token.WaitTimeout(5000 * time.Millisecond)
	if !sent {
		// Timeout occurred
		log.Errorf("Timeout occurred while trying to publish to topic '%s'", topic)
		return
	}
}

func (t *MqttTrigger) constructStartRequest(message string) *StartRequest {
	//TODO how to handle reply to, reply feature
	req := &StartRequest{}
	data := make(map[string]interface{})
	data["message"] = message
	req.Data = data
	return req
}

// StartRequest describes a request for starting a ProcessInstance
type StartRequest struct {
	ProcessURI  string                 `json:"flowUri"`
	Data        map[string]interface{} `json:"data"`
	Interceptor *support.Interceptor   `json:"interceptor"`
	Patch       *support.Patch         `json:"patch"`
	ReplyTo     string                 `json:"replyTo"`
}
