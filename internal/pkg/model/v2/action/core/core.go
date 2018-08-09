package core

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	mservice "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/action/service"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
)

var log = logger.GetLogger("action-mashling")

func ExecuteMashling(payload interface{}, routes []types.Route, services []types.Service) (code int, output interface{}, err error) {
	// Create services map
	serviceMap := make(map[string]types.Service)
	for _, service := range services {
		serviceMap[service.Name] = service
	}

	// Route to be executed once it is identified by the conditional evaluation.
	var routeToExecute *types.Route

	// Setup conditional VM with defaults.
	vmDefaults := make(map[string]interface{})
	if payload != nil {
		vmDefaults["payload"] = payload
	}
	vmDefaults["async"] = false
	// Add ENV flags to the vmDefaults
	envFlags := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envFlags[pair[0]] = pair[1]
	}
	vmDefaults["env"] = envFlags
	vm, err := mservice.NewVM(vmDefaults)
	if err != nil {
		return -1, nil, err
	}

	// Evaluate route conditions to select which one to execute.
	for _, route := range routes {
		var truthiness bool
		truthiness, err = evaluateTruthiness(route.Condition, vm)
		if err != nil {
			continue
		}
		if truthiness {
			log.Info("route identified via conditional evaluation to true: ", route.Condition)
			routeToExecute = &route
			break
		}
	}
	// Contains all elements of request: right now just payload, environment flags and service instances.
	executionContext := make(map[string]interface{})
	executionContext["payload"] = &payload
	executionContext["env"] = envFlags

	// Execute the identified route if it exists and handle the async option.
	if routeToExecute != nil {
		if routeToExecute.Async {
			log.Info("executing route asynchronously")
			vmDefaults["async"] = true
			asyncVM, vmerr := mservice.NewVM(vmDefaults)
			if vmerr != nil {
				return -1, nil, vmerr
			}
			go executeRoute(routeToExecute, serviceMap, &executionContext, asyncVM)
			vm.SetPrimitiveInVM("async", true)
		} else {
			err = executeRoute(routeToExecute, serviceMap, &executionContext, vm)
		}
		if err != nil {
			log.Error("error executing route: ", err)
		}
	} else {
		log.Info("no route to execute, continuing to reply handler")
	}

	if routeToExecute != nil {
		for _, response := range routeToExecute.Responses {
			var truthiness bool
			truthiness, err = evaluateTruthiness(response.Condition, vm)
			if err != nil {
				continue
			}
			if truthiness {
				output, oErr := translateMappings(&executionContext, map[string]interface{}{"code": response.Output.Code})
				if oErr != nil {
					return -1, nil, oErr
				}
				var code int
				codeElement, ok := output["code"]
				if ok {
					switch cv := codeElement.(type) {
					case float64:
						code = int(cv)
					case int:
						code = cv
					case string:
						code, err = strconv.Atoi(cv)
						if err != nil {
							log.Info("unable to format extracted code string from response output", cv)
						}
					}
				}
				if ok && code != 0 {
					log.Info("Code identified in response output: ", code)
				} else {
					log.Info("Code contents is not found or not an integer, default response code is 200")
					code = 200
				}
				// Translate data mappings
				var data interface{}
				nestedData, ok := response.Output.Data.(map[string]interface{})
				if ok {
					data, oErr = translateMappings(&executionContext, nestedData)
					if oErr != nil {
						return -1, nil, oErr
					}
				} else {
					interimData, dErr := translateMappings(&executionContext, map[string]interface{}{"data": response.Output.Data})
					if dErr != nil {
						return -1, nil, dErr
					}
					data, ok = interimData["data"]
					if !ok {
						return -1, nil, errors.New("cannot extract data from response output")
					}
				}
				return code, data, err
			}
		}
	}
	return 0, nil, err
}

