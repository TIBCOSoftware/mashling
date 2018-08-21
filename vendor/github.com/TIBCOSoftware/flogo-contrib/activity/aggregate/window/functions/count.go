package functions

func AddSampleCount(a, b interface{}) interface{} {

	if a == nil {
		return 1
	}

	return a.(int) + 1
}
