package function

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime/debug"
	"strings"

	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/funcexprtype"
)

var logrus = logger.GetLogger("function")

type Func interface {
	Eval(inputScope, outputScope data.Scope) ([]interface{}, error)
	String() string
}

type FunctionExp struct {
	Name   string       `json:"name"`
	Params []*Parameter `json:"params"`
}

type Parameter struct {
	Function *FunctionExp `json:"function"`
	Type     funcexprtype.Type    `json:"type"`
	Value    interface{}  `json:"value"`
}

func (p *Parameter) UnmarshalJSON(paramData []byte) error {
	ser := &struct {
		Function *FunctionExp `json:"function"`
		Type     funcexprtype.Type    `json:"type"`
		Value    interface{}  `json:"value"`
	}{}

	if err := json.Unmarshal(paramData, ser); err != nil {
		return err
	}

	p.Function = ser.Function
	p.Type = ser.Type

	v, err := ConvertToValue(ser.Value, ser.Type)
	if err != nil {
		return err
	}

	p.Value = v

	return nil
}

func (p *Parameter) IsEmtpy() bool {
	if p.Function != nil {
		if p.Function.Name == "" && p.Type == 0 && p.Value == nil && len(p.Function.Params) <= 0 {
			return true
		}
	} else {
		if p.Type == 0 && p.Value == nil {
			return true
		}
	}

	return false
}

func (p *Parameter) IsFunction() bool {
	if funcexprtype.FUNCTION == p.Type {
		return true
	}
	return false
}

func (f *FunctionExp) Eval() ([]interface{}, error) {

	values, err := f.callFunction(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	returns := []interface{}{}
	for _, v := range values {
		returns = append(returns, convertType(v))
	}
	return returns, err
}

func (f *FunctionExp) EvalWithScope(inputScope data.Scope, resolver data.Resolver) ([]interface{}, error) {

	values, err := f.callFunction(nil, inputScope, resolver)
	if err != nil {
		logrus.Errorf("Execution failed for function [%s] error - %+v", f.Name, err.Error())
		return nil, err
	}
	returns := []interface{}{}
	for _, v := range values {
		returns = append(returns, convertType(v))
	}
	return returns, err
}

func (f *FunctionExp) EvalWithData(data interface{}, inputScope data.Scope, resolver data.Resolver) ([]interface{}, error) {

	values, err := f.callFunction(data, inputScope, resolver)
	if err != nil {
		return nil, err
	}
	returns := []interface{}{}
	for _, v := range values {
		returns = append(returns, convertType(v))
	}
	return returns, err
}

func convertType(value reflect.Value) interface{} {
	return value.Interface()
}

func (f *FunctionExp) getRealFunction() (Function, error) {
	return GetFunction(f.Name)
}

func (f *FunctionExp) getMethod() (reflect.Value, error) {
	var ptr reflect.Value
	s, err := f.getRealFunction()
	if err != nil {
		return reflect.Value{}, err
	}

	value := reflect.ValueOf(s)
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(s))
		temp := ptr.Elem()
		temp.Set(value)
	}

	method := value.MethodByName("Eval")
	if method.IsValid() {
		logrus.Debug("valid")
	} else {
		logrus.Debug("invalid")
		method = ptr.MethodByName("Eval")
		if !method.IsValid() {
			logrus.Debug("invalid also, ", f.Name)
			return reflect.Value{}, errors.New("Method invalid..")

		}
	}

	return method, nil
}

func convertToFunctionName(name string) string {
	if name != "" {
		return strings.Title(name)
	}
	return name
}

func (f *FunctionExp) callFunction(fdata interface{}, inputScope data.Scope, resolver data.Resolver) (results []reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			logrus.Debugf("StackTrace: %s", debug.Stack())
		}
	}()

	method, err := f.getMethod()
	if err != nil {
		return nil, err
	}

	inputs := []reflect.Value{}
	for i, p := range f.Params {
		if p.IsFunction() {
			result, err := p.Function.callFunction(fdata, inputScope, resolver)
			if err != nil {
				return nil, err
			}

			for _, v := range result {
				inputs = append(inputs, v)
			}
		} else {
			if !p.IsEmtpy() {
				if p.Type == funcexprtype.REF {
					logrus.Debug("Mapping ref field should done before eval function.")
					var field *ref.MappingRef
					switch p.Value.(type) {
					case string:
						field = ref.NewMappingRef(p.Value.(string))
					case *ref.MappingRef:
						field = p.Value.(*ref.MappingRef)
					}
					//TODO
					if inputScope == nil {
						p.Value = field.GetRef()
					} else {

						v, err := field.Eval(inputScope, resolver)
						if err != nil {
							return nil, err
						}
						p.Value = v

					}

				} else if p.Type == funcexprtype.ARRAYREF {
					logrus.Debug("Mapping ref field should done before eval function.")
					var field *ref.ArrayRef
					switch p.Value.(type) {
					case string:
						//field = &arrayref.MappingRef{p.Value.(string)}
					case *ref.ArrayRef:
						field = p.Value.(*ref.ArrayRef)
					}
					//TODO
					if inputScope == nil {
						p.Value = field.GetRef()
					} else {
						v, err := field.EvalFromData(fdata)
						if err != nil {
							return nil, err
						}
						p.Value = v
					}
				}
				if p.Value != nil {
					inputs = append(inputs, reflect.ValueOf(p.Value))
				} else {
					t := method.Type().In(i)
					funcStr := method.Type().String()
					if strings.Contains(funcStr, "...") {
						parameterNum := method.Type().NumIn()
						if parameterNum > 1 {
							//2. Variadic as latest parameter
							//func(name string, id int, ids ...string)
							if i == parameterNum-1 {
								inputs = append(inputs, reflect.Zero(t.Elem()))
							} else {
								inputs = append(inputs, reflect.Zero(t))
							}
						} else {
							//1. only one variadic parameter
							//func(...string)
							inputs = append(inputs, reflect.Zero(t.Elem()))
						}
					} else {
						inputs = append(inputs, reflect.Zero(t))
					}
				}
			}
		}
	}

	logrus.Debugf("Input Parameters: %+v", inputs)
	values := method.Call(inputs)
	return f.extractErrorFromValues(values)
}

func (f *FunctionExp) extractErrorFromValues(values []reflect.Value) ([]reflect.Value, error) {
	tempValues := []reflect.Value{}

	var err error
	for _, value := range values {
		if value.Type().Name() == "error" {
			if value.Interface() != nil {
				err = value.Interface().(error)
			}
		} else {
			tempValues = append(tempValues, value)
		}
	}

	return tempValues, err
}