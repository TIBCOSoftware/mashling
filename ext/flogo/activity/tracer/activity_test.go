/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package tracer

import (
	"context"
	"io/ioutil"
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
	opentracing.SetGlobalTracer(&opentracing.NoopTracer{})

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	span := opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	tc.SetInput(ivTracing, ctx)

	act.Eval(tc)

	//check result attr
	outSpan := tc.GetOutput(ovSpan)
	if outSpan == nil {
		t.Error("span is nil")
	}
	tracing := tc.GetOutput(ovTracing)
	if tracing == nil {
		t.Error("tracing is nil")
	}

	act = NewActivity(getActivityMetadata())
	tc = test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivSpan, outSpan)

	act.Eval(tc)
}
