package v2

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/cache"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v1"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/schema"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/pkg/files"
)

// LoadGateway loads a V2 Gateway app instance.
func LoadGateway(configuration []byte) (*Gateway, error) {
	gw := &Gateway{}
	gateway := &types.Schema{}
	gw.SchemaVersion = Version
	var flogoJSON []byte
	key, err := files.ChecksumContents(configuration)
	if err != nil {
		return gw, err
	}
	if cache.Enabled && cache.Cache.InCache(key) {
		flogoJSON, err = cache.Cache.LoadFromCache(key)
		if err != nil {
			return gw, err
		}
		err = json.Unmarshal(configuration, gateway)
		if err != nil {
			return gw, err
		}
		log.Println("[mashling] Post processed configuration contents found in cache")
	} else {
		log.Println("[mashling] Post processed configuration contents *not* found in cache, processing now...")
		err = schema.Validate(configuration)
		if err != nil {
			return gw, err
		}
		err = json.Unmarshal(configuration, gateway)
		if err != nil {
			return gw, err
		}
		// Load remote Flogo flow references if they exist.
		for index, service := range gateway.Gateway.Services {
			if service.Type != "flogoFlow" {
				continue
			}
			if reference, ok := service.Settings["reference"].(string); ok {
				action, ferr := v1.LoadFlogoFlow(reference)
				if ferr != nil {
					return gw, ferr
				}
				service.Settings["definition"] = action
				gateway.Gateway.Services[index] = service
			} else if definition, ok := service.Settings["definition"].(map[string]interface{}); ok {
				rawAction, jerr := json.Marshal(definition)
				if jerr != nil {
					return gw, jerr
				}
				service.Settings["definition"] = json.RawMessage(rawAction)
				gateway.Gateway.Services[index] = service
			}
		}
		gerrs, ferr := validateConfigurationContents(gateway)
		if ferr != nil {
			gw.ErrorDetails = gerrs
			return gw, ferr
		}
		if len(gerrs) > 0 {
			gw.ErrorDetails = gerrs
			return gw, errors.New("error validating contents of configuration")
		}
		flogoJSON, err = Translate(gateway)
		if err != nil {
			return gw, err
		}
		if cache.Enabled {
			err = cache.Cache.WriteToCache(key, flogoJSON)
			if err != nil {
				return gw, err
			}
			log.Println("[mashling] Post processed configuration contents written to cache")
		}
	}
	gw.MashlingConfig = *gateway
	jsonParser := json.NewDecoder(bytes.NewReader(flogoJSON))
	gw.FlogoApp = app.Config{}
	err = jsonParser.Decode(&gw.FlogoApp)
	if err != nil {
		return gw, err
	}

	gw.FlogoEngine, err = engine.New(&gw.FlogoApp)
	if err != nil {
		return gw, err
	}
	return gw, nil
}

func validateConfigurationContents(gateway *types.Schema) ([]gwerrors.Error, error) {
	// Right now we just handle missing dependencies.
	var deps []string
	var gerrs []gwerrors.Error
	// Check trigger types for satisfied dependencies.
	for _, trigger := range gateway.Gateway.Triggers {
		if _, exists := registry.SupportedImports[trigger.Type]; !exists {
			deps = append(deps, trigger.Type)
		}
	}
	for _, service := range gateway.Gateway.Services {
		if service.Type == "flogoActivity" {
			if ref, ok := service.Settings["ref"].(string); ok {
				if _, exists := registry.SupportedImports[ref]; !exists {
					deps = append(deps, ref)
				}
			}
		} else if service.Type == "flogoFlow" {
			if definition, ok := service.Settings["definition"].(json.RawMessage); ok {
				missingDeps, ferr := v1.IdentifyMissingFlogoDependencies(definition)
				if ferr != nil {
					return gerrs, ferr
				}
				if missingDeps != nil {
					deps = append(deps, missingDeps...)
				}
			} else {
				return gerrs, errors.New("missing Flogo flow definition")
			}
		}
	}
	if deps != nil {
		gerrs = append(gerrs, &gwerrors.MissingDependency{MissingDependencies: deps})
	}
	return gerrs, nil
}
