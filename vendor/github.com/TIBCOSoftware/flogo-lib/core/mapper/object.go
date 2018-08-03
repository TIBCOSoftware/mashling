package mapper

import (
	"errors"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper"
)

func MapObject(mappingObject interface{}, scope data.Scope, resolver data.Resolver) (result interface{}, err error) {

	complexObject, ok := mappingObject.(data.ComplexObject)

	if ok {
		mappingObject = complexObject.Value
	}

	switch t := mappingObject.(type) {
	case map[string]interface{}:
		return dm(t, scope, resolver)
	case []interface{}:
		return da(t, scope, resolver)
	default:
		return nil, errors.New("unsupported mapping object")
	}
}

func dm(mappingObject map[string]interface{}, scope data.Scope, resolver data.Resolver) (result map[string]interface{}, err error) {

	result = make(map[string]interface{}, len(mappingObject))

	for key, value := range mappingObject {

		switch t := value.(type) {
		case string:

			if strings.HasPrefix(t, "{{") {
				result[key], err = evalExpr(t, scope, resolver)
			} else {
				result[key] = t
			}
		case map[string]interface{}:
			result[key], err = dm(t, scope, resolver)
		case []interface{}:
			result[key], err = da(t, scope, resolver)
		default:
			result[key] = t
		}
	}

	return result, err
}

func da(mappingArr []interface{}, scope data.Scope, resolver data.Resolver) (result []interface{}, err error) {

	result = make([]interface{}, len(mappingArr))

	for idx, value := range mappingArr {

		switch t := value.(type) {
		case string:

			if strings.HasPrefix(t, "{{") {
				result[idx], err = evalExpr(t, scope, resolver)
			} else {
				result[idx] = t
			}
		case map[string]interface{}:
			result[idx], err = dm(t, scope, resolver)
		case []interface{}:
			result[idx], err = da(t, scope, resolver)
		default:
			result[idx] = t
		}
	}

	return result, err
}

func evalExpr(exprString string, scope data.Scope, resolver data.Resolver) (interface{}, error) {

	exprStr := exprString[2 : len(exprString)-2]
	//support just assign for now
	//todo: optimization - trim whitspace when first creating object mapping
	return exprmapper.GetMappingValue(strings.TrimSpace(exprStr), scope, resolver)
}
