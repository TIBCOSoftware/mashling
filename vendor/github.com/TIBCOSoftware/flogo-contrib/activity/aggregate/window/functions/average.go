package functions

func AggregateBlocksAvg(blocks []interface{}, start int, size int) interface{} {

	switch blocks[0].(type) {
	case int:
		return avgInt(blocks, size)
	case float64:
		return avgFloat(blocks, size)
	case []int:
		return avgIntArray(blocks, size)
	case []float64:
		return avgFloatArray(blocks, size)
	}

	//todo handle unsupported type
	return 0
}

func avgInt(blocks []interface{}, size int) interface{} {
	total := 0
	for _, block := range blocks {
		total += block.(int)
	}
	return total / (len(blocks) * size)
}

func avgFloat(blocks []interface{}, size int) interface{} {
	total := 0.0
	for _, block := range blocks {
		total += block.(float64)
	}
	return total / float64(len(blocks) * size)
}

func avgIntArray(blocks []interface{}, size int) interface{} {

	firstBlock := blocks[0].([]int)
	result := make([]int, len(firstBlock))

	for _, block := range blocks {
		arrBlock := block.([]int)
		for i, val := range arrBlock {
			result[i] += val
		}
	}

	for i, val := range result {
		result[i] = val / (len(blocks)*size)
	}

	return result
}

func avgFloatArray(blocks []interface{}, size int) interface{} {
	firstBlock := blocks[0].([]float64)
	result := make([]float64, len(firstBlock))

	for _, block := range blocks {
		arrBlock := block.([]float64)
		for i, val := range arrBlock {
			result[i] += val
		}
	}

	for i, val := range result {
		result[i] = val / float64(len(blocks) * size)
	}

	return result
}

func AggregateSingleAvg(a interface{}, count int) interface{} {
	switch x := a.(type) {
	case int:
		return x / count
	case float64:
		return x / float64(count)
	case []int:
		ret := make([]int, len(x))
		copy(ret, x)
		for idx, value := range ret {
			ret[idx] = value / count
		}
		return ret
	case []float64:
		ret := make([]float64, len(x))
		copy(ret, x)
		for idx, value := range ret {
			ret[idx] = value / float64(count)
		}
		return ret
	}

	//todo handle unsupported type
	return 0
}