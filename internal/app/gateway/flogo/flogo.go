//go:generate go run registry/generate/registry_generator.go

package flogo

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/extension"
)

// ResetGlobalContext sets a new extension provider. It is only used in dev mode
// for hot reloading of configuration files.
func ResetGlobalContext() {
	flow.SetExtensionProvider(extension.New())
}
