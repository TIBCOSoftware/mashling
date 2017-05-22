package app

import (
	"encoding/json"
	"strings"

	"fmt"
	api "github.com/TIBCOSoftware/flogo-cli/app"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling-lib/model"
	"github.com/TIBCOSoftware/mashling-lib/types"
	"github.com/TIBCOSoftware/mashling-lib/util"
	"path"
)

// CreateMashling creates a gateway application from the specified json gateway descriptor
func CreateMashling(env env.Project, gatewayJson string, appDir string, appName string, vendorDir string) error {

	descriptor, err := model.ParseGatewayDescriptor(gatewayJson)
	if err != nil {
		return err
	}

	if appName != "" {
		altJson := strings.Replace(gatewayJson, `"`+descriptor.Gateway.Name+`"`, `"`+appName+`"`, 1)
		altDescriptor, err := model.ParseGatewayDescriptor(altJson)

		//see if we can get away with simple replace so we don't reorder the existing json
		if err == nil && altDescriptor.Gateway.Name == appName {
			gatewayJson = altJson
		} else {
			//simple replace didn't work so we have to unmarshal & re-marshal the supplied json
			var appObj map[string]interface{}
			err := json.Unmarshal([]byte(gatewayJson), &appObj)
			if err != nil {
				return err
			}

			appObj["name"] = appName

			updApp, err := json.MarshalIndent(appObj, "", "  ")
			if err != nil {
				return err
			}
			gatewayJson = string(updApp)
		}

		descriptor.Gateway.Name = appName
	} else {
		appName = descriptor.Gateway.Name
	}

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

	//translate the gateway model to the flogo model
	for _, link := range descriptor.Gateway.EventLinks {
		triggerName := link.Trigger

		successPaths := link.SuccessPaths
		for _, path := range successPaths {
			handlerName := path.Handler

			flogoTrigger, err := model.CreateFlogoTrigger(configNamedMap, triggerNamedMap[triggerName], handlerNamedMap[handlerName])
			if err != nil {
				return err
			}

			flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)

			flogoAction, err := model.CreateFlogoFlowAction(handlerNamedMap[handlerName])
			if err != nil {
				return err
			}

			flogoAppActions = append(flogoAppActions, flogoAction)
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
	bytes, err := json.MarshalIndent(flogoApp, "", "\t")
	if err != nil {
		return nil
	}

	flogoJson := string(bytes)

	err = api.CreateApp(SetupNewProjectEnv(), flogoJson, appDir, appName, vendorDir)
	if err != nil {
		return err
	}

	fmt.Println("Generated mashling Artifacts.")
	fmt.Println("Building mashling Artifacts.")

	options := &api.BuildOptions{SkipPrepare: false, PrepareOptions: &api.PrepareOptions{OptimizeImports: false, EmbedConfig: false}}
	api.BuildApp(SetupExistingProjectEnv(appDir), options)

	err = fgutil.CreateFileFromString(path.Join(appDir, util.Gateway_Definition_File_Name), gatewayJson)
	if err != nil {
		return err
	}

	fmt.Println("Mashling gateway successfully built!")

	return nil
}
