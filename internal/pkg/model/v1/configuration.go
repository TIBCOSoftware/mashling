package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	mashling "github.com/TIBCOSoftware/mashling/cli/app"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/cache"
	gwerrors "github.com/TIBCOSoftware/mashling/internal/pkg/model/errors"
	"github.com/TIBCOSoftware/mashling/lib/types"
	"github.com/TIBCOSoftware/mashling/pkg/files"
)

var (
	githubRawContent = "https://raw.githubusercontent.com"
)

// LoadGateway loads a V1 Gateway app instance.
func LoadGateway(configuration []byte) (*Gateway, error) {
	gw := &Gateway{}
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
		log.Println("[mashling] Post processed configuration contents found in cache")
	} else {
		log.Println("[mashling] Post processed configuration contents *not* found in cache, processing now...")
		isValidJSON, verr := mashling.IsValidGateway(string(configuration))
		if verr != nil {
			return gw, verr
		}
		if !isValidJSON {
			return gw, errors.New("invalid gateway schema")
		}

		mashed := &types.Microgateway{}
		err = json.Unmarshal(configuration, mashed)
		if err != nil {
			return gw, err
		}

		// Cycle through and change refs to definitions
		for index, handler := range mashed.Gateway.EventHandlers {
			if handler.Reference != "" {
				action, ferr := LoadFlogoFlow(handler.Reference)
				if ferr != nil {
					return gw, ferr
				}
				handler.Definition = action
				handler.Reference = ""
				mashed.Gateway.EventHandlers[index] = handler
			}
		}
		gerrs, perr := validateConfigurationContents(mashed)
		if len(gerrs) > 0 {
			gw.ErrorDetails = gerrs
			return gw, errors.New("error validating contents of configuration")
		}
		if perr != nil {
			return gw, perr
		}
		flogoJSON, err = Translate(mashed)
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

func validateConfigurationContents(gateway *types.Microgateway) ([]gwerrors.Error, error) {
	// Right now we just handle missing dependencies.
	var deps []string
	var gerrs []gwerrors.Error
	var err error
	// Check trigger types for satisfied dependencies.
	for _, trigger := range gateway.Gateway.Triggers {
		if _, exists := registry.SupportedImports[trigger.Type]; !exists {
			deps = append(deps, trigger.Type)
		}
	}
	for _, handler := range gateway.Gateway.EventHandlers {
		if handler.Definition != nil {
			missingDeps, ferr := IdentifyMissingFlogoDependencies(handler.Definition)
			if ferr != nil {
				return gerrs, ferr
			}
			if missingDeps != nil {
				deps = append(deps, missingDeps...)
			}
		}
	}
	if deps != nil {
		gerrs = append(gerrs, &gwerrors.MissingDependency{MissingDependencies: deps})
	}
	return gerrs, err
}
