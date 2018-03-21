package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry/triggers"
	condition "github.com/TIBCOSoftware/mashling/lib/conditions"
	"github.com/TIBCOSoftware/mashling/lib/types"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

// Translate translates mashling gateway JSON config to a Flogo app.
func Translate(descriptor *types.Microgateway) ([]byte, error) {
	flogoAppTriggers := []*ftrigger.Config{}
	flogoAppActions := []*faction.Config{}

	//1. load the configuration, if provided.
	configNamedMap := make(map[string]types.Config)
	for _, config := range descriptor.Gateway.Configurations {
		configNamedMap[config.Name] = config
	}

	triggerNamedMap := make(map[string]types.Trigger)
	for _, trigger := range descriptor.Gateway.Triggers {
		triggerNamedMap[trigger.Name] = trigger
	}

	handlerNamedMap := make(map[string]types.EventHandler)
	for _, evtHandler := range descriptor.Gateway.EventHandlers {
		handlerNamedMap[evtHandler.Name] = evtHandler
	}

	createdHandlers := make(map[string]bool)

	//new map to maintain existing trigger and its settings, to be used in comparing one trigger definition with another
	createdTriggersMap := make(map[string]*ftrigger.Config)

	//translate the gateway model to the flogo model
	for _, link := range descriptor.Gateway.EventLinks {
		triggerNames := link.Triggers

		for _, triggerName := range triggerNames {
			dispatches := link.Dispatches

			//create trigger sections for flogo
			/**
			TODO handle condition parsing and setting the condition in the trigger.
			//get the condition if available
			condition := path.If
			//create the trigger using the condition
			.......
			*/
			flogoTrigger, isNew, terr := createFlogoTrigger(configNamedMap, triggerNamedMap[triggerName], handlerNamedMap, dispatches, createdTriggersMap)
			if terr != nil {
				return nil, terr
			}

			//	check if the trigger is a new trigger or a "logically" same trigger.
			if *isNew {
				//looks like a new trigger has been added
				flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
			} else {
				//looks like an existing trigger with matching settings is found and is now modified with a new handler
				for index, v := range flogoAppTriggers {
					if v.Name == flogoTrigger.Name {
						// Found the old trigger entry in the list!
						//remove it..
						flogoAppTriggers = append(flogoAppTriggers[:index], flogoAppTriggers[index+1:]...)
						//...and add the modified trigger to the list
						flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
					}
				}
			}

			//create unique handler actions
			for _, dispatch := range dispatches {
				handlerName := dispatch.Handler

				if !createdHandlers[handlerName] {
					//not already created, so create it
					flogoAction, aerr := createFlogoFlowAction(handlerNamedMap[handlerName])
					if aerr != nil {
						return nil, aerr
					}

					flogoAppActions = append(flogoAppActions, flogoAction)
					createdHandlers[handlerName] = true
				}
			}
		}

	}

	flogoApp := app.Config{
		Name:        descriptor.Gateway.Name,
		Type:        util.Flogo_App_Type,
		Version:     descriptor.Gateway.Version,
		Description: descriptor.Gateway.Description,
		Triggers:    flogoAppTriggers,
		Actions:     flogoAppActions,
	}

	//create flogo PP JSON
	flogoJSON, err := json.MarshalIndent(flogoApp, "", "\t")
	if err != nil {
		return flogoJSON, nil
	}

	return flogoJSON, nil
}

