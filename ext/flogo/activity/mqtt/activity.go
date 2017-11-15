package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	opentracing "github.com/opentracing/opentracing-go"
	ctx "golang.org/x/net/context"
)

const (
	ivContent  = "content"
	ivTopic    = "topic"
	ivBroker   = "broker"
	ivID       = "id"
	ivUser     = "user"
	ivPassword = "password"
	ivQOS      = "qos"
	ivTracing  = "tracing"

	ovTracing = "tracing"
)

var (
	errorTopicIsNotAString            = errors.New("topic is not a string")
	errorInvalidEmptyTopicToPublishTo = errors.New("Invalid empty topic to publish to")
	errorQOSIsNotANumber              = errors.New("qos is not a number")
)

var log = logger.GetLogger("activity-tibco-rest")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	var span opentracing.Span
	if tracing := context.GetInput(ivTracing); tracing != nil {
		span = opentracing.SpanFromContext(tracing.(ctx.Context))
	}

	if span != nil {
		span = opentracing.StartSpan(
			context.TaskName(),
			opentracing.ChildOf(span.Context()))
		context.SetOutput(ovTracing, opentracing.ContextWithSpan(ctx.Background(), span))
		defer span.Finish()
	}

	setTag := func(key string, value interface{}) {
		if span != nil {
			span.SetTag(key, value)
		}
	}

	logError := func(format string, a ...interface{}) {
		str := fmt.Sprintf(format, a...)
		setTag("error", str)
		log.Error(str)
	}

	opts := mqtt.NewClientOptions()
	if broker, ok := context.GetInput(ivBroker).(string); ok {
		setTag("broker", broker)
		opts.AddBroker(broker)
	}
	if id, ok := context.GetInput(ivID).(string); ok {
		opts.SetClientID(id)
	}
	if user, ok := context.GetInput(ivUser).(string); ok {
		opts.SetUsername(user)
	}
	if password, ok := context.GetInput(ivPassword).(string); ok {
		opts.SetPassword(password)
	}
	opts.SetCleanSession(true)
	opts.SetStore(mqtt.NewMemoryStore())

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logError("Connection to mqtt broker failed %v", token.Error())
	}
	defer client.Disconnect(250)

	content := ""
	switch v := context.GetInput(ivContent).(type) {
	case int, int64, float64, bool, json.Number:
		content = fmt.Sprintf("%v", v)
	case string:
		content = v
	default:
		var data []byte
		data, err = json.Marshal(v)
		if err != nil {
			logError("Invalid content %v", err)
			break
		}
		content = string(data)
	}
	setTag("content", content)

	topic, ok := context.GetInput(ivTopic).(string)
	if !ok {
		logError(errorTopicIsNotAString.Error())
		return false, errorTopicIsNotAString
	}
	if len(topic) == 0 {
		logError(errorInvalidEmptyTopicToPublishTo.Error())
		return false, errorInvalidEmptyTopicToPublishTo
	}
	setTag("topic", topic)

	qos, ok := context.GetInput(ivQOS).(float64)
	if !ok {
		logError(errorQOSIsNotANumber.Error())
		return false, errorQOSIsNotANumber
	}

	token := client.Publish(topic, byte(qos), false, content)
	sent := token.WaitTimeout(5000 * time.Millisecond)
	if !sent {
		logError("Timeout occurred while trying to publish to topic '%s'", topic)
		return false, err
	}

	return true, nil
}
