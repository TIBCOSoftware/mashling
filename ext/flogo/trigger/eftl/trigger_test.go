/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package eftl

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/mashling/commons/lib/eftl"
	"github.com/mashling/commons/lib/util"

	opentracing "github.com/opentracing/opentracing-go"
)

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
	"name": "tibco-eftl",
	"settings": {
	  "url": "ws://localhost:9191/channel",
	  "id": "eftl",
	  "user": "",
	  "password": ""
	},
	"handlers": [
	  {
		"actionId": "test",
		"settings": {
		  "dest": "test_start"
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

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	return nil, nil
}
func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	return nil, nil
}
func TestInit(t *testing.T) {
	jsonMetadata := getJsonMetadata()

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
	// TODO: Init no longer exists.
	// runner := &TestRunner{t: t}
	// tgr.Init(runner)
	tgr.Metadata()
}

func TestGetLocalIP(t *testing.T) {
	ip := getLocalIP()
	if ip == "" {
		t.Error("failed to get local ip")
	}
}

func TestConfigureTracer(t *testing.T) {
	jsonMetadata := getJsonMetadata()

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
	// TODO: Init no longer exists.
	// runner := &TestRunner{t: t}
	// tgr.Init(runner)

	eftl := tgr.(*Trigger)
	eftl.config.Settings["tracer"] = "zipkin"
	eftl.config.Settings["tracerEndpoint"] = "http://localhost:9411/api/v1/spans"
	eftl.config.Settings["tracerDebug"] = true
	eftl.config.Settings["tracerSameSpan"] = true
	eftl.config.Settings["tracerID128Bit"] = true
	err = eftl.tracer.ConfigureTracer(eftl.config.Settings, "localhost", "test")
	if err != nil {
		t.Fatal(err)
	}

	eftl.config.Settings["tracer"] = "appdash"
	eftl.config.Settings["tracerEndpoint"] = "localhost:7701"
	err = eftl.tracer.ConfigureTracer(eftl.config.Settings, "localhost", "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestEndpoint(t *testing.T) {
	_, err := net.Dial("tcp", "127.0.0.1:9191")
	if err != nil {
		t.Log("EFTL message broker is not available, skipping test...")
		return
	}

	jsonMetadata := getJsonMetadata()

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

	// TODO: Init no longer exists.
	// runner := &TestRunner{t: t}
	// tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	connection, err := eftl.Connect("ws://localhost:9191/channel", nil, nil)
	if err != nil {
		t.Errorf("connect failed: %s", err)
	}
	defer connection.Disconnect()

	content := `{"message": "Test message payload!"}`

	t.Log("---- doing first publish ----")
	err = connection.Publish(eftl.Message{
		"_dest":   "test_start",
		"content": []byte(content),
	})
	if err != nil {
		t.Errorf("publish failed: %s", err)
	}

	t.Log("---- doing second publish ----")
	err = connection.Publish(eftl.Message{
		"_dest":   "test_start",
		"content": []byte(content),
	})
	if err != nil {
		t.Errorf("publish failed: %s", err)
	}
}

func TestHandler(t *testing.T) {
	tracer := &opentracing.NoopTracer{}
	span := Span{
		Span: tracer.StartSpan("topic"),
	}
	defer span.Finish()

	handler := &OptimizedHandler{
		defaultActionID: "action_1",
		dispatches: []*Dispatch{
			{
				actionID:  "action_2",
				condition: "${trigger.content.value == A}",
			},
			{
				actionID:  "action_3",
				condition: "${env.value == A}",
			},
			{
				actionID:  "action_4",
				condition: "${#}",
			},
			{
				actionID:  "action_5",
				condition: "${trigger.header.stuff == A}",
			},
		},
	}

	action, _ := handler.GetActionID("null", span)
	if action != "action_1" {
		t.Error("expected action_1")
	}

	content := `{"value": "A"}`
	action, _ = handler.GetActionID(content, span)
	if action != "action_2" {
		t.Error("expected action_2")
	}

	os.Setenv("value", "A")
	action, _ = handler.GetActionID("null", span)
	if action != "action_3" {
		t.Error("expected action_3")
	}
}

const testCreateHandlersConfig string = `{
	"name": "tibco-mqtt",
	"settings": {
	  "url": "ws://localhost:9191/channel",
	  "id": "eftl",
	  "user": "",
	  "password": ""
	},
	"handlers": [
	  {
		"actionId": "action_1",
		"settings": {
		  "dest": "topic_1"
		}
	  },
		  {
		"actionId": "action_2",
		"settings": {
		  "dest": "topic_1",
				  "Condition": "${trigger.content.value == A}"
		}
	  },
		  {
		"actionId": "action_3",
		"settings": {
		  "dest": "topic_2"
		}
	  }
	]
  }`

func TestCreateHandlers(t *testing.T) {
	jsonMetadata := getJsonMetadata()
	// New  factory
	md := trigger.NewMetadata(jsonMetadata)
	f := NewFactory(md)

	// New Trigger
	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testCreateHandlersConfig), config)
	if err != nil {
		t.Error(err)
	}
	tgr := f.New(config).(*Trigger)

	handlers := tgr.CreateHandlers()
	if handlers["topic_1"].defaultActionID != "action_1" {
		t.Error("default action for topic_1 should be action_1")
	}
	if len(handlers["topic_1"].dispatches) != 1 {
		t.Error("there should be 1 dispatches for topic_1")
	}
	if handlers["topic_2"].defaultActionID != "action_3" {
		t.Error("default action for topic_2 should be action_3")
	}
	if len(handlers["topic_2"].dispatches) != 0 {
		t.Error("there should be 0 dispatches for topic_2")
	}
}

