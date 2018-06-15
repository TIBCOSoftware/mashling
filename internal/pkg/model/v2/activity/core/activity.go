package core

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/mashling/internal/pkg/logger"
	mservice "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/activity/service"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/pkg/strings"
)

var log = logger.GetLogger("mashling-activity-core")
var initializedServices = map[string]map[string]mservice.Service{}
var initializedRoutes = map[string][]types.Route{}
var envFlags = map[string]string{}

// MashlingCore is a stub for your Activity implementation
type MashlingCore struct {
	metadata *activity.Metadata
}

func WarmUp() error {
	// Capture env variable once
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envFlags[pair[0]] = pair[1]
	}
	return nil
}

// WarmUpServices sets up re-used service data structures once to reduce per invocation overhead. This should not be called concurrently.
func WarmUpServices(mashlingInstance string, services []types.Service) error {
	// Create services map
	serviceMap := make(map[string]mservice.Service)
	for _, service := range services {
		serviceInstance, err := mservice.Initialize(service)
		if err != nil {
			return err
		}
		serviceMap[service.Name] = serviceInstance
	}
	initializedServices[mashlingInstance] = serviceMap
	return nil
}

// WarmUpRoutes sets up re-used route data structures once to reduce per invocation overhead. This should not be called concurrently.
func WarmUpRoutes(mashlingIdentifier string, routes []types.Route) error {
	initializedRoutes[mashlingIdentifier] = routes
	return nil
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
		log.Debug("Executing mashling-core with empty payload.")
	} else {
		log.Debug("Executing mashling-core with payload: ", payload)
	}
	identifier := context.GetInput("identifier").(string)
	instance := context.GetInput("instance").(string)

	log.Debug("Executing mashling-core with identifier: ", identifier)
	log.Debug("Executing mashling-core with instance: ", instance)

	// Route to be executed once it is identified by the conditional evaluation.
	var routeToExecute *types.Route

	// Setup conditional VM with defaults.
	vmDefaults := make(map[string]interface{})
	if payload != nil {
		vmDefaults["payload"] = payload
	}
	vmDefaults["async"] = false
	vmDefaults["env"] = envFlags
	vm, err := mservice.NewVM(vmDefaults)
	if err != nil {
		return false, err
	}

	// Evaluate route conditions to select which one to execute.
	for _, route := range initializedRoutes[identifier] {
		var truthiness bool
		truthiness, err = evaluateTruthiness(route.Condition, vm)
		if err != nil {
			continue
		}
		if truthiness {
			log.Debug("route identified via conditional evaluation to true: ", route.Condition)
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
			log.Debug("executing route asynchronously")
			vmDefaults["async"] = true
			asyncVM, vmerr := mservice.NewVM(vmDefaults)
			if vmerr != nil {
				return false, vmerr
			}
			go executeRoute(instance, routeToExecute, &executionContext, asyncVM)
			vm.SetPrimitiveInVM("async", true)
		} else {
			err = executeRoute(instance, routeToExecute, &executionContext, vm)
		}
		if err != nil {
			log.Error("error executing route: ", err)
		}
	} else {
		log.Debug("no route to execute, continuing to reply handler")
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
				output, oErr := translateMappings(&executionContext, map[string]interface{}{"code": response.Output.Code})
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
							log.Debug("unable to format extracted code string from response output", cv)
						}
					}
				}
				if ok && code != 0 {
					log.Debug("Code identified in response output: ", code)
				} else {
					log.Debug("Code contents is not found or not an integer, default response code is 200")
					code = 200
				}
				// Translate data mappings
				var data interface{}
				nestedData, ok := response.Output.Data.(map[string]interface{})
				if ok {
					data, oErr = translateMappings(&executionContext, nestedData)
					if oErr != nil {
						return false, oErr
					}
				} else {
					interimData, dErr := translateMappings(&executionContext, map[string]interface{}{"data": response.Output.Data})
					if dErr != nil {
						return false, dErr
					}
					data, ok = interimData["data"]
					if !ok {
						return false, errors.New("cannot extract data from response output")
					}
				}
				replyHandler.Reply(code, data, nil)
				return true, err
			}
		}
	}
	log.Debug("no response conditions evaluated to true")
	return true, err
}

func executeRoute(instance string, route *types.Route, executionContext *map[string]interface{}, vm *mservice.VM) (err error) {
	for _, step := range route.Steps {
		var truthiness bool
		truthiness, err = evaluateTruthiness(step.Condition, vm)
		if err != nil {
			return err
		}
		if truthiness {
			err = invokeService(instance, step.Service, executionContext, step.Input, vm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func evaluateTruthiness(condition string, vm *mservice.VM) (truthy bool, err error) {
	if condition == "" {
		log.Debug("condition was empty and thus evaluates to true")
		return true, nil
	}
	truthy, err = vm.EvaluateToBool(condition)
	if err != nil {
		log.Debugf("condition evaluation causes error so is false: %s", condition)
		return false, err
	}
	log.Debugf("condition evaluated to %t: %s", truthy, condition)
	return truthy, err
}

func invokeService(instance string, serviceName string, executionContext *map[string]interface{}, input map[string]interface{}, vm *mservice.VM) (err error) {
	serviceInstance := initializedServices[instance][serviceName]
	log.Debug("invoking service: ", serviceName)
	values, mErr := translateMappings(executionContext, input)
	if mErr != nil {
		return mErr
	}
	response, err := serviceInstance.Execute(values)
	if err != nil {
		return err
	}
	mappedResp := map[string]mservice.Response{"response": response}
	err = vm.SetInVM(serviceName, mappedResp)
	if err != nil {
		return err
	}
	(*executionContext)[serviceName] = mappedResp
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
