package data

import (
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

type Resolver interface {
	Resolve(toResolve string, scope Scope) (value interface{}, err error)
}

var resolver = &BasicResolver{}

func GetBasicResolver() Resolver {
	return resolver
}

type BasicResolver struct {
}

func (r *BasicResolver) Resolve(toResolve string, scope Scope) (value interface{}, err error) {

	var details *ResolutionDetails

	if strings.HasPrefix(toResolve, "${") {
		details, err = GetResolutionDetailsOld(toResolve)
	} else if strings.HasPrefix(toResolve, "$") {
		details, err = GetResolutionDetails(toResolve[1:])
	} else {

		if scope == nil {
			//todo is this what we should do in this circumstance? or throw an error?
			return toResolve, nil
		}

		return SimpleScopeResolve(toResolve, scope)
	}

	if err != nil {
		return nil, err
	}

	if details == nil {
		return nil, fmt.Errorf("unable to resolve '%s'", toResolve)
	}

	var exists bool

	switch details.ResolverName {
	case "property":
		// Property resolution
		provider := GetPropertyProvider()
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
	case ".":
		//Current scope resolution
		attr, exists := scope.GetAttr(details.Property)
		if !exists {
			return nil, fmt.Errorf("failed to resolve current scope: '%s', not found in scope", details.Property)
		}
		value = attr.Value()
	default:
		return nil, fmt.Errorf("unsupported resolver: %s", details.ResolverName)
	}

	if details.Path != "" {
		value, err = PathGetValue(value, details.Path)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	return value, nil
}

func SimpleScopeResolve(toResolve string, scope Scope) (value interface{}, err error) {
	//idx := strings.Index(toResolve, ".")
	idx := strings.IndexFunc(toResolve, isSep)

	if idx != -1 {
		attr, found := scope.GetAttr(toResolve[:idx])
		if !found {
			return nil, fmt.Errorf("could not resolve '%s'", toResolve)
		}
		value, err := PathGetValue(attr.Value(), toResolve[idx:])
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		return value, nil

	} else {
		attr, found := scope.GetAttr(toResolve)
		if !found {
			return nil, fmt.Errorf("could not resolve '%s'", toResolve)
		}

		return attr.Value(), nil
	}
}

type ResolutionDetails struct {
	ResolverName string
	Item         string
	Property     string
	Path         string
}

func GetResolutionDetails(toResolve string) (*ResolutionDetails, error) {

	//todo optimize, maybe tokenize first
	details := &ResolutionDetails{}

	bracketIdx := strings.Index(toResolve, "]")
	exprLen := len(toResolve)
	if bracketIdx == (exprLen-1) && (strings.HasPrefix(toResolve, "property") || strings.HasPrefix(toResolve, "env")) {
		//$property[] or $env[] resolution
		itemIdx := strings.Index(toResolve, "[")
		details.ResolverName = toResolve[:itemIdx]
		details.Property = toResolve[itemIdx+1 : exprLen-1]
		return details, nil
	}

	dotIdx := strings.Index(toResolve, ".")
	if dotIdx == -1 {
		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
	}

	itemIdx := strings.Index(toResolve[:dotIdx], "[")

	if itemIdx != -1 {
		details.Item = toResolve[itemIdx+1 : dotIdx-1]
		details.ResolverName = toResolve[:itemIdx]
	} else {
		//For the case to get current scope atribute data
		if strings.HasPrefix(toResolve, "$.") || strings.HasPrefix(toResolve, ".") {
			details.ResolverName = toResolve[:dotIdx+1]
		} else {
			details.ResolverName = toResolve[:dotIdx]
		}
		//special case for activity without brackets
		if strings.HasPrefix(toResolve, "activity") {
			nextDot := strings.Index(toResolve[dotIdx+1:], ".") + dotIdx + 1
			details.Item = toResolve[dotIdx+1 : nextDot]
			dotIdx = nextDot
		}
	}

	pathIdx := strings.IndexFunc(toResolve[dotIdx+1:], isSep)

	if pathIdx != -1 {
		pathStart := pathIdx + dotIdx + 1
		details.Path = toResolve[pathStart:]
		details.Property = toResolve[dotIdx+1 : pathStart]
	} else {
		details.Property = toResolve[dotIdx+1:]
	}

	return details, nil
}

func GetResolutionDetailsOld(toResolve string) (*ResolutionDetails, error) {

	//todo optimize, maybe tokenize first

	closeIdx := strings.Index(toResolve, "}")

	if len(toResolve) < 4 || closeIdx == -1 {
		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
	}

	details := &ResolutionDetails{}

	dotIdx := strings.Index(toResolve, ".")

	if dotIdx == -1 {
		return nil, fmt.Errorf("invalid resolution expression [%s]", toResolve)
	}

	details.ResolverName = toResolve[2:dotIdx]

	if details.ResolverName == "activity" {
		nextDot := strings.Index(toResolve[dotIdx+1:], ".") + dotIdx + 1
		details.Item = toResolve[dotIdx+1 : nextDot]
		dotIdx = nextDot
	}
	details.Property = toResolve[dotIdx+1 : closeIdx]

	if closeIdx+1 < len(toResolve) {
		details.Path = toResolve[closeIdx+1:]
	}

	return details, nil
}

func isSep(r rune) bool {
	return r == '.' || r == '['
}

func GetValueWithResolver(valueMap map[string]interface{}, key string) (interface{}, bool) {

	val, exists := valueMap[key]

	if !exists {
		return nil, false
	}

	strVal, ok := val.(string)

	if ok {
		if strVal == "" {
			return "", true
		}

		if strVal[0] == '$' {

			v, err := GetBasicResolver().Resolve(strVal, nil)
			if err != nil {
				if strings.HasPrefix(err.Error(), "unsupported resolver") {
					return val, true
				}
				//todo double check this case
				return val, true
			}

			return v, true
		} else {
			return val, true
		}
	}

	return val, true
}
