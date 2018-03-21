package model

import (
	"sync"
)

var (
	modelsMu sync.RWMutex
	models   = make(map[string]*FlowModel)
)

// Register registers the specified flow model
func Register(model *FlowModel) {
	modelsMu.Lock()
	defer modelsMu.Unlock()

	if model == nil {
		panic("model.Register: model cannot be nil")
	}

	id := model.Name()

	if _, dup := models[id]; dup {
		panic("model.Register: model " + id + " already registered")
	}

	models[id] = model
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
