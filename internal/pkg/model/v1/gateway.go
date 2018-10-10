package v1

import (
	"encoding/json"
	"log"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/mashling/internal/pkg/consul"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/internal/pkg/services"
	"github.com/TIBCOSoftware/mashling/internal/pkg/swagger"
	"github.com/mashling/commons/lib/types"
)

// Gateway contains all data needed to run a v1 mashling gateway app.
type Gateway struct {
	FlogoApp       app.Config
	FlogoEngine    engine.Engine
	MashlingConfig interface{}
	SchemaVersion  string
	ErrorDetails   []gwerrors.Error
	pingEnabled    bool
	PingService    services.PingService
}

// Init initializes the Gateway.
func (g *Gateway) Init(pingEnabled bool, pingPort string) error {
	log.Println("[mashling] Initializing Flogo engine...")

	g.pingEnabled = pingEnabled
	//Initialize ping service if it is enabled
	if g.pingEnabled {
		//construct ping response
		pingResponse := services.PingResponse{
			Version:        g.Version(),
			Appversion:     g.AppVersion(),
			Appdescription: g.Description()}

		g.PingService = services.GetPingService()
		g.PingService.Init(pingPort, pingResponse)
	}
	return g.FlogoEngine.Init(true)
}

// Start starts the Gateway.
func (g *Gateway) Start() error {
	log.Println("[mashling] Starting Flogo engine...")
	err := g.FlogoEngine.Start()
	if err != nil {
		return err
	}

	//Start ping service if it is enabled
	if g.pingEnabled {
		log.Println("[mashling] Starting Ping service...")
		err = g.PingService.Start()
		if err != nil {
			g.Stop()
			return err
		}
	}
	return nil
}

// Stop stops the Gateway.
func (g *Gateway) Stop() error {
	log.Println("[mashling] Stoppping Flogo engine...")
	if g.pingEnabled {
		log.Println("[mashling] Stoppping Ping service...")
		g.PingService.Stop()
	}
	return g.FlogoEngine.Stop()
}

// Version returns the current schema version used to configure the Gateway.
func (g *Gateway) Version() string {
	return g.SchemaVersion
}

// AppVersion returns the version specified in the Gateway configuration.
func (g *Gateway) AppVersion() string {
	return g.FlogoApp.Version
}

// Name returns the name specified in the Gateway configuration.
func (g *Gateway) Name() string {
	return g.FlogoApp.Name
}

// Description returns the description specified in the Gateway configuration.
func (g *Gateway) Description() string {
	return g.FlogoApp.Description
}

// Errors returns the associated slice of ErrorDetails.
func (g *Gateway) Errors() []gwerrors.Error {
	return g.ErrorDetails
}

// Configuration returns the user provided Mashling configuration file contents.
func (g *Gateway) Configuration() interface{} {
	return g.MashlingConfig
}

// Swagger returns Swagger 2.0 docs based off of the triggers defined for this Gateway.
func (g *Gateway) Swagger(hostname string, triggerName string) ([]byte, error) {
	gConf := g.MashlingConfig.(types.Microgateway)
	var endpoints []swagger.Endpoint
	for _, trigger := range gConf.Gateway.Triggers {
		if triggerName == "" || triggerName == trigger.Name {
			if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
				var endpoint swagger.Endpoint
				endpoint.Name = trigger.Name
				endpoint.Description = trigger.Description
				err := json.Unmarshal(trigger.Settings, &endpoint)
				if err != nil {
					return nil, err
				}
				var beginDelim, endDelim rune
				switch trigger.Type {
				case "github.com/TIBCOSoftware/flogo-contrib/trigger/rest":
					beginDelim = ':'
					endDelim = '/'
				case "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger":
					beginDelim = '{'
					endDelim = '}'
				default:
					beginDelim = '{'
					endDelim = '}'
				}
				endpoint.BeginDelim = beginDelim
				endpoint.EndDelim = endDelim
				endpoints = append(endpoints, endpoint)
			}
		}
	}
	return swagger.Generate(hostname, g.Name(), g.Description(), g.AppVersion(), endpoints)
}

// ConsulServiceDefinition returns Consul compatible service definitions.
func (g *Gateway) ConsulServiceDefinition() ([]consul.ServiceDefinition, error) {

	gConf := g.MashlingConfig.(types.Microgateway)

	//load the configuration, if provided.
	configNamedMap := make(map[string]types.Config)
	for _, config := range gConf.Gateway.Configurations {
		configNamedMap[config.Name] = config
	}

	var consulServiceDefinitions []consul.ServiceDefinition
	for _, trigger := range gConf.Gateway.Triggers {
		if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
			var consulServiceDefinition consul.ServiceDefinition

			var mtSettings interface{}
			if err := json.Unmarshal([]byte(trigger.Settings), &mtSettings); err != nil {
				return nil, err
			}

			//resolve any configuration references if the "config" param is set in the settings
			mashTriggerSettings := mtSettings.(map[string]interface{})
			mashTriggerSettingsUsable := mtSettings.(map[string]interface{})
			for k, v := range mashTriggerSettings {
				mashTriggerSettingsUsable[k] = v
			}

			if configNamedMap != nil && len(configNamedMap) > 0 {
				//inherit the configuration settings if the trigger uses configuration reference
				err := resolveConfigurationReference(configNamedMap, trigger, mashTriggerSettingsUsable)
				if err != nil {
					return nil, err
				}
			}

			settings, err := json.MarshalIndent(&mashTriggerSettingsUsable, "", "    ")
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(settings, &consulServiceDefinition)
			if err != nil {
				return nil, err
			}
			consulServiceDefinition.Name = trigger.Name
			consulServiceDefinitions = append(consulServiceDefinitions, consulServiceDefinition)
		}

	}

	return consulServiceDefinitions, nil
}
