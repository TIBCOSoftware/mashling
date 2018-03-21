package action

import (
	"errors"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"sync"
)

var (
	factoryMu sync.Mutex
	factories = make(map[string]Factory)
	actionMu  sync.Mutex
	actions   = make(map[string]Action)
)

func GetFactory(ref string) Factory {
	factoryMu.Lock()
	defer factoryMu.Unlock()

	return factories[ref]
}

func RegisterFactory(ref string, factory Factory) error {
	factoryMu.Lock()
	defer factoryMu.Unlock()

	logger.Debugf("Registering action factory: '%s'", ref)

	if len(ref) == 0 {
		return errors.New("RegisterFactory: ref is empty")
	}

	if factory == nil {
		return errors.New("RegisterFactory: factory is nil")
	}

	if factories[ref] != nil {
		return fmt.Errorf("RegisterFactory: already registered factory for ref '%s'", ref)
	}

	factories[ref] = factory

	return nil
}

func Factories() map[string]Factory {
	factoryMu.Lock()
	defer factoryMu.Unlock()

	factoriesCopy := make(map[string]Factory, len(factories))

	for k, v := range factories {
		factoriesCopy[k] = v
	}

	return factoriesCopy
}

func Get(id string) Action {
	actionMu.Lock()
	defer actionMu.Unlock()

	return actions[id]
}

func Register(id string, act Action) error {
	actionMu.Lock()
	defer actionMu.Unlock()

	if len(id) == 0 {
		return fmt.Errorf("error registering action, id is empty")
	}

	if act == nil {
		return fmt.Errorf("error registering action for id '%s', action is nil", id)
	}

	if actions[id] != nil {
		return fmt.Errorf("Error registering action, action already registered for id '%s'", id)
	}

	// copy on write to avoid synchronization on access
	actionsCopy := make(map[string]Action, len(actions)+1)

	for k, v := range actions {
		actionsCopy[k] = v
	}

	actionsCopy[id] = act

	actions = actionsCopy

	return nil
}

func Actions() map[string]Action {
	actionMu.Lock()
	defer actionMu.Unlock()

	actionsCopy := make(map[string]Action, len(actions))

	for id, act := range actions {
		actionsCopy[id] = act
	}

	return actionsCopy
}
