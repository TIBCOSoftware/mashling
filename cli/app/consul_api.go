/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/mashling/lib/model"
	"github.com/TIBCOSoftware/mashling/lib/types"
)

const (
	registerURI   = "http://localhost:8500/v1/agent/service/register"
	deRegisterURI = "http://localhost:8500/v1/agent/service/deregister/"
)

type consulServiceDef struct {
	Name    string `json:"Name"`
	Port    string `json:"port"`
	Address string `json:"address"`
}

func generateConsulDef(gatewayJSON string) ([]consulServiceDef, error) {

	triggers, err := generateFlogoTriggers(gatewayJSON)
	if err != nil {
		return nil, err
	}

	var consulServices = make([]consulServiceDef, len(triggers))

	for i, trigger := range triggers {

		if trigger.Ref == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Ref == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
			var consulDef consulServiceDef

			settings, err := json.MarshalIndent(&trigger.Settings, "", "    ")
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(settings, &consulDef)
			if err != nil {
				return nil, err
			}
			consulServices[i].Name = trigger.Name
			consulServices[i].Port = consulDef.Port

		}
	}

	return consulServices, nil
}

//RegisterWithConsul registers suplied gateway json with consul
func RegisterWithConsul(gatewayJSON string, consulToken string, consulDefDir string) error {

	consulServices, err := generateConsulDef(gatewayJSON)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	for _, content := range consulServices {

		port, _ := strconv.Atoi(content.Port)

		contentMap := map[string]interface{}{
			"Name":    content.Name,
			"Address": "127.0.0.1",
			"Port":    port,
		}

		contentPayload, err := json.MarshalIndent(&contentMap, "", "    ")
		if err != nil {
			return err
		}

		statusCode, err := callConsulService(registerURI, []byte(contentPayload))

		if err != nil {
			return err
		}

		if statusCode != http.StatusOK {
			return fmt.Errorf("registration failed : status code %v", statusCode)
		}
	}
	return nil
}

//DeregisterFromConsul removes suplied gateway json from consul
func DeregisterFromConsul(gatewayJSON string, consulToken string, consulDefDir string) error {

	consulServices, err := generateConsulDef(gatewayJSON)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	for _, content := range consulServices {

		fullURI := deRegisterURI + content.Name

		statusCode, err := callConsulService(fullURI, []byte(""))

		if err != nil {
			return err
		}

		if statusCode != http.StatusOK {
			return fmt.Errorf("deregistration failed : status code %v", statusCode)
		}
	}
	return nil
}

func callConsulService(uri string, payload []byte) (int, error) {

	client := &http.Client{}
	r, _ := http.NewRequest("PUT", uri, bytes.NewReader([]byte(payload)))
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, err
}

//generateFlogoJson generates flogo json
func generateFlogoTriggers(gatewayJSON string) ([]*ftrigger.Config, error) {

	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		return nil, err
	}

	flogoAppTriggers := []*ftrigger.Config{}

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

	createdTriggersMap := make(map[string]*ftrigger.Config)

	//translate the gateway model to the flogo model
	for _, link := range descriptor.Gateway.EventLinks {
		triggerNames := link.Triggers

		for _, triggerName := range triggerNames {
			dispatches := link.Dispatches

			flogoTrigger, isNew, err := model.CreateFlogoTrigger(configNamedMap, triggerNamedMap[triggerName], handlerNamedMap, dispatches, createdTriggersMap)
			if err != nil {
				return nil, err
			}

			if *isNew {
				//looks like a new trigger has been added
				flogoAppTriggers = append(flogoAppTriggers, flogoTrigger)
			} else {
				//looks like an existing trigger with matching settings is found and modified with a new handler
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
					createdHandlers[handlerName] = true
				}
			}
		}

	}
	return flogoAppTriggers, nil
}
