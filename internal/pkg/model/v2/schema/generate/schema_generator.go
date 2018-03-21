// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/jeffreybozek/jsonschema"
)

func main() {
	schema := jsonschema.Reflect(&types.Schema{})
	schemaJSON, err := json.MarshalIndent(schema, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	err = ioutil.WriteFile("schema.json", schemaJSON, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
