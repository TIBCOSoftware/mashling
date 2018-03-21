package action

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for the Activity
type Metadata struct {
	ID      string
	Options map[string]*data.Attribute
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

// MarshalJSON overrides the default MarshalJSON for TaskEnv
func (md *Metadata) MarshalJSON() ([]byte, error) {

	options := make([]*data.Attribute, 0, len(md.Options))

	for _, value := range md.Options {
		options = append(options, value)
	}

	return json.Marshal(&struct {
		Name    string            `json:"name"`
		Ref     string            `json:"ref"`
		Options []*data.Attribute `json:"options"`
	}{
		Name:    md.ID,
		Options: options,
	})
}

// UnmarshalJSON overrides the default UnmarshalJSON for TaskEnv
func (md *Metadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Name    string            `json:"name"`
		Ref     string            `json:"ref"`
		Options []*data.Attribute `json:"options"`
	}{}

	if err := json.Unmarshal(b, ser); err != nil {
		return err
	}

	if len(ser.Ref) > 0 {
		md.ID = ser.Ref
	} else {
		// Added for backwards compatibility
		// TODO remove and add a proper error once the BC is removed
		md.ID = ser.Name
	}

	md.Options = make(map[string]*data.Attribute, len(ser.Options))

	for _, attr := range ser.Options {
		md.Options[attr.Name()] = attr
	}

	return nil
}
