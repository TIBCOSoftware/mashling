package data

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// CoerceToValue coerce a value to the specified type
func CoerceToValue(value interface{}, dataType Type) (interface{}, error) {

	var coerced interface{}
	var err error

	switch dataType {
	case TypeAny:
		coerced, err = CoerceToAny(value)
	case TypeString:
		coerced, err = CoerceToString(value)
	case TypeInteger:
		coerced, err = CoerceToInteger(value)
	case TypeLong:
		coerced, err = CoerceToLong(value)
	case TypeDouble:
		coerced, err = CoerceToDouble(value)
	case TypeBoolean:
		coerced, err = CoerceToBoolean(value)
	case TypeObject:
		coerced, err = CoerceToObject(value)
	case TypeComplexObject:
		coerced, err = CoerceToComplexObject(value)
	case TypeArray:
		coerced, err = CoerceToArrayIfNecessary(value)
	case TypeParams:
		coerced, err = CoerceToParams(value)
	}

	if err != nil {
		return nil, err
	}

	return coerced, nil
}

// CoerceToString coerce a value to a string
func CoerceToString(val interface{}) (string, error) {

	switch t := val.(type) {
	case string:
		return t, nil
	case int:
		return strconv.Itoa(t), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case json.Number:
		return t.String(), nil
	case bool:
		return strconv.FormatBool(t), nil
	case nil:
		return "", nil
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return "", fmt.Errorf("unable to Coerce %#v to string", t)
		}
		return string(b), nil
	}
}

// CoerceToInteger coerce a value to an integer
func CoerceToInteger(val interface{}) (int, error) {
	switch t := val.(type) {
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float64:
		return int(t), nil
	case json.Number:
		i, err := t.Int64()
		return int(i), err
	case string:
		return strconv.Atoi(t)
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to coerce %#v to integer", val)
	}
}

// CoerceToInteger coerce a value to an integer
func CoerceToLong(val interface{}) (int64, error) {
	switch t := val.(type) {
	case int:
		return int64(t), nil
	case int64:
		return t, nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case json.Number:
		return t.Int64()
	case string:
		return strconv.ParseInt(t, 10, 64)
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("Unable to coerce %#v to integer", val)
	}
}

// Deprecated: Use CoerceToDouble()
func CoerceToNumber(val interface{}) (float64, error) {
	return CoerceToDouble(val)
}

// CoerceToDouble coerce a value to a double/float64
func CoerceToDouble(val interface{}) (float64, error) {
	switch t := val.(type) {
	case int:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float64:
		return t, nil
	case json.Number:
		return t.Float64()
	case string:
		return strconv.ParseFloat(t, 64)
	case bool:
		if t {
			return 1.0, nil
		}
		return 0.0, nil
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("Unable to coerce %#v to float", val)
	}
}

// CoerceToBoolean coerce a value to a boolean
func CoerceToBoolean(val interface{}) (bool, error) {
	switch t := val.(type) {
	case bool:
		return t, nil
	case int, int64:
		return t != 0, nil
	case float64:
		return t != 0.0, nil
	case json.Number:
		i, err := t.Int64()
		return i != 0, err
	case string:
		return strconv.ParseBool(t)
	case nil:
		return false, nil
	default:
		str, err := CoerceToString(val)
		if err != nil {
			return false, fmt.Errorf("unable to coerce %#v to bool", val)
		}
		return strconv.ParseBool(str)
	}
}

// CoerceToObject coerce a value to an object
func CoerceToObject(val interface{}) (map[string]interface{}, error) {

	switch t := val.(type) {
	case map[string]interface{}:
		return t, nil
	case string:
		m := make(map[string]interface{})
		if t != "" {
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", val)
			}
		}
		return m, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", val)
	}
}

// CoerceToArray coerce a value to an array of empty interface values
func CoerceToArray(val interface{}) ([]interface{}, error) {

	switch t := val.(type) {
	case []interface{}:
		return t, nil

	case []map[string]interface{}:
		var a []interface{}
		for _, v := range t {
			a = append(a, v)
		}
		return a, nil
	case string:
		a := make([]interface{}, 0)
		if t != "" {
			err := json.Unmarshal([]byte(t), &a)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", val)
			}
		}
		return a, nil
	case nil:
		return nil, nil
	default:
		s := reflect.ValueOf(val)
		if s.Kind() == reflect.Slice {
			a := make([]interface{}, s.Len())

			for i := 0; i < s.Len(); i++ {
				a[i] = s.Index(i).Interface()
			}
			return a, nil
		}
		return nil, fmt.Errorf("unable to coerce %#v to []interface{}", val)
	}
}

