//go:generate go-bindata -pkg pattern -o assets.go DefaultHttpPattern.json DefaultChannelPattern.json

package pattern

import (
	"encoding/json"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

type PatternDefinition struct {
	Dispatch types.Dispatch  `json:"dispatch"`
	Services []types.Service `json:"services"`
}

func Load(pattern string) (*PatternDefinition, error) {
	patternJSON, err := Asset(pattern + ".json")
	if err != nil {
		return nil, err
	}
	pDef := &PatternDefinition{}
	err = json.Unmarshal(patternJSON, pDef)
	if err != nil {
		return nil, err
	}
	return pDef, nil
}
