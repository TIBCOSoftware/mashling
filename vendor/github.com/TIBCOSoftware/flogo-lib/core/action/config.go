package action

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Config is the configuration for the Action
type Config struct {
	Id       string          `json:"id"`
	Ref      string          `json:"ref"`
	Data     json.RawMessage `json:"data"`
	Metadata *ConfigMetadata `json:"metadata"`
}

// Metadata is the configuration metadata for the Action
type ConfigMetadata struct {
	Input  []*data.Attribute `json:"input"`
	Output []*data.Attribute `json:"output"`
}

func GetConfigInputMetadata(act Action) []*data.Attribute {

	if act.Config() != nil {
		if act.Config().Metadata != nil {
			return act.Config().Metadata.Input
		}
	}

	return nil
}

func GetConfigOutputMetadata(act Action) []*data.Attribute {

	if act.Config() != nil {
		if act.Config().Metadata != nil {
			return act.Config().Metadata.Output
		}
	}

	return nil
}