// CoerceToArrayIfNecessary coerce a value to an array if it isn't one already
func CoerceToArrayIfNecessary(val interface{}) (interface{}, error) {

	if val == nil {
		return nil, nil
	}

	rt := reflect.TypeOf(val).Kind()

	if rt == reflect.Array || rt == reflect.Slice {
		return val, nil
	}

	switch t := val.(type) {
	case string:
		a := make([]interface{}, 0)
		if t != "" {
			err := json.Unmarshal([]byte(t), &a)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to map[string]interface{}", val)
			}
		}
		return a, nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to []interface{}", val)
	}
}

// CoerceToAny coerce a value to generic value
func CoerceToAny(val interface{}) (interface{}, error) {

	switch t := val.(type) {

	case json.Number:
		if strings.Contains(t.String(), ".") {
			return t.Float64()
		} else {
			return t.Int64()
		}
	default:
		return val, nil
	}
}

// CoerceToParams coerce a value to params
func CoerceToParams(val interface{}) (map[string]string, error) {

	switch t := val.(type) {
	case map[string]string:
		return t, nil
	case string:
		m := make(map[string]string)
		if t != "" {
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce %#v to params", val)
			}
		}
		return m, nil
	case map[string]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mVal, err := CoerceToString(v)
			if err != nil {
				return nil, err
			}
			m[k] = mVal
		}
		return m, nil
	case map[interface{}]string:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey, err := CoerceToString(k)
			if err != nil {
				return nil, err
			}
			m[mKey] = v
		}
		return m, nil
	case map[interface{}]interface{}:

		var m = make(map[string]string, len(t))
		for k, v := range t {

			mKey, err := CoerceToString(k)
			if err != nil {
				return nil, err
			}

			mVal, err := CoerceToString(v)
			if err != nil {
				return nil, err
			}
			m[mKey] = mVal
		}
		return m, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to map[string]string", val)
	}
}

// CoerceToObject coerce a value to an complex object
func CoerceToComplexObject(val interface{}) (*ComplexObject, error) {
	//If the val is nil then just return empty struct
	var emptyComplexObject = &ComplexObject{Value: "{}"}
	if val == nil {
		return emptyComplexObject, nil
	}
	switch t := val.(type) {
	case string:
		if val == "" {
			return emptyComplexObject, nil
		} else {
			complexObject := &ComplexObject{}
			err := json.Unmarshal([]byte(t), complexObject)
			if err != nil {
				return nil, err

			}
			return handleComplex(complexObject), nil
		}
	case map[string]interface{}:
		v, err := json.Marshal(val)
		if err != nil {
			return nil, err
		}
		complexObject := &ComplexObject{}
		err = json.Unmarshal(v, complexObject)
		if err != nil {
			return nil, err
		}
		return handleComplex(complexObject), nil
	case *ComplexObject:
		return handleComplex(val.(*ComplexObject)), nil
	default:
		return nil, fmt.Errorf("unable to coerce %#v to complex object", val)
	}
}

func handleComplex(complex *ComplexObject) *ComplexObject {
	if complex != nil {
		if complex.Value == "" {
			complex.Value = "{}"
		}
	}
	return complex
}

//var mapHelper *MapHelper = &MapHelper{}
//
//func GetMapHelper() *MapHelper {
//	return mapHelper
//}
//
//type MapHelper struct {
//}
//
//func (h *MapHelper) GetInt(data map[string]interface{}, key string) (int, bool) {
//	mapVal, exists := data[key]
//	if exists {
//		value, ok := mapVal.(int)
//
//		if ok {
//			return value, true
//		}
//	}
//
//	return 0, false
//}
//
//func (h *MapHelper) GetString(data map[string]interface{}, key string) (string, bool) {
//	mapVal, exists := data[key]
//	if exists {
//		value, ok := mapVal.(string)
//
//		if ok {
//			return value, true
//		}
//	}
//
//	return "", false
//}
//
//func (h *MapHelper) GetBool(data map[string]interface{}, key string) (bool, bool) {
//	mapVal, exists := data[key]
//	if exists {
//		value, ok := mapVal.(bool)
//
//		if ok {
//			return value, true
//		}
//	}
//
//	return false, false
//}
//
//func (h *MapHelper) ToAttributes(data map[string]interface{}, metadata []*Attribute, ignoreExtras bool) []*Attribute {
//
//	size := len(metadata)
//	if !ignoreExtras {
//		size = len(data)
//	}
//	attrs := make([]*Attribute, 0, size)
//
//	metadataMap := make(map[string]*Attribute)
//	for _, attr := range metadata {
//		metadataMap[attr.Name()] = attr
//	}
//
//	//todo do special handling for complex_object metadata (merge or ref it)
//	for key, value := range data {
//		mdAttr, exists := metadataMap[key]
//
//		if !exists {
//			if !ignoreExtras {
//				//todo handle error
//				attr, _ := NewAttribute(key, TypeAny, value)
//				attrs = append(attrs, attr)
//			}
//		} else {
//			attr, _ := NewAttribute(key, mdAttr.Type(), value)
//			attrs = append(attrs, attr)
//		}
//	}
//
//	return attrs
//}
