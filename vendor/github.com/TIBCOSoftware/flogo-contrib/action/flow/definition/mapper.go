package definition

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
	"sync"
)

// MapperDef represents a Mapper, which is a collection of mappings
type MapperDef struct {
	Mappings []*data.MappingDef
}

type MapperFactory interface {
	// NewMapper creates a new Mapper from the specified MapperDef
	NewMapper(mapperDef *MapperDef) data.Mapper

	// NewActivityInputMapper creates a new Activity Input Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	NewActivityInputMapper(task *Task, mapperDef *data.MapperDef) data.Mapper

	// NewActivityOutputMapper creates a new Activity Output Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	NewActivityOutputMapper(task *Task, mapperDef *data.MapperDef) data.Mapper

	// GetDefaultTaskOutputMapper get the default Activity Output Mapper for the
	// specified Task
	GetDefaultActivityOutputMapper(task *Task) data.Mapper

	// NewTaskInputMapper creates a new Input Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	//Deprecated
	NewTaskInputMapper(task *Task, mapperDef *MapperDef) data.Mapper

	// NewTaskOutputMapper creates a new Output Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	//Deprecated
	NewTaskOutputMapper(task *Task, mapperDef *MapperDef) data.Mapper

	// GetDefaultTaskOutputMapper get the default Output Mapper for the
	// specified Task
	//Deprecated
	GetDefaultTaskOutputMapper(task *Task) data.Mapper
}

var mapperFactory MapperFactory

func SetMapperFactory(factory MapperFactory) {
	mapperFactory = factory

	baseFactory, ok := interface{}(factory).(mapper.Factory)
	if ok {
		mapper.SetFactory(baseFactory)
	}
}

func GetMapperFactory() MapperFactory {

	//temp hack until we consolidate mapper definition
	if mapperFactory == nil {
		mapperFactory = &BasicMapperFactory{baseFactory: mapper.GetFactory()}
	}

	return mapperFactory
}

type BasicMapperFactory struct {
	baseFactory mapper.Factory
}

func (mf *BasicMapperFactory) NewMapper(mapperDef *MapperDef) data.Mapper {
	return mf.baseFactory.NewMapper(&data.MapperDef{Mappings: mapperDef.Mappings}, GetDataResolver())
}

func (mf *BasicMapperFactory) NewActivityInputMapper(task *Task, mapperDef *data.MapperDef) data.Mapper {
	id := task.definition.name + "." + task.id + ".input"
	return mf.baseFactory.NewUniqueMapper(id, mapperDef, GetDataResolver())
}

func (mf *BasicMapperFactory) NewActivityOutputMapper(task *Task, mapperDef *data.MapperDef) data.Mapper {
	id := task.definition.name + "." + task.id + ".output"
	return mf.baseFactory.NewUniqueMapper(id, mapperDef, nil)
}

func (mf *BasicMapperFactory) GetDefaultActivityOutputMapper(task *Task) data.Mapper {
	act := task.activityCfg.Activity
	attrNS := "_A." + task.ID() + "."

	if act.Metadata().DynamicIO {

		return &DefaultActivityOutputMapper{attrNS: attrNS, task: task}
		////todo validate dynamic on instantiation
		//dynamic, _ := act.(activity.DynamicIO)
		//dynamicIO, _ := dynamic.IOMetadata(&DummyTaskCtx{task: task})
		////todo handler error
		//if dynamicIO != nil {
		//	return &DefaultActivityOutputMapper{attrNS: attrNS, outputMetadata: dynamicIO.Output}
		//}
	}

	return &DefaultActivityOutputMapper{attrNS: attrNS, outputMetadata: act.Metadata().Output}
}

// BasicMapper is a simple object holding and executing mappings
type DefaultActivityOutputMapper struct {
	attrNS string
	//activityMetadata *activity.Metadata
	outputMetadata map[string]*data.Attribute

	task *Task

	mutex sync.Mutex


}

func (m *DefaultActivityOutputMapper) Apply(inputScope data.Scope, outputScope data.Scope) error {

	m.mutex.Lock()
	if m.outputMetadata == nil {
		act := m.task.activityCfg.Activity
		if act.Metadata().DynamicIO {
			//todo validate dynamic on instantiation
			dynamic, _ := act.(activity.DynamicIO)
			dynamicIO, _ := dynamic.IOMetadata(&DummyTaskCtx{task: m.task})
			//todo handler error
			if dynamicIO != nil {
				m.outputMetadata = dynamicIO.Output
			} else {
				m.outputMetadata =act.Metadata().Output
			}
		}

	}
	m.mutex.Unlock()

	oscope := outputScope.(data.MutableScope)

	for _, attr := range m.outputMetadata {

		oAttr, _ := inputScope.GetAttr(attr.Name())

		if oAttr != nil {
			oscope.AddAttr(m.attrNS+attr.Name(), attr.Type(), oAttr.Value())
		}
	}

	return nil
}

//Deprecated
func (mf *BasicMapperFactory) NewTaskInputMapper(task *Task, mapperDef *MapperDef) data.Mapper {
	id := task.definition.name + "." + task.id + ".input"
	return mf.baseFactory.NewUniqueMapper(id, &data.MapperDef{Mappings: mapperDef.Mappings}, GetDataResolver())
}

//Deprecated
func (mf *BasicMapperFactory) NewTaskOutputMapper(task *Task, mapperDef *MapperDef) data.Mapper {
	id := task.definition.name + "." + task.id + ".output"
	return mf.baseFactory.NewUniqueMapper(id, &data.MapperDef{Mappings: mapperDef.Mappings}, nil)
}

//Deprecated
func (mf *BasicMapperFactory) GetDefaultTaskOutputMapper(task *Task) data.Mapper {
	return &DefaultTaskOutputMapper{task: task}
}

// BasicMapper is a simple object holding and executing mappings
//Deprecated
type DefaultTaskOutputMapper struct {
	task *Task
}

//Deprecated
func (m *DefaultTaskOutputMapper) Apply(inputScope data.Scope, outputScope data.Scope) error {

	oscope := outputScope.(data.MutableScope)

	act := activity.Get(m.task.ActivityConfig().Ref())

	attrNS := "_A." + m.task.ID() + "."

	for _, attr := range act.Metadata().Output {

		oAttr, _ := inputScope.GetAttr(attr.Name())

		if oAttr != nil {
			oscope.AddAttr(attrNS+attr.Name(), attr.Type(), oAttr.Value())
		}
	}

	return nil
}

//Temporary hack for determining dynamic default outputs

type DummyTaskCtx struct {
	task *Task
}

func (*DummyTaskCtx) ActivityHost() activity.Host {
	return nil
}

func (ctx *DummyTaskCtx) Name() string {
	return ctx.task.Name()
}

func (ctx *DummyTaskCtx) GetSetting(setting string) (value interface{}, exists bool) {
	val, found := ctx.task.ActivityConfig().GetSetting(setting)
	if found {
		return val.Value(), true
	}

	return nil, false
}

func (*DummyTaskCtx) GetInitValue(key string) (value interface{}, exists bool) {
	return nil, false
}

func (*DummyTaskCtx) GetInput(name string) interface{} {
	return ""
}

func (*DummyTaskCtx) GetOutput(name string) interface{} {
	return nil
}

func (*DummyTaskCtx) SetOutput(name string, value interface{}) {
}

func (*DummyTaskCtx) GetSharedTempData() map[string]interface{} {
	return nil
}

func (ctx *DummyTaskCtx) TaskName() string {
	return ctx.task.Name()
}

func (*DummyTaskCtx) FlowDetails() activity.FlowDetails {
	return nil
}
