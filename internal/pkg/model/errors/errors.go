package errors

import (
	"bytes"
	"fmt"
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

// UndefinedReference is an error for a reference to another configuration component that does not exist.
type UndefinedReference struct {
	ReferenceType  string
	ReferencedFrom string
	Reference      string
}

//Type is an error when an undefined reference is found.
func (e *UndefinedReference) Type() string {
	return "Undefined reference found"
}

//Details returns the undefined reference details.
func (e *UndefinedReference) Details() string {
	return fmt.Sprintf("%s reference %s from %s is undefined", e.ReferenceType, e.Reference, e.ReferencedFrom)
}
