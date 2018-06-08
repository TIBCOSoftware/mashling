package service

import (
	"github.com/pointlander/injectsec"
)

var maker = injectsec.NewDetectorMaker()

// SQLD is a SQL injection attack detector
type SQLD struct {
	values map[string]interface{}
	Attack bool `json:"attack"`
}

// InitializeSQLD creates a SQLD service
func InitializeSQLD(settings map[string]interface{}) (service *SQLD, err error) {
	service = &SQLD{}
	err = service.UpdateRequest(settings)
	return
}

// Execute executes the SQLD service
func (s *SQLD) Execute() (err error) {
	detector := maker.Make()

	var testMap func(a map[string]interface{}) bool
	testMap = func(a map[string]interface{}) bool {
		for _, v := range a {
			switch element := v.(type) {
			case map[string]interface{}:
				if testMap(element) {
					return true
				}
			case string:
				if detector.Detect(element) {
					return true
				}
			}
		}

		return false
	}

	test := func(key string) bool {
		if a, ok := s.values[key]; ok {
			switch b := a.(type) {
			case map[string]interface{}:
				if testMap(b) {
					s.Attack = true
					return true
				}
			case map[string]string:
				for _, v := range b {
					if detector.Detect(v) {
						s.Attack = true
						return true
					}
				}
			}
		}

		return false
	}

	if test("pathParams") {
		return
	}
	if test("queryParams") {
		return
	}
	if test("content") {
		return
	}

	return
}

// UpdateRequest updates the SQLD service
func (s *SQLD) UpdateRequest(values map[string]interface{}) (err error) {
	a := values["payload"]
	if a == nil {
		return
	}
	b := *a.(*interface{})
	c, ok := b.(map[string]interface{})
	if !ok {
		return
	}
	s.values = c

	return
}
