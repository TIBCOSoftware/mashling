package Core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	mservice "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/activity/service"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
)

var log = logger.GetLogger("activity-mashling-core")

// MashlingCore is a stub for your Activity implementation
type MashlingCore struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MashlingCore{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MashlingCore) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MashlingCore) Eval(context activity.Context) (done bool, err error) {
	// github.com/TIBCOSoftware/flogo-lib/core/mapper/object.go has fmt.Printf statement commented out to stop flow params from being written to screen.
	payload := context.GetInput("mashlingPayload")
	if payload == nil {
		log.Info("Executing mashling-core with empty payload.")
	} else {
		log.Info("Executing mashling-core with payload: ", payload)
	}
	identifier := context.GetInput("identifier").(string)
	instance := context.GetInput("instance").(string)
	rawRoutes := context.GetInput("routes").([]interface{})
	rawServices := context.GetInput("services").([]interface{})

	log.Info("Executing mashling-core with identifier: ", identifier)
	log.Info("Executing mashling-core with instance: ", instance)
	log.Debug("Executing mashling-core with routes: ", rawRoutes)
	log.Debug("Executing mashling-core with services: ", rawServices)

	// Parse routes
	var routes []types.Route
	var routesJSON json.RawMessage
	routesJSON, err = json.Marshal(rawRoutes)
	if err != nil {
		log.Error("error loading routes")
		return false, err
	}
	err = json.Unmarshal(routesJSON, &routes)
	if err != nil {
		log.Error("error parsing routes")
		return false, err
	}

	// Parse services
	var services []types.Service
	var servicesJSON json.RawMessage
	servicesJSON, err = json.Marshal(rawServices)
	if err != nil {
		log.Error("error loading services")
		return false, err
	}
	err = json.Unmarshal(servicesJSON, &services)
	if err != nil {
		log.Error("error parsing services")
		return false, err
	}
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
		return false, err
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
				return false, vmerr
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

	replyHandler := context.FlowDetails().ReplyHandler()

	if replyHandler != nil && routeToExecute != nil {
		for _, response := range routeToExecute.Responses {
			var truthiness bool
			truthiness, err = evaluateTruthiness(response.Condition, vm)
			if err != nil {
				continue
			}
			if truthiness {
				output, oErr := translateMappings(&executionContext, response.Output)
				if oErr != nil {
					return false, oErr
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
					replyHandler.Reply(code, output, nil)
				} else {
					log.Info("Code contents is not found or not an integer, default response is 200")
					replyHandler.Reply(200, output, nil)
				}
				return true, err
			}
		}
	}
	log.Info("no response conditions evaluated to true")
	return true, err
}

func executeRoute(route *types.Route, services map[string]types.Service, executionContext *map[string]interface{}, vm *mservice.VM) (err error) {
	for _, step := range route.Steps {
		var truthiness bool
		truthiness, err = evaluateTruthiness(step.Condition, vm)
		if err != nil {
			return err
		}
		if truthiness {
			err = invokeService(services[step.Service], executionContext, step.Input, step.Output, vm)
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

func invokeService(serviceDef types.Service, executionContext *map[string]interface{}, input map[string]interface{}, output map[string]interface{}, vm *mservice.VM) (err error) {
	log.Info("invoking service type: ", serviceDef.Type)
	serviceInstance, err := mservice.Initialize(serviceDef)
	if err != nil {
		return err
	}
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
	err = vm.SetInVM(serviceDef.Name, serviceInstance)
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