func createFlogoTrigger(configDefinitions map[string]types.Config, trigger types.Trigger, namedHandlerMap map[string]types.EventHandler,
	dispatches []types.Dispatch, createdTriggersMap map[string]*ftrigger.Config) (*ftrigger.Config, *bool, error) {

	var flogoTrigger ftrigger.Config
	flogoTrigger.Name = trigger.Name
	flogoTrigger.Id = trigger.Name
	flogoTrigger.Ref = trigger.Type
	var mtSettings interface{}
	if err := json.Unmarshal([]byte(trigger.Settings), &mtSettings); err != nil {
		return nil, nil, err
	}

	//resolve any configuration references if the "config" param is set in the settings
	mashTriggerSettings := mtSettings.(map[string]interface{})
	mashTriggerSettingsUsable := mtSettings.(map[string]interface{})
	for k, v := range mashTriggerSettings {
		mashTriggerSettingsUsable[k] = v
	}

	if configDefinitions != nil && len(configDefinitions) > 0 {
		//inherit the configuration settings if the trigger uses configuration reference
		err := resolveConfigurationReference(configDefinitions, trigger, mashTriggerSettingsUsable)
		if err != nil {
			return nil, nil, err
		}
	}

	//check if the trigger has valid settings required
	triggerMD, err := GetLocalTriggerMetadata(trigger.Type)
	if err != nil {
		return nil, nil, err
	}
	//2. check if the trigger metadata contains the settings
	triggerSettings := make(map[string]interface{})
	handlerSettings := make(map[string]map[string]interface{})

	for key, value := range mashTriggerSettingsUsable {
		if util.IsValidTriggerSetting(triggerMD, key) {
			triggerSettings[key] = value
		}
	}

	isNew := true

	//check if the trigger specifies a boolean setting key named 'optimize'
	if util.CheckTriggerOptimization(mashTriggerSettingsUsable) {
		log.Printf("[mashling] Trigger specifies %v property setting true\n", util.Gateway_Trigger_Optimize_Property)

		//2.1 check if a trigger having the same settings is already created
		//2.2 organize the trigger names as a list so that they can be sorted alphabetically. Golang maps are unordered and the iteration order is not guaranteed across multiple iterations.
		triggerNames := make([]string, len(createdTriggersMap))
		i := 0
		for k := range createdTriggersMap {
			triggerNames[i] = k
			i++
		}
		sort.Strings(triggerNames)

		//iterate over the list of trigger names, now sorted alphabetically.
		for _, name := range triggerNames {
			createdTrigger := createdTriggersMap[name]
			if reflect.DeepEqual(createdTrigger.Settings, triggerSettings) {
				//looks like we found an existing trigger that has the same settings. No need to create a new trigger object. just create a new handler on the existing trigger
				log.Printf("[mashling] Found a trigger having same settings %v %v\n", name, triggerSettings)
				flogoTrigger = *createdTrigger
				isNew = false
				break
			} else {
				log.Printf("[mashling] Current trigger %v did not match settings of trigger %v %v\n", flogoTrigger.Name, name, triggerSettings)
			}
		}
	}

	//3. check if the trigger handler metadata contain the settings
	handlers := []*ftrigger.HandlerConfig{}
	var handler types.EventHandler
	for _, dispatch := range dispatches {
		handler = namedHandlerMap[dispatch.Handler]

		handlerSettings[handler.Name] = make(map[string]interface{})
		for key, value := range mashTriggerSettingsUsable {
			if util.IsValidTriggerHandlerSetting(triggerMD, key) {
				handlerSettings[handler.Name][key] = value
			}
		}

		//check if any condition is specified & if the Condition setting is part of the trigger metadata
		if dispatch.If != "" {
			//check if the trigger metadata supports Condition handler setting
			if util.IsValidTriggerHandlerSetting(triggerMD, util.Flogo_Trigger_Handler_Setting_Condition) {
				//check if the condition is valid.
				condition.ValidateOperatorInExpression(dispatch.If)
				//set the condition on the trigger as is. the trigger should parse and interpret it.
				handlerSettings[handler.Name][util.Flogo_Trigger_Handler_Setting_Condition] = dispatch.If
			} else {
				log.Printf("[mashling] The trigger [%v] does not support [%v] handler setting. skippng the condition logic.\n", trigger.Type, util.Flogo_Trigger_Handler_Setting_Condition)
			}
		}
		flogoTrigger.Settings = triggerSettings
		flogoHandler := ftrigger.HandlerConfig{
			ActionId: handler.Name,
			Settings: handlerSettings[handler.Name],
		}

		//Add autoIdReply & useReplyHandler settings only for valid trigger (i.e http trigger. For kafka these settings are invalid).
		if util.IsValidTriggerHandlerSetting(triggerMD, util.Gateway_Trigger_Handler_UseReplyHandler) {
			flogoHandler.Settings[util.Gateway_Trigger_Handler_UseReplyHandler] = util.Gateway_Trigger_Handler_UseReplyHandler_Default
		}
		if util.IsValidTriggerHandlerSetting(triggerMD, util.Gateway_Trigger_Handler_AutoIdReply) {
			flogoHandler.Settings[util.Gateway_Trigger_Handler_AutoIdReply] = util.Gateway_Trigger_Handler_AutoIdReply_Default
		}

		handlers = append(handlers, &flogoHandler)
	}

	flogoTrigger.Handlers = append(flogoTrigger.Handlers, handlers...)

	if isNew {
		log.Printf("[mashling] Adding a new trigger with settings %v %v\n", flogoTrigger.Name, triggerSettings)
		createdTriggersMap[flogoTrigger.Name] = &flogoTrigger
	}

	return &flogoTrigger, &isNew, nil
}

