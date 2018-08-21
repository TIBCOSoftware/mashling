package app

import (
	"encoding/json"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"io/ioutil"
	"regexp"
	"strings"
)

// App is the configuration for the App
type Config struct {
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Version     string             `json:"version"`
	Description string             `json:"description"`

	Properties  []*data.Attribute  `json:"properties"`
	Channels    []string           `json:"channels"`
	Triggers    []*trigger.Config  `json:"triggers"`
	Resources   []*resource.Config `json:"resources"`
	Actions     []*action.Config   `json:"actions"`
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
	return LoadConfig("")
}

func LoadConfig(flogoJson string) (*Config, error) {
	if flogoJson == "" {
		configPath := config.GetFlogoConfigPath()

		flogo, err := os.Open(configPath)
		if err != nil {
			return nil, err
		}

		file, err := ioutil.ReadAll(flogo)
		if err != nil {
			return nil, err
		}

		updated, err := preprocessConfig(file)
		if err != nil {
			return nil, err
		}

		app := &Config{}
		err = json.Unmarshal(updated, &app)
		if err != nil {
			return nil, err
		}
		return app, nil
	} else {
		updated, err := preprocessConfig([]byte(flogoJson))
		if err != nil {
			return nil, err
		}

		app := &Config{}
		err = json.Unmarshal(updated, &app)
		if err != nil {
			return nil, err
		}
		return app, nil
	}
}

func preprocessConfig(appJson []byte) ([]byte, error) {

	// For now decode secret values
	re := regexp.MustCompile("SECRET:[^\\\\\"]*")
	for _, match := range re.FindAll(appJson, -1) {
		encodedValue := string(match[7:])
		decodedValue, err := data.GetSecretValueHandler().DecodeValue(encodedValue)
		if err != nil {
			return nil, err
		}
		appstring := strings.Replace(string(appJson), string(match), decodedValue, -1)
		appJson = []byte(appstring)
	}

	return appJson, nil
}

func GetProperties(properties []*data.Attribute) (map[string]interface{}, error) {

	props := make(map[string]interface{})
	if properties != nil {
		for _, property := range properties {
			pValue := property.Value()
			strValue, ok := pValue.(string)
			if ok {
				if strValue != "" && strValue[0] == '$' {
					// Needs resolution
					resolvedValue, err := data.GetBasicResolver().Resolve(strValue, nil)
					if err != nil {
						return props, err
					}
					pValue = resolvedValue
				}
			}
			value, err := data.CoerceToValue(pValue, property.Type())
			if err != nil {
				return props, err
			}
			props[property.Name()] = value
		}
		return props, nil
	}

	return props, nil
}

//used for old action config

//func FixUpApp(cfg *Config) {
//
//	if cfg.Resources != nil || cfg.Actions == nil {
//		//already new app format
//		return
//	}
//
//	idToAction := make(map[string]*action.Config)
//	for _, act := range cfg.Actions {
//		idToAction[act.Id] = act
//	}
//
//	for _, trg := range cfg.Triggers {
//		for _, handler := range trg.Handlers {
//
//			oldAction := idToAction[handler.ActionId]
//
//			newAction := &action.Config{Ref: oldAction.Ref}
//
//			if oldAction != nil {
//				newAction.Mappings = oldAction.Mappings
//			} else {
//				if handler.ActionInputMappings != nil {
//					newAction.Mappings = &data.IOMappings{}
//					newAction.Mappings.Input = handler.ActionInputMappings
//					newAction.Mappings.Output = handler.ActionOutputMappings
//				}
//			}
//
//			newAction.Data = oldAction.Data
//			newAction.Metadata = oldAction.Metadata
//
//			handler.Action = newAction
//		}
//	}
//
//	cfg.Actions = nil
//}
