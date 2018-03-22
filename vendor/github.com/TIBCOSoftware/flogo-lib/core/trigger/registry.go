package trigger

import (
	"fmt"
)

var (
	triggerFactories = make(map[string]Factory)
)

func RegisterFactory(ref string, f Factory) error {

	if len(ref) == 0 {
		return fmt.Errorf("'ref' must be specified when registering a trigger factory")
	}

	if f == nil {
		return fmt.Errorf("cannot register 'nil' trigger factory")
	}

	if triggerFactories[ref] != nil {
		return fmt.Errorf("trigger factory already registered for ref '%s'", ref)
	}

	triggerFactories[ref] = f

	return nil
}

func GetFactory(ref string) Factory {
	return triggerFactories[ref]
}

func Factories() map[string]Factory {
	//todo return copy
	return triggerFactories
}
