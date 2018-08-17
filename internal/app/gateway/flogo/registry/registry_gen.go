// Package registry handles Flogo trigger and activity registration. DO NOT EDIT!! This file was auto-generated by mashling at
// 2018-08-17 09:06:22.006314608 -0400 EDT m=+0.002396474
package registry

import (
	vendorgithubcomTIBCOSoftwareflogocontribactivityactreply "github.com/TIBCOSoftware/flogo-contrib/activity/actreply"
	vendorgithubcomTIBCOSoftwareflogocontribactivityactreturn "github.com/TIBCOSoftware/flogo-contrib/activity/actreturn"
	vendorgithubcomTIBCOSoftwareflogocontribactivityaggregate "github.com/TIBCOSoftware/flogo-contrib/activity/aggregate"
	vendorgithubcomTIBCOSoftwareflogocontribactivityapp "github.com/TIBCOSoftware/flogo-contrib/activity/app"
	vendorgithubcomTIBCOSoftwareflogocontribactivityawsiot "github.com/TIBCOSoftware/flogo-contrib/activity/awsiot"
	vendorgithubcomTIBCOSoftwareflogocontribactivityawssns "github.com/TIBCOSoftware/flogo-contrib/activity/awssns"
	vendorgithubcomTIBCOSoftwareflogocontribactivitycoap "github.com/TIBCOSoftware/flogo-contrib/activity/coap"
	vendorgithubcomTIBCOSoftwareflogocontribactivitycouchbase "github.com/TIBCOSoftware/flogo-contrib/activity/couchbase"
	vendorgithubcomTIBCOSoftwareflogocontribactivitycounter "github.com/TIBCOSoftware/flogo-contrib/activity/counter"
	vendorgithubcomTIBCOSoftwareflogocontribactivityerror "github.com/TIBCOSoftware/flogo-contrib/activity/error"
	vendorgithubcomTIBCOSoftwareflogocontribactivitykafkapub "github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub"
	vendorgithubcomTIBCOSoftwareflogocontribactivitylambda "github.com/TIBCOSoftware/flogo-contrib/activity/lambda"
	vendorgithubcomTIBCOSoftwareflogocontribactivitylog "github.com/TIBCOSoftware/flogo-contrib/activity/log"
	vendorgithubcomTIBCOSoftwareflogocontribactivitymapper "github.com/TIBCOSoftware/flogo-contrib/activity/mapper"
	vendorgithubcomTIBCOSoftwareflogocontribactivitymongodb "github.com/TIBCOSoftware/flogo-contrib/activity/mongodb"
	vendorgithubcomTIBCOSoftwareflogocontribactivityrest "github.com/TIBCOSoftware/flogo-contrib/activity/rest"
	vendorgithubcomTIBCOSoftwareflogocontribactivitysubflow "github.com/TIBCOSoftware/flogo-contrib/activity/subflow"
	vendorgithubcomTIBCOSoftwareflogocontribactivitytwilio "github.com/TIBCOSoftware/flogo-contrib/activity/twilio"
	vendorgithubcomTIBCOSoftwareflogocontribactivitywsmessage "github.com/TIBCOSoftware/flogo-contrib/activity/wsmessage"
	vendorgithubcomTIBCOSoftwareflogocontribtriggercli "github.com/TIBCOSoftware/flogo-contrib/trigger/cli"
	vendorgithubcomTIBCOSoftwareflogocontribtriggercoap "github.com/TIBCOSoftware/flogo-contrib/trigger/coap"
	vendorgithubcomTIBCOSoftwareflogocontribtriggerkafkasub "github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub"
	vendorgithubcomTIBCOSoftwareflogocontribtriggerlambda "github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	vendorgithubcomTIBCOSoftwareflogocontribtriggermqtt "github.com/TIBCOSoftware/flogo-contrib/trigger/mqtt"
	vendorgithubcomTIBCOSoftwareflogocontribtriggerrest "github.com/TIBCOSoftware/flogo-contrib/trigger/rest"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	extflogoactivityeftl "github.com/TIBCOSoftware/mashling/ext/flogo/activity/eftl"
	extflogoactivitylambda "github.com/TIBCOSoftware/mashling/ext/flogo/activity/lambda"
	extflogoactivitymqtt "github.com/TIBCOSoftware/mashling/ext/flogo/activity/mqtt"
	extflogoactivityping "github.com/TIBCOSoftware/mashling/ext/flogo/activity/ping"
	extflogoactivityrest "github.com/TIBCOSoftware/mashling/ext/flogo/activity/rest"
	extflogoactivitytracer "github.com/TIBCOSoftware/mashling/ext/flogo/activity/tracer"
	extflogotriggereftl "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/eftl"
	extflogotriggergorillamuxtrigger "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger"
	extflogotriggergrpc "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/grpc"
	extflogotriggerkafkasubrouter "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/kafkasubrouter"
	extflogotriggermqtt "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/mqtt"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry/activities"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry/triggers"
	// Force github.com/TIBCOSoftware/mashling/pkg/flogo/action to be loaded and init run.
	_ "github.com/TIBCOSoftware/mashling/pkg/flogo/action"
	// Force github.com/TIBCOSoftware/flogo-contrib/action/flow to be loaded and init run.
	_ "github.com/TIBCOSoftware/flogo-contrib/action/flow"
)

