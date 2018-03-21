package twilio

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/sfreiberg/gotwilio"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-twilio")

const (
	ivAcctSID   = "accountSID"
	ivAuthToken = "authToken"
	ivFrom      = "from"
	ivTo        = "to"
	ivMessage   = "message"
)

// TwilioActivity is a Twilio Activity implementation
type TwilioActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new TwilioActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TwilioActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *TwilioActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *TwilioActivity) Eval(context activity.Context) (done bool, err error) {

	accountSID := context.GetInput(ivAcctSID).(string)
	authToken := context.GetInput(ivAuthToken).(string)
	from := context.GetInput(ivFrom).(string)
	to := context.GetInput(ivTo).(string)
	message := context.GetInput(ivMessage).(string)

	twilio := gotwilio.NewTwilioClient(accountSID, authToken)

	resp, _, err := twilio.SendSMS(from, to, message, "", "")

	if err != nil {
		log.Error("Error sending SMS:", err)
	}

	log.Debug("Response:", resp)

	return true, nil
}
