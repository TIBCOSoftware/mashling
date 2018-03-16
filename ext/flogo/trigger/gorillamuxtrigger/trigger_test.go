/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package gorillamuxtrigger

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var testJsonMetadata = getJSONMetadata()

func getJSONMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
	 "name": "animals_rest_trigger",
	 "id": "animals_rest_trigger",
	 "ref": "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger",
	 "settings": {
		 "port": "9096"
	 },
	 "output": null,
	 "handlers": [
		 {
			 "actionId": "mammals_handler",
			 "settings": {
				 "Condition": "${trigger.content.name in (ELEPHANT,CAT)}",
				 "autoIdReply": "false",
				 "method": "PUT",
				 "path": "/pets",
				 "useReplyHandler": "false"
			 },
			 "output": null,
			 "outputs": null
		 },
		 {
			 "actionId": "birds_handler",
			 "settings": {
				 "Condition": "${trigger.content.name == SPARROW}",
				 "autoIdReply": "false",
				 "method": "PUT",
				 "path": "/pets",
				 "useReplyHandler": "false"
			 },
			 "output": null,
			 "outputs": null
		 },
		 {
			 "actionId": "animals_handler",
			 "settings": {
				 "autoIdReply": "false",
				 "method": "PUT",
				 "path": "/pets",
				 "useReplyHandler": "false"
			 },
			 "output": null,
			 "outputs": null
		 },
		 {
			 "actionId": "animals_get_handler",
			 "settings": {
				 "autoIdReply": "false",
				 "method": "GET",
				 "path": "/pets/{petId}",
				 "useReplyHandler": "false"
			 },
			 "output": null,
			 "outputs": null
		 }
	 ],
	 "outputs": null
 }`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestInitOk(t *testing.T) {
	// New  factory
	f := &RestFactory{}

	config := trigger.Config{}
	tgr := f.New(&config)

	runner := &TestRunner{}

	json.Unmarshal([]byte(testConfig), &config)
	tgr.Init(runner)
}

//Run the specified Action
func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	return nil, nil
}

func TestHandlerOk(t *testing.T) {

	// New  factory
	md := trigger.NewMetadata(testJsonMetadata)
	f := NewFactory(md)

	config := trigger.Config{}
	tgr := f.New(&config)
	runner := &TestRunner{}
	json.Unmarshal([]byte(testConfig), &config)

	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()

	uri := "http://localhost:9096/pets/6"

	req, err := http.NewRequest("GET", uri, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		t.Fail()
	}
}

func BenchmarkHandlerOk(b *testing.B) {

	// New  factory
	md := trigger.NewMetadata(testJsonMetadata)
	f := NewFactory(md)

	config := trigger.Config{}
	tgr := f.New(&config)
	runner := &TestRunner{}
	json.Unmarshal([]byte(testConfig), &config)

	tgr.Init(runner)
	tgr.Start()
	defer tgr.Stop()

	uri := "http://localhost:9096/pets/6"

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	for i := 0; i < b.N; i++ {
		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			b.Fail()
		}
	}

}
