package trigger

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"sync"
)

var (
	triggersMu sync.Mutex
	triggers   = make(map[string]TriggerDeprecated)
	reg        = &registry{}
)

type Registry interface {
	AddFactory(ref string, f Factory) error
	GetFactories() map[string]Factory
	AddInstance(id string, instance *TriggerInstance) error
	GetInstance(id string) *TriggerInstance
}

type registry struct {
	factories map[string]Factory
	instances map[string]*TriggerInstance
}

func GetRegistry() Registry {
	return reg
}

func RegisterFactory(ref string, f Factory) error {
	return reg.AddFactory(ref, f)
}

func (r *registry) AddFactory(ref string, f Factory) error {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	logger.Debugf("Registering trigger factory: '%s'", ref)

	if len(ref) == 0 {
		return fmt.Errorf("registry.RegisterFactory: ref is empty")
	}

	if f == nil {
		return fmt.Errorf("registry.RegisterFactory: factory is nil")
	}

	// copy on write to avoid synchronization on access
	newFs := make(map[string]Factory, len(r.factories))

	for k, v := range r.factories {
		newFs[k] = v
	}

	if newFs[ref] != nil {
		return fmt.Errorf("registry.RegisterFactory: already registered factory for ref '%s'", ref)
	}

	newFs[ref] = f

	r.factories = newFs

	return nil
}

func Factories() map[string]Factory {
	return reg.GetFactories()
}

// GetFactories returns a copy of the factories map
func (r *registry) GetFactories() map[string]Factory {

	newFs := make(map[string]Factory, len(r.factories))

	for k, v := range r.factories {
		newFs[k] = v
	}

	return newFs
}

func RegisterInstance(id string, inst *TriggerInstance) error {
	return reg.AddInstance(id, inst)
}

func (r *registry) AddInstance(id string, inst *TriggerInstance) error {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	if len(id) == 0 {
		return fmt.Errorf("registry.RegisterInstance: id is empty")
	}

	if inst == nil {
		return fmt.Errorf("registry.RegisterInstance: instance is nil")
	}

	// copy on write to avoid synchronization on access
	newInst := make(map[string]*TriggerInstance, len(r.instances))

	for k, v := range r.instances {
		newInst[k] = v
	}

	if newInst[id] != nil {
		return fmt.Errorf("registry.RegisterInstance: already registered instance for id '%s'", id)
	}

	newInst[id] = inst

	r.instances = newInst

	return nil
}

// Register registers the specified trigger
func Register(trigger TriggerDeprecated) {
	triggersMu.Lock()
	defer triggersMu.Unlock()

	if trigger == nil {
		panic("trigger.Register: trigger is nil")
	}
	id := trigger.Metadata().ID

	if _, dup := triggers[id]; dup {
		panic("trigger.Register: Register called twice for trigger " + id)
	}
	// copy on write to avoid synchronization on access
	newTriggers := make(map[string]TriggerDeprecated, len(triggers))

	for k, v := range triggers {
		newTriggers[k] = v
	}

	newTriggers[id] = trigger
	triggers = newTriggers
}

// Triggers gets all the registered triggers
func Triggers() []TriggerDeprecated {

	var curTriggers = triggers

	list := make([]TriggerDeprecated, 0, len(curTriggers))

	for _, value := range curTriggers {
		list = append(list, value)
	}

	return list
}

// Instance gets specified trigger instance
func Instance(id string) *TriggerInstance {
	return reg.GetInstance(id)
}

// GetInstances gets specified trigger instance
func (r registry) GetInstance(id string) *TriggerInstance {
	return r.instances[id]
}

// Get gets specified trigger
func Get(id string) TriggerDeprecated {
	//var curTriggers = triggers
	return triggers[id]
}

func GetTriggerInstanceInfo() []TriggerInstanceInfo {
	currentInstances := reg.instances
	list := make([]TriggerInstanceInfo, 0, len(currentInstances))
	for id, triggerInstance := range currentInstances {
		list = append(list, TriggerInstanceInfo{
			Name:   id,
			Status: triggerInstance.Status,
			Error:  triggerInstance.Error,
		})
	}
	return list
}