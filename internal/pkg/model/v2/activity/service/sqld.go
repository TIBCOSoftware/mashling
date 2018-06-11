package service

import (
	"github.com/pointlander/injectsec"
)

var maker = injectsec.NewDetectorMaker()

// SQLD is a SQL injection attack detector
type SQLD struct {
	values map[string]interface{}
	Attack float32 `json:"attack"`
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

	var testMap func(a map[string]interface{}) (err error)
	testMap = func(a map[string]interface{}) (err error) {
		for _, v := range a {
			switch element := v.(type) {
			case map[string]interface{}:
				err = testMap(element)
				if err != nil {
					return
				}
			case string:
				probability, err := detector.Detect(element)
				if probability > s.Attack {
					s.Attack = probability
				}
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	test := func(key string) (err error) {
		if a, ok := s.values[key]; ok {
			switch b := a.(type) {
			case map[string]interface{}:
				err = testMap(b)
			case map[string]string:
				for _, v := range b {
					probability, err := detector.Detect(v)
					if probability > s.Attack {
						s.Attack = probability
					}
					if err != nil {
						return err
					}
				}
			}
		}

		return
	}

	err = test("pathParams")
	if err != nil {
		return
	}
	err = test("queryParams")
	if err != nil {
		return
	}
	err = test("content")
	if err != nil {
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
