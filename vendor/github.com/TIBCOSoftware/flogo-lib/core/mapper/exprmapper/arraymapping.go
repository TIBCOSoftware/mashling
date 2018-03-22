package exprmapper

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	flogojson "github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var arraylog = logger.GetLogger("array-mapping")

const (
	PRIMITIVE = "primitive"
	FOREACH   = "foreach"
	NEWARRAY  = "NEWARRAY"
)

type ArrayMapping struct {
	From   interface{}     `json:"from"`
	To     string          `json:"to"`
	Type   string          `json:"type"`
	Fields []*ArrayMapping `json:"fields,omitempty"`
}

func (a *ArrayMapping) Validate() error {
	//Validate root from/to field
	if a.From == nil {
		return fmt.Errorf("The array mapping validation failed for the mapping [%s]. Ensure valid array is mapped in the mapper. ", a.To)
	}

	if a.To == "" || len(a.To) <= 0 {
		return fmt.Errorf("The array mapping validation failed for the mapping [%s]. Ensure valid array is mapped in the mapper. ", a.From)
	}

	if a.Type == FOREACH {
		//Validate root from/to field
		if a.From == NEWARRAY {
			//Make sure no array ref fields exist
			for _, field := range a.Fields {
				if field.Type == FOREACH {
					return field.Validate()
				}
				stringVal, ok := field.From.(string)
				if ok && ref.IsArrayMapping(stringVal) {
					return fmt.Errorf("The array mapping validation failed, due to invalid new array mapping [%s]", stringVal)
				}

			}
		} else {
			for _, field := range a.Fields {
				if field.Type == FOREACH {
					return field.Validate()
				}
			}
		}

	}

	return nil
}

func (a *ArrayMapping) DoArrayMapping(inputScope, outputScope data.Scope, resolver data.Resolver) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			logger.Debugf("StackTrace: %s", debug.Stack())
		}
	}()
	switch a.Type {
	case PRIMITIVE:
		mappingDef := a.mappingDef()
		err := Map(mappingDef, inputScope, outputScope, resolver)
		if err != nil {
			return err
		}
	case FOREACH:
		//First Level
		toRef := ref.NewMappingRef(a.To)
		var fromValue interface{}
		var err error

		//TODO this might never be call.. try to delete
		stringVal, _ := a.From.(string)
		if expression.IsExpression(stringVal) {
			exp, err := expression.NewExpression(stringVal).GetExpression()
			if err != nil {
				log.Errorf("New expression from %s error: %s", stringVal, err.Error())
				return err
			}

			fromValue, err = exp.EvalWithScope(inputScope, resolver)
			if err != nil {
				log.Errorf("Eval expression from scope error: %s", err.Error())
				return err
			}

		} else if expression.IsFunction(stringVal) {
			log.Debugf("The mapping ref is a function")
			function, err := expression.NewFunctionExpression(stringVal).GetFunction()
			if err != nil {
				log.Errorf("New function from %s error: %s", stringVal, err.Error())
				return err
			}
			log.Debugf("Function is:%+v", function)
			funcValue, err := function.EvalWithScope(inputScope, resolver)
			if err != nil {
				log.Errorf("Eval function error %s", err.Error())
				return err
			}

			if funcValue != nil && len(funcValue) == 1 {
				fromValue = funcValue[0]
			} else if funcValue != nil && len(funcValue) > 1 {
				fromValue = funcValue[0]
			}

		} else {
			if strings.EqualFold(stringVal, NEWARRAY) {
				log.Debugf("Init a new array for field", a.To)
				fromValue = make([]interface{}, 1)
			} else {
				fromRef := ref.NewMappingRef(stringVal)
				fromValue, err = fromRef.GetValue(inputScope, resolver)
				if err != nil {
					return err
				}
			}
		}
		//Check if fields is empty for primitive array mapping
		if a.Fields == nil || len(a.Fields) <= 0 {
			//Set value directlly to MapTo field
			return setValueToOutputScopde(a.To, outputScope, fromValue, resolver)
		}

		//Loop array
		fromArrayvalues, ok := fromValue.([]interface{})
		if !ok {
			return fmt.Errorf("Failed to get array value from [%s], due to error- [%s] value not an array", a.From, a.From)
		}

		toValue, err := toRef.GetValueFromOutputScope(outputScope)
		if err != nil {
			return err
		}
		toValue = toInterface(toValue)
		objArray := make([]interface{}, len(fromArrayvalues))
		for i, _ := range objArray {
			objArray[i] = make(map[string]interface{})
		}

		mappingField, err := toRef.GetFields()
		if err != nil {
			return fmt.Errorf("Get fields from mapping string error, due to [%s]", err.Error())
		}
		if mappingField != nil && len(mappingField.Fields) > 0 {
			vv, err := flogojson.SetFieldValue(objArray, toValue, mappingField)
			if err != nil {
				return err
			}
			log.Debugf("Set Value return as %+v", vv)
		} else {
			toValue = objArray
		}

		if err != nil {
			return err
		}

		for i, arrayV := range fromArrayvalues {
			err = a.runArrayMap(arrayV, objArray[i], a.Fields, inputScope, outputScope, resolver)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		fieldName, err := toRef.GetFieldName()
		if err != nil {
			return err
		}
		//Get Value from fields
		toFieldName, err := toRef.GetFieldName()
		if err != nil {
			return err
		}

		if len(mappingField.Fields) > 0 {
			return SetAttribute(fieldName, toValue, outputScope)
		}
		return SetAttribute(fieldName, getFieldValue(toValue, toFieldName), outputScope)
	}
	return nil
}

