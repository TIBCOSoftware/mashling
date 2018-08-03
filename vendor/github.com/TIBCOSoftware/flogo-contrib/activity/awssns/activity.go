package awssns

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// List of input and output variables names
const (
	ConfAWSAccessKeyID     = "accessKey"
	ConfAWSSecretAccessKey = "secretKey"
	ConfAWSDefaultRegion   = "region"
	ConfSMSType            = "smsType"
	ConfSMSFrom            = "from"
	ConfSMSTo              = "to"
	ConfSMSMessage         = "message"
	OUTMessageID           = "messageId"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-awssns")

// AWSSNS Structure for the AWSNS activity
type AWSSNS struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AWSSNS{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AWSSNS) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AWSSNS) Eval(context activity.Context) (done bool, err error) {

	if context.GetInput(ConfAWSAccessKeyID) == nil || context.GetInput(ConfAWSSecretAccessKey) == nil || context.GetInput(ConfAWSDefaultRegion) == nil || context.GetInput(ConfSMSFrom) == nil || context.GetInput(ConfSMSTo) == nil || context.GetInput(ConfSMSMessage) == nil {
		log.Error("Required variables have not been set !")
		return false, fmt.Errorf("required variables have not been set")
	}

	AWSAccessKeyID := context.GetInput(ConfAWSAccessKeyID).(string)
	AWSSecretAccessKey := context.GetInput(ConfAWSSecretAccessKey).(string)
	AWSDefaultRegion := context.GetInput(ConfAWSDefaultRegion).(string)
	SMSFrom := context.GetInput(ConfSMSFrom).(string)
	SMSTo := context.GetInput(ConfSMSTo).(string)
	SMSMessage := context.GetInput(ConfSMSMessage).(string)
	SMSType := context.GetInput(ConfSMSType).(string)

	log.Debug("Setting credentials")
	var snsCreds = credentials.NewStaticCredentials(AWSAccessKeyID, AWSSecretAccessKey, "")
	log.Debug("Creating session")
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: snsCreds,
		Region:      aws.String(AWSDefaultRegion),
	}))
	log.Debug("Session created")

	log.Debug("Creating service")
	svc := sns.New(sess)
	log.Debug("Service created")

	log.Debug("Setting SMS parameters")
	params := &sns.PublishInput{
		Message:     aws.String(SMSMessage),
		PhoneNumber: aws.String(SMSTo),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"AWS.SNS.SMS.SenderID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(SMSFrom),
			},
			"AWS.SNS.SMS.SMSType": {
				DataType:    aws.String("String"),
				StringValue: aws.String(SMSType),
			},
		},
	}
	log.Debug("SMS parameters set.")
	log.Debug("Publishing SMS")
	resp, err := svc.Publish(params)

	if err != nil {
		log.Errorf(err.Error())
		return false, err
	}
	context.SetOutput(OUTMessageID, *resp.MessageId)
	log.Infof("Message sent [%s]", *resp.MessageId)
	return true, nil
}
