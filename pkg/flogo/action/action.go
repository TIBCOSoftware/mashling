package action

import (
	"context"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

const (
	MashlingActionRef = "github.com/TIBCOSoftware/mashling/pkg/flogo/action"
)

type MashlingAction struct {
	ioMetadata *data.IOMetadata
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

	return mAction, nil
}

func (s *MashlingAction) Metadata() *action.Metadata {
	return metadata
}

func (s *MashlingAction) IOMetadata() *data.IOMetadata {
	return s.ioMetadata
}

func (s *MashlingAction) Run(context context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	return nil, nil
}
