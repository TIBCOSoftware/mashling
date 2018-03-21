package v1

import (
	"log"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
)

// Gateway contains all data needed to run a v1 mashling gateway app.
type Gateway struct {
	FlogoApp      app.Config
	FlogoEngine   engine.Engine
	SchemaVersion string
	ErrorDetails  []gwerrors.Error
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

// Description returns the description specified in the Gateway configuration.
func (g *Gateway) Description() string {
	return g.FlogoApp.Description
}

// Errors returns the associated slice of ErrorDetails.
func (g *Gateway) Errors() []gwerrors.Error {
	return g.ErrorDetails
}
