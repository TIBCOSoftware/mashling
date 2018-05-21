package errors

import (
	"bytes"
)

// Error are errors specific to gateway parsing and creation.
type Error interface {
	Type() string
	Details() string
}

// MissingDependency is an error for missing Flogo activity or trigger dependencies.
type MissingDependency struct {
	MissingDependencies []string
}

//Type is an error when missing dependencies are found.
func (e *MissingDependency) Type() string {
	return "Missing dependencies found"
}

//Details returns the missing dependencies found.
func (e *MissingDependency) Details() string {
	var buffer bytes.Buffer

	for _, dep := range e.MissingDependencies {
		buffer.WriteString(dep)
		buffer.WriteString("\n")
	}
	return buffer.String()
}
