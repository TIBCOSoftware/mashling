package pingreply

import (
	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-reply")

const (
	ivCode = "code"
	ivData = "data"
)

// ReplyActivity is an Activity that is used to reply via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type ReplyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new ReplyActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ReplyActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ReplyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ReplyActivity) Eval(ctx activity.Context) (done bool, err error) {

	replyCode := ctx.GetInput(ivCode).(int)
	replyData := ctx.GetInput(ivData)

	result := replyData.(map[string]interface{})

	if strings.Compare(result["response"].(string), "Detailed") == 0 {

		pingDataVar := util.PingDataPntr.GetData()
		dataBytes := []byte(`{
			"response":"success",
			"mashlingCliRevision":"` + pingDataVar.MashlingCliRev + `",
			"MashlingCliLocalRev":"` + pingDataVar.MashlingCliLocalRev + `",
			"MashlingCliVersion":"` + pingDataVar.MashlingCliVersion + `",
			"SchemaVersion":"` + pingDataVar.SchemaVersion + `",
			"AppVersion":"` + pingDataVar.AppVersion + `",
			"FlogolibRev":"` + pingDataVar.FlogolibRev + `",
			"MashlingRev":"` + pingDataVar.MashlingRev + `",
			"AppDescrption":"` + pingDataVar.AppDescrption + `"
			}`)

		var replyJSON interface{}
		err = json.Unmarshal(dataBytes, &replyJSON)

		replyData = replyJSON
	}

	log.Debugf("Code :'%d', Data: '%+v'", replyCode, replyData)

	replyHandler := ctx.FlowDetails().ReplyHandler()

	//todo support replying with error

	if replyHandler != nil {

		replyHandler.Reply(replyCode, replyData, nil)
	}

	return true, nil
}
