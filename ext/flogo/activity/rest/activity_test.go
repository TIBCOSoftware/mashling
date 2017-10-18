/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	opentracing "github.com/opentracing/opentracing-go"
)

const reqPostStr string = `{
  "name": "my pet"
}
`

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

var petID string

func TestSimplePost(t *testing.T) {
	opentracing.SetGlobalTracer(&opentracing.NoopTracer{})

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivMethod, "POST")
	tc.SetInput(ivURI, "http://petstore.swagger.io/v2/pet")
	tc.SetInput(ivContent, reqPostStr)

	span := opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	tc.SetInput(ivTracing, ctx)

	//eval
	act.Eval(tc)
	val := tc.GetOutput(ovResult)

	fmt.Printf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	fmt.Println("petID:", petID)

	tracing := tc.GetOutput(ovTracing)
	if tracing == nil {
		t.Error("tracing is nil")
	}
}

func TestSimpleGet(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivMethod, "GET")
	tc.SetInput(ivURI, "http://petstore.swagger.io/v2/pet/"+petID)

	//eval
	act.Eval(tc)

	val := tc.GetOutput(ovResult)
	fmt.Printf("result: %v\n", val)
}

func TestParamGet(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput(ivMethod, "GET")
	tc.SetInput(ivURI, "http://petstore.swagger.io/v2/pet/:id")

	pathParams := map[string]string{
		"id": petID,
	}
	tc.SetInput(ivPathParams, pathParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput(ovResult)
	fmt.Printf("result: %v\n", val)
}

func TestSimpleGetQP(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivMethod, "GET")
	tc.SetInput(ivURI, "http://petstore.swagger.io/v2/pet/findByStatus")

	queryParams := map[string]string{
		"status": "ava",
	}
	tc.SetInput(ivQueryParams, queryParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput(ovResult)
	fmt.Printf("result: %v\n", val)
}

func TestBuildURI(t *testing.T) {

	uri := "http://localhost:7070/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI2(t *testing.T) {

	uri := "https://127.0.0.1:7070/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id":  "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI3(t *testing.T) {

	uri := "http://localhost/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI4(t *testing.T) {

	uri := "https://127.0.0.1/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id":  "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}
