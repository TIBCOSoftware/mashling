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
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

const tempRepoName = "sampleRepo"

func doGitClone(path, ref string) error {
	cmd := exec.Command("git", "clone", "https://"+ref, tempRepoName)
	cmd.Dir = path
	return cmd.Run()
}

//GetGithubResource used to get github files present in given path
func GetGithubResource(gitHubPath string, resourceFile string) ([]byte, error) {

	tmp, err := ioutil.TempDir("", "github_resource")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)

	tokens := strings.Split(gitHubPath, "/")
	gitRepoPath := gitHubPath
	gitCloneFlag := false

	if len(tokens) == 0 {
		fmt.Println("Invalid github path")
		return nil, nil
	}

	for i := 0; i < len(tokens); i++ {
		err := doGitClone(tmp, gitRepoPath)
		if err == nil {
			gitCloneFlag = true
			break
		}
		index := strings.LastIndex(gitRepoPath, "/")
		if index < 0 {
			gitCloneFlag = false
			break
		}
		gitRepoPath = gitRepoPath[0:index]
	}

	if !gitCloneFlag {
		fmt.Println("Provided github refference is Invalid ", gitHubPath)
		return nil, nil
	}

	resourceFilePath := strings.Replace(gitHubPath, gitRepoPath, "", -1)

	return ioutil.ReadFile(filepath.Join(tmp, tempRepoName, resourceFilePath, resourceFile))
}

//GetTriggerMetadata returns trigger.json for supplied trigger github path
func GetTriggerMetadata(gitHubPath string) (*ftrigger.Metadata, error) {
	goPathVendor := filepath.Join(os.Getenv("GOPATH"), "vendor")
	triggerMetadata := &ftrigger.Metadata{}
	if _, err := os.Stat(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name)); os.IsNotExist(err) {
		fmt.Println("creating trigger.json ", gitHubPath, Gateway_Trigger_Metadata_JSON_Name)
		if _, err := os.Stat(goPathVendor); os.IsNotExist(err) {
			os.Mkdir(goPathVendor, os.ModePerm)
		}
		data, err := GetGithubResource(gitHubPath, Gateway_Trigger_Metadata_JSON_Name)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, triggerMetadata)

		os.Create(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name))
		err = ioutil.WriteFile(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name), data, os.ModePerm)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("reading trigger.json ", gitHubPath, Gateway_Trigger_Metadata_JSON_Name)
		data, err := ioutil.ReadFile(filepath.Join(goPathVendor, gitHubPath, Gateway_Trigger_Metadata_JSON_Name))
		if err != nil {
			return nil, err
		}
		json.Unmarshal(data, triggerMetadata)
	}

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

var PingDataPntr = &PingDataDet{}

type PingDataDet struct {
	MashlingCliRev      string
	MashlingCliLocalRev string
	MashlingCliVersion  string
	SchemaVersion       string
	AppVersion          string
	FlogolibRev         string
	MashlingRev         string
	AppDescrption       string
}
