package app

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"os"
)

// App is the configuration for the App
type Config struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Properties  map[string]interface{} `json:"properties"`
	Triggers    []*trigger.Config      `json:"triggers"`
	Actions     []*action.Config       `json:"actions"`
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
