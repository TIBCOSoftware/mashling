package app

import (
	"encoding/json"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

// App is the configuration for the App
type Config struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
	Triggers    []*trigger.Config      `json:"triggers"`
	Resources   []*resource.Config     `json:"resources"`

	//for backwards compatibility
	Actions []*action.Config `json:"actions"`
}

// defaultConfigProvider implementation of ConfigProvider
type defaultConfigProvider struct {
}

// ConfigProvider interface to implement to provide the app configuration
type ConfigProvider interface {
	GetApp() (*Config, error)
}

// DefaultSerializer returns the default App Serializer
func DefaultConfigProvider() ConfigProvider {
	return &defaultConfigProvider{}
}

// GetApp returns the app configuration
func (d *defaultConfigProvider) GetApp() (*Config, error) {

	configPath := config.GetFlogoConfigPath()

	flogo, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(flogo)
	app := &Config{}
	err = jsonParser.Decode(&app)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func FixUpApp(cfg *Config) {

	if cfg.Resources != nil || cfg.Actions == nil {
		//already new app format
		return
	}

	idToAction := make(map[string]*action.Config)
	for _, act := range cfg.Actions {
		idToAction[act.Id] = act
	}

	for _, trg := range cfg.Triggers {
		for _, handler := range trg.Handlers {

			oldAction := idToAction[handler.ActionId]

			newAction := &action.Config{Ref: oldAction.Ref}

			if oldAction != nil {
				newAction.Mappings = oldAction.Mappings
			} else {
				if handler.ActionInputMappings != nil {
					newAction.Mappings = &data.IOMappings{}
					newAction.Mappings.Input = handler.ActionInputMappings
					newAction.Mappings.Output = handler.ActionOutputMappings
				}
			}

			newAction.Data = oldAction.Data
			newAction.Metadata = oldAction.Metadata

			handler.Action = newAction
		}
	}

	cfg.Actions = nil
}
