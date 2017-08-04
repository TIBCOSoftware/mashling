package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"reflect"
	"sort"

	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	condition "github.com/TIBCOSoftware/mashling-lib/conditions"
	"github.com/TIBCOSoftware/mashling-lib/types"
	"github.com/TIBCOSoftware/mashling-lib/util"
)

// ParseGatewayDescriptor parse the application descriptor
func ParseGatewayDescriptor(appJson string) (*types.Microgateway, error) {
	descriptor := &types.Microgateway{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

func CreateFlogoTrigger(configDefinitions map[string]types.Config, trigger types.Trigger, namedHandlerMap map[string]types.EventHandler,
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

	//substitute for any ENV variable values referenced in the settings. the expressions will be in the format ${ENV.HOST_NAME} where HOST_NAME is the env property
	err := resolveEnvironmentProperties(mashTriggerSettingsUsable)
	if err != nil {
		return nil, nil, err
	}

	//check if the trigger has valid settings required
	//1. get the trigger resource from github
	triggerMD, err := util.GetTriggerMetadata(trigger.Type)
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
		fmt.Sprintf("Trigger specifies %v property setting true", util.Gateway_Trigger_Optimize_Property)

		//2.1 check if a trigger having the same settings is already created
		//2.2 organize the trigger names as a list so that they can be sorted alphabetically. Golang maps are unordered and the iteration order is not guaranteed across multiple iterations.
		triggerNames := make([]string, len(createdTriggersMap))
		i := 0
		for k, _ := range createdTriggersMap {
			triggerNames[i] = k
			i++
		}
		sort.Strings(triggerNames)

		//iterate over the list of trigger names, now sorted alphabetically.
		for _, name := range triggerNames {
			createdTrigger := createdTriggersMap[name]
			if reflect.DeepEqual(createdTrigger.Settings, triggerSettings) {
				//looks like we found an existing trigger that has the same settings. No need to create a new trigger object. just create a new handler on the existing trigger
				fmt.Sprintf("Found a trigger having same settings %v %v ", name, triggerSettings)
				flogoTrigger = *createdTrigger
				isNew = false
				break
			} else {
				fmt.Sprintf("Current trigger %v did not match settings of trigger %v %v", flogoTrigger.Name, name, triggerSettings)
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
				fmt.Sprintf("The trigger [%v] does not support [%v] handler setting. skippng the condition logic.", trigger.Type, util.Flogo_Trigger_Handler_Setting_Condition)
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
		fmt.Sprintf("Adding a new trigger with settings %v %v ", flogoTrigger.Name, triggerSettings)
		createdTriggersMap[flogoTrigger.Name] = &flogoTrigger
	}

	return &flogoTrigger, &isNew, nil
}

func CreateFlogoFlowAction(handler types.EventHandler) (*faction.Config, error) {
	flogoAction := types.FlogoAction{}
	reference := handler.Reference
	gatewayAction := faction.Config{}

	// handle param name-values provided as part of the handler
	if handler.Params != nil {
		var handlerParams interface{}
		if err := json.Unmarshal([]byte(handler.Params), &handlerParams); err != nil {
			return nil, err
		}
		if handlerParams != nil {
			handlerParamsMap := handlerParams.(map[string]interface{})
			//substitute for any ENV variable values referenced in the params. the expressions will be in the format ${ENV.HOST_NAME} where HOST_NAME is the env property
			err := resolveEnvironmentProperties(handlerParamsMap)
			if err != nil {
				return nil, err
			}
		}
	}

	//TODO use the params as inputs to the flow action. flogo tunables can be used to validate the params and, if valid, the values can be injected into the flow

	if reference != "" {
		//reference is provided, get the referenced resource inline. the provided path should be the git path e.g. github.com/<userid>/resources/app.json
		index := strings.LastIndex(reference, "/")

		if index < 0 {
			return nil, errors.New("Invalid URL reference. Pls provide the github path to mashling flow json")
		}
		gitHubPath := reference[0:index]

		resourceFile := reference[index+1 : len(reference)]

		data, err := util.GetGithubResource(gitHubPath, resourceFile)

		var flogoFlowDef *app.Config
		err = json.Unmarshal(data, &flogoFlowDef)
		if err != nil {
			return nil, err
		}

		actions := flogoFlowDef.Actions
		if len(actions) != 1 {
			return nil, errors.New("Please make sure that the pattern flow has only one action")
		}

		action := actions[0]
		action.Id = handler.Name
		gatewayAction = faction.Config{
			Id:   handler.Name,
			Data: action.Data,
			Ref:  action.Ref,
		}

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

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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

func resolveEnvironmentProperties(settings map[string]interface{}) error {
	for k, v := range settings {
		value := v.(string)
		valid, propertyName := util.ValidateEnvPropertySettingExpr(&value)
		if !valid {
			continue
		}
		//lets get the env property value
		propertyNameStr := *propertyName
		propertyValue, found := os.LookupEnv(propertyNameStr)
		if !found {
			return errors.New(fmt.Sprintf("ENV property [%v] referenced by the gateway is not set.", propertyNameStr))
		}
		settings[k] = propertyValue
	}
	return nil
}
