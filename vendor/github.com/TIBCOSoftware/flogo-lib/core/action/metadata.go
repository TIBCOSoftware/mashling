package action

import (
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// Metadata is the metadata for the Activity
type Metadata struct {
	ID       string `json:"ref"`
	Version  string
	Settings map[string]*data.Attribute
	Input    map[string]*data.Attribute
	Output   map[string]*data.Attribute

	Async    bool `json:"async"`
	Passthru bool `json:"passthru"`
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

// UnmarshalJSON overrides the default UnmarshalJSON for TaskEnv
func (md *Metadata) UnmarshalJSON(b []byte) error {

	ser := &struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Ref     string `json:"ref"`

		Settings []*data.Attribute `json:"settings"`
		Input    []*data.Attribute `json:"input"`
		Output   []*data.Attribute `json:"output"`
		Async    bool              `json:"async"`
		Passthru bool              `json:"passthru"`
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

	md.Version = ser.Version

	md.Async = ser.Async
	md.Passthru = ser.Passthru

	md.Settings = make(map[string]*data.Attribute, len(ser.Settings))
	md.Input = make(map[string]*data.Attribute, len(ser.Input))
	md.Output = make(map[string]*data.Attribute, len(ser.Output))

	for _, attr := range ser.Settings {
		md.Settings[attr.Name()] = attr
	}

	if len(ser.Input) > 0 {
		for _, attr := range ser.Input {
			md.Input[attr.Name()] = attr
		}
	}

	if len(ser.Output) > 0 {
		for _, attr := range ser.Output {
			md.Output[attr.Name()] = attr
		}
	}

	return nil
}
