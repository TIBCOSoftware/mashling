package functions

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

func AddSampleSum(a, b interface{}) interface{} {

	if a == nil {
		return b
	} else if b == nil {
		return a
	}

	switch x := a.(type) {
	case int:
		return x + b.(int)
	case float64:
		return x + b.(float64)
	case []int:
		y := b.([]int)
		for idx, value := range x {
			x[idx] = value + y[idx]
		}
		return x
	case []float64:
		y := b.([]float64)
		for idx, value := range x {
			x[idx] = value + y[idx]
		}
		return x
	}

	logger.Errorf("unknown")
	panic("invalid input")
}

func AggregateBlocksSum(blocks []interface{}, start int,  size int) interface{} {

	switch blocks[0].(type) {
	case int:
		return sumInt(blocks)
	case float64:
		return sumFloat(blocks)
	case []int:
		return sumIntArray(blocks)
	case []float64:
		return sumFloatArray(blocks)
	}

	//todo handle unsupported type
	return 0
}

func sumInt(blocks []interface{}) interface{} {
	total := 0
	for _, block := range blocks {
		total += block.(int)
	}
	return total
}

func sumFloat(blocks []interface{}) interface{} {
	total := 0.0
	for _, block := range blocks {
		total += block.(float64)
	}
	return total
}

func sumIntArray(blocks []interface{}) interface{} {

	firstBlock := blocks[0].([]int)
	total := make([]int, len(firstBlock))

	for _, block := range blocks {
		arrBlock := block.([]int)
		for i, val := range arrBlock {
			total[i] += val
		}
	}

	return total
}

func sumFloatArray(blocks []interface{}) interface{} {
	firstBlock := blocks[0].([]float64)
	total := make([]float64, len(firstBlock))

	for _, block := range blocks {
		arrBlock := block.([]float64)
		for i, val := range arrBlock {
			total[i] += val
		}
	}

	return total
}


