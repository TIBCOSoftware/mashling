package action

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Config is the configuration for the Action
type Config struct {
	Ref      string           `json:"ref"`
	Data     json.RawMessage  `json:"data"`
	Mappings *data.IOMappings `json:"mappings"`

	// Deprecated: No longer used
	Id string `json:"id"`
	// Deprecated: No longer used
	Metadata *data.IOMetadata `json:"metadata"`
}
