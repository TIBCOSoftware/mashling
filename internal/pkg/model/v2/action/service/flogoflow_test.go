package service

import (
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

const flogoflow = `{
 "flow": {
  "explicitReply": true,
  "type": 1,
  "attributes": [],
  "rootTask": {
   "id": 1,
   "type": 1,
   "tasks": [],
   "links": [],
   "attributes": []
  }
 }
}`

func TestFlogoFlow(t *testing.T) {
	logger.SetLogLevel(logger.ErrorLevel)
	ff := flow.ActionFactory{}
	ff.Init()

	var data interface{}
	err := json.Unmarshal([]byte(flogoflow), &data)
	if err != nil {
		t.Fatal(err)
	}
	service := types.Service{
		Type: "flogoFlow",
		Settings: map[string]interface{}{
			"reference": "http://example.com/test.json",
			"definition": map[string]interface{}{
				"ref":  "github.com/TIBCOSoftware/flogo-contrib/action/flow",
				"data": data,
			},
		},
	}
	instance, err := Initialize(service)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
