package types

import "encoding/json"
import "github.com/TIBCOSoftware/flogo-lib/core/action"

type Microgateway struct {
	Gateway Gateway `json:"gateway"`
}

type Gateway struct {
	Name           string         `json:"name"`
	Version        string         `json:"version"`
	Description    string         `json:"description,omitempty"`
	Configurations []Config       `json:"configurations"`
	Triggers       []Trigger      `json:"triggers"`
	EventHandlers  []EventHandler `json:"event_handlers"`
	EventLinks     []EventLink    `json:"event_links"`
}

type Config struct {
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Description string          `json:"description,omitempty"`
	Settings    json.RawMessage `json:"settings"`
}

type Trigger struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Type        string          `json:"type"`
	Settings    json.RawMessage `json:"settings"`
}

type EventHandler struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Reference   string          `json:"reference,omitempty"`
	Params      json.RawMessage `json:"params,omitempty"`
	Definition  json.RawMessage `json:"definition,omitempty"`
}

type EventLink struct {
	Triggers   []string   `json:"triggers"`
	Dispatches []Dispatch `json:"dispatches"`
}

type Path struct {
	If      string `json:"if,omitempty"`
	Handler string `json:"handler,omitempty"`
}

type Dispatch struct {
	Path
}

type FlogoAction struct {
	action.Config
}
