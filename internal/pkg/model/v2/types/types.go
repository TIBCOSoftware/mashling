package types

// Schema contains schema version and configuration information for a gateway instance.
type Schema struct {
	Version string  `json:"mashling_schema" jsonschema:"required"`
	Gateway Gateway `json:"gateway" jsonschema:"required"`
}

// Gateway contains the runtime behavior of a gateway instance.
type Gateway struct {
	Name         string     `json:"name" jsonschema:"required"`
	Version      string     `json:"version" jsonschema:"required"`
	DisplayName  string     `json:"display_name,omitempty"`
	DisplayImage string     `json:"display_image,omitempty"`
	Description  string     `json:"description,omitempty"`
	Triggers     []Trigger  `json:"triggers" jsonschema:"required"`
	Dispatches   []Dispatch `json:"dispatches" jsonschema:"required"`
	Services     []Service  `json:"services,omitempty"`
	Policies     []Policy   `json:"policies,omitempty"`
}

// Trigger contains the event listener definitions and configurations.
type Trigger struct {
	Name        string                 `json:"name" jsonschema:"required"`
	Type        string                 `json:"type" jsonschema:"required"`
	Description string                 `json:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty" jsonschema:"additionalProperties"`
	Handlers    []Handler              `json:"handlers" jsonschema:"required"`
}

// Dispatch links events on a trigger to execution flows.
type Dispatch struct {
	Name   string  `json:"name" jsonschema:"required"`
	Routes []Route `json:"routes" jsonschema:"required"`
}

// Route conditionally defines an execution flow.
type Route struct {
	Condition string     `json:"if,omitempty"`
	Async     bool       `json:"async,omitempty"`
	Policies  []string   `json:"policies,omitempty"`
	Steps     []Step     `json:"steps" jsonschema:"required"`
	Responses []Response `json:"responses,omitempty"`
}

// Step conditionally defines a step in a route's execution flow.
type Step struct {
	Condition string                 `json:"if,omitempty"`
	Service   string                 `json:"service" jsonschema:"required"`
	Input     map[string]interface{} `json:"input,omitempty" jsonschema:"additionalProperties"`
}

// Response defines response handling rules.
type Response struct {
	Condition string `json:"if,omitempty"`
	Error     bool   `json:"error" jsonschema:"required"`
	Output    Output `json:"output,omitempty" jsonschema:"required"`
}

// Output defines response output values back to a trigger event.
type Output struct {
	Code int         `json:"code" jsonschema:"required"`
	Data interface{} `json:"data" jsonschema:"additionalProperties"`
}

// Service defines a functional target that may be invoked by a step in an execution flow.
type Service struct {
	Name        string                 `json:"name" jsonschema:"required"`
	Type        string                 `json:"type" jsonschema:"required"`
	Description string                 `json:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty" jsonschema:"additionalProperties"`
}

// Policy defines an invocation rule that may be applied to a route.
type Policy struct {
	Name        string                 `json:"name" jsonschema:"required"`
	Type        string                 `json:"type" jsonschema:"required"`
	Description string                 `json:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty" jsonschema:"additionalProperties"`
}

// Handler maps a trigger and settings to a specific dispatch
type Handler struct {
	Dispatch string                 `json:"dispatch" jsonschema:"required"`
	Settings map[string]interface{} `json:"settings,omitempty" jsonschema:"additionalProperties"`
}
