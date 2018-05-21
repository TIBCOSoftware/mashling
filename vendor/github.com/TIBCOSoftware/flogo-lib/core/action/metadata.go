package action

import (
	"encoding/json"
)

// Metadata is the metadata for the Activity
type Metadata struct {
	ID    string `json:"ref"`
	Async bool   `json:"async"`
}

// NewMetadata creates the metadata object from its json representation
func NewMetadata(jsonMetadata string) *Metadata {
	md := &Metadata{}
	err := json.Unmarshal([]byte(jsonMetadata), md)
	if err != nil {
		panic("Unable to parse activity metadata: " + err.Error())
	}

	return md
}