// SupportedImports contains all Flogo imports supported in the code base.
var SupportedImports map[string]bool

// init registers the triggers and activities.
func init() {
	err := RegisterUniverse()
	if err != nil {
		panic(err.Error())
	}
}

// RegisterUniverse registers all the activities and triggers that we know how to handle.
func RegisterUniverse() error {
	SupportedImports = make(map[string]bool)
	var err error
	// Register triggers
	var tr []byte
	var trmd *trigger.Metadata
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/trigger/mqtt
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/trigger/mqtt"] = true
	tr, err = triggers.Asset("ext/flogo/trigger/mqtt/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, extflogotriggermqtt.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/trigger/cli
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/trigger/cli"] = true
	tr, err = triggers.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/trigger/cli/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, vendorgithubcomTIBCOSoftwareflogocontribtriggercli.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/trigger/coap
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/trigger/coap"] = true
	tr, err = triggers.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/trigger/coap/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, vendorgithubcomTIBCOSoftwareflogocontribtriggercoap.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/trigger/lambda
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"] = true
	tr, err = triggers.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/trigger/lambda/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, vendorgithubcomTIBCOSoftwareflogocontribtriggerlambda.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/trigger/mqtt
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/trigger/mqtt"] = true
	tr, err = triggers.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/trigger/mqtt/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, vendorgithubcomTIBCOSoftwareflogocontribtriggermqtt.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/trigger/rest
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/trigger/rest"] = true
	tr, err = triggers.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/trigger/rest/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, vendorgithubcomTIBCOSoftwareflogocontribtriggerrest.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/trigger/grpc
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/trigger/grpc"] = true
	tr, err = triggers.Asset("ext/flogo/trigger/grpc/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, extflogotriggergrpc.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/trigger/kafkasubrouter
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/trigger/kafkasubrouter"] = true
	tr, err = triggers.Asset("ext/flogo/trigger/kafkasubrouter/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, extflogotriggerkafkasubrouter.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub"] = true
	tr, err = triggers.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/trigger/kafkasub/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, vendorgithubcomTIBCOSoftwareflogocontribtriggerkafkasub.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/trigger/eftl
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/trigger/eftl"] = true
	tr, err = triggers.Asset("ext/flogo/trigger/eftl/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, extflogotriggereftl.NewFactory(trmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger"] = true
	tr, err = triggers.Asset("ext/flogo/trigger/gorillamuxtrigger/trigger.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	trmd = trigger.NewMetadata(string(tr))
	trigger.RegisterFactory(trmd.ID, extflogotriggergorillamuxtrigger.NewFactory(trmd))

	// Register activities
	var act []byte
	var actmd *activity.Metadata
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/counter
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/counter"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/counter/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitycounter.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitykafkapub.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/wsmessage
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/wsmessage"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/wsmessage/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitywsmessage.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/activity/lambda
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/activity/lambda"] = true
	act, err = activities.Asset("ext/flogo/activity/lambda/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(extflogoactivitylambda.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/awssns
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/awssns"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/awssns/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityawssns.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/activity/ping
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/activity/ping"] = true
	act, err = activities.Asset("ext/flogo/activity/ping/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(extflogoactivityping.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/activity/rest
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/activity/rest"] = true
	act, err = activities.Asset("ext/flogo/activity/rest/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(extflogoactivityrest.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/activity/tracer
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/activity/tracer"] = true
	act, err = activities.Asset("ext/flogo/activity/tracer/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(extflogoactivitytracer.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/actreply
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/actreply"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/actreply/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityactreply.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/actreturn
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/actreturn"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/actreturn/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityactreturn.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/twilio
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/twilio"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/twilio/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitytwilio.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/activity/eftl
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/activity/eftl"] = true
	act, err = activities.Asset("ext/flogo/activity/eftl/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(extflogoactivityeftl.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/mashling/ext/flogo/activity/mqtt
	SupportedImports["github.com/TIBCOSoftware/mashling/ext/flogo/activity/mqtt"] = true
	act, err = activities.Asset("ext/flogo/activity/mqtt/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(extflogoactivitymqtt.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/coap
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/coap"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/coap/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitycoap.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/couchbase
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/couchbase"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/couchbase/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitycouchbase.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/lambda
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/lambda"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/lambda/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitylambda.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/mapper
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/mapper"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/mapper/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitymapper.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/mongodb
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/mongodb"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/mongodb/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitymongodb.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/aggregate
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/aggregate"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/aggregate/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityaggregate.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/app
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/app"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/app/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityapp.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/log
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/log"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/log/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitylog.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/rest
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/rest"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/rest/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityrest.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/subflow
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/subflow"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/subflow/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivitysubflow.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/awsiot
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/awsiot"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/awsiot/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityawsiot.NewActivity(actmd))
	// Register github.com/TIBCOSoftware/flogo-contrib/activity/error
	SupportedImports["github.com/TIBCOSoftware/flogo-contrib/activity/error"] = true
	act, err = activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/error/activity.json")
	if err != nil {
		// Asset was not found.
		return err
	}
	actmd = activity.NewMetadata(string(act))
	activity.Register(vendorgithubcomTIBCOSoftwareflogocontribactivityerror.NewActivity(actmd))
	return nil
}
