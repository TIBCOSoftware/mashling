package v2

import (
	"log"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/mashling/internal/pkg/consul"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/internal/pkg/services"
	"github.com/TIBCOSoftware/mashling/internal/pkg/swagger"
)

// Gateway contains all data needed to run a v2 mashling gateway app.
type Gateway struct {
	FlogoApp       app.Config
	FlogoEngine    engine.Engine
	MashlingConfig interface{}
	SchemaVersion  string
	ErrorDetails   []gwerrors.Error
	PingService    services.PingService
}

// Init initializes the Gateway.
func (g *Gateway) Init(pingPort string) error {
	log.Println("[mashling] Initializing Flogo engine...")

	g.PingService = services.GetPingService()
	g.PingService.Init(pingPort, g.Version(), g.AppVersion(), g.Description())

	return g.FlogoEngine.Init(true)
}

// Start starts the Gateway.
func (g *Gateway) Start() error {
	log.Println("[mashling] Starting Flogo engine...")
	g.PingService.Start()
	return g.FlogoEngine.Start()
}

// Stop stops the Gateway.
func (g *Gateway) Stop() error {
	log.Println("[mashling] Stoppping Flogo engine...")
	g.PingService.Stop()
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
	gConf := g.MashlingConfig.(types.Schema)
	var endpoints []swagger.Endpoint
	for _, trigger := range gConf.Gateway.Triggers {
		if triggerName == "" || triggerName == trigger.Name {
			if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
				for _, handler := range trigger.Handlers {
					var endpoint swagger.Endpoint
					endpoint.Name = trigger.Name
					endpoint.Method = handler.Settings["method"].(string)
					endpoint.Path = handler.Settings["path"].(string)
					endpoint.Description = trigger.Description
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
	}
	return swagger.Generate(hostname, g.Name(), g.Description(), g.AppVersion(), endpoints)
}

// ConsulServiceDefinition returns Consul compatible service definitions.
func (g *Gateway) ConsulServiceDefinition() ([]consul.ServiceDefinition, error) {
	gConf := g.MashlingConfig.(types.Schema)
	var consulServiceDefinitions []consul.ServiceDefinition
	for _, trigger := range gConf.Gateway.Triggers {
		if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
			var consulServiceDefinition consul.ServiceDefinition
			consulServiceDefinition.Name = trigger.Name
			consulServiceDefinition.Port = trigger.Settings["port"].(string)
			consulServiceDefinitions = append(consulServiceDefinitions, consulServiceDefinition)
		}
	}
	return consulServiceDefinitions, nil
}
