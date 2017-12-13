/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/env"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

func GetGithubResource(gitHubPath string, resourceFile string) ([]byte, error) {
	gbProject := env.NewGbProjectEnv()
	tmp, err := ioutil.TempDir("", "github_resource")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	err = os.Mkdir(tmp+"/src", 0755)
	if err != nil {
		return nil, err
	}
	gbProject.Init(tmp)

	resourceDir := gbProject.GetVendorSrcDir()
	resourcePath := resourceDir + "/" + gitHubPath + "/" + resourceFile

	gbProject.InstallDependency(gitHubPath, "")

	return ioutil.ReadFile(resourcePath)
}

func GetTriggerMetadata(gitHubPath string) (*ftrigger.Metadata, error) {
	gbProject := env.NewGbProjectEnv()

	gbProject.Init(os.Getenv("GOPATH"))

	resourceDir := gbProject.GetVendorSrcDir()
	triggerPath := resourceDir + "/" + gitHubPath + "/" + Gateway_Trigger_Metadata_JSON_Name

	gbProject.InstallDependency(gitHubPath, "")
	data, err := ioutil.ReadFile(triggerPath)
	if err != nil {
		return nil, err
	}
	triggerMetadata := &ftrigger.Metadata{}
	json.Unmarshal(data, triggerMetadata)
	return triggerMetadata, nil
}

func IsValidTriggerSetting(metadata *ftrigger.Metadata, property string) bool {
	settings := metadata.Settings
	for key := range settings {
		if key == property {
			return true
		}
	}

	return false
}

func IsValidTriggerHandlerSetting(metadata *ftrigger.Metadata, property string) bool {
	settings := metadata.Handler.Settings

	for _, element := range settings {
		if element.Name() == property {
			return true
		}
	}

	return false
}

func ValidateTriggerConfigExpr(expression *string) (bool, *string) {
	if expression == nil {
		return false, nil
	}

	exprValue := *expression
	if strings.HasPrefix(exprValue, Gateway_Trigger_Config_Prefix) && strings.HasSuffix(exprValue, Gateway_Trigger_Config_Suffix) {
		//get name of the config
		str := exprValue[len(Gateway_Trigger_Config_Prefix) : len(exprValue)-1]
		return true, &str
	} else {
		return false, &exprValue
	}
}

func CheckTriggerOptimization(triggerSettings map[string]interface{}) bool {
	if val, ok := triggerSettings[Gateway_Trigger_Optimize_Property]; ok {
		optimize, err := strconv.ParseBool(val.(string))
		if err != nil {
			//check if its a boolean
			optimize, found := val.(bool)
			if !found {
				return found
			}
			return optimize
		}
		return optimize
	} else {
		return Gateway_Trigger_Optimize_Property_Default
	}
}

func validateEnvPropertySettingExpr(expression *string) (bool, *string) {
	if expression == nil {
		return false, nil
	}

	exprValue := *expression
	if strings.HasPrefix(exprValue, Gateway_Trigger_Setting_Env_Prefix) && strings.HasSuffix(exprValue, Gateway_Trigger_Setting_Env_Suffix) {
		//get name of the property
		str := exprValue[len(Gateway_Trigger_Setting_Env_Prefix) : len(exprValue)-1]
		return true, &str
	}
	return false, &exprValue
}

// ResolveEnvironmentProperties resolves environment properties mentioned in the settings map.
func ResolveEnvironmentProperties(settings map[string]interface{}) error {
	for k, v := range settings {
		value := v.(string)
		valid, propertyName := validateEnvPropertySettingExpr(&value)
		if !valid {
			continue
		}
		//lets get the env property value
		propertyNameStr := *propertyName
		propertyValue, found := os.LookupEnv(propertyNameStr)
		if !found {
			return fmt.Errorf("environment property [%v] is not set", propertyNameStr)
		}
		settings[k] = propertyValue
	}
	return nil
}
