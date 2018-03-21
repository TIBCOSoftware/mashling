package flow

import (
	"encoding/json"
)

type Flavor struct {
	// The flow is embedded and uncompressed
	Flow json.RawMessage `json:"flow"`
	// The flow is a URI
	FlowCompressed json.RawMessage `json:"flowCompressed"`
	// The flow is a URI
	FlowURI json.RawMessage `json:"flowURI"`
}
