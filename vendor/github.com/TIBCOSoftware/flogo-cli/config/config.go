package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
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
	Name        string                `json:"name"`
	Type        string                `json:"type"`
	Version     string                `json:"version"`
	Description string                `json:"description"`
	AppModel    string                `json:"appModel,omitempty"`
	Triggers    []*trigger.Config     `json:"triggers"`
	Resources   []*ResourceDescriptor `json:"resources"`
	//deprecated
	Actions []*ActionDescriptor `json:"actions"`
}

type ResourceDescriptor struct {
	ID         string          `json:"id"`
	Compressed bool            `json:"compressed"`
	Data       json.RawMessage `json:"data"`
}

type ResourceData struct {
	Tasks        []*Task          `json:"tasks"`
	Links        []*Task          `json:"links"`
	ErrorHandler *ErrorHandlerRep `json:"errorHandler"`
}

// TaskOld is part of the flow structure
type TaskOld struct {
	Ref   string     `json:"activityRef"`
	Tasks []*TaskOld `json:"tasks"`
}

type Task struct {
	Activity *struct {
		Ref string `json:"ref"`
	} `json:"activity"`
}

type ErrorHandlerRep struct {
	Tasks []*Task `json:"tasks"`
}

type ActionDescriptor struct {
	ID   string `json:"id"`
	Ref  string `json:"ref"`
	Data *struct {
		Flow *struct {
			RootTask         *TaskOld `json:"rootTask"`
			ErrorHandlerTask *TaskOld `json:"errorHandlerTask"`
		} `json:"flow"`
	} `json:"data"`
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

func ExtractAllDependencies(appjson string) ([]*Dependency, error) {
	var deps []*Dependency

	flogoApp := &FlogoAppDescriptor{}
	jsonParser := json.NewDecoder(strings.NewReader(appjson))
	err := jsonParser.Decode(&flogoApp)
	if err != nil {
		return deps, err
	}
	deps = append(deps, extractTrigersDependency(flogoApp.Triggers)...)
	resourceDeps, err := extractResourceDependency(flogoApp.Resources)
	if err != nil {
		return deps, err
	}
	deps = append(deps, resourceDeps...)

	deps = append(deps, ExtractDependenciesActionOld(flogoApp.Actions)...)
	return deps, nil
}

func extractTrigersDependency(triggers []*trigger.Config) []*Dependency {
	var deps []*Dependency

	if triggers != nil && len(triggers) > 0 {
		for _, t := range triggers {
			deps = append(deps, &Dependency{ContribType: TRIGGER, Ref: t.Ref})
			if t.Handlers != nil {
				for _, t := range t.Handlers {
					if t.Action != nil {
						deps = append(deps, &Dependency{ContribType: ACTION, Ref: t.Action.Ref})
					}
				}
			}
		}
	}
	return deps
}
func extractResourceDependency(resources []*ResourceDescriptor) ([]*Dependency, error) {
	var deps []*Dependency

	if resources != nil && len(resources) > 0 {
		for _, t := range resources {
			if t.Compressed {

			}

			var flowDefBytes []byte

			if t.Compressed {
				decodedBytes, err := decodeAndUnzip(string(t.Data))
				if err != nil {
					return deps, fmt.Errorf("error decoding compressed resource with id '%s', %s", t.ID, err.Error())
				}

				flowDefBytes = decodedBytes
			} else {
				flowDefBytes = t.Data
			}

			var defRep *ResourceData
			err := json.Unmarshal(flowDefBytes, &defRep)
			if err != nil {
				return deps, fmt.Errorf("error marshalling flow resource with id '%s', %s", t.ID, err.Error())
			}

			if defRep.Tasks != nil {
				for _, task := range defRep.Tasks {
					deps = append(deps, &Dependency{ContribType: ACTIVITY, Ref: task.Activity.Ref})
				}
			}

			//Error handler
			if defRep.ErrorHandler != nil {
				for _, task := range defRep.ErrorHandler.Tasks {
					deps = append(deps, &Dependency{ContribType: ACTIVITY, Ref: task.Activity.Ref})
				}
			}

		}
	}
	return deps, nil
}

type depHolder struct {
	deps []*Dependency
}

// ExtractDependencies extracts dependencies from from application descriptor
func ExtractDependenciesActionOld(actions []*ActionDescriptor) []*Dependency {
	dh := &depHolder{}

	for _, action := range actions {
		dh.deps = append(dh.deps, &Dependency{ContribType: ACTION, Ref: action.Ref})

		if action.Data != nil && action.Data.Flow != nil {
			extractDepsFromTaskOld(action.Data.Flow.RootTask, dh)
			//Error handle flow
			if action.Data.Flow.ErrorHandlerTask != nil {
				extractDepsFromTaskOld(action.Data.Flow.ErrorHandlerTask, dh)
			}
		}
	}

	return dh.deps
}

// extractDepsFromTask extract dependencies from a TaskOld and is children
func extractDepsFromTaskOld(TaskOld *TaskOld, dh *depHolder) {

	if TaskOld.Ref != "" {
		dh.deps = append(dh.deps, &Dependency{ContribType: ACTIVITY, Ref: TaskOld.Ref})
	}

	for _, childTask := range TaskOld.Tasks {
		extractDepsFromTaskOld(childTask, dh)
	}
}

func decodeAndUnzip(encoded string) ([]byte, error) {

	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	return unzip(decoded)
}

func unzip(compressed []byte) ([]byte, error) {

	buf := bytes.NewBuffer(compressed)
	r, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	jsonAsBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return jsonAsBytes, nil
}
