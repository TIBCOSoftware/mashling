package function

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("function-registry")

type Function interface {
	GetName() string
	GetCategory() string
}

var (
	functionsMu sync.Mutex
	functions   = make(map[string]Function)
)

func Registry(f Function) {
	functionsMu.Lock()

	defer functionsMu.Unlock()

	if f == nil {
		log.Errorf("Cannot rregistry nil function")
		return
	}

	log.Debugf("Registry function name %s tag %s", f.GetName(), f.GetCategory())

	var registeName string
	if f.GetCategory() != "" && len(strings.TrimSpace(f.GetCategory())) > 0 {
		registeName = f.GetCategory() + "." + f.GetName()
	} else {
		registeName = f.GetName()
	}
	functions[strings.ToLower(registeName)] = f
}

func GetFunction(name string) (Function, error) {
	name = strings.ToLower(name)
	f, ok := functions[name]
	if ok {
		return f, nil
	}
	for k, _ := range functions {
		log.Debugf("function %s", k)
	}
	return nil, fmt.Errorf("No function %s found", name)
}

func GetFunctionByTag(name string, tag string) (Function, error) {
	regName := strings.ToLower(getRegisteName(name, tag))
	f, ok := functions[regName]
	if ok {
		return f, nil
	}

	for k, _ := range functions {
		log.Debugf("function %s", k)
	}
	return nil, fmt.Errorf("No function name %s tag %s found", name, tag)
}

func ListAllFunctions() []string {
	var keys []string
	for k, _ := range functions {
		keys = append(keys, k)
	}
	return keys
}

func getRegisteName(name, tag string) string {
	var registeName string
	if tag != "" && len(strings.TrimSpace(tag)) > 0 {
		registeName = tag + "." + name
	} else {
		registeName = name
	}
	return registeName
}
