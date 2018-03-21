package action

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Action is an action to perform as a result of a trigger
type Action interface {

	//Config get the Action's config
	Config() *Config

	//Metadata get the Action's metadata
	Metadata() *Metadata

	// Run this Action
	Run(context context.Context, inputs []*data.Attribute, options map[string]interface{}, handler ResultHandler) error
}

// Factory is used to create new instances for an action
type Factory interface {
	New(config *Config) Action
}

// Runner runs actions
type Runner interface {
	//DEPRECATED
	Run(context context.Context, action Action, uri string, options interface{}) (code int, data interface{}, err error)

	//Run the specified Action
	RunAction(ctx context.Context, act Action, options map[string]interface{}) (results map[string]*data.Attribute, err error)
}

// ResultHandler used to handle results from the Action
type ResultHandler interface {

	HandleResult(results map[string]*data.Attribute, err error)

	Done()
}
