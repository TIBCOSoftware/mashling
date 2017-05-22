package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling-lib/types"
	"github.com/TIBCOSoftware/mashling-lib/util"
	"os"
	"strings"
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

func CreateFlogoTrigger(configDefinitions map[string]types.Config, trigger types.Trigger, handler types.EventHandler) (*ftrigger.Config, error) {
	var flogoTrigger ftrigger.Config
	flogoTrigger.Name = trigger.Name
	flogoTrigger.Id = trigger.Name
	flogoTrigger.Ref = trigger.Type
	var mtSettings interface{}
	if err := json.Unmarshal([]byte(trigger.Settings), &mtSettings); err != nil {
		return nil, err
	}

	//resolve any configuration references if the "config" param is set in the settings
	mashTriggerSettings := mtSettings.(map[string]interface{})
	mashTriggerSettingsUsable := mtSettings.(map[string]interface{})
	for k, v := range mashTriggerSettings {
		mashTriggerSettingsUsable[k] = v
	}

	if configRef, ok := mashTriggerSettings[util.Gateway_Trigger_Config_Ref_Key]; ok {
		//get the configuration details
		//the expression would be e.g. ${configurations.kafkaConfig}
		configExpr := configRef.(string)
		valid, configName := util.ValidateTriggerConfigExpr(&configExpr)
		if !valid {
			return nil, fmt.Errorf("Invalid Configuration reference specified in the Trigger settings [%v]", configName)
		}
		//lets get the config object details
		configNameStr := *configName

		if configObject, ok := configDefinitions[configNameStr]; ok {
			if configObject.Type != trigger.Type {
				return nil, fmt.Errorf("Mismatch in the Configuration reference [%v] and the Trigger type [%v]", configObject.Type, trigger.Type)
			}

			var configObjSettings interface{}
			if err := json.Unmarshal([]byte(configObject.Settings), &configObjSettings); err != nil {
				return nil, err
			}
			configSettingsMap := configObjSettings.(map[string]interface{})
			//delete the "config" key from the the Usable trigger settings map
			delete(mashTriggerSettingsUsable, util.Gateway_Trigger_Config_Ref_Key)
			//copy from the config settings into the usable trigger settings map
			for k, v := range configSettingsMap {
				mashTriggerSettingsUsable[k] = v
			}

		}
	}

	//check if the trigger has valid settings required
	//1. get the trigger resource from github
	triggerMD, err := util.GetTriggerMetadata(trigger.Type)
	if err != nil {
		return nil, err
	}
	//2. check if the trigger metadata contains the settings
	triggerSettings := make(map[string]interface{})
	handlerSettings := make(map[string]interface{})

	for key, value := range mashTriggerSettingsUsable {
		if util.IsValidTriggerSetting(triggerMD, key) {
			triggerSettings[key] = value
		}

		if util.IsValidTriggerHandlerSetting(triggerMD, key) {
			handlerSettings[key] = value
		}
	}
	//3. check if the trigger handler metadata contain the settings

	flogoTrigger.Settings = triggerSettings
	flogoHandler := ftrigger.HandlerConfig{
		ActionId: handler.Name,
		Settings: handlerSettings,
	}

	handlers := []*ftrigger.HandlerConfig{}
	handlers = append(handlers, &flogoHandler)

	flogoHandler.Settings[util.Gateway_Trigger_Handler_UseReplyHandler] = util.Gateway_Trigger_Handler_UseReplyHandler_Default
	flogoHandler.Settings[util.Gateway_Trigger_Handler_AutoIdReply] = util.Gateway_Trigger_Handler_AutoIdReply_Default
	flogoTrigger.Handlers = handlers

	return &flogoTrigger, nil
}

func CreateFlogoFlowAction(handler types.EventHandler) (*faction.Config, error) {
	flogoAction := types.FlogoAction{}
	reference := &handler.Reference
	gatewayAction := faction.Config{}

	if reference != nil {
		//reference is provided, get the referenced resource inline. the provided path should be the git path e.g. github.com/<userid>/resources/app.json
		referenceString := *reference

		index := strings.LastIndex(referenceString, "/")

		if index < 0 {
			return nil, errors.New("Invalid URL reference. Pls provide the github path to mashling flow json")
		}
		gitHubPath := referenceString[0:index]

		resourceFile := referenceString[index+1 : len(referenceString)]

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
