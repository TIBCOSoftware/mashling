package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/imdario/mergo"
)

// FlogoFlow is a Flogo flow service.
type FlogoFlow struct {
	Request FlogoFlowRequest `json:"request"`
	Action  action.Action    `json:"action"`
}

// FlogoFlowRequest is a flogo flow service request.
type FlogoFlowRequest struct {
	Definition map[string]interface{} `json:"definition"`
	Reference  string                 `json:"reference"`
	Inputs     map[string]interface{} `json:"inputs"`
}

// FlogoFlowResponse is a flogo flow service response.
type FlogoFlowResponse struct {
	Done    bool                   `json:"done"`
	Error   string                 `json:"error"`
	Outputs map[string]interface{} `json:"outputs"`
}

// InitializeFlogoFlow initializes a FlogoFlow service with provided settings.
func InitializeFlogoFlow(settings map[string]interface{}) (flogoFlowService *FlogoFlow, err error) {
	flogoFlowService = &FlogoFlow{}
	// req := FlogoFlowRequest{}
	// req.Inputs = make(map[string]interface{})
	// flogoFlowService.Request = req
	flogoFlowService.Request, err = flogoFlowService.createRequest(settings)
	if err != nil {
		return flogoFlowService, nil
	}
	var flowAction action.Action
	cfg := &action.Config{}
	rawData, err := json.Marshal(flogoFlowService.Request.Definition["data"])
	if err != nil {
		return flogoFlowService, err
	}
	cfg.Data = rawData
	cfg.Id = flogoFlowService.Request.Reference
	cfg.Ref = flogoFlowService.Request.Definition["ref"].(string)
	ff := flow.ActionFactory{}
	flowAction, err = ff.New(cfg)
	if err != nil {
		return flogoFlowService, err
	}
	flogoFlowService.Action = flowAction
	return flogoFlowService, err
}

func (f *FlogoFlow) createRequest(settings map[string]interface{}) (FlogoFlowRequest, error) {
	request := FlogoFlowRequest{}
	for k, v := range settings {
		switch k {
		case "definition":
			definition, ok := v.(map[string]interface{})
			if !ok {
				return request, errors.New("invalid type for definition")
			}
			request.Definition = definition
		case "reference":
			reference, ok := v.(string)
			if !ok {
				return request, errors.New("invalid type for reference")
			}
			request.Reference = reference
		case "inputs":
			inputs, ok := v.(map[string]interface{})
			if !ok {
				return request, errors.New("invalid type for inputs")
			}
			request.Inputs = inputs
			if err := mergo.Merge(&request.Inputs, f.Request.Inputs); err != nil {
				return request, errors.New("unable to merge inputs values")
			}
		default:
			// ignore and move on.
		}
	}
	if err := mergo.Merge(&request, f.Request); err != nil {
		return request, errors.New("unable to merge request values")
	}
	return request, nil
}

// Execute invokes this FlogoActivity service.
func (f *FlogoFlow) Execute(requestValues map[string]interface{}) (Response, error) {
	response := FlogoFlowResponse{}
	request, err := f.createRequest(requestValues)
	if err != nil {
		return response, err
	}
	var attrs []*data.Attribute
	mAttrs := make(map[string]*data.Attribute)
	if request.Inputs != nil {

		for k, v := range request.Inputs {
			attr, dErr := data.NewAttribute(k, data.TypeAny, v)
			if dErr != nil {
				response.Error = dErr.Error()
				return response, dErr
			}
			attrs = append(attrs, attr)
			mAttrs[k] = attr
			attr, dErr = data.NewAttribute("_T."+k, data.TypeAny, v)
			if dErr != nil {
				response.Error = dErr.Error()
				return response, dErr
			}
			attrs = append(attrs, attr)
			mAttrs["_T."+k] = attr
		}
	}
	ctx := trigger.NewContext(context.Background(), attrs)
	r := runner.NewDirect()
	outputData, err := r.Execute(ctx, f.Action, mAttrs)
	outputs := make(map[string]interface{})
	for _, v := range outputData {
		outputs[v.Name()] = v.Value()
	}
	response.Done = true
	if err != nil {
		response.Error = err.Error()
	}
	response.Outputs = outputs
	return response, err
}
