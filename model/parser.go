package model

import (
	"encoding/json"
	"errors"
	"github.com/TIBCOSoftware/flogo-cli/env"
	"github.com/TIBCOSoftware/flogo-lib/app"
	faction "github.com/TIBCOSoftware/flogo-lib/core/action"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling-lib/types"
	"io/ioutil"
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

func CreateFlogoTrigger(trigger types.Trigger, handler types.EventHandler) (*ftrigger.Config, error) {
	var flogoTrigger ftrigger.Config
	flogoTrigger.Name = trigger.Name
	flogoTrigger.Id = trigger.Name
	flogoTrigger.Ref = trigger.Type
	var ftSettings interface{}
	if err := json.Unmarshal([]byte(trigger.Settings), &ftSettings); err != nil {
		return nil, err
	}
	flogoTrigger.Settings = ftSettings.(map[string]interface{})
	flogoHandler := ftrigger.HandlerConfig{
		ActionId: handler.Name,
		Settings: ftSettings.(map[string]interface{}),
	}

	handlers := []*ftrigger.HandlerConfig{}
	handlers = append(handlers, &flogoHandler)

	flogoHandler.Settings["useReplyHandler"] = "false"
	flogoHandler.Settings["autoIdReply"] = "false"
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

		gbProject := env.NewGbProjectEnv()

		gbProject.Init(os.Getenv("GOPATH"))

		resourceDir := gbProject.GetVendorSrcDir()
		resourcePath := resourceDir + "/" + gitHubPath + "/" + resourceFile

		gbProject.InstallDependency(gitHubPath, "")

		data, err := ioutil.ReadFile(resourcePath)
		if err != nil {
			return nil, err
		}

		err = gbProject.UninstallDependency(gitHubPath)
		if err != nil {
			return nil, err
		}

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
