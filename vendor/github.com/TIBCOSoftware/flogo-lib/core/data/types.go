package data

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Type denotes a data type
type Type int

const (
	TypeAny           Type = iota
	TypeString
	TypeInteger
	TypeLong
	TypeDouble
	TypeBoolean
	TypeObject
	TypeComplexObject
	TypeArray
	TypeParams

)

var types = [...]string{
	"any",
	"string",
	"integer",
	"long",
	"double",
	"boolean",
	"object",
	"complexObject",
	"array",
	"params",
}

func (t Type) String() string {
	return types[t]
}

// ToTypeEnum get the data type that corresponds to the specified name
func ToTypeEnum(typeStr string) (Type, bool) {

	switch strings.ToLower(typeStr) {
	case "any":
		return TypeAny, true
	case "string":
		return TypeString, true
	case "integer", "int":
		return TypeInteger, true
	case "long":
		return TypeLong, true
	case "double", "number":
		return TypeDouble, true
	case "boolean", "bool":
		return TypeBoolean, true
	case "object":
		return TypeObject, true
	case "complexobject", "complex_object":
		return TypeComplexObject, true
	case "array":
		return TypeArray, true
	case "params":
		return TypeParams, true
	default:
		return TypeAny, false
	}
}

// GetType get the Type of the supplied value
func GetType(val interface{}) (Type, error) {

	switch t := val.(type) {
	case string:
		return TypeString, nil
	case int, int32:
		return TypeInteger, nil
	case int64:
		return TypeLong, nil
	case float64:
		return TypeDouble, nil
	case json.Number:
		if strings.Contains(t.String(), ".") {
			return TypeDouble, nil
		} else {
			return TypeLong, nil
		}
	case bool:
		return TypeBoolean, nil
	case map[string]interface{}:
		return TypeObject, nil
	case ComplexObject:
		return TypeComplexObject, nil
	case []interface{}:
		return TypeArray, nil
	case map[string]string:
		return TypeParams, nil
	default:
		return TypeAny, fmt.Errorf("unable to determine type of %#v", t)
	}
}

func IsSimpleType(val interface{}) bool {

	switch val.(type) {
	case string, int, int32, float32, float64, json.Number, bool:
		return true
	default:
		return false
	}
}
