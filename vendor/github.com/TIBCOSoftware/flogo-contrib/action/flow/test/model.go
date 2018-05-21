package test

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/model/simple"
)

func init() {
	model.Register(NewTestModel())
}

func NewTestModel() *model.FlowModel {
	m := model.New("test")
	m.RegisterFlowBehavior(&simple.FlowBehavior{})
	m.RegisterDefaultTaskBehavior("basic", &simple.TaskBehavior{})

	return m
}
