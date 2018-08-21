package function

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime/debug"
	"strings"

	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/funcexprtype"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/ref"
	"github.com/TIBCOSoftware/flogo-lib/logger"
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
	Function *FunctionExp      `json:"function"`
	Type     funcexprtype.Type `json:"type"`
	Value    interface{}       `json:"value"`
}

func (p *Parameter) UnmarshalJSON(paramData []byte) error {
	ser := &struct {
		Function *FunctionExp      `json:"function"`
		Type     funcexprtype.Type `json:"type"`
		Value    interface{}       `json:"value"`
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

func (f *FunctionExp) Eval() (interface{}, error) {
	value, err := f.callFunction(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return convertType(value), err
}

func (f *FunctionExp) EvalWithScope(inputScope data.Scope, resolver data.Resolver) (interface{}, error) {

	value, err := f.callFunction(nil, inputScope, resolver)
	if err != nil {
		logrus.Errorf("Execution failed for function [%s] error - %+v", f.Name, err.Error())
		return nil, err
	}
	return convertType(value), err
}

func (f *FunctionExp) EvalWithData(data interface{}, inputScope data.Scope, resolver data.Resolver) (interface{}, error) {
	value, err := f.callFunction(data, inputScope, resolver)
	if err != nil {
		return nil, err
	}
	return convertType(value), err
}

func HandleToSingleOutput(values interface{}) interface{} {
	if values != nil {
		switch t := values.(type) {
		case []interface{}:
			return t[0]
		default:
			return t
		}
	}
	return nil
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

func (f *FunctionExp) callFunction(fdata interface{}, inputScope data.Scope, resolver data.Resolver) (results reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%+v", r)
			logrus.Debugf("StackTrace: %s", debug.Stack())
		}
	}()

	method, err := f.getMethod()
	if err != nil {
		return reflect.Value{}, err
	}

	inputs := []reflect.Value{}
	for i, p := range f.Params {
		if p.IsFunction() {
			result, err := p.Function.callFunction(fdata, inputScope, resolver)
			if err != nil {
				return reflect.Value{}, err
			}
			inputs = append(inputs, result)
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
							return reflect.Value{}, err
						}
						p.Value = v
					}

				} else if p.Type == funcexprtype.ARRAYREF {
					logrus.Debug("Mapping ref field should done before eval function.")
					var field *ref.ArrayRef
					switch p.Value.(type) {
					case string:
						field = ref.NewArrayRef(p.Value.(string))
					case *ref.ArrayRef:
						field = p.Value.(*ref.ArrayRef)
					}
					if inputScope == nil {
						p.Value = field.GetRef()
					} else {
						if fdata == nil {
							//Array mapping should not go here for today, take is as get current scope.
							//TODO how to know it is array mapping or get current scope
							ref := ref.NewMappingRef(field.GetRef())
							v, err := ref.Eval(inputScope, resolver)
							if err != nil {
								return reflect.Value{}, err
							}
							p.Value = v

						} else {
							v, err := field.EvalFromData(fdata)
							if err != nil {
								return reflect.Value{}, err
							}
							p.Value = v

						}

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
	args, err := ensureArguments(method, inputs)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("Function '%s' argument validation failed due to error %s", f.Name, err.Error())
	}
	values := method.Call(args)
	return f.extractErrorFromValues(values)
}

func ensureArguments(method reflect.Value, in []reflect.Value) ([]reflect.Value, error) {

	var retInputs []reflect.Value
	methodType := method.Type()
	n := method.Type().NumIn()
	for i := 0; i < n; i++ {
		if xt, targ := in[i].Type(), methodType.In(i); !xt.AssignableTo(targ) {
			v, err := convertArgs(targ, in[i])
			if err != nil {
				return nil, fmt.Errorf("argument type mismatch. Can not convert type %s to type %s. ", xt.String(), targ.String())
			}
			retInputs = append(retInputs, reflect.ValueOf(v))
		} else {
			retInputs = append(retInputs, in[i])
		}
	}

	if methodType.IsVariadic() {
		m := len(in) - n
		elem := methodType.In(n - 1).Elem()
		for j := 0; j < m; j++ {
			x := in[n+j]
			if xt := x.Type(); !xt.AssignableTo(elem) {
				v, err := convertArgs(elem, x)
				if err != nil {
					return nil, fmt.Errorf("argument type mismatch. Can not convert type %s to type %s. ", xt.String(), elem.String())
				}
				retInputs = append(retInputs, reflect.ValueOf(v))
			} else {
				retInputs = append(retInputs, x)
			}
		}
	}
	return retInputs, nil
}

func convertArgs(argType reflect.Type, in reflect.Value) (interface{}, error) {
	var v interface{}
	var err error
	switch argType.Kind() {
	case reflect.Bool:
		v, err = data.CoerceToBoolean(in.Interface())
	case reflect.Interface:
		v, err = data.CoerceToAny(in.Interface())
	case reflect.Int:
		v, err = data.CoerceToInteger(in.Interface())
	case reflect.Int64:
		v, err = data.CoerceToLong(in.Interface())
	case reflect.String:
		v, err = data.CoerceToString(in.Interface())
	case reflect.Float64:
		v, err = data.CoerceToDouble(in.Interface())
	default:
		v = in.Interface()
	}
	return v, err

}

func (f *FunctionExp) extractErrorFromValues(values []reflect.Value) (reflect.Value, error) {
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

	if len(tempValues) > 1 {
		return tempValues[0], fmt.Errorf("Not support function multiple returns")
	}
	return tempValues[0], err
}
