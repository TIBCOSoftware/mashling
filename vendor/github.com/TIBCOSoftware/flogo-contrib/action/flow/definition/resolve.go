package definition

import (
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

var resolver = &FlowResolver{}

func GetDataResolver() data.Resolver {
	return resolver
}

type FlowResolver struct {
}

func (r *FlowResolver) Resolve(toResolve string, scope data.Scope) (value interface{}, err error){

	var details *data.ResolutionDetails

	if strings.HasPrefix(toResolve,"${") {
		details,err = data.GetResolutionDetailsOld(toResolve)
	} else if strings.HasPrefix(toResolve, "$") {
		details,err = data.GetResolutionDetails(toResolve[1:])
	} else {
		return data.SimpleScopeResolve(toResolve, scope)
	}

	if err != nil {
		return nil, err
	}

	if details == nil {
		return nil, fmt.Errorf("unable to determine resolver for %s", toResolve)
	}

	var exists bool

	switch details.ResolverName {
	case "property":
		// Property resolution
		provider := data.GetPropertyProvider()
		value, exists = provider.GetProperty(details.Property + details.Path) //should we add the path and reset it to ""
		if !exists {
			err := fmt.Errorf("failed to resolve Property: '%s', ensure that property is configured in the application", details.Property)
			logger.Error(err.Error())
			return nil, err
		}
	case "env":
		// Environment resolution
		value, exists = os.LookupEnv(details.Property + details.Path)
		if !exists {
			err := fmt.Errorf("failed to resolve Environment Variable: '%s', ensure that variable is configured", details.Property)
			logger.Error(err.Error())
			return "", err
		}
	case "activity":
		attr, exists := scope.GetAttr("_A." + details.Item + "." +details.Property)
		if !exists {
			return nil, fmt.Errorf("failed to resolve activity attr: '%s', not found in flow", details.Property)
		}
		value = attr.Value()
	case "trigger":
		attr, exists := scope.GetAttr("_T." + details.Property)
		if !exists {
			return nil, fmt.Errorf("failed to resolve trigger attr: '%s', not found in flow", details.Property)
		}
		value = attr.Value()
	case "flow":
		attr, exists := scope.GetAttr(details.Property)
		if !exists {
			return nil, fmt.Errorf("failed to resolve flow attr: '%s', not found in flow", details.Property)
		}
		value = attr.Value()
	case "current":
		attr, exists := scope.GetAttr("$current." + details.Property)
		if !exists {
			return nil, fmt.Errorf("failed to resolve current working data: '%s', not found in scope", details.Property)
		}
		value = attr.Value()
	default:
		return nil, fmt.Errorf("unsupported resolver: %s", details.ResolverName)
	}

	if details.Path != "" {
		value, err = data.PathGetValue(value, details.Path)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return value, nil
}