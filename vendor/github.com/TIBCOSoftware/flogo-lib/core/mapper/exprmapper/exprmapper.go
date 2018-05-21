package exprmapper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"reflect"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("mapper")

const (
	MAP_TO_INPUT = "$INPUT"
)

func Map(mapping *data.MappingDef, inputScope, outputScope data.Scope, resolver data.Resolver) error {
	mappingValue, err := GetMappingValue(mapping.Value, inputScope, resolver)
	if err != nil {
		return err
	}
	err = setValueToOutputScopde(mapping.MapTo, outputScope, mappingValue, resolver)
	if err != nil {
		err = fmt.Errorf("Set value %+v to output [%s] error - %s", mappingValue, mapping.MapTo, err.Error())
		log.Error(err)
		return err
	}
	log.Debugf("Set value %+v to %s Done", mappingValue, mapping.MapTo)
	return nil
}

func GetMappingValue(mappingV interface{}, inputScope data.Scope, resolver data.Resolver) (interface{}, error) {
	if mappingV == nil || reflect.TypeOf(mappingV).Kind() != reflect.String {
		return mappingV, nil
	}

	mappingValue := mappingV.(string)
	expressionType := expression.GetExpressionType(mappingValue)
	if expressionType == expression.TERNARY_EXPRESSION {
		exp, err := expression.NewExpression(mappingValue).GetTernaryExpression()
		if err != nil {
			return nil, fmt.Errorf("Parsing ternary expression [%s] error - %s", mappingValue, err.Error())
		}

		funcValue, err := exp.EvalWithScope(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", mappingValue, err.Error())
		}
		log.Debugf("Ternary expression value: %+v", funcValue)
		return funcValue, nil
	} else if expressionType == expression.EXPRESSION {
		exp, err := expression.NewExpression(mappingValue).GetExpression()
		if err != nil {
			return nil, fmt.Errorf("Parsing expression [%s] error - %s", mappingValue, err.Error())
		}

		funcValue, err := exp.EvalWithScope(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", mappingValue, err.Error())
		}
		log.Debugf("Expression value: %+v", funcValue)
		return funcValue, nil

	} else if expressionType == expression.FUNCTION {
		log.Debugf("The mapping ref is a function")
		function, err := expression.NewFunctionExpression(mappingValue).GetFunction()
		if err != nil {
			return nil, fmt.Errorf("Parsing function [%s] error - %s", mappingValue, err.Error())
		}
		funcValue, err := function.EvalWithScope(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", mappingValue, err.Error())
		}

		if funcValue != nil && len(funcValue) == 1 {
			return funcValue[0], nil

		} else if funcValue != nil && len(funcValue) > 1 {
			return funcValue, nil
		}

	} else if !isMappingRef(mappingValue) {
		log.Debugf("Mapping value is literal set directly to field")
		log.Debugf("Mapping ref %s and value %+v", mappingValue, mappingValue)
		return mappingValue, nil
	} else {

		mappingref := ref.NewMappingRef(mappingValue)
		mappingValue, err := mappingref.GetValue(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Get value from ref [%s] error - %s", mappingref.GetRef(), err.Error())

		}
		log.Debugf("Mapping ref %s and value %+v", mappingValue, mappingValue)
		return mappingValue, nil
	}

	return nil, nil
}

func setValueToOutputScopde(mapTo string, outputScope data.Scope, value interface{}, resolver data.Resolver) error {
	toMappingRef := ref.NewMappingRef(mapTo)
	actRootField, err := toMappingRef.GetActivtyRootField()
	if err != nil {
		return err
	}
	if field.HasSpecialFields(mapTo) {
		fields, err := field.GetAllspecialFields(mapTo)
		if err != nil {
			return fmt.Errorf("Get fields from field %s error, due to [%s]", mapTo, err.Error())
		}
		if len(fields) == 1 {
			//No complex mapping exist
			return SetAttribute(actRootField, value, outputScope)
		} else if len(fields) > 1 {
			//Complex mapping
			return settValueToComplexObject(toMappingRef, actRootField, outputScope, value)
		}
		return fmt.Errorf("No field name found for mapTo [%s]", mapTo)
	}

	if strings.HasPrefix(mapTo, "$") || strings.Index(mapTo, ".") > 0 {
		return settValueToComplexObject(toMappingRef, actRootField, outputScope, value)
	}
	return SetAttribute(mapTo, value, outputScope)
}

func settValueToComplexObject(toMappingRef *ref.MappingRef, fieldName string, outputScope data.Scope, value interface{}) error {
	complexVlaueIn, err := toMappingRef.GetValueFromOutputScope(outputScope)
	if err != nil {
		return err
	}
	fields, err := toMappingRef.GetFields()
	if err != nil {
		return err
	}

	log.Debugf("Set value %+v to fields %s", value, fields)
	complexValue, err2 := json.SetFieldValue(value, complexVlaueIn, fields)
	if err2 != nil {
		return err2
	}

	return SetAttribute(fieldName, complexValue, outputScope)
}

func isMappingRef(mappingref string) bool {
	if mappingref == "" || !strings.HasPrefix(mappingref, "$") {
		return false
	}
	return true
}

func SetAttribute(fieldName string, value interface{}, outputScope data.Scope) error {
	//Set Attribute value back to attribute
	attribute, exist := outputScope.GetAttr(fieldName)
	if exist {
		switch attribute.Type() {
		case data.TypeComplexObject:
			complexObject := attribute.Value().(*data.ComplexObject)
			newComplexObject := &data.ComplexObject{Metadata: complexObject.Metadata, Value: value}
			outputScope.SetAttrValue(fieldName, newComplexObject)
		default:
			outputScope.SetAttrValue(fieldName, value)
		}

	} else {
		return errors.New("Cannot found attribute " + fieldName + " at output scope")
	}
	return nil
}


func RemovePrefixInput(str string) string {
	if str != "" && strings.HasPrefix(str, MAP_TO_INPUT) {
		//Remove $INPUT for mapTo
		newMapTo := str[len(MAP_TO_INPUT):]
		if strings.HasPrefix(newMapTo, ".") {
			newMapTo = newMapTo[1:]
		}
		str = newMapTo
	}
	return str
}

