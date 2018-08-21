package activity

import "fmt"

type Factory func(config *Config) (Activity, error)


var (
	activityFactories = make(map[string]Factory)
)

func RegisterFactory(ref string, f Factory) error {

	if len(ref) == 0 {
		return fmt.Errorf("'ref' must be specified when registering a activity factory")
	}

	if f == nil {
		return fmt.Errorf("cannot register 'nil' activity factory")
	}

	if activityFactories[ref] != nil {
		return fmt.Errorf("activity factory already registered for ref '%s'", ref)
	}

	activityFactories[ref] = f

	return nil
}

func GetFactory(ref string) Factory {
	return activityFactories[ref]
}
