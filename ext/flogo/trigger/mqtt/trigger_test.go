package mqtt

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	opentracing "github.com/opentracing/opentracing-go"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "name": "tibco-mqtt",
  "settings": {
    "broker": "tcp://127.0.0.1:1883",
    "id": "flogoEngine",
    "user": "",
    "password": "",
    "store": "",
    "qos": "0",
    "cleansess": "false"
  },
  "handlers": [
    {
      "actionId": "device_info",
      "settings": {
        "topic": "test_start"
      }
    }
  ]
}`

var _ action.Runner = &TestRunner{}

type TestRunner struct {
	t *testing.T
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string,
	options interface{}) (code int, data interface{}, err error) {
	tr.t.Logf("Ran Action: %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(context context.Context, actionID string, inputGenerator action.InputGenerator,
	options map[string]interface{}) (results map[string]interface{}, err error) {
	tr.t.Logf("Ran Action: %v", actionID)
	return nil, nil
}

func TestInit(t *testing.T) {

	// New  factory
	md := trigger.NewMetadata(jsonMetadata)
	f := NewFactory(md)

	// New Trigger
	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	if err != nil {
		t.Error(err)
	}
	tgr := f.New(config)

	runner := &TestRunner{t: t}

	tgr.Init(runner)
}

func TestEndpoint(t *testing.T) {

	_, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Log("MQTT message broker is not available, skipping test...")
		return
	}

	// New  factory
	md := trigger.NewMetadata(jsonMetadata)
	f := NewFactory(md)

	// New Trigger
	config := &trigger.Config{}
	err = json.Unmarshal([]byte(testConfig), config)
	if err != nil {
		t.Error(err)
	}
	tgr := f.New(config)

	runner := &TestRunner{t: t}

	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("flogo_test")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	t.Log("---- doing first publish ----")

	token := client.Publish("test_start", 0, false, `{"message": "Test message payload!"}`)
	token.Wait()

	duration2 := time.Duration(2) * time.Second
	time.Sleep(duration2)

	t.Log("---- doing second publish ----")

	token = client.Publish("test_start", 0, false, `{"message": "Test message payload!"}`)
	token.Wait()

	duration5 := time.Duration(5) * time.Second
	time.Sleep(duration5)

	client.Disconnect(250)
	t.Log("Sample Publisher Disconnected")
}

func TestHandler(t *testing.T) {
	tracer := &opentracing.NoopTracer{}
	span := Span{
		Span: tracer.StartSpan("topic"),
	}
	defer span.Finish()

	handler := &OptimizedHandler{
		defaultActionId: "action_1",
		dispatches: []*Dispatch{
			{
				actionId:  "action_2",
				condition: "${trigger.content.value == A}",
			},
			{
				actionId:  "action_3",
				condition: "${env.value == A}",
			},
		},
	}

	action := handler.GetActionID("null", span)
	if action != "action_1" {
		t.Error("expected action_1")
	}

	content := `{"value": "A"}`
	action = handler.GetActionID(content, span)
	if action != "action_2" {
		t.Error("expected action_2")
	}

	os.Setenv("value", "A")
	action = handler.GetActionID("null", span)
	if action != "action_3" {
		t.Error("expected action_3")
	}
}

const testCreateHandlersConfig string = `{
  "name": "tibco-mqtt",
  "settings": {
    "broker": "tcp://127.0.0.1:1883",
    "id": "flogoEngine",
    "user": "",
    "password": "",
    "store": "",
    "qos": "0",
    "cleansess": "false"
  },
  "handlers": [
    {
      "actionId": "action_1",
      "settings": {
        "topic": "topic_1"
      }
    },
		{
      "actionId": "action_2",
      "settings": {
        "topic": "topic_1",
				"Condition": "${trigger.content.value == A}"
      }
    },
		{
      "actionId": "action_3",
      "settings": {
        "topic": "topic_2"
      }
    }
  ]
}`

func TestCreateHandlers(t *testing.T) {
	// New  factory
	md := trigger.NewMetadata(jsonMetadata)
	f := NewFactory(md)

	// New Trigger
	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testCreateHandlersConfig), config)
	if err != nil {
		t.Error(err)
	}
	tgr := f.New(config).(*MqttTrigger)

	handlers := tgr.CreateHandlers()
	if handlers["topic_1"].defaultActionId != "action_1" {
		t.Error("default action for topic_1 should be action_1")
	}
	if len(handlers["topic_1"].dispatches) != 1 {
		t.Error("there should be 1 dispatches for topic_1")
	}
	if handlers["topic_2"].defaultActionId != "action_3" {
		t.Error("default action for topic_2 should be action_3")
	}
	if len(handlers["topic_2"].dispatches) != 0 {
		t.Error("there should be 0 dispatches for topic_2")
	}
}
