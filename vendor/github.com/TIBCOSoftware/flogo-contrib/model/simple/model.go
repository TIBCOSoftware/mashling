package simple

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/model/simple/behaviors"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("flowmodel-simple")

const (
	MODEL_NAME = "tibco-simple"
)

func init() {
	model.Register(New())
}

func New() *model.FlowModel {
	m := model.New(MODEL_NAME)
	m.RegisterFlowBehavior(&behaviors.Flow{})
	m.RegisterDefaultTaskBehavior("basic", &behaviors.Task{})
	m.RegisterTaskBehavior("iterator", &behaviors.IteratorTask{})

	return m
}
