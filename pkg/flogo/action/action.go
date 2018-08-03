package action

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/activity/core"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

const (
	MashlingActionRef = "github.com/TIBCOSoftware/mashling/pkg/flogo/action"
)

type MashlingAction struct {
	ioMetadata *data.IOMetadata
	routes     []types.Route
	services   []types.Service
	identifier string
	instance   string
}

type Data struct {
	Routes     []types.Route   `json:"routes"`
	Services   []types.Service `json:"services"`
	Identifier string          `json:"identifier"`
	Instance   string          `json:"Instance"`
}

//todo fix this
var metadata = &action.Metadata{ID: "MashlingActionRef"}

func init() {
	action.RegisterFactory(MashlingActionRef, &Factory{})
}

type Factory struct {
}

func (f *Factory) Init() error {
	return nil
}

func (f *Factory) New(config *action.Config) (action.Action, error) {

	mAction := &MashlingAction{}
	var actionData Data
	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("failed to load mashling data: '%s' error '%s'", config.Id, err.Error())
	}
	mAction.routes = actionData.Routes
	mAction.services = actionData.Services
	mAction.identifier = actionData.Identifier
	mAction.instance = actionData.Instance
	return mAction, nil
}

func (m *MashlingAction) Metadata() *action.Metadata {
	return metadata
}

func (m *MashlingAction) IOMetadata() *data.IOMetadata {
	return m.ioMetadata
}

func (m *MashlingAction) Run(context context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	code, data, err := Core.ExecuteMashling(nil, m.identifier, m.instance, m.routes, m.services)
	return nil, err
}
