package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/TIBCOSoftware/mashling/pkg/strings"
	"github.com/dop251/goja"
)

// JS is a JS service.
type JS struct {
	Request  JSRequest  `json:"request"`
	Response JSResponse `json:"response"`
}

// JSRequest is a JS service request.
type JSRequest struct {
	Script     string                 `json:"script"`
	Parameters map[string]interface{} `json:"parameters"`
}

// JSResponse is a JS service response.
type JSResponse struct {
	Error        bool                   `json:"error"`
	ErrorMessage string                 `json:"errorMessage"`
	Result       map[string]interface{} `json:"result"`
}

// Execute invokes this JS service.
func (j *JS) Execute() (err error) {
	j.Response = JSResponse{}
	result := make(map[string]interface{})
	vm, err := NewVM(nil)
	if err != nil {
		j.Response.Error = true
		j.Response.ErrorMessage = err.Error()
		return err
	}
	vm.SetInVM("parameters", j.Request.Parameters)
	vm.SetInVM("result", result)
	_, err = vm.vm.RunScript("JSServiceScript", j.Request.Script)
	if err != nil {
		j.Response.Error = true
		j.Response.ErrorMessage = err.Error()
		return err
	}
	err = vm.GetFromVM("result", &result)
	if err != nil {
		j.Response.Error = true
		j.Response.ErrorMessage = err.Error()
		return err
	}
	j.Response.Result = result
	return err
}

// InitializeJS initializes a JS service with provided settings.
func InitializeJS(settings map[string]interface{}) (j *JS, err error) {
	j = &JS{}
	req := JSRequest{}
	req.Parameters = make(map[string]interface{})
	for k, v := range settings {
		switch k {
		case "script":
			script, ok := v.(string)
			if !ok {
				return j, errors.New("invalid type for script")
			}
			req.Script = script
		case "parameters":
			parameters, ok := v.(map[string]interface{})
			if !ok {
				return j, errors.New("invalid type for headers")
			}
			req.Parameters = parameters
		default:
			// ignore and move on.
		}
		j.Request = req
	}
	return j, err
}

// VM represents a VM object.
type VM struct {
	vm *goja.Runtime
}

// NewVM initializes a new VM with defaults.
func NewVM(defaults map[string]interface{}) (vm *VM, err error) {
	vm = &VM{}
	vm.vm = goja.New()
	_, err = vm.vm.RunScript("AssignFunc", objectAssignFunc)
	if err != nil {
		return vm, err
	}
	for k, v := range defaults {
		if v != nil {
			vm.vm.Set(k, v)
		}
	}
	return vm, err
}

// EvaluateToBool evaluates a string condition within the context of the VM.
func (vm *VM) EvaluateToBool(condition string) (truthy bool, err error) {
	if condition == "" {
		return true, nil
	}
	var res goja.Value
	res, err = vm.vm.RunString(condition)
	if err != nil {
		return false, err
	}
	truthy, ok := res.Export().(bool)
	if !ok {
		err = errors.New("condition does not evaluate to bool")
		return false, err
	}
	return truthy, err
}

// SetInVM sets the object name and value in the VM.
func (vm *VM) SetInVM(name string, object interface{}) (err error) {
	var valueJSON json.RawMessage
	var vmObject map[string]interface{}
	valueJSON, err = json.Marshal(object)
	if err != nil {
		return err
	}
	err = json.Unmarshal(valueJSON, &vmObject)
	if err != nil {
		return err
	}
	vm.vm.Set(name, vmObject)
	return err
}

// GetFromVM extracts the current object value from the VM.
func (vm *VM) GetFromVM(name string, object interface{}) (err error) {
	var valueJSON json.RawMessage
	var vmObject map[string]interface{}
	vm.vm.ExportTo(vm.vm.Get(name), &vmObject)

	valueJSON, err = json.Marshal(vmObject)
	if err != nil {
		return err
	}
	err = json.Unmarshal(valueJSON, object)
	if err != nil {
		return err
	}
	return err
}

// SetPrimitiveInVM sets primitive value in VM.
func (vm *VM) SetPrimitiveInVM(name string, primitive interface{}) {
	vm.vm.Set(name, primitive)
}

// RunTranslationMappings maps objects in the VM to new values that may be
// literals or other data already in the VM.
func (vm *VM) RunTranslationMappings(objectRoot string, mappings map[string]interface{}) (err error) {
	if len(mappings) == 0 {
		return err
	}
	assignObjFmt := "assign(%s, \"%s\", %s);\n"
	assignStrFmt := "assign(%s, \"%s\", \"%s\");\n"
	var transformation bytes.Buffer
	for k, v := range mappings {
		switch value := v.(type) {
		case string:
			if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
				// this is a variable
				value = strings.Replace(value, "${", "", 1)
				value = util.TrimSuffix(value, "}")
				transformation.WriteString(fmt.Sprintf(assignObjFmt, objectRoot, k, value))
			} else {
				transformation.WriteString(fmt.Sprintf(assignStrFmt, objectRoot, k, value))
			}
		default:
			// convert to raw JSON and insert.
			valueJSON, jerr := json.Marshal(value)
			if jerr != nil {
				return jerr
			}
			transformation.WriteString(fmt.Sprintf(assignObjFmt, objectRoot, k, string(valueJSON)))
		}
	}
	_, err = vm.vm.RunScript(objectRoot+"Translations", transformation.String())
	return err
}

const objectAssignFunc = `function assign(obj, keyPath, value) {
   keyPath = keyPath.split('.');
   lastKeyIndex = keyPath.length-1;
   for (var i = 0; i < lastKeyIndex; ++ i) {
     key = keyPath[i];
     if (!(key in obj))
       obj[key] = {}
     obj = obj[key];
   }
   obj[keyPath[lastKeyIndex]] = value;
}`
