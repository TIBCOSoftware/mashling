package mapper

import (
	"fmt"

	"encoding/json"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var mapplerLog = logger.GetLogger("basic-mapper")

type Factory interface {
	// NewMapper creates a new data.Mapper from the specified data.MapperDef
	NewMapper(mapperDef *data.MapperDef, resolver data.Resolver) data.Mapper

	// NewUniqueMapper creates a unique data.Mapper from the specified data.MapperDef
	// the ID can be used to facilitate use precompiled mappers
	NewUniqueMapper(ID string, mapperDef *data.MapperDef, resolver data.Resolver) data.Mapper
}

var factory Factory

func SetFactory(factory Factory) {
	factory = factory
}

func GetFactory() Factory {

	if factory == nil {
		factory = &BasicMapperFactory{}
	}

	return factory
}

type BasicMapperFactory struct {
}

func (mf *BasicMapperFactory) NewMapper(mapperDef *data.MapperDef, resolver data.Resolver) data.Mapper {
	return NewBasicMapper(mapperDef, resolver)
}

func (mf *BasicMapperFactory) NewUniqueMapper(ID string, mapperDef *data.MapperDef, resolver data.Resolver) data.Mapper {
	return NewBasicMapper(mapperDef, resolver)
}

// BasicMapper is a simple object holding and executing mappings
type BasicMapper struct {
	mappings []*data.MappingDef
	resolver data.Resolver
}

// NewBasicMapper creates a new BasicMapper with the specified mappings
func NewBasicMapper(mapperDef *data.MapperDef, resolver data.Resolver) data.Mapper {

	var mapper BasicMapper
	mapper.mappings = mapperDef.Mappings

	if resolver == nil {
		mapper.resolver = &data.BasicResolver{}
	} else {
		mapper.resolver = resolver
	}

	return &mapper
}

// Mappings gets the mappings of the BasicMapper
func (m *BasicMapper) Mappings() []*data.MappingDef {
	return m.mappings
}

// Apply executes the mappings using the values from the input scope
// and puts the results in the output scope
//
// return error
func (m *BasicMapper) Apply(inputScope data.Scope, outputScope data.Scope) error {

	if err := m.UpdateMapping(); err != nil {
		return fmt.Errorf("Update mapping ref error %s", err.Error())
	}

	//todo validate types
	for _, mapping := range m.mappings {

		switch mapping.Type {
		case data.MtAssign:

			toResolve, ok := mapping.Value.(string)
			if !ok {
				return fmt.Errorf("invalid assign value: %v", mapping.Value)
			}

			var val interface{}
			var err error

			if m.resolver != nil {
				val, err = m.resolver.Resolve(toResolve, inputScope)
				if err != nil {
					return err
				}
			}

			assignExpr := NewAssignExpr(mapping.MapTo, val)
			_, err = assignExpr.Eval(outputScope)
			if err != nil {
				return err
			}

		case data.MtLiteral:
			assignExpr := NewAssignExpr(mapping.MapTo, mapping.Value)

			_, err := assignExpr.Eval(outputScope)

			if err != nil {
				return err
			}
		case data.MtObject:

			val, err := MapObject(mapping.Value, inputScope, m.resolver)
			if err != nil {
				return err
			}

			err = outputScope.SetAttrValue(mapping.MapTo, val)
			if err != nil {
				return err
			}
		case data.MtExpression:
			err := exprmapper.Map(mapping, inputScope, outputScope, m.resolver)
			if err != nil {
				return fmt.Errorf("Expression mapping failed, due to %s", err.Error())
			}
		case data.MtArray:
			//ArrayMapping
			mapplerLog.Debugf("Array mapping value %s", mapping.Value)
			//Array mapping value must be string
			arrayMapping, err := exprmapper.ParseArrayMapping(mapping.Value)
			if err != nil {
				return fmt.Errorf("Array mapping structure error -  %s", err.Error())
			}

			if err := arrayMapping.Validate(); err != nil {
				return err
			}
			if err = arrayMapping.DoArrayMapping(inputScope, outputScope, m.resolver); err != nil {
				return fmt.Errorf("Do array mapping error - %s", err.Error())
			}

		}

	}

	return nil
}

func (m *BasicMapper) UpdateMapping() error {
	var newMappingDefs []*data.MappingDef
	for _, mapping := range m.mappings {
		var mappingDef *data.MappingDef
		//Remove all $INPUT for mapTo include array mapping
		if mapping.MapTo != "" && strings.HasPrefix(mapping.MapTo, exprmapper.MAP_TO_INPUT) {
			mappingDef = &data.MappingDef{Type: mapping.Type, Value: mapping.Value, MapTo: exprmapper.RemovePrefixInput(mapping.MapTo)}
		} else {
			mappingDef = mapping
		}

		switch mappingDef.Type {
		//Array mapping
		case data.MtArray:
			//Update Array Mapping
			arrayMapping, err := exprmapper.ParseArrayMapping(mapping.Value)
			if err != nil {
				return fmt.Errorf("Array mapping structure error -  %s", err.Error())
			}

			arrayMapping.RemovePrefixForMapTo()
			v, err := json.Marshal(arrayMapping)
			if err != nil {
				return err
			}
			mappingDef.Value = string(v)
		}
		mapplerLog.Debugf("Updated mapping def %+v", mappingDef)
		newMappingDefs = append(newMappingDefs, mappingDef)
	}
	m.mappings = newMappingDefs
	return nil
}
