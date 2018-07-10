package gorgonia

import (
	"github.com/pkg/errors"
	"gorgonia.org/tensor"
)

// BinaryXent is a convenience function for doing binary crossentropy stuff.
// The formula is as below:
// 		-(y * logprob) +  (1-y)(1-logprob)
func BinaryXent(output, target *Node) (retVal *Node, err error) {
	var one *Node
	var logO, omt, omo, tLogO *Node

	// which constant one to use?
	var dt tensor.Dtype
	if dt, err = dtypeOf(output.t); err != nil {
		return nil, errors.Wrapf(err, dtypeExtractionFail, output.t)
	}

	switch dt {
	case Float64:
		one = onef64
	case Float32:
		one = onef32
	default:
		return nil, errors.Errorf(nyiFail, "BinaryXEnt", dt)
	}

	if logO, err = Log(output); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	if omt, err = Sub(one, target); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	if omo, err = Sub(one, output); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	if tLogO, err = HadamardProd(target, logO); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	if retVal, err = Log(omo); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	if retVal, err = HadamardProd(omt, retVal); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	if retVal, err = Add(tLogO, retVal); err != nil {
		return nil, errors.Wrap(err, operationError)
	}

	return Neg(retVal)
}

// Dropout is a convenience function to implement dropout.
// It uses randomly zeroes out a *Tensor with a probability drawn from
// a uniform distribution
func Dropout(x *Node, prob float64) (retVal *Node, err error) {
	if prob == 0.0 {
		return x, nil
	}

	var dt tensor.Dtype
	if dt, err = dtypeOf(x.t); err != nil {
		return nil, errors.Wrap(err, dtypeOfFail)
	}

	var opp, pr Value // opp = 1 per p
	switch dt {
	case Float64:
		opp, _ = anyToScalar(1.0 / prob)
		pr, _ = anyToScalar(prob)
	case Float32:
		opp, _ = anyToScalar(float32(1.0 / prob))
		pr, _ = anyToScalar(float32(prob))
	default:
		return nil, errors.Errorf(nyiTypeFail, "Dropout()", dt)
	}

	p := NewConstant(pr)
	c := NewConstant(opp)

	m := UniformRandomNode(x.g, dt, 0, 1, x.shape...)
	if retVal, err = Gt(m, p, true); err != nil {
		return nil, errors.Wrap(err, "Greater Than failed")
	}

	if retVal, err = HadamardProd(x, retVal); err != nil {
		return nil, errors.Wrap(err, mulFail)
	}

	return HadamardDiv(retVal, c)
}

// Rectify is a convenience function for creating rectified linear units activation functions.
// This function uses >=, which is the canonical version. If you want to use >, you can create
// your own by just following this.
func Rectify(x *Node) (retVal *Node, err error) {
	var zero *Node
	var dt tensor.Dtype

	// which zero to use?
	if dt, err = dtypeOf(x.t); err != nil {
		return nil, errors.Wrap(err, dtypeOfFail)
	}
	switch dt {
	case Float64:
		zero = zerof64
	case Float32:
		zero = zerof32
	default:
		return nil, errors.Errorf(nyiFail, "ReLu", dt)
	}

	cmp := newElemBinOp(gteOpType, x, zero)
	cmp.retSame = true

	if retVal, err = ApplyOp(cmp, x, zero); err != nil {
		return nil, errors.Wrap(err, applyOpFail)
	}

	return HadamardProd(x, retVal)
}

