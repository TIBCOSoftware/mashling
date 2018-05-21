package data

import (
	"encoding/json"
	"errors"
)

// MappingType is an enum for possible MappingDef Types
type MappingType int

const (
	// MtAssign denotes an attribute to attribute assignment
	MtAssign MappingType = 1

	// MtLiteral denotes a literal to attribute assignment
	MtLiteral MappingType = 2

	// MtExpression denotes a expression execution to perform mapping
	MtExpression MappingType = 3

	// MtObject denotes a object construction mapping
	MtObject MappingType = 4

	MtArray MappingType = 5
)

// MappingDef is a simple structure that defines a mapping
type MappingDef struct {
	//Type the mapping type
	Type MappingType

	//Value the mapping value to execute to determine the result (rhs)
	Value interface{}

	//Result the name of attribute to place the result of the mapping in (lhs)
	MapTo string
}

// Mapper is a constructs that maps values from one scope to another
type Mapper interface {
	Apply(inputScope Scope, outputScope Scope) error
}

// MapperDef represents a Mapper, which is a collection of mappings
type MapperDef struct {
	//todo possibly add optional lang/mapper type so we can fast fail on unsupported mappings/mapper combo
	Mappings []*MappingDef
}

type IOMappings struct {
	Input  []*MappingDef `json:"input,omitempty"`
	Output []*MappingDef `json:"output,omitempty"`
}

func (md *MappingDef) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Type  interface{} `json:"type"`
		Value interface{} `json:"value"`
		MapTo string      `json:"mapTo"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	md.MapTo = ser.MapTo
	md.Value = ser.Value
	intType, err := ConvertMappingType(ser.Type)
	if err == nil {
		md.Type = intType
	}
	return err
}

func ConvertMappingType(mapType interface{}) (MappingType, error) {
	strType, _ := CoerceToString(mapType)
	switch strType {
	case "assign", "1":
		return MtAssign, nil
	case "literal", "2":
		return MtLiteral, nil
	case "expression", "3":
		return MtExpression, nil
	case "object", "4":
		return MtObject, nil
	case "array", "5":
		return MtArray, nil
	default:
		return 0, errors.New("unsupported mapping type: " + strType)
	}
}
