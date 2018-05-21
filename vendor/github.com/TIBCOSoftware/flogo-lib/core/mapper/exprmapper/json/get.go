package json

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"encoding/json"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("json")

func GetFieldValueP(jsonData string, path string) (interface{}, error) {
	jsonParsed, err := ParseJSON([]byte(jsonData))
	if err != nil {
		return nil, err

	}
	return getFieldValueP(&JSONData{container: jsonParsed, rw: sync.RWMutex{}}, path)
}

func GetFieldValueFromInP(data interface{}, path string) (interface{}, error) {
	var jsonParsed *Container
	var err error

	if reflect.TypeOf(data).Kind() == reflect.String {
		jsonParsed, err = ParseJSON([]byte(data.(string)))
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		jsonParsed, err = ParseJSON(b)
	}

	if err != nil {
		return nil, err

	}
	return getFieldValueP(&JSONData{container: jsonParsed, rw: sync.RWMutex{}}, path)
}

func GetFieldValueFromIn(data interface{}, mappingField *field.MappingField) (interface{}, error) {
	var jsonParsed *Container
	var err error

	if reflect.TypeOf(data).Kind() == reflect.String {
		jsonParsed, err = ParseJSON([]byte(data.(string)))
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		jsonParsed, err = ParseJSON(b)
	}

	if err != nil {
		return nil, err

	}
	return getFieldValue(&JSONData{container: jsonParsed, rw: sync.RWMutex{}}, mappingField)
}

func getFieldValue(jsonData *JSONData, mappingField *field.MappingField) (interface{}, error) {
	var res interface{}
	var err error
	if mappingField.HasArray && mappingField.HasSpecialField {
		res, err = handleGetSpecialFields(jsonData, mappingField.Fields)
	} else if mappingField.HasArray {
		data, err := getFiledContainer(jsonData, strings.Join(mappingField.Fields, "."))
		if data != nil {
			return data.Data(), err
		}
		return nil, err
	} else if mappingField.HasSpecialField {
		res, err = handleGetSpecialFields(jsonData, mappingField.Fields)
	} else {
		data, err := getFiledContainer(jsonData, strings.Join(mappingField.Fields, "."))
		return data.Data(), err
	}
	return res, err
}

func getFieldValueP(jsonData *JSONData, path string) (interface{}, error) {
	var res interface{}
	var err error
	if field.HasArray(path) && field.HasSpecialFields(path) {
		fields, errs := field.GetAllspecialFields(path)
		if errs != nil {
			return nil, errs
		}
		res, err = handleGetSpecialFields(jsonData, fields)
	} else if field.HasArray(path) {
		data, err := getFiledContainer(jsonData, path)
		return data.Data(), err
	} else if field.HasSpecialFields(path) {
		fields, errs := field.GetAllspecialFields(path)
		if errs != nil {
			return nil, errs
		}
		res, err = handleGetSpecialFields(jsonData, fields)
	} else {
		data, err := getFiledContainer(jsonData, path)
		return data.Data(), err
	}

	return res, err
}

func getFiledContainer(jsonData *JSONData, path string) (*Container, error) {
	jsonData.rw.RLock()

	defer jsonData.rw.RUnlock()
	container := jsonData.container
	if strings.Index(path, "[") >= 0 {
		//Array

		arrayFieldName := getFieldName(path)
		index, err := getFieldSliceIndex(path)
		if err != nil {
			return nil, fmt.Errorf("Get array index error, due to [%s]", err.Error())
		}

		var jsonField *Container
		if arrayFieldName == "" {
			//Root level array handling
			jsonField, err = container.ArrayElement(index)
		} else {
			log.Debug("Field Name:", arrayFieldName, " and index: ", index)
			jsonField, err = container.ArrayElementP(index, arrayFieldName)

		}
		if err != nil {
			return nil, err
		}

		restPath := getRestArrayFieldName(path)
		if restPath == "" {
			// value := jsonField.Data()
			// log.Debug("Type :", reflect.TypeOf(value))
			// log.Debug("Value :", value)
			return jsonField, nil
		}
		if strings.Index(restPath, "[") > 0 {
			return getFiledContainer(&JSONData{container: jsonField, rw: sync.RWMutex{}}, restPath)
		}
		filed := jsonField.Path(restPath)
		return filed, nil

	}
	// value := container.Path(path).Data()
	// log.Debug("Type :", reflect.TypeOf(value))
	// log.Debug("Value :", value)
	return container.Path(path), nil
}

func getRestArrayFieldName(fieldName string) string {
	if strings.Index(fieldName, "]") >= 0 {
		closeBracketIndex := strings.Index(fieldName, "]")
		if len(fieldName) == closeBracketIndex+1 {
			return ""
		}
		return fieldName[closeBracketIndex+2:]
	}
	return fieldName
}

func getFieldName(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		return fieldName[0:strings.Index(fieldName, "[")]
	}

	return fieldName
}

func getFieldSliceIndex(fieldName string) (int, error) {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		i, err := strconv.Atoi(index)

		if err != nil {
			return -2, nil
		}
		return i, nil
	}

	return -1, nil
}

func getNameInsideBrancket(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		index := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		return index
	}

	return ""
}
