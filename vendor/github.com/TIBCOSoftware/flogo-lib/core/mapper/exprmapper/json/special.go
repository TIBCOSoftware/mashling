package json

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

func getArrayFieldName(fields []string) ([]string, int, int) {
	var tmpFields []string
	index := -1
	var arrayIndex int
	for i, field := range fields {
		if strings.Index(field, "[") >= 0 && strings.Index(field, "]") >= 0 {
			tmpFields = append(tmpFields, getFieldName(field))
			arrayIndex, _ = getFieldSliceIndex(field)
			index = i
			break
		} else {
			tmpFields = append(tmpFields, field)
		}
	}
	return tmpFields, index, arrayIndex
}

func hasArrayFieldInArray(fields []string) bool {
	for _, field := range fields {
		if strings.Index(field, "[") >= 0 && strings.HasSuffix(field, "]") {
			//Make sure the index are integer
			_, err := strconv.Atoi(getNameInsideBrancket(field))
			if err == nil {
				return true
			}
		}
	}
	return false
}

func handleArrayWithSpecialFields(value interface{}, jsonData *JSONData, fields []string) (interface{}, error) {

	log.Debugf("All fields %+v", fields)
	jsonData.rw.Lock()
	defer jsonData.rw.Unlock()

	container := jsonData.container
	if hasArrayFieldInArray(fields) {
		arrayFields, fieldNameindex, arrayIndex := getArrayFieldName(fields)
		//No array field found
		if fieldNameindex == -1 {
			if arrayIndex == -2 {
				//Append
				err := container.ArrayAppend(value, arrayFields...)
				if err != nil {
					return nil, err
				}
			} else {
				//set to exist index array
				size, err := container.ArrayCount(arrayFields...)
				if err != nil {
					return nil, err
				}
				if arrayIndex > size-1 {
					err := container.ArrayAppend(value, arrayFields...)
					if err != nil {
						return nil, err
					}
				} else {
					array := container.S(arrayFields...)
					_, err := array.SetIndex(value, arrayIndex)
					if err != nil {
						return nil, err
					}
				}
			}
		} else {
			restFields := fields[fieldNameindex+1:]
			if container.Exists(arrayFields...) {
				_, err := container.ArrayElement(arrayIndex, arrayFields...)
				if err != nil {
					return nil, err
				}
				return handleArrayWithSpecialFields(value, &JSONData{container: container, rw: sync.RWMutex{}}, restFields)
			}
			//Create new one
			array, err := container.ArrayOfSize(arrayIndex+1, arrayFields...)
			if err != nil {
				return nil, err
			}

			if hasArrayFieldInArray(restFields) {
				return handleArrayWithSpecialFields(value, &JSONData{container: array, rw: sync.RWMutex{}}, restFields)
			}
			newObject, err := ParseJSON([]byte("{}"))
			_, err = newObject.Set(value, restFields...)
			log.Debugf("new object %s", newObject.String())
			if err != nil {
				return nil, err
			}
			maps := &map[string]interface{}{}
			err = json.Unmarshal(newObject.Bytes(), maps)
			if err != nil {
				return nil, err
			}
			_, err = array.SetIndex(maps, arrayIndex)

		}
	} else {
		_, err := jsonData.container.Set(value, fields...)
		if err != nil {
			return nil, err
		}
	}

	return jsonData.container.object, nil
}

func handleGetSpecialFields(jsonData *JSONData, fields []string) (interface{}, error) {

	log.Debugf("All fields %+v", fields)
	jsonData.rw.Lock()
	defer jsonData.rw.Unlock()

	container := jsonData.container
	if hasArrayFieldInArray(fields) {
		arrayFields, fieldNameindex, arrayIndex := getArrayFieldName(fields)
		//No array field found
		if fieldNameindex == -1 {
			return container.S(arrayFields...).Data(), nil
		}
		restFields := fields[fieldNameindex+1:]
		specialField, err := container.ArrayElement(arrayIndex, arrayFields...)
		if err != nil {
			return nil, err
		}
		log.Debugf("Array element value %s", specialField)
		if hasArrayFieldInArray(restFields) {
			return handleGetSpecialFields(&JSONData{container: specialField, rw: sync.RWMutex{}}, restFields)
		}
		return specialField.S(restFields...).Data(), nil
	}
	log.Debugf("No array found for array %+v and size %d", fields, len(fields))
	return container.S(fields...).Data(), nil
}
