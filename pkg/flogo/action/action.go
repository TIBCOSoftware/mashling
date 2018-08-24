package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/action/core"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/action/pattern"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

const (
	MashlingActionRef = "github.com/TIBCOSoftware/mashling/pkg/flogo/action"
)

var (
	defaultManager *MashlingManager
)

type MashlingAction struct {
	mashlingURI   string
	metadata      *action.Metadata
	ioMetadata    *data.IOMetadata
	dispatch      types.Dispatch
	services      []types.Service
	pattern       string
	configuration map[string]interface{}
}

type Data struct {
	MashlingURI   string                 `json:"mashlingURI"`
	Dispatch      json.RawMessage        `json:"dispatch"`
	Services      json.RawMessage        `json:"services"`
	Pattern       string                 `json:"pattern"`
	Configuration map[string]interface{} `json:"configuration"`
}

type MashlingManager struct {
	resMashlings map[string]*Data
}

func init() {
	action.RegisterFactory(MashlingActionRef, &Factory{})
	defaultManager := &MashlingManager{}
	defaultManager.resMashlings = make(map[string]*Data)
	resource.RegisterManager("mashling", defaultManager)
}

func (mm *MashlingManager) LoadResource(config *resource.Config) error {

	mashlingDefBytes := config.Data

	var mashlingDefinition *Data
	err := json.Unmarshal(mashlingDefBytes, &mashlingDefinition)
	if err != nil {
		return fmt.Errorf("error marshalling mashling definition resource with id '%s', %s", config.ID, err.Error())
	}

	mm.resMashlings[config.ID] = mashlingDefinition
	return nil
}

func (mm *MashlingManager) GetResource(id string) interface{} {
	return mm.resMashlings[id]
}

type Factory struct {
}

func (f *Factory) Init() error {
	return nil
}

func (f *Factory) New(config *action.Config) (action.Action, error) {
	mAction := &MashlingAction{}
	mAction.metadata = &action.Metadata{ID: config.Id}
	var actionData *Data
	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("failed to load mashling data: '%s' error '%s'", config.Id, err.Error())
	}
	if actionData.MashlingURI != "" {
		// Load action data from resources
		resData, err := resource.Get(actionData.MashlingURI)
		if err != nil {
			return nil, fmt.Errorf("failed to load mashling URI data: '%s' error '%s'", config.Id, err.Error())
		}
		actionData = resData.(*Data)
	}
	// Extract configuration
	mAction.configuration = actionData.Configuration
	// Extract pattern
	mAction.pattern = actionData.Pattern
	if mAction.pattern == "" {
		// Parse routes
		var dispatch types.Dispatch
		err = json.Unmarshal(actionData.Dispatch, &dispatch)
		if err != nil {
			return nil, err
		}
		// Parse services
		var services []types.Service
		err = json.Unmarshal(actionData.Services, &services)
		if err != nil {
			return nil, err
		}
		mAction.dispatch = dispatch
		mAction.services = services
	} else {
		pDef, err := pattern.Load(mAction.pattern)
		if err != nil {
			return nil, err
		}
		mAction.dispatch = pDef.Dispatch
		mAction.services = pDef.Services
	}

	return mAction, nil
}

func (m *MashlingAction) Metadata() *action.Metadata {
	return m.metadata
}

func (m *MashlingAction) IOMetadata() *data.IOMetadata {
	return m.ioMetadata
}

func (m *MashlingAction) Run(context context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	payload := make(map[string]interface{})
	for k, v := range inputs {
		payload[k] = v.Value()
	}

	code, mData, err := core.ExecuteMashling(payload, m.configuration, m.dispatch.Routes, m.services)
	output := make(map[string]*data.Attribute)
	codeAttr, err := data.NewAttribute("code", data.TypeInteger, code)
	if err != nil {
		return nil, err
	}
	output["code"] = codeAttr
	dataAttr, err := data.NewAttribute("data", data.TypeObject, mData)
	if err != nil {
		return nil, err
	}
	output["data"] = dataAttr
	return output, err
}
