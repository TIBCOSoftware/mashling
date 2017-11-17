/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package types

import "encoding/json"
import "github.com/TIBCOSoftware/flogo-lib/core/action"

type Microgateway struct {
	MashlingSchema string  `json:"mashling_schema"`
	Gateway        Gateway `json:"gateway"`
}

type Gateway struct {
	Name           string         `json:"name"`
	Version        string         `json:"version"`
	DisplayName    string         `json:"display_name"`
	DisplayImage   string         `json:"display_image"`
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
	If          string          `json:"if,omitempty"`
	Handler     string          `json:"handler,omitempty"`
	InputParams json.RawMessage `json:"inputParams,omitempty"`
}

type Dispatch struct {
	Path
}

type FlogoAction struct {
	action.Config
}
