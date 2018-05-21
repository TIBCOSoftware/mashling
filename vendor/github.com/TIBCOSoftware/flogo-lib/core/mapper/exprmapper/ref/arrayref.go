package ref

import (
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/json"
)

type ArrayRef struct {
	ref string
}

func (m *ArrayRef) GetRef() string{
	return m.ref
}

func NewArrayRef(ref string) *ArrayRef {
	return &ArrayRef{ref:ref}
}


func (m *ArrayRef) EvalFromData(data interface{}) (interface{}, error) {
	log.Debugf("Eval mapping field and ref: %s", m.ref)
	//
	value, err := m.getValueFromRef(data, m.ref)
	if err != nil {
		log.Errorf("Get From from ref error %+v", err)
	}

	log.Debugf("Mapping ref eval result: %p", value)
	return value, err
}

func (m *ArrayRef) Eval(inputScope, outputScope data.Scope) (interface{}, error) {

	return nil, fmt.Errorf("Array ref not support eval")

}

func (m *ArrayRef) getValueFromRef(object interface{}, ref string) (interface{}, error) {
	reference := GetFieldNameFromArrayRef(ref)
	return json.GetFieldValueFromInP(object, reference)
}

func GetFieldNameFromArrayRef(arrayRef string) string {
	if arrayRef != "" {
		if IsArrayMapping(arrayRef) {
			return arrayRef[2:]
		}
	}
	return arrayRef
}

func IsArrayMapping(ref string) bool {
	if ref != "" {
		return strings.HasPrefix(ref, "$.") || strings.HasPrefix(ref, "$$")
	}
	return false
}
