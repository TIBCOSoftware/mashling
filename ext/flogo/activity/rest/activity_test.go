/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package rest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

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

func TestMain(m *testing.M) {
	opentracing.SetGlobalTracer(&opentracing.NoopTracer{})

	database := make([]map[string]interface{}, 0, 10)
	http.HandleFunc("/v2/pet", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			var pet map[string]interface{}
			err = json.Unmarshal(body, &pet)
			if err != nil {
				panic(err)
			}
			pet["id"] = len(database)
			database = append(database, pet)
			body, err = json.Marshal(pet)
			if err != nil {
				panic(err)
			}

			_, err = w.Write(body)
			if err != nil {
				panic(err)
			}
		}
	})

	http.HandleFunc("/v2/pet/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			parts := strings.Split(r.URL.Path, "/")
			id, err := strconv.Atoi(parts[3])
			if err != nil {
				panic(err)
			}
			data, err := json.Marshal(database[id])
			if err != nil {
				panic(err)
			}
			_, err = w.Write(data)
			if err != nil {
				panic(err)
			}
		}
	})

	http.HandleFunc("/v2/pet/findByStatus", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			query := r.URL.Query()
			if query["status"][0] != "ava" {
				panic("invalid status")
			}
			data, err := json.Marshal(database[0])
			if err != nil {
				panic(err)
			}
			_, err = w.Write(data)
			if err != nil {
				panic(err)
			}
		}
	})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		http.Serve(listener, nil)
	}()

	os.Exit(m.Run())
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

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivMethod, "POST")
	tc.SetInput(ivURI, "http://localhost:8080/v2/pet")
	tc.SetInput(ivContent, reqPostStr)

	span := opentracing.StartSpan("test")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	tc.SetInput(ivTracing, ctx)

	//eval
	act.Eval(tc)
	val := tc.GetOutput(ovResult)

	t.Logf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	t.Log("petID:", petID)
	if petID != "0" {
		t.Fatal("invalid pet id")
	}

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
	tc.SetInput(ivURI, "http://localhost:8080/v2/pet/"+petID)

	//eval
	act.Eval(tc)

	val := tc.GetOutput(ovResult)
	t.Logf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	t.Log("petID:", petID)
	if petID != "0" {
		t.Fatal("invalid pet id")
	}
}

func TestParamGet(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput(ivMethod, "GET")
	tc.SetInput(ivURI, "http://localhost:8080/v2/pet/:id")

	pathParams := map[string]string{
		"id": petID,
	}
	tc.SetInput(ivPathParams, pathParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput(ovResult)
	t.Logf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	t.Log("petID:", petID)
	if petID != "0" {
		t.Fatal("invalid pet id")
	}
}

func TestSimpleGetQP(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivMethod, "GET")
	tc.SetInput(ivURI, "http://localhost:8080/v2/pet/findByStatus")

	queryParams := map[string]string{
		"status": "ava",
	}
	tc.SetInput(ivQueryParams, queryParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput(ovResult)
	t.Logf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	t.Log("petID:", petID)
	if petID != "0" {
		t.Fatal("invalid pet id")
	}
}

func TestGetContentType(t *testing.T) {
	contentType := getContentType("test")
	if contentType != contentTypeTextPlain {
		t.Error("content type should be ", contentTypeTextPlain)
	}

	contentType = getContentType(1)
	if contentType != contentTypeTextPlain {
		t.Error("content type should be ", contentTypeTextPlain)
	}

	contentType = getContentType(make([]int, 1))
	if contentType != contentTypeApplicationJSON {
		t.Error("content type should be ", contentTypeApplicationJSON)
	}
}

func TestMethodIsValid(t *testing.T) {
	if !methodIsValid(methodDELETE) {
		t.Error("method should be valid")
	}
}

func TestBuildURI(t *testing.T) {

	uri := "http://localhost:7070/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	t.Log(newURI)
	if newURI != "http://localhost:7070/flow/1234" {
		t.Fatal("invalid uri")
	}
}

func TestBuildURI2(t *testing.T) {

	uri := "https://127.0.0.1:7070/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id":  "1234",
	}

	newURI := BuildURI(uri, params)

	t.Log(newURI)
	if newURI != "https://127.0.0.1:7070/flow/1234/test" {
		t.Fatal("invalid uri")
	}
}

func TestBuildURI3(t *testing.T) {

	uri := "http://localhost/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	t.Log(newURI)
	if newURI != "http://localhost/flow/1234" {
		t.Fatal("invalid uri")
	}
}

func TestBuildURI4(t *testing.T) {

	uri := "https://127.0.0.1/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id":  "1234",
	}

	newURI := BuildURI(uri, params)

	t.Log(newURI)
	if newURI != "https://127.0.0.1/flow/1234/test" {
		t.Fatal("invalid uri")
	}
}

func TestGetCerts(t *testing.T) {
	_, err := getCerts("./certs")
	if err != nil {
		t.Error(err)
	}
}
