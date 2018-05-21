package ref

import (
	ejson "encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json/field"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("MappingRef")

type MappingRef struct {
	ref string
}

func NewMappingRef(ref string) *MappingRef {
	//Compatible TriggerData, the $TriggerData might in function or expression
	if strings.Index(ref, "$TriggerData") >= 0 {
		return &MappingRef{ref: strings.Replace(ref, "$TriggerData", "$flow", -1)}
	} else {
		return &MappingRef{ref: ref}
	}
}

func (m *MappingRef) GetRef() string {
	return m.ref
}

func (m *MappingRef) Eval(inputScope data.Scope, resovler data.Resolver) (interface{}, error) {
	log.Debugf("Eval mapping field and ref: %s", m.ref)

	if inputScope == nil {
		return nil, errors.New("Input scope cannot nil while eval mapping ref")
	}
	value, err := m.GetValue(inputScope, resovler)
	if err != nil {
		log.Errorf("Get From from ref error %+v", err)
	}

	log.Debugf("Mapping ref eval result: %+v", value)
	return value, err

}

func (m *MappingRef) GetValue(inputScope data.Scope, resovler data.Resolver) (interface{}, error) {

	inStruct, err := m.getValueFromAttribute(inputScope, resovler)
	if err != nil {
		return nil, err
	}
	mappingFiled, err := m.GetFields()
	if err != nil {
		return nil, err
	}

	if mappingFiled == nil || len(mappingFiled.Fields) <= 0 {
		value, err := toInfterface(inStruct)
		if err != nil {
			value = inStruct
		}
		return value, nil
	}
	mappingValue, err := json.GetFieldValueFromIn(inStruct, mappingFiled)
	if err != nil {
		return nil, err
	}
	return mappingValue, nil
}

func toInfterface(value interface{}) (interface{}, error) {

	var paramMap interface{}

	if value == nil {
		return paramMap, nil
	}

	switch t := value.(type) {
	case string:
		err := ejson.Unmarshal([]byte(t), &paramMap)
		if err != nil {
			return nil, err
		}
		return paramMap, nil
	default:
		return value, nil
	}
	return paramMap, nil
}

func (m *MappingRef) getValueFromAttribute(inputscope data.Scope, resolver data.Resolver) (interface{}, error) {

	log.Debugf("Get value from attribute ref [%s]", m.ref)
	if strings.HasPrefix(m.ref, "${") {
		//backward compatible for old version OSS
		return resolver.Resolve(m.ref, inputscope)
	}

	resolutionDetails, err := data.GetResolutionDetails(m.ref)
	if err != nil {
		return nil, fmt.Errorf("Get activity name and root field error, %s", err.Error())
	}
	//Only need activity and field name
	resolutionDetails.Path = ""
	//var newRef string
	var newRef string
	if resolutionDetails.ResolverName != "$." {
		if resolutionDetails.Item != "" {
			newRef = resolutionDetails.ResolverName + "[" + resolutionDetails.Item + "]" + "." + resolutionDetails.Property
		} else {
			newRef = resolutionDetails.ResolverName + "." + resolutionDetails.Property
		}
	} else {
		newRef = resolutionDetails.ResolverName + resolutionDetails.Property
	}

	log.Debugf("Activity and root field name is: %s", newRef)
	value, err := resolver.Resolve(newRef, inputscope)
	if err != nil {
		return nil, err
	}
	var relvalue interface{}
	if value != nil {
		switch t := value.(type) {
		case *data.ComplexObject:
			relvalue = t.Value
		default:
			relvalue = value
		}
	}

	return relvalue, nil
}

func (m *MappingRef) GetValueFromOutputScope(outputtscope data.Scope) (interface{}, error) {
	fieldName, err := m.GetFieldName()
	if err != nil {
		return nil, err
	}
	log.Debugf("GetValueFromOutputScope field name %s", fieldName)

	attribute, exist := outputtscope.GetAttr(fieldName)
	log.Debugf("GetValueFromOutputScope field name %s and exist %t ", fieldName, exist)

	if exist {
		switch attribute.Type() {
		case data.TypeComplexObject:
			complexObject := attribute.Value().(*data.ComplexObject)
			object := complexObject.Value
			//Convert the object to exist struct.
			//TODO return interface rather than string
			if object == nil {
				return "{}", nil
			}
			return object, nil
		default:
			return attribute.Value(), nil
		}

	}
	return nil, fmt.Errorf("Cannot found attribute %s", fieldName)
}

func (m *MappingRef) GetActivtyRootField() (string, error) {
	if field.HasSpecialFields(m.ref) {
		fields, err := field.GetAllspecialFields(m.ref)
		if err != nil {
			return "", err
		}
		activityNameRef := fields[0]
		if strings.HasPrefix(activityNameRef, "$") {
			activityName := activityNameRef[1:]
			activityName = getActivityName(activityName)
			fieldName := "_A." + activityName + "." + getFieldName(fields[1])
			return fieldName, nil
		}
		return getFieldName(fields[0]), nil
	}

	if strings.HasPrefix(m.ref, "$") {
		log.Debugf("Mapping ref %s", m.ref)
		mappingFields := strings.Split(m.ref, ".")
		//Which might like $A3.
		//field := mappingFields[1]
		var activityID string
		if strings.HasPrefix(m.ref, "$") {
			activityID = mappingFields[0][1:]
		} else {
			activityID = mappingFields[0]
		}

		//fieldName := "{" + activityID + "." + getFieldName(mappingFields[1]) + "}"
		fieldName := "_A." + getActivityName(activityID) + "." + getFieldName(mappingFields[1])

		log.Debugf("Field name now is: %s", fieldName)
		return fieldName, nil
	} else if strings.Index(m.ref, ".") > 0 {
		log.Debugf("Mapping ref %s", m.ref)
		mappingFields := strings.Split(m.ref, ".")
		log.Debugf("Field name now is: %s", mappingFields[0])
		return getFieldName(mappingFields[0]), nil
	} else {
		return m.ref, nil
	}
}

func (m *MappingRef) GetFields() (*field.MappingField, error) {
	hasArray := field.HasArray(m.ref)
	if field.HasSpecialFields(m.ref) {
		fields, err := field.GetAllspecialFields(m.ref)
		if err != nil {
			return nil, err
		}
		activityNameRef := fields[0]
		if strings.HasPrefix(activityNameRef, "$") {
			if strings.HasSuffix(fields[1], "]") {
				//Root element is an array
				arrayIndexPart := getArrayIndexPart(fields[1])
				fields[1] = arrayIndexPart
				return &field.MappingField{HasArray: hasArray, HasSpecialField: true, Fields: fields[1:]}, nil
			} else {
				return &field.MappingField{HasArray: hasArray, HasSpecialField: true, Fields: fields[2:]}, nil
			}
		}
		return &field.MappingField{HasArray: hasArray, HasSpecialField: true, Fields: fields[1:]}, nil
	}

	if strings.HasPrefix(m.ref, "$") {
		fieldArray := strings.Split(m.ref, ".")
		if strings.HasSuffix(fieldArray[1], "]") {
			//Root element is an array
			arrayIndexPart := getArrayIndexPart(fieldArray[1])
			fieldArray[1] = arrayIndexPart
			return &field.MappingField{HasArray: hasArray, HasSpecialField: false, Fields: fieldArray[1:]}, nil
		} else {
			return &field.MappingField{HasArray: hasArray, HasSpecialField: false, Fields: fieldArray[2:]}, nil
		}
	} else if strings.Index(m.ref, ".") >= 0 {
		return &field.MappingField{HasArray: hasArray, HasSpecialField: false, Fields: strings.Split(m.ref, ".")[1:]}, nil
	} else {
		//Only attribute name no field name
		return &field.MappingField{HasArray: hasArray, HasSpecialField: false, Fields: []string{}}, nil
	}
}

func (m *MappingRef) GetFieldName() (string, error) {
	if field.HasSpecialFields(m.ref) {
		fields, err := field.GetAllspecialFields(m.ref)
		if err != nil {
			return "", err
		}
		activityNameRef := fields[0]
		if strings.HasPrefix(activityNameRef, "$") {
			return getFieldName(fields[1]), nil
		}
		return getFieldName(fields[0]), nil
	}

	if strings.HasPrefix(m.ref, "$") || strings.Index(m.ref, ".") > 0 {
		log.Debugf("Mapping ref %s", m.ref)
		mappingFields := strings.Split(m.ref, ".")
		if strings.HasPrefix(m.ref, "$") {
			return getFieldName(mappingFields[1]), nil

		}
		log.Debugf("Field name now is: %s", mappingFields[0])
		return getFieldName(mappingFields[0]), nil

	}
	return getFieldName(m.ref), nil
}

func (m *MappingRef) GetActivityId() (string, error) {

	dotIndex := strings.Index(m.ref, ".")

	if dotIndex == -1 {
		return "", fmt.Errorf("invalid resolution expression [%s]", m.ref)
	}

	firstItemIndex := strings.Index(m.ref[:dotIndex], "[")

	if firstItemIndex != -1 {
		return m.ref[firstItemIndex+1 : dotIndex-1], nil
	}
	return "", nil
}

//
//func GetResolutionDetails(toResolve string) (*string, error) {
//
//
//	dotIdx := strings.Index(toResolve, ".")
//
//	if dotIdx == -1 {
//		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
//	}
//
//	details := &ResolutionDetails{}
//	itemIdx := strings.Index(toResolve[:dotIdx], "[")
//
//	if itemIdx != -1 {
//		details.Item = toResolve[itemIdx+1:dotIdx-1]
//		details.ResolverName = toResolve[:itemIdx]
//	} else {
//		details.ResolverName = toResolve[:dotIdx]
//
//		//special case for activity without brackets
//		if strings.HasPrefix(toResolve, "activity") {
//			nextDot := strings.Index(toResolve[dotIdx+1:], ".") + dotIdx + 1
//			details.Item = toResolve[dotIdx+1:nextDot]
//			dotIdx = nextDot
//		}
//	}
//
//	pathIdx := strings.IndexFunc(toResolve[dotIdx+1:], isSep)
//
//	if pathIdx != -1 {
//		pathStart := pathIdx + dotIdx + 1
//		details.Path = toResolve[pathStart:]
//		details.Property = toResolve[dotIdx+1:pathStart]
//	} else {
//		details.Property = toResolve[dotIdx+1:]
//	}
//
//	return details, nil
//}

func getFieldName(fieldname string) string {
	if strings.Index(fieldname, "[") > 0 && strings.Index(fieldname, "]") > 0 {
		return fieldname[:strings.Index(fieldname, "[")]
	}
	return fieldname
}

func getActivityName(fieldname string) string {
	//$activity[name]
	startIndex := strings.Index(fieldname, "[")
	endIndex := strings.Index(fieldname, "]")
	if startIndex >= 0 {
		return fieldname[startIndex+1 : endIndex]
	} else {
		return fieldname
	}
}

//getArrayIndexPart get array part of the string. such as name[0] return [0]
func getArrayIndexPart(fieldName string) string {
	if strings.Index(fieldName, "[") >= 0 {
		return fieldName[strings.Index(fieldName, "[") : strings.Index(fieldName, "]")+1]
	}
	return ""
}
