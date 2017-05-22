package util

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-cli/env"
	ftrigger "github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"io/ioutil"
	"os"
)

func GetGithubResource(gitHubPath string, resourceFile string) ([]byte, error) {
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
	return data, nil
}

func GetTriggerMetadata(gitHubPath string) (*ftrigger.Metadata, error) {
	gbProject := env.NewGbProjectEnv()

	gbProject.Init(os.Getenv("GOPATH"))

	resourceDir := gbProject.GetVendorSrcDir()
	triggerPath := resourceDir + "/" + gitHubPath + "/trigger.json"

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
		if element.Name == property {
			return true
		}
	}

	return false
}
