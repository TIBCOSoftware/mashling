package functions

func AddSampleMax(a, b interface{}) interface{} {

	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	switch x := a.(type) {
	case int:
		if x > b.(int) {
			return a
		}
		return b
	case float64:
		if x > b.(float64) {
			return a
		}
		return b
	case []int:
		y := b.([]int)
		for idx, value := range x {
			if  y[idx] > value {
				x[idx] = y[idx]
			}
		}
		return x
	case []float64:
		y := b.([]float64)
		for idx, value := range x {
			if  y[idx] > value {
				x[idx] = y[idx]
			}
		}
		return x
	}

	panic("invalid input")
}

func AggregateBlocksMax(blocks []interface{}, start int, size int) interface{} {

	switch blocks[0].(type) {
	case int:
		return maxInt(blocks)
	case float64:
		return maxFloat(blocks)
	case []int:
		return maxIntArray(blocks)
	case []float64:
		return maxFloatArray(blocks)
	}

	//todo handle unsupported type
	return 0
}

func AggregateBlocksCount(blocks []interface{}, start int, size int) interface{} {
	return len(blocks)
}

func maxInt(blocks []interface{}) interface{} {
	max := blocks[0].(int)

	for _, block := range blocks {
		if block.(int) > max {
			max = block.(int)
		}
	}
	return max
}

func maxFloat(blocks []interface{}) interface{} {

	max := blocks[0].(float64)

	for _, block := range blocks {
		if block.(float64) > max {
			max = block.(float64)
		}
	}
	return max
}

func maxIntArray(blocks []interface{}) interface{} {

	firstBlock := blocks[0].([]int)
	var max []int
	copy(max, firstBlock)

	for _, block := range blocks {
		arrBlock := block.([]int)
		for i, val := range arrBlock {
			if val > max[i] {
				max[i] = val
			}
		}
	}

	return max
}

func maxFloatArray(blocks []interface{}) interface{} {
	firstBlock := blocks[0].([]float64)
	var max []float64
	copy(max, firstBlock)

	for _, block := range blocks {
		arrBlock := block.([]float64)
		for i, val := range arrBlock {
			if val > max[i] {
				max[i] = val
			}
		}
	}

	return max
}
