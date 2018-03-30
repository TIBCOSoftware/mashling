package v1

import (
	"encoding/json"
	"log"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/mashling/internal/pkg/consul"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/internal/pkg/swagger"
	"github.com/TIBCOSoftware/mashling/lib/types"
)

// Gateway contains all data needed to run a v1 mashling gateway app.
type Gateway struct {
	FlogoApp       app.Config
	FlogoEngine    engine.Engine
	MashlingConfig interface{}
	SchemaVersion  string
	ErrorDetails   []gwerrors.Error
}

// Init initializes the Gateway.
func (g *Gateway) Init() error {
	log.Println("[mashling] Initializing Flogo engine...")
	return g.FlogoEngine.Init(true)
}

// Start starts the Gateway.
func (g *Gateway) Start() error {
	log.Println("[mashling] Starting Flogo engine...")
	return g.FlogoEngine.Start()
}

// Stop stops the Gateway.
func (g *Gateway) Stop() error {
	log.Println("[mashling] Stoppping Flogo engine...")
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
	var consulServiceDefinitions []consul.ServiceDefinition
	for _, trigger := range gConf.Gateway.Triggers {
		if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
			var consulServiceDefinition consul.ServiceDefinition

			settings, err := json.MarshalIndent(&trigger.Settings, "", "    ")
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(settings, &consulServiceDefinition)
			if err != nil {
				return nil, err
			}
			consulServiceDefinition.Name = trigger.Name
		}
	}
	return consulServiceDefinitions, nil
}
