package mqtt

import (
	"context"
	"io/ioutil"
	"net"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	opentracing "github.com/opentracing/opentracing-go"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {
	_, err := net.Dial("tcp", "127.0.0.1:1883")
	if err != nil {
		t.Log("MQTT message broker is not available, skipping test...")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivContent, `{"test": "hello world"}`)
	tc.SetInput(ivTopic, "test")
	tc.SetInput(ivBroker, "tcp://localhost:1883")
	tc.SetInput(ivID, "flogo")
	tc.SetInput(ivUser, "")
	tc.SetInput(ivPassword, "")
	tc.SetInput(ivQOS, float64(0))

	span := opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	tc.SetInput(ivTracing, ctx)
	defer span.Finish()

	act.Eval(tc)

	//check result attr
	tracing := tc.GetOutput(ovTracing)
	if tracing == nil {
		t.Error("tracing is nil")
	}
}
