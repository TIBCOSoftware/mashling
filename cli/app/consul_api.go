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
	"os/exec"
	"strconv"

	"github.com/TIBCOSoftware/mashling/lib/model"
	"github.com/TIBCOSoftware/mashling/lib/types"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

const (
	registerURI   = "/v1/agent/service/register"
	deRegisterURI = "/v1/agent/service/deregister/"
)

type consulServiceDef struct {
	Name    string `json:"Name"`
	Port    string `json:"port"`
	Address string `json:"address"`
}

type triggerDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type"`
	Settings    map[string]interface{} `json:"settings"`
}

/***
generateConsulDef generates consul service definition from supplied gateway.json
***/
func generateConsulDef(gatewayJSON string) ([]consulServiceDef, error) {

	triggers, err := generateTriggers(gatewayJSON)
	if err != nil {
		return nil, err
	}

	var consulServices = make([]consulServiceDef, len(triggers))

	for i, trigger := range triggers {

		if trigger.Type == "github.com/TIBCOSoftware/flogo-contrib/trigger/rest" || trigger.Type == "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/gorillamuxtrigger" {
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

//RegisterWithConsul registers suplied gateway.json services with consul
func RegisterWithConsul(gatewayJSON string, consulToken string, consulDefDir string, consulAddress string) error {

	consulServices, err := generateConsulDef(gatewayJSON)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	var localIP = getLocalIP()

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	for _, content := range consulServices {

		port, _ := strconv.Atoi(content.Port)

		checkMap := map[string]interface{}{
			"tcp":      localIP + ":" + content.Port,
			"interval": "30s",
			"timeout":  "1s",
		}

		contentMap := map[string]interface{}{
			"Name":    content.Name,
			"Address": localIP,
			"Port":    port,
			"check":   checkMap,
		}

		contentPayload, err := json.MarshalIndent(&contentMap, "", "    ")
		if err != nil {
			return err
		}

		fullURI := "http://" + consulAddress + registerURI

		if len(consulDefDir) != 0 {

			err := os.Chdir(consulDefDir)
			if err != nil {
				return err
			}

			file, err := os.Create(content.Name + ".json")
			defer file.Close()

			if err != nil {
				return err
			}

			serviceMap := map[string]interface{}{
				"service": contentMap,
			}

			serviceContentPayload, err := json.MarshalIndent(&serviceMap, "", "    ")
			if err != nil {
				return err
			}
			_, dataErr := file.Write(serviceContentPayload)
			if dataErr != nil {
				return dataErr
			}

		} else {

			statusCode, err := callConsulService(fullURI, []byte(contentPayload), consulToken)

			if err != nil {
				return err
			}

			if statusCode != http.StatusOK {
				return fmt.Errorf("registration failed : status code %v", statusCode)
			}
		}
	}

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	fmt.Println("===================================")
	fmt.Println("Successfully registered with consul")
	fmt.Println("===================================")
	return nil
}

//DeregisterFromConsul removes suplied gateway.json services from consul
func DeregisterFromConsul(gatewayJSON string, consulToken string, consulDefDir string, consulAddress string) error {

	consulServices, err := generateConsulDef(gatewayJSON)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate consul payload \n\n")
		return err
	}

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	for _, content := range consulServices {

		fullURI := "http://" + consulAddress + deRegisterURI + content.Name

		if len(consulDefDir) != 0 {

			err := os.Chdir(consulDefDir)
			if err != nil {
				return err
			}

			err = os.Remove(content.Name + ".json")
			if err != nil {
				return err
			}

		} else {
			statusCode, err := callConsulService(fullURI, []byte(""), consulToken)

			if err != nil {
				return err
			}

			if statusCode != http.StatusOK {
				return fmt.Errorf("deregistration failed : status code %v", statusCode)
			}
		}
	}

	if len(consulDefDir) != 0 {
		err = reloadConsul(consulToken)
		if err != nil {
			return err
		}
	}

	fmt.Println("======================================")
	fmt.Println("Successfully de-registered with consul")
	fmt.Println("======================================")
	return nil
}

/**
callConsulService Performs PUT API call on consul agent
**/
func callConsulService(uri string, payload []byte, consulToken string) (int, error) {

	client := &http.Client{}
	r, _ := http.NewRequest("PUT", uri, bytes.NewReader([]byte(payload)))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Consul-Token", consulToken)

	resp, err := client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, err
}

/**
reloadConsul used to reload consul services
**/
func reloadConsul(consulSecurityToken string) error {

	command := exec.Command("consul", "reload", "-token="+consulSecurityToken)
	err := command.Run()

	if err != nil {
		return fmt.Errorf("command error output [%v]", err)
	}

	return nil
}

/**
generateTriggers generates array of triggers from supplied gatewayjson
**/
func generateTriggers(gatewayJSON string) ([]triggerDef, error) {

	descriptor, err := model.ParseGatewayDescriptor(gatewayJSON)
	if err != nil {
		return nil, err
	}

	//1. load the configuration, if provided.
	configNamedMap := make(map[string]types.Config)
	for _, config := range descriptor.Gateway.Configurations {
		configNamedMap[config.Name] = config
	}

	var triggers = make([]triggerDef, len(descriptor.Gateway.Triggers))

	for i, trigger := range descriptor.Gateway.Triggers {

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

		if configNamedMap != nil && len(configNamedMap) > 0 {
			//inherit the configuration settings if the trigger uses configuration reference
			err := resolveConfigurationReference(configNamedMap, trigger, mashTriggerSettingsUsable)
			if err != nil {
				return nil, err
			}
		}

		triggers[i].Name = trigger.Name
		triggers[i].Description = trigger.Description
		triggers[i].Type = trigger.Type
		triggers[i].Settings = mashTriggerSettingsUsable

	}

	return triggers, nil
}

/**
resolveConfigurationReference resolves config settings in triggers
**/
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
