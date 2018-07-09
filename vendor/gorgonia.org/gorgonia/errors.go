package gorgonia

import (
	"fmt"

	"github.com/pkg/errors"
)

// NoOpError is an error returned when an operation does nothing.
type NoOpError interface {
	NoOp() bool
}

type noopError struct{}

func (e noopError) NoOp() bool    { return true }
func (e noopError) Error() string { return "NoOp" }

// errNoStabilization is an error used internally for when there is no stabilization mechanism is found.
type errNoStabilization interface {
	error
	noStabilization() bool
}

// nostabilizationErr is used internally to communicate that there isn't any stabilization possible
type noStabilizationErr struct{}

func (noStabilizationErr) Error() string         { return "No stabilization mechanism found" }
func (noStabilizationErr) noStabilization() bool { return true }

// noIncrErr is an error used internally when a Value cannot be incremented
type noIncrErr struct {
	v Value
}

func (noIncrErr) Error() string  { return incrErr }
func (e noIncrErr) Value() Value { return e.v }

// oomError represents an Out of tensor.Memory error. It is typically used for CUDA related machine work
type oomError struct {
	res       int64
	allocated int64
}

func (e oomError) Reserved() int64  { return e.res }
func (e oomError) Allocated() int64 { return e.allocated }
func (e oomError) Error() string    { return fmt.Sprintf("allocated/reserved: %v/%v", e.allocated, e.res) }

// AutoDiffError is an error which should be passed if the function is not differentiable. This is useful for Op implementations
type AutoDiffError struct{}

func (err AutoDiffError) Error() string { return "AutoDiffError" }

// vmContextualError is an error that is used to wrap errors that arise from the VM
type vmContextualError struct {
	error
	node  *Node // which node was it processing
	instr int   // what instruction ID it was
}

func (err vmContextualError) Node() *Node        { return err.node }
func (err vmContextualError) Value() Value       { return err.node.Value() }
func (err vmContextualError) InstructionID() int { return err.instr }

func nyi(what string, implFor interface{}) error {
	return errors.Errorf(nyiFail, what, implFor)
}

func nondiffErr(op Op) error {
	return errors.Errorf("%s is a non-differentiable function", op)
}

// checkErrSetDeriv sets the deriv if the error is a Valuer. Helper function for linalg operations
func checkErrSetDeriv(err error, dv *dualValue) error {
	if ver, ok := err.(Valuer); ok {
		return dv.SetDeriv(ver.Value())
	}
	return err
}
