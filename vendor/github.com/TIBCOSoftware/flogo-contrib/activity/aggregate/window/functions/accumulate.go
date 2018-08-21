package functions


func AddSampleAccum(a, b interface{}) interface{} {

	var accum []interface{}

	if a == nil {
		accum = make([]interface{}, 0)
	} else {
		accum = a.([]interface{})
	}

	accum = append(accum, b)

	return accum
}


func AggregateBlocksAccumulate(blocks []interface{}, start int, size int) interface{} {

	accum := make([]interface{}, 0, len(blocks))

	for i := 0; i < len(blocks); i++ {

		accum = append(accum, blocks[(start+i)%len(blocks)])
	}

	return accum
}