// Im2Col converts a BCHW image block to columns. The kernel, pad and stride parameter must be shape of size 2, no more no less
// This poor naming scheme clearly comes from matlab
func Im2Col(n *Node, kernel, pad, stride tensor.Shape) (retVal *Node, err error) {
	if kernel.Dims() != 2 {
		return nil, errors.Errorf("kernel shape is supposed to have a dim of 2")
	}
	if pad.Dims() != 2 {
		return nil, errors.Errorf("pad is supposed to have a dim of 2")
	}
	if stride.Dims() != 2 {
		return nil, errors.Errorf("strides is supposed to have a dim of 2")
	}

	if kernel[0] <= 0 || kernel[1] <= 0 {
		return nil, errors.Errorf("cannot have negative or 0 in kernel shape")
	}

	if stride[0] <= 0 || stride[1] <= 0 {
		return nil, errors.Errorf("cannot have negative or 0 in stride: %v", stride)
	}

	if pad[0] < 0 || pad[1] < 0 {
		return nil, errors.Errorf("cannot have negative padding")
	}
	op := makeIm2ColOp(kernel[0], kernel[1], pad[0], pad[1], stride[0], stride[1])
	return ApplyOp(op, n)
}

// Conv2d is a simple 2D convoution, to be used for CPU computation only. If CuDNN is used, use the CUDAConv2D function.
// These are the properties the inputs must fulfil:
//
// im: must have 4D shape. Expected format is BCHW (batch, channel, height, width)
// filter: must have 4D shape: (batch, kernel, height, width)
// kernelShape: shape of the filter kernel
// pad: len(pad) == 2
// stride: len(stride) == 2
func Conv2d(im, filter *Node, kernelShape tensor.Shape, pad, stride []int) (retVal *Node, err error) {
	// checks
	for _, s := range stride {
		if s <= 0 {
			return nil, errors.Errorf("Cannot use strides of less than or equal 0: %v", stride)
		}
	}

	for _, p := range pad {
		if p < 0 {
			return nil, errors.Errorf("Cannot use padding of less than 0: %v", pad)
		}
	}

	var colIm *Node
	if colIm, err = Im2Col(im, kernelShape, pad, stride); err != nil {
		return
	}

	layer := filter.Shape()[0]
	kernel := filter.Shape()[1]
	row := filter.Shape()[2]
	col := filter.Shape()[3]

	var flattened *Node
	if flattened, err = Reshape(filter, tensor.Shape{layer, kernel * row * col}); err != nil {
		return
	}

	// extract patch
	batch := colIm.Shape()[0]
	m := colIm.Shape()[1]
	n := colIm.Shape()[2]
	z := colIm.Shape()[3]

	var patch, colImLayer *Node
	if patch, err = Reshape(colIm, tensor.Shape{batch * m * n, z}); err != nil {
		return
	}

	op := linAlgBinOp{
		āBinaryOperator: matMulOperator,
		transA:          false,
		transB:          true,
	}

	if colImLayer, err = ApplyOp(op, patch, flattened); err != nil {
		return
	}

	// now reshape and transpose the values back into the original order
	var res *Node
	if res, err = Reshape(colImLayer, tensor.Shape{batch, m, n, layer}); err != nil {
		return
	}
	return Transpose(res, 0, 3, 1, 2)
}

// Conv1d is a 1D convlution. It relies on Conv2D
func Conv1d(in, filter *Node, kernel, pad, stride int) (*Node, error) {
	return Conv2d(in, filter, tensor.Shape{1, kernel}, []int{0, pad}, []int{1, stride})
}

func MaxPool2D(x *Node, kernel tensor.Shape, pad, stride []int) (*Node, error) {
	xShape := x.Shape()
	h, w := xShape[2], xShape[3]
	kh, kw := kernel[0], kernel[1]
	ph, pw := pad[0], pad[1]

	// check shape
	if xShape.Dims() != 4 {
		return nil, errors.Errorf("Expected input to have a shape with dimension 4")
	}
	if kernel.Dims() != 2 {
		return nil, errors.Errorf("Expected kernel to have a shape of dimension 2")
	}

	if h-kh == 0 && ph == 0 {
		// error
		return nil, errors.New("Impossible height/kernel/pad combination")
	}

	if w-kw == 0 && pw == 0 {
		// error
		return nil, errors.New("Impossible width/kernel/pad combination")
	}

	op := newMaxPoolOp(xShape, kernel, pad, stride)
	return ApplyOp(op, x)
}