const testJSONMessage = `{
	  "replyTo": "abc123",
	  "pathParams": {
		  "param": "a"
	  },
	  "queryParams": {
		  "param": "b"
	  }
  }`

func TestConstructJSONStartRequest(t *testing.T) {
	tracer := &opentracing.NoopTracer{}
	span := Span{
		Span: tracer.StartSpan("topic"),
	}
	defer span.Finish()

	jsonMetadata := getJsonMetadata()
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

	// TODO: Init no longer exists.
	// runner := &TestRunner{t: t}
	// tgr.Init(runner)

	replyTo, params := tgr.(*Trigger).constructStartRequest([]byte(testJSONMessage), span)
	if params == nil {
		t.Fatal("params is nil")
	}
	if replyTo == "" {
		t.Fatal("replyTo is an empty string")
	}
	pathParams := params["pathParams"]
	if pathParams == nil {
		t.Fatal("pathParams is nil")
	}
	if pathParams.(map[string]string)["param"] != "a" {
		t.Fatal("param should be a")
	}
	queryParams := params["queryParams"]
	if queryParams == nil {
		t.Fatal("queryParams is nil")
	}
	if queryParams.(map[string]string)["param"] != "b" {
		t.Fatal("param should be b")
	}
	content := params["content"]
	if content == nil {
		t.Fatal("content is nil")
	}
	tracing := params["tracing"]
	if tracing == nil {
		t.Fatal("tracing is nil")
	}

	if content.(map[string]interface{})["replyTo"] != nil {
		t.Fatal("replyTo should be nil in content")
	}
	if content.(map[string]interface{})["pathParams"] != nil {
		t.Fatal("pathParams should be nil in content")
	}
	if content.(map[string]interface{})["queryParams"] != nil {
		t.Fatal("queryParams should be nil in content")
	}
}

const testXMLMessage = `<?xml version="1.0"?>

 <test replyTo="abc123">
  <message>hello world</message>
  <pathParams>
   <item key="param" value="a"/>
  </pathParams>
  <message>hello world</message>
  <queryParams>
   <item key="param" value="b"/>
  </queryParams>
  <message>hello world</message>
 </test>`

func TestConstructXMLStartRequest(t *testing.T) {
	tracer := &opentracing.NoopTracer{}
	span := Span{
		Span: tracer.StartSpan("topic"),
	}
	defer span.Finish()

	jsonMetadata := getJsonMetadata()
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
	// TODO: Init no longer exists.
	//runner := &TestRunner{t: t}

	// tgr.Init(runner)

	replyTo, params := tgr.(*Trigger).constructStartRequest([]byte(testXMLMessage), span)
	if params == nil {
		t.Fatal("params is nil")
	}
	if replyTo == "" {
		t.Fatal("replyTo is an empty string")
	}
	pathParams := params["pathParams"]
	if pathParams == nil {
		t.Fatal("pathParams is nil")
	}
	if pathParams.(map[string]string)["param"] != "a" {
		t.Fatal("param should be a")
	}
	queryParams := params["queryParams"]
	if queryParams == nil {
		t.Fatal("queryParams is nil")
	}
	if queryParams.(map[string]string)["param"] != "b" {
		t.Fatal("param should be b")
	}
	content := params["content"]
	if content == nil {
		t.Fatal("content is nil")
	}
	tracing := params["tracing"]
	if tracing == nil {
		t.Fatal("tracing is nil")
	}

	if content, ok := content.(map[string]interface{}); ok {
		getRoot := func() map[string]interface{} {
			body := content[util.XMLKeyBody]
			if body == nil {
				return nil
			}
			for _, e := range body.([]interface{}) {
				element, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				name, ok := element[util.XMLKeyType].(string)
				if !ok || name != util.XMLTypeElement {
					continue
				}
				return element
			}
			return nil
		}
		root := getRoot()

		if root["replyTo"] != nil {
			t.Fatal("replyTo should be nil in content")
		}

		find := func(target string) bool {
			rootBody, ok := root[util.XMLKeyBody].([]interface{})
			if !ok {
				return false
			}
			for _, e := range rootBody {
				element, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				name, ok := element[util.XMLKeyName].(string)
				if !ok || name != target {
					continue
				}
				return true
			}

			return false
		}
		if find("pathParams") {
			t.Fatal("pathParams should be nil in content")
		}
		if find("queryParams") {
			t.Fatal("queryParams should be nil in content")
		}
	}
}
