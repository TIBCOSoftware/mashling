package engine

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/config"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/runner"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

// Config is the configuration for the engine
type Config struct {
	LogLevel         string
	RunnerConfig     *RunnerConfig
	ValidateTriggers bool
	Services         map[string]*util.ServiceConfig
}

type serEngineConfig struct {
	LogLevel     string                `json:"loglevel"`
	DtValidation bool                  `json:"disableTriggerValidation,omitempty"`
	RunnerConfig *RunnerConfig         `json:"actionRunner"`
	Services     []*util.ServiceConfig `json:"services"`
}

// RunnerConfig is the configuration for the engine level runner
type RunnerConfig struct {
	Type   string               `json:"type"`
	Pooled *runner.PooledConfig `json:"pooled,omitempty"`
}

// TriggersConfig is the triggers configuration for the engine
type TriggersConfig struct {
	Triggers map[string]*trigger.Config
}

type serTriggersConfig struct {
	Triggers []*trigger.Config `json:"triggers"`
}

// DefaultConfig returns the default engine configuration
func DefaultConfig() *Config {

	var engineConfig Config

	engineConfig.LogLevel = config.GetLogLevel()
	engineConfig.RunnerConfig = defaultRunnerConfig()

	return &engineConfig
}

// DefaultTriggersConfig returns the default triggers configuration
func DefaultTriggersConfig() *TriggersConfig {

	var triggersConfig TriggersConfig
	triggersConfig.Triggers = make(map[string]*trigger.Config)

	return &triggersConfig
}

// MarshalJSON marshals the EngineConfig to JSON
func (ec *Config) MarshalJSON() ([]byte, error) {

	var services []*util.ServiceConfig

	for _, value := range ec.Services {
		services = append(services, value)
	}

	return json.Marshal(&serEngineConfig{
		LogLevel:     ec.LogLevel,
		DtValidation: !ec.ValidateTriggers,
		RunnerConfig: ec.RunnerConfig,
		Services:     services,
	})
}

// UnmarshalJSON unmarshals EngineConfog from JSON
func (ec *Config) UnmarshalJSON(data []byte) error {

	ser := &serEngineConfig{}
	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	ec.LogLevel = ser.LogLevel
	ec.ValidateTriggers = !ser.DtValidation

	if ser.RunnerConfig != nil {
		ec.RunnerConfig = ser.RunnerConfig
	} else {
		ec.RunnerConfig = defaultRunnerConfig()
	}

	if ser.Services != nil {
		ec.Services = make(map[string]*util.ServiceConfig)

		for _, value := range ser.Services {
			ec.Services[value.Name] = value
		}
	}

	return nil
}

// MarshalJSON marshals the EngineConfig to JSON
func (tc *TriggersConfig) MarshalJSON() ([]byte, error) {

	var triggers []*trigger.Config

	for _, value := range tc.Triggers {
		triggers = append(triggers, value)
	}

	return json.Marshal(&serTriggersConfig{
		Triggers: triggers,
	})
}

// UnmarshalJSON unmarshals EngineConfog from JSON
func (tc *TriggersConfig) UnmarshalJSON(data []byte) error {

	ser := &serTriggersConfig{}
	if err := json.Unmarshal(data, ser); err != nil {
		return err
	}

	tc.Triggers = make(map[string]*trigger.Config)

	for _, value := range ser.Triggers {
		tc.Triggers[value.Name] = value
	}

	return nil
}

// LoadConfigFromFile loads the engine Config from the specified JSON file
func LoadConfigFromFile(fileName string) *Config {

	if len(fileName) == 0 {
		panic("file name cannot be empty")
	}

	configFile, _ := os.Open(fileName)

	if configFile != nil {

		engineConfig := &Config{}

		decoder := json.NewDecoder(configFile)
		decodeErr := decoder.Decode(engineConfig)
		if decodeErr != nil {
			err := fmt.Errorf("Error decoding %s - %s", fileName, decodeErr.Error())
			panic(err)
		}

		// Quick fix until we refactor the runnerConfig in flogo.json
		if engineConfig.RunnerConfig != nil {
			engineConfig.RunnerConfig.Pooled = NewPooledConfig()
		}

		engineConfig.LogLevel = config.GetLogLevel()

		return engineConfig
	}

	return nil
}

// LoadConfigFromJSON loads the engine Config from the specified JSON file
func LoadConfigFromJSON(configJSON string) *Config {

	engineConfig := &Config{}
	decodeErr := json.Unmarshal([]byte(configJSON), engineConfig)
	if decodeErr != nil {
		err := fmt.Errorf("Error decoding %s - %s", "engineConfig", decodeErr.Error())
		panic(err)
	}

	// Quick fix until we refactor the runnerConfig in flogo.json
	if engineConfig.RunnerConfig != nil {
		engineConfig.RunnerConfig.Pooled = NewPooledConfig()
	}

	engineConfig.LogLevel = config.GetLogLevel()

	return engineConfig
}

// LoadTriggersConfigFromFile loads the triggers Config from the specified JSON file
func LoadTriggersConfigFromFile(fileName string) *TriggersConfig {

	if len(fileName) == 0 {
		panic("file name cannot be empty")
	}

	configFile, _ := os.Open(fileName)

	if configFile != nil {

		triggersConfig := &TriggersConfig{}

		decoder := json.NewDecoder(configFile)
		decodeErr := decoder.Decode(triggersConfig)
		if decodeErr != nil {
			err := fmt.Errorf("Error decoding %s - %s", fileName, decodeErr.Error())
			panic(err)
		}

		return triggersConfig
	}

	return nil
}

// LoadTriggersConfigFromJSON loads the engine Config from the specified JSON file
func LoadTriggersConfigFromJSON(configJSON string) *TriggersConfig {

	triggersConfig := &TriggersConfig{}
	decodeErr := json.Unmarshal([]byte(configJSON), triggersConfig)
	if decodeErr != nil {
		err := fmt.Errorf("Error decoding %s - %s", "triggersConfig", decodeErr.Error())
		panic(err)
	}

	return triggersConfig
}

func defaultRunnerConfig() *RunnerConfig {
	return &RunnerConfig{Type: "pooled", Pooled: NewPooledConfig()}
}

//NewPooledConfig creates a new Pooled config, looks for environment variables to override default values
func NewPooledConfig() *runner.PooledConfig {
	return &runner.PooledConfig{NumWorkers: config.GetRunnerWorkers(), WorkQueueSize: config.GetRunnerQueueSize()}
}