func createFlogoFlowAction(handler types.EventHandler) (*faction.Config, error) {
	flogoAction := types.FlogoAction{}
	reference := handler.Reference
	gatewayAction := faction.Config{}

	// handle param name-values provided as part of the handler
	// if handler.Params != nil {
	// 	var handlerParams interface{}
	// 	if err := json.Unmarshal([]byte(handler.Params), &handlerParams); err != nil {
	// 		return nil, err
	// 	}
	// 	if handlerParams != nil {
	// 		handlerParamsMap := handlerParams.(map[string]interface{})
	// 		//substitute for any ENV variable values referenced in the params. the expressions will be in the format ${ENV.HOST_NAME} where HOST_NAME is the env property
	// 		err := resolveEnvironmentProperties(handlerParamsMap)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }

	//TODO use the params as inputs to the flow action. flogo tunables can be used to validate the params and, if valid, the values can be injected into the flow

	if reference != "" {
		return nil, errors.New("only inlined Flogo action definitions allowed")
	} else if handler.Definition != nil {
		//definition is provided inline
		err := json.Unmarshal([]byte(handler.Definition), &flogoAction)
		if err != nil {
			return nil, err
		}
		gatewayAction = faction.Config{
			Id:   handler.Name,
			Data: flogoAction.Data,
			Ref:  flogoAction.Ref,
		}
	}

	return &gatewayAction, nil
}

func resolveConfigurationReference(configDefinitions map[string]types.Config, trigger types.Trigger, settings map[string]interface{}) error {
	if configRef, ok := settings[util.Gateway_Trigger_Config_Ref_Key]; ok {
		//get the configuration details
		//the expression would be e.g. ${configurations.kafkaConfig}
		configExpr := configRef.(string)
		valid, configName := util.ValidateTriggerConfigExpr(&configExpr)
		if !valid {
			return fmt.Errorf("Invalid Configuration reference specified in the Trigger settings [%v]", configName)
		}
		//lets get the config object details
		configNameStr := *configName

		if configObject, ok := configDefinitions[configNameStr]; ok {
			if configObject.Type != trigger.Type {
				return fmt.Errorf("Mismatch in the Configuration reference [%v] and the Trigger type [%v]", configObject.Type, trigger.Type)
			}

			var configObjSettings interface{}
			if err := json.Unmarshal([]byte(configObject.Settings), &configObjSettings); err != nil {
				return err
			}
			configSettingsMap := configObjSettings.(map[string]interface{})
			//delete the "config" key from the the Usable trigger settings map
			delete(settings, util.Gateway_Trigger_Config_Ref_Key)
			//copy from the config settings into the usable trigger settings map, if the key does NOT exist in the trigger already.
			//this is to ensure that the individual trigger can override a property defined in a "common" configuration
			for k, v := range configSettingsMap {
				if _, ok := settings[k]; !ok {
					settings[k] = v
				}
			}
		}
	}
	return nil
}

// GetLocalTriggerMetadata extracts trigger metadata from the local json definition.
func GetLocalTriggerMetadata(gitHubPath string) (*ftrigger.Metadata, error) {
	// Look for local first
	triggerPath := strings.Replace(gitHubPath, "github.com/TIBCOSoftware/mashling/", "", 1) + "/trigger.json"
	fmt.Println(triggerPath)
	data, err := triggers.Asset(triggerPath)
	if err != nil {
		// Look in vendor now
		triggerPath := "vendor/" + gitHubPath + "/trigger.json"
		data, err = triggers.Asset(triggerPath)
		if err != nil {
			return nil, err
		}
	}
	triggerMetadata := &ftrigger.Metadata{}
	json.Unmarshal(data, triggerMetadata)
	return triggerMetadata, nil
}
