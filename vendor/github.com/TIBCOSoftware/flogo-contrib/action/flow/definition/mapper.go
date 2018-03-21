package definition

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
)

// MapperDef represents a Mapper, which is a collection of mappings
type MapperDef struct {
	Mappings []*data.MappingDef
}

type MapperFactory interface {
	// NewMapper creates a new Mapper from the specified MapperDef
	NewMapper(mapperDef *MapperDef) data.Mapper

	// NewTaskInputMapper creates a new Input Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	NewTaskInputMapper(task *Task, mapperDef *MapperDef) data.Mapper

	// NewTaskOutputMapper creates a new Output Mapper from the specified MapperDef
	// for the specified Task, method to facilitate pre-compiled mappers
	NewTaskOutputMapper(task *Task, mapperDef *MapperDef) data.Mapper

	// GetDefaultTaskOutputMapper get the default Output Mapper for the
	// specified Task
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

func (mf *BasicMapperFactory) NewTaskInputMapper(task *Task, mapperDef *MapperDef) data.Mapper {
	id := task.definition.name + "." + task.id + ".input"
	return mf.baseFactory.NewUniqueMapper(id, &data.MapperDef{Mappings: mapperDef.Mappings}, GetDataResolver())
}

func (mf *BasicMapperFactory) NewTaskOutputMapper(task *Task, mapperDef *MapperDef) data.Mapper {
	id := task.definition.name + "." + task.id + ".output"
	return mf.baseFactory.NewUniqueMapper(id, &data.MapperDef{Mappings: mapperDef.Mappings}, nil)
}

func (mf *BasicMapperFactory) GetDefaultTaskOutputMapper(task *Task) data.Mapper {
	return &DefaultOutputMapper{task: task}
}

// BasicMapper is a simple object holding and executing mappings
type DefaultOutputMapper struct {
	task *Task
}

func (m *DefaultOutputMapper) Apply(inputScope data.Scope, outputScope data.Scope) error {

	oscope := outputScope.(data.MutableScope)

	act := activity.Get(m.task.ActivityRef())

	attrNS := "_A." + m.task.ID() + "."

	for _, attr := range act.Metadata().Output {

		oAttr, _ := inputScope.GetAttr(attr.Name())

		if oAttr != nil {
			oscope.AddAttr(attrNS+attr.Name(), attr.Type(), oAttr.Value())
		}
	}

	return nil
}
