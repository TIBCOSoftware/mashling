package service

import (
	"errors"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

// Service encapsulates everything necessary to execute a step against a target.
type Service interface {
	Execute() (err error)
	UpdateRequest(values map[string]interface{}) (err error)
}

// Initialize sets up the service based off of the service definition.
func Initialize(serviceDef types.Service) (service Service, err error) {
	switch sType := serviceDef.Type; sType {
	case "http":
		return InitializeHTTP(serviceDef.Settings)
	case "js":
		return InitializeJS(serviceDef.Settings)
	case "flogoActivity":
		return InitializeFlogoActivity(serviceDef.Settings)
	case "flogoFlow":
		return InitializeFlogoFlow(serviceDef.Settings)
	case "sqld":
		return InitializeSQLD(serviceDef.Settings)
	case "anomaly":
		return InitializeAnomaly(serviceDef.Settings)
	case "jwt":
		return InitializeJWT(serviceDef.Settings)
	default:
		return nil, errors.New("unknown service type")
	}
}