func (a *ArrayMapping) mappingDef() *data.MappingDef {
	return &data.MappingDef{MapTo: a.To, Value: a.From, Type: data.MtExpression}
}

func (a *ArrayMapping) runArrayMap(fromValue, value interface{}, fields []*ArrayMapping, inputScope, outputScope data.Scope, resolver data.Resolver) error {
	for _, field := range fields {
		if field.Type == PRIMITIVE {
			fValue, err := GetValueFromArrayRef(fromValue, field.From, inputScope, resolver)
			if err != nil {
				return err
			}
			log.Debugf("Array mapping from %s 's value %+v", field.From, fValue)
			err = field.DoMap(fValue, value, ref.GetFieldNameFromArrayRef(field.To), inputScope, outputScope, resolver)
			if err != nil {
				return err
			}
		} else if field.Type == FOREACH {
			var fromArrayvalues []interface{}
			if strings.EqualFold(field.From.(string), NEWARRAY) {
				log.Debugf("Init a new array for field", field.To)
				fromArrayvalues = make([]interface{}, 1)
			} else {
				fValue, err := GetValueFromArrayRef(fromValue, field.From, inputScope, resolver)
				if err != nil {
					return err
				}
				var ok bool
				fromArrayvalues, ok = fValue.([]interface{})
				if !ok {
					return fmt.Errorf("Failed to get array value from [%s], due to error- value not an array", fValue)
				}
			}

			toValue := toInterface(value)
			objArray := make([]interface{}, len(fromArrayvalues))
			for i, _ := range objArray {
				objArray[i] = make(map[string]interface{})
			}

			_, err := flogojson.SetFieldValueP(objArray, toValue, ref.GetFieldNameFromArrayRef(field.To))
			if err != nil {
				return err
			}
			//Check if fields is empty for primitive array mapping
			if field.Fields == nil || len(field.Fields) <= 0 {
				for f, v := range fromArrayvalues {
					objArray[f] = v
				}
				continue
			}

			for i, arrayV := range fromArrayvalues {
				err = a.runArrayMap(arrayV, objArray[i], field.Fields, inputScope, outputScope, resolver)
				if err != nil {
					return err
				}
			}
		}

	}

	return nil

}

