package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
)

var flowActions sync.Map

// FlogoFlow is a Flogo flow service.
type FlogoFlow struct {
	Request  FlogoFlowRequest  `json:"request"`
	Response FlogoFlowResponse `json:"response"`
}

// FlogoFlowRequest is a flogo flow service request.
type FlogoFlowRequest struct {
	Definition map[string]interface{} `json:"definition"`
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
	req := FlogoFlowRequest{}
	req.Inputs = make(map[string]interface{})
	for k, v := range settings {
		switch k {
		case "definition":
			definition, ok := v.(map[string]interface{})
			if !ok {
				return flogoFlowService, errors.New("invalid type for definition")
			}
			req.Definition = definition
		case "inputs":
			inputs, ok := v.(map[string]interface{})
			if !ok {
				return flogoFlowService, errors.New("invalid type for inputs")
			}
			req.Inputs = inputs
		default:
			// ignore and move on.
		}
		flogoFlowService.Request = req
	}
	return flogoFlowService, err
}

// Execute invokes this FlogoActivity service.
func (f *FlogoFlow) Execute() (err error) {
	// Ignore IDs and do everything by ref?
	var flowAction action.Action
	flowActionStored, exists := flowActions.Load(f.Request.Definition["ref"].(string))
	if !exists {
		cfg := &action.Config{}
		f.Response = FlogoFlowResponse{}
		rawData, err := json.Marshal(f.Request.Definition["data"])
		if err != nil {
			return err
		}
		cfg.Data = rawData
		cfg.Id = f.Request.Definition["ref"].(string)
		cfg.Ref = f.Request.Definition["ref"].(string)

		ff := flow.FlowFactory{}
		flowAction = ff.New(cfg)
		flowActions.Store(f.Request.Definition["ref"].(string), flowAction)
	} else {
		flowAction = flowActionStored.(action.Action)
	}

	ctx := context.Background()

	if f.Request.Inputs != nil {

		var attrs []*data.Attribute

		for k, v := range f.Request.Inputs {
			attr, _ := data.NewAttribute(k, data.ANY, v)
			attrs = append(attrs, attr)
		}
		ctx = trigger.NewContext(context.Background(), attrs)
	}

	r := runner.NewDirect()
	outputData, err := r.RunAction(ctx, flowAction, nil)
	outputs := make(map[string]interface{})
	for _, v := range outputData {
		outputs[v.Name()] = v.Value()
	}
	f.Response = FlogoFlowResponse{}
	f.Response.Done = true
	if err != nil {
		f.Response.Error = err.Error()
	}
	f.Response.Outputs = outputs
	return err
}