func executeRoute(route *types.Route, services map[string]types.Service, executionContext *map[string]interface{}, vm *mservice.VM) (err error) {
	for _, step := range route.Steps {
		var truthiness bool
		truthiness, err = evaluateTruthiness(step.Condition, vm)
		if err != nil {
			return err
		}
		if truthiness {
			err = invokeService(services[step.Service], executionContext, step.Input, vm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func evaluateTruthiness(condition string, vm *mservice.VM) (truthy bool, err error) {
	if condition == "" {
		log.Info("condition was empty and thus evaluates to true")
		return true, nil
	}
	truthy, err = vm.EvaluateToBool(condition)
	if err != nil {
		log.Infof("condition evaluation causes error so is false: %s", condition)
		return false, err
	}
	log.Infof("condition evaluated to %t: %s", truthy, condition)
	return truthy, err
}

func invokeService(serviceDef types.Service, executionContext *map[string]interface{}, input map[string]interface{}, vm *mservice.VM) (err error) {
	log.Info("invoking service type: ", serviceDef.Type)
	serviceInstance, err := mservice.Initialize(serviceDef)
	if err != nil {
		return err
	}
	defer func() {
		vmErr := vm.SetInVM(serviceDef.Name, serviceInstance)
		if vmErr != nil {
			err = vmErr
		}
	}()
	(*executionContext)[serviceDef.Name] = &serviceInstance
	values, mErr := translateMappings(executionContext, input)
	if mErr != nil {
		return mErr
	}
	err = serviceInstance.UpdateRequest(values)
	if err != nil {
		return err
	}
	err = serviceInstance.Execute()
	if err != nil {
		return err
	}
	return nil
}

func translateMappings(executionContext *map[string]interface{}, mappings map[string]interface{}) (values map[string]interface{}, err error) {
	values = make(map[string]interface{})
	if len(mappings) == 0 {
		return values, err
	}
	for fullKey, v := range mappings {
		var convertedValue interface{}
		switch value := v.(type) {
		case string:
			if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
				// this is a variable so we need to evaluate it.
				value = strings.Replace(value, "${", "", 1)
				convertedValue, err = getValueFromDotNotation(*executionContext, util.TrimSuffix(value, "}"))
				if err != nil {
					return values, err
				}
			} else {
				convertedValue = value
			}
		default:
			convertedValue = value
		}
		values[fullKey] = convertedValue
	}
	return expandMap(values), err
}

func getValueFromDotNotation(rootObject interface{}, fullPropertyName string) (interface{}, error) {
	dotNames := strings.Split(fullPropertyName, ".")
	var err error
	for _, subPropertyName := range dotNames {
		rootObject, err = getProperty(rootObject, subPropertyName)
		if err != nil {
			return nil, err
		}
		if rootObject == nil {
			return nil, nil
		}
	}
	return rootObject, nil
}

func getProperty(obj interface{}, property string) (interface{}, error) {
	objKind := reflect.TypeOf(obj).Kind()
	// Check for pointer
	if objKind == reflect.Ptr {
		obj = reflect.ValueOf(obj).Elem().Interface()
		objKind = reflect.TypeOf(obj).Kind()
	}
	// Check if plain map
	if objKind == reflect.Map {
		val := reflect.ValueOf(obj)
		valueOf := val.MapIndex(reflect.ValueOf(property))
		if valueOf == reflect.Zero(reflect.ValueOf(property).Type()) {
			return nil, nil
		}
		index := val.MapIndex(reflect.ValueOf(property))
		if !index.IsValid() {
			return nil, nil
		}
		return index.Interface(), nil
	}
	if !(objKind == reflect.Struct || objKind == reflect.Ptr) {
		return nil, errors.New("can only get property fields from struct interfaces")
	}
	property = strings.Title(property)
	var objValue reflect.Value
	if objKind == reflect.Ptr {
		objValue = reflect.ValueOf(obj).Elem()
	} else {
		objValue = reflect.ValueOf(obj)
	}
	propertyField := objValue.FieldByName(property)
	if !propertyField.IsValid() {
		return nil, fmt.Errorf("%s type has no property named %s", objKind, property)
	}
	return propertyField.Interface(), nil
}

// Turn dot notation map into nested map structure.
func expandMap(m map[string]interface{}) map[string]interface{} {
	var tree = make(map[string]interface{})
	for key, value := range m {
		keys := strings.Split(key, ".")
		subTree := tree
		for _, treeKey := range keys[:len(keys)-1] {
			subTreeNew, ok := subTree[treeKey]
			if !ok {
				subTreeNew = make(map[string]interface{})
				subTree[treeKey] = subTreeNew
			}
			subTree = subTreeNew.(map[string]interface{})
		}
		subTree[keys[len(keys)-1]] = value
	}
	return tree
}