func (a *ArrayMapping) DoMap(fromValue, value interface{}, to string, inputScope, outputScope data.Scope, resolver data.Resolver) error {
	switch a.Type {
	case PRIMITIVE:
		_, err := flogojson.SetFieldValueP(fromValue, value, to)
		if err != nil {
			return err
		}
	case FOREACH:
		fValue, err := GetValueFromArrayRef(fromValue, a.From, inputScope, resolver)
		if err != nil {
			return err
		}
		tValue, err := GetValueFromArrayRef(value, a.To, inputScope, resolver)
		if err != nil {
			return err
		}
		err = a.runArrayMap(fValue, tValue, a.Fields, inputScope, outputScope, resolver)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ArrayMapping) RemovePrefixForMapTo() {
	if a == nil {
		return
	}

	a.To = RemovePrefixInput(a.To)

	if a.Type == FOREACH {
		//Validate root from/to field
		if a.From == NEWARRAY {
			//Make sure no array ref fields exist
			for _, field := range a.Fields {
				if field.Type == FOREACH {
					field.RemovePrefixForMapTo()
				} else {
					field.To = RemovePrefixInput(field.To)
				}
			}

		} else {
			for _, field := range a.Fields {
				if field.Type == FOREACH {
					field.RemovePrefixForMapTo()
				} else {
					field.To = RemovePrefixInput(field.To)
				}
			}
		}

	}
}

func ParseArrayMapping(arrayDatadata interface{}) (*ArrayMapping, error) {
	amapping := &ArrayMapping{}
	switch t := arrayDatadata.(type) {
	case string:
		err := json.Unmarshal([]byte(t), amapping)
		if err != nil {
			return nil, err
		}
	case interface{}:
		s, err := data.CoerceToString(t)
		if err != nil {
			return nil, fmt.Errorf("Convert array mapping value to string error, due to [%s]", err.Error())
		}
		err = json.Unmarshal([]byte(s), amapping)
		if err != nil {
			return nil, err
		}
	}
	return amapping, nil
}

func toInterface(data interface{}) interface{} {

	switch t := data.(type) {
	case string:
		if strings.EqualFold("{}", t) {
			return make(map[string]interface{})
		}
	default:
		if t == nil {
			//TODO maybe consider other types as well
			return make(map[string]interface{})
		}
	}
	return data
}

func getFieldValue(value interface{}, fieldName string) interface{} {
	switch t := value.(type) {
	case map[string]interface{}:
		return t[fieldName]
	default:
		return value
	}
	return value
}

func GetValueFromArrayRef(object interface{}, expressionRef interface{}, inputScope data.Scope, resolver data.Resolver) (interface{}, error) {

	var fromValue interface{}
	var err error

	stringVal, ok := expressionRef.(string)

	if !ok {
		//Non string value
		return expressionRef, nil
	}
	if expression.IsTernaryExpression(stringVal) {
		exp, err := expression.NewExpression(stringVal).GetTernaryExpression()
		if err != nil {
			return nil, fmt.Errorf("Parsing ternary expression [%s] error - %s", stringVal, err.Error())
		}

		funcValue, err := exp.EvalWithScope(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", stringVal, err.Error())
		}
		log.Debugf("Ternary expression value: %+v", funcValue)
		return funcValue, nil
	} else if expression.IsExpression(stringVal) {
		exp, err := expression.NewExpression(stringVal).GetExpression()
		if err != nil {
			log.Errorf("Parsing expression from [%s] error - %s", stringVal, err.Error())
			return nil, err
		}

		fromValue, err = exp.EvalWithScope(inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", stringVal, err.Error())
		}

	} else if expression.IsFunction(stringVal) {
		log.Debugf("The mapping ref is a function")
		function, err := expression.NewFunctionExpression(stringVal).GetFunction()
		if err != nil {
			log.Errorf("Parsing function from [%s] error - %s", stringVal, err.Error())
			return nil, err
		}
		funcValue, err := function.EvalWithData(object, inputScope, resolver)
		if err != nil {
			return nil, fmt.Errorf("Execution failed for mapping [%s] due to error - %s", stringVal, err.Error())
		}

		if funcValue != nil && len(funcValue) == 1 {
			fromValue = funcValue[0]
		} else if funcValue != nil && len(funcValue) > 1 {
			fromValue = funcValue[0]
		}

	} else {
		if ref.IsArrayMapping(stringVal) {
			reference := ref.GetFieldNameFromArrayRef(stringVal)
			fromValue, err = flogojson.GetFieldValueFromInP(object, reference)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(stringVal, "$") {
			fromRef := ref.NewMappingRef(stringVal)
			fromValue, err = fromRef.GetValue(inputScope, resolver)
			if err != nil {
				return nil, fmt.Errorf("Get value from [%s] failed, due to error - %s", stringVal, err.Error())
			}
		} else {
			fromValue = expressionRef
		}

	}

	return fromValue, err

}
