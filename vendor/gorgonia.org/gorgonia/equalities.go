package gorgonia

import "gorgonia.org/tensor"

func scalarEq(a, b Scalar) bool {
	switch at := a.(type) {
	case *F64:
		if bt, ok := b.(*F64); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	case *F32:
		if bt, ok := b.(*F32); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	case *I:
		if bt, ok := b.(*I); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	case *I32:
		if bt, ok := b.(*I32); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	case *I64:
		if bt, ok := b.(*I64); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	case *U8:
		if bt, ok := b.(*U8); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	case *B:
		if bt, ok := b.(*B); ok {
			if at == bt {
				return true
			}
			return *at == *bt
		}
		return false
	}
	return false
}

func scalarClose(a, b Scalar) bool {
	switch at := a.(type) {
	case *F64:
		if bt, ok := b.(*F64); ok {
			return closeF64(float64(*at), float64(*bt))
		}
		return false
	case *F32:
		if bt, ok := b.(*F32); ok {
			return closeF32(float32(*at), float32(*bt))
		}
		return false
	default:
		return scalarEq(a, b)
	}
}

func closeF32(a, b float32) bool {
	const EPSILON32 = 1e-5
	if (a-b) < EPSILON32 && (b-a) < EPSILON32 {
		return true
	}
	return false
}

func closeF64(a, b float64) bool {
	const EPSILON64 = 1e-10
	if (a-b) < EPSILON64 && (b-a) < EPSILON64 {
		return true
	}
	return false
}

func tensorClose(a, b tensor.Tensor) bool {
	aDt := a.Dtype()
	bDt := b.Dtype()
	if aDt != bDt {
		return false
	}

	switch aDt {
	case tensor.Float64:
		aFs := a.Data().([]float64)
		bFs := b.Data().([]float64)
		if len(aFs) != len(bFs) {
			return false
		}
		aFs = aFs[:len(aFs)]
		bFs = bFs[:len(aFs)]
		for i, v := range aFs {
			if !closeF64(v, bFs[i]) {
				return false
			}
		}
		return true
	case tensor.Float32:
		aFs := a.Data().([]float32)
		bFs := b.Data().([]float32)
		if len(aFs) != len(bFs) {
			return false
		}
		aFs = aFs[:len(aFs)]
		bFs = bFs[:len(aFs)]
		for i, v := range aFs {
			if !closeF32(v, bFs[i]) {
				return false
			}
		}
		return true
	default:
		return a.Eq(b)
	}

}

/*
func axesEq(a, b axes) bool {
	if len(a) != len(b) {
		return false
	}

	for i, s := range a {
		if b[i] != s {
			return false
		}
	}
	return true
}

// yes it's exactly the same as axesEq
func coordEq(a, b coordinates) bool {
	if len(a) != len(b) {
		return false
	}

	for i, s := range a {
		if b[i] != s {
			return false
		}
	}
	return true
}
*/

func constEq(a, b constant) (ok bool) {
	switch at := a.(type) {
	case constantScalar:
		var bt constantScalar
		if bt, ok = b.(constantScalar); !ok {
			return
		}

		return bt == at
	case constantTensor:
		var bt constantTensor
		if bt, ok = b.(constantTensor); !ok {
			return
		}
		return at.v.Eq(bt.v)
	default:
		panic("Not yet implemented")
	}
}

// fastest comparisons to least fastest
func nodeEq(a, b *Node) bool {
	if a == b {
		return true
	}

	if a.isInput() {
		if !b.isInput() {
			return false
		}
		return a.name == b.name
	}

	if b.isInput() {
		return false
	}

	// hashcode is good for comparing Op (TODO: benchmark this vs reflect.DeepEq)
	if a.op.Hashcode() != b.op.Hashcode() {
		return false
	}

	if len(a.children) != len(b.children) {
		return false
	}

	if a.t != b.t {
		return false
	}

	if !a.shape.Eq(b.shape) {
		return false
	}

	return true
}
