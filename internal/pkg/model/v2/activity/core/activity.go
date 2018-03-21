package Core

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	mservice "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/activity/service"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
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

	// Execute the identified route if it exists and handle the async option.
	if routeToExecute != nil {
		if routeToExecute.Async {
			log.Info("executing route asynchronously")
			vmDefaults["async"] = true
			asyncVM, vmerr := mservice.NewVM(vmDefaults)
			if vmerr != nil {
				return false, vmerr
			}
			go executeRoute(routeToExecute, serviceMap, asyncVM)
			vm.SetPrimitiveInVM("async", true)
		} else {
			err = executeRoute(routeToExecute, serviceMap, vm)
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
				output := make(map[string]interface{})
				err = vm.SetInVM("output", output)
				if err != nil {
					return false, err
				}
				err = vm.RunTranslationMappings("output", response.Output)
				if err != nil {
					return false, err
				}
				err = vm.GetFromVM("output", &output)
				if err != nil {
					return false, err
				}
				var code float64
				var ok bool
				var codeElement interface{}
				if codeElement, ok = output["code"]; ok {
					code, ok = codeElement.(float64)
				}
				if ok && code != 0 {
					log.Info("Code identified in response output: ", int(code))
					replyHandler.Reply(int(code), output, nil)
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

func executeRoute(route *types.Route, services map[string]types.Service, vm *mservice.VM) (err error) {
	for _, step := range route.Steps {
		var truthiness bool
		truthiness, err = evaluateTruthiness(step.Condition, vm)
		if err != nil {
			return err
		}
		if truthiness {
			err = invokeService(services[step.Service], step.Input, step.Output, vm)
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

func invokeService(serviceDef types.Service, input map[string]interface{}, output map[string]interface{}, vm *mservice.VM) (err error) {
	log.Info("invoking service type: ", serviceDef.Type)
	serviceInstance, err := mservice.Initialize(serviceDef)
	if err != nil {
		return err
	}
	err = vm.SetInVM(serviceDef.Name, serviceInstance)
	if err != nil {
		return err
	}
	err = vm.RunTranslationMappings(serviceDef.Name+".request", input)
	if err != nil {
		return err
	}
	err = vm.GetFromVM(serviceDef.Name, serviceInstance)
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
	err = vm.RunTranslationMappings(serviceDef.Name+".response", output)
	if err != nil {
		return err
	}
	return nil
}
