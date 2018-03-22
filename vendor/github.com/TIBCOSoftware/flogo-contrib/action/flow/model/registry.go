package model

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/util"
	"sync"
)

var (
	modelsMu     sync.RWMutex
	models       = make(map[string]*FlowModel)
	defaultModel *FlowModel
)

// Register registers the specified flow model
func Register(flowModel *FlowModel) {
	modelsMu.Lock()
	defer modelsMu.Unlock()

	if flowModel == nil {
		panic("model.Register: model cannot be nil")
	}

	id := flowModel.Name()

	if _, dup := models[id]; dup {
		panic("model.Register: model " + id + " already registered")
	}

	models[id] = flowModel
	util.RegisterModelValidator(id, flowModel)
}

// Registered gets all the registered flow models
func Registered() []*FlowModel {

	modelsMu.RLock()
	defer modelsMu.RUnlock()

	list := make([]*FlowModel, 0, len(models))

	for _, value := range models {
		list = append(list, value)
	}

	return list
}

// Get gets specified FlowModel
func Get(id string) *FlowModel {
	return models[id]
}

// Register registers the specified flow model
func RegisterDefault(model *FlowModel) {
	modelsMu.Lock()
	defer modelsMu.Unlock()

	if model == nil {
		panic("model.RegisterDefault: model cannot be nil")
	}

	id := model.Name()

	if _, dup := models[id]; !dup {
		models[id] = model
	}

	defaultModel = model
}

func Default() *FlowModel {
	return defaultModel
}
