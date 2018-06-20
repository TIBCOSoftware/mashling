package service

import (
	"fmt"
	"os"

	"github.com/pointlander/injectsec"
	"github.com/pointlander/injectsec/gru"
)

var maker *injectsec.DetectorMaker

func init() {
	var err error
	maker, err = injectsec.NewDetectorMaker()
	if err != nil {
		panic(err)
	}
}

// SQLD is a SQL injection attack detector
type SQLD struct {
	values       map[string]interface{}
	Attack       float32                `json:"attack"`
	AttackValues map[string]interface{} `json:"attackValues"`
	Maker        *injectsec.DetectorMaker
}

// InitializeSQLD creates a SQLD service
func InitializeSQLD(settings map[string]interface{}) (service *SQLD, err error) {
	service = &SQLD{
		AttackValues: make(map[string]interface{}),
	}
	err = service.UpdateRequest(settings)
	return
}

// Execute executes the SQLD service
func (s *SQLD) Execute() (err error) {
	var detector *gru.Detector
	if s.Maker != nil {
		detector = s.Maker.Make()
	} else {
		detector = maker.Make()
	}

	var testMap func(a, values map[string]interface{}) (err error)
	testMap = func(a, values map[string]interface{}) (err error) {
		for k, v := range a {
			switch element := v.(type) {
			case []interface{}:
				valuesList := make([]interface{}, 0, len(element))
				for _, item := range element {
					switch element := item.(type) {
					case map[string]interface{}:
						childValues := make(map[string]interface{})
						err = testMap(element, childValues)
						if err != nil {
							return
						}
						valuesList = append(valuesList, childValues)
					case string:
						probability, err := detector.Detect(element)
						valuesList = append(valuesList, float64(probability))
						if probability > s.Attack {
							s.Attack = probability
						}
						if err != nil {
							return err
						}
					}
				}
				values[k] = valuesList
			case map[string]interface{}:
				childValues := make(map[string]interface{})
				err = testMap(element, childValues)
				if err != nil {
					return
				}
				values[k] = childValues
			case string:
				probability, err := detector.Detect(element)
				values[k] = float64(probability)
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
			case []interface{}:
				valuesList := make([]interface{}, 0, len(b))
				for _, item := range b {
					switch element := item.(type) {
					case map[string]interface{}:
						childValues := make(map[string]interface{})
						err = testMap(element, childValues)
						if err != nil {
							return
						}
						valuesList = append(valuesList, childValues)
					case string:
						probability, err := detector.Detect(element)
						valuesList = append(valuesList, float64(probability))
						if probability > s.Attack {
							s.Attack = probability
						}
						if err != nil {
							return err
						}
					}
				}
				s.AttackValues[key] = valuesList
			case map[string]interface{}:
				values := make(map[string]interface{})
				err = testMap(b, values)
				s.AttackValues[key] = values
			case map[string]string:
				values := make(map[string]interface{})
				for _, v := range b {
					probability, err := detector.Detect(v)
					values[v] = float64(probability)
					if probability > s.Attack {
						s.Attack = probability
					}
					if err != nil {
						return err
					}
				}
				s.AttackValues[key] = values
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
	file := values["file"]
	if file != nil {
		if file, ok := file.(string); ok {
			var in *os.File
			in, err = os.Open(file)
			if err != nil {
				return
			}
			defer in.Close()
			s.Maker, err = injectsec.NewDetectorMakerWithWeights(in)
			if err != nil {
				return
			}
		} else {
			err = fmt.Errorf("file should be a string")
			return
		}
	}

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
