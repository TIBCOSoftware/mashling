package activity

type Config struct {
	Ref         string                 `json:"ref"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	InputAttrs  map[string]interface{} `json:"input,omitempty"`
	OutputAttrs map[string]interface{} `json:"output,omitempty"`
}

func (c *Config) FixUp(metadata *Metadata)  {

	//fixup settings

	//generate input

	//fix up outputs
}
