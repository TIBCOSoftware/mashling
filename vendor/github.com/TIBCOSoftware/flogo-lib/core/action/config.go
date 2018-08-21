package action

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Config is the configuration for the Action
type Config struct {
	//inline action
	Ref      string                 `json:"ref"`
	Settings map[string]interface{} `json:"settings"`
	Data     json.RawMessage        `json:"data"`

	//referenced action
	Id string `json:"id"`

	// Deprecated: No longer used
	Metadata *data.IOMetadata `json:"metadata"`
}

//do we need a call that will "fix up" the config, coerce to the right attr, using the metadata?

