package action

import (
	"fmt"
)

var (
	actionFactories = make(map[string]Factory)

	// Deprecated: No longer used
	actions = make(map[string]Action)
)

func RegisterFactory(ref string, f Factory) error {

	if len(ref) == 0 {
		return fmt.Errorf("'ref' must be specified when registering an action factory")
	}

	if f == nil {
		return fmt.Errorf("cannot register 'nil' action factory")
	}

	if actionFactories[ref] != nil {
		return fmt.Errorf("action factory already registered for ref '%s'", ref)
	}

	actionFactories[ref] = f

	return nil
}

func GetFactory(ref string) Factory {
	return actionFactories[ref]
}

func Factories() map[string]Factory {
	//todo return copy
	return actionFactories
}

// Deprecated: No longer used
func Get(id string) Action {

	return actions[id]
}

// Deprecated: No longer used
func Register(id string, act Action) error {

	if len(id) == 0 {
		return fmt.Errorf("error registering action, id is empty")
	}

	if act == nil {
		return fmt.Errorf("error registering action for id '%s', action is nil", id)
	}

	if actions[id] != nil {
		return fmt.Errorf("error registering action, action already registered for id '%s'", id)
	}

	actions[id] = act

	return nil
}

// Deprecated: No longer used
func Actions() map[string]Action {

	actionsCopy := make(map[string]Action, len(actions))

	for id, act := range actions {
		actionsCopy[id] = act
	}

	return actionsCopy
}
