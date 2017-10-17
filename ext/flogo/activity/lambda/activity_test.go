package lambda

import (
	"encoding/json"
	"testing"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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

type Hello struct {
	Name string `json:"name"`
}

func TestLambdaInvokeWithSecurity(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("arn", "<arn:..")
	tc.SetInput("region", "us-east-1")
	tc.SetInput("accessKey", "<access_key>")
	tc.SetInput("secretKey", "<secret>")

	payLoad := Hello{
		Name: "Matt",
	}
	b, _ := json.Marshal(payLoad)
	tc.SetInput("payload", string(b))

	//eval
	_, err := act.Eval(tc)

	if err == nil {
		t.Fail()
	}
}
