package functions

func AddSampleMin(a, b interface{}) interface{} {

	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	switch x := a.(type) {
	case int:
		if x < b.(int) {
			return a
		}
		return b
	case float64:
		if x < b.(float64) {
			return a
		}
		return b
	case []int:
		y := b.([]int)
		for idx, value := range x {
			if  y[idx] < value {
				x[idx] = y[idx]
			}
		}
		return x
	case []float64:
		y := b.([]float64)
		for idx, value := range x {
			if  y[idx] < value {
				x[idx] = y[idx]
			}
		}
		return x
	}

	panic("invalid input")
}


func AggregateBlocksMin(blocks []interface{}, start int, size int) interface{} {

	switch blocks[0].(type) {
	case int:
		return minInt(blocks)
	case float64:
		return minFloat(blocks)
	case []int:
		return minIntArray(blocks)
	case []float64:
		return minFloatArray(blocks)
	}

	//todo handle unsupported type
	return 0
}

func minInt(blocks []interface{}) interface{} {
	min := blocks[0].(int)

	for _, block := range blocks {
		if block.(int) < min {
			min = block.(int)
		}
	}
	return min
}

func minFloat(blocks []interface{}) interface{} {

	min := blocks[0].(float64)

	for _, block := range blocks {
		if block.(float64) < min {
			min = block.(float64)
		}
	}
	return min
}

func minIntArray(blocks []interface{}) interface{} {

	firstBlock := blocks[0].([]int)
	var min []int
	copy(min, firstBlock)

	for _, block := range blocks {
		arrBlock := block.([]int)
		for i, val := range arrBlock {
			if val < min[i] {
				min[i] = val
			}
		}
	}

	return min
}

func minFloatArray(blocks []interface{}) interface{} {
	firstBlock := blocks[0].([]float64)
	var min []float64
	copy(min, firstBlock)

	for _, block := range blocks {
		arrBlock := block.([]float64)
		for i, val := range arrBlock {
			if val < min[i] {
				min[i] = val
			}
		}
	}

	return min
}
