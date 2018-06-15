package service

import (
	"encoding/json"
	"errors"

	"github.com/dop251/goja"
	"github.com/imdario/mergo"
)

// JS is a JS service.
type JS struct {
	Request JSRequest `json:"request"`
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
func (j *JS) Execute(requestValues map[string]interface{}) (Response, error) {
	response := JSResponse{}
	request, err := j.createRequest(requestValues)
	if err != nil {
		return response, err
	}
	result := make(map[string]interface{})
	vm, err := NewVM(nil)
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return response, err
	}
	vm.SetInVM("parameters", request.Parameters)
	vm.SetInVM("result", result)
	_, err = vm.vm.RunScript("JSServiceScript", request.Script)
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return response, err
	}
	err = vm.GetFromVM("result", &result)
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return response, err
	}
	response.Result = result
	return response, err
}

// InitializeJS initializes a JS service with provided settings.
func InitializeJS(settings map[string]interface{}) (j *JS, err error) {
	j = &JS{}
	// req := JSRequest{}
	// req.Parameters = make(map[string]interface{})
	// j.Request = req
	j.Request, err = j.createRequest(settings)
	return j, err
}

func (j *JS) createRequest(settings map[string]interface{}) (JSRequest, error) {
	request := JSRequest{}
	for k, v := range settings {
		switch k {
		case "script":
			script, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for script")
			}
			request.Script = script
		case "parameters":
			parameters, ok := v.(map[string]interface{})
			if !ok {
				return request, errors.New("invalid type for headers")
			}
			request.Parameters = parameters
			if err := mergo.Merge(&request.Parameters, j.Request.Parameters); err != nil {
				return request, errors.New("unable to merge parameters values")
			}
		default:
			// ignore and move on.
		}
		if err := mergo.Merge(&request, j.Request); err != nil {
			return request, errors.New("unable to merge request values")
		}
	}
	return request, nil
}

// VM represents a VM object.
type VM struct {
	vm *goja.Runtime
}

// NewVM initializes a new VM with defaults.
func NewVM(defaults map[string]interface{}) (vm *VM, err error) {
	vm = &VM{}
	vm.vm = goja.New()
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
