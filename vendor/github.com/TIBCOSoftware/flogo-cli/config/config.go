package config

import (
	"encoding/json"
	"reflect"
	"regexp"
)

type ContribType int

const (
	ACTION ContribType = 1 + iota
	TRIGGER
	ACTIVITY
	FLOW_MODEL
	REF

	FileDescriptor string = "flogo.json"
	FileImportsGo  string = "imports.go"
)

var ctStr = [...]string{
	"all",
	"action",
	"trigger",
	"activity",
	"flow-model",
}

func (m ContribType) String() string { return ctStr[m] }

func ToContribType(name string) ContribType {
	switch name {
	case "action":
		return ACTION
	case "trigger":
		return TRIGGER
	case "activity":
		return ACTIVITY
	case "flow-model":
		return FLOW_MODEL
	case "all":
		return 0
	}

	return -1
}

// FlogoAppDescriptor is the descriptor for a Flogo application
type FlogoAppDescriptor struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AppModel    string `json:"appModel,omitempty"`
	Triggers []*TriggerDescriptor `json:"triggers"`
}

// TriggerDescriptor is the config descriptor for a Trigger
type TriggerDescriptor struct {
	ID  string `json:"id"`
	Ref string `json:"ref"`
}

type TriggerMetadata struct {
	Name string `json:"name"`
	Ref  string `json:"ref"`
	Shim string `json:"shim"`
}

//FlogoPaletteDescriptor a package: just change to a list of references
type FlogoPaletteDescriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Extensions []Dependency `json:"extensions"`
}

type Descriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Dependency struct {
	ContribType ContribType
	Ref         string
}

func (d *Dependency) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ContribType string `json:"type"`
		Ref         string `json:"ref"`
	}{
		ContribType: d.ContribType.String(),
		Ref:         d.Ref,
	})
}

func (d *Dependency) UnmarshalJSON(data []byte) error {
	ser := &struct {
		ContribType string `json:"type"`
		Ref         string `json:"ref"`
	}{}

	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	d.Ref = ser.Ref
	d.ContribType = ToContribType(ser.ContribType)

	return nil
}

type depHolder struct {
	deps []*Dependency
}

func ExtractAllDependencies(appJson string) ([]*Dependency) {
	dh := &depHolder{}
	var descriptor interface{}
	//Should be valid app json
	json.Unmarshal([]byte(appJson), &descriptor)
	//Find all "ref" values in the model
	traverse(descriptor, dh)
	return dh.deps
}

func traverse(data interface{}, dh *depHolder ) {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		d := reflect.ValueOf(data)
		tmpData := make([]interface{}, d.Len())
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for _, v := range tmpData {
			traverse(v, dh)
		}
	} else if reflect.ValueOf(data).Kind() == reflect.Map {
		d := reflect.ValueOf(data)
		for _, k := range d.MapKeys() {
			match, _ := regexp.MatchString("(ref|activityRef)", k.String())
			if match {
				refVal := d.MapIndex(k).Interface()
				dh.deps = append(dh.deps, &Dependency{ContribType: REF, Ref: refVal.(string)})
			} else {
				if d.MapIndex(k).Interface() != nil {
					typeOfValue := reflect.TypeOf(d.MapIndex(k).Interface()).Kind()
					if typeOfValue == reflect.Map || typeOfValue == reflect.Slice {
						traverse(d.MapIndex(k).Interface(), dh)
					}
				}
			}
		}
	}
}
