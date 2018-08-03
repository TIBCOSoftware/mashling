package aggregator

import "sync"

type BlockAverage struct {
	windowSize   int
	values       []float64
	nextValueIdx int
	mutex        *sync.Mutex
}

func init() {
	RegisterFactory("block_avg", NewBlockAverage)
}

func (ba *BlockAverage) Add(value float64) (bool, float64) {

	ba.mutex.Lock()
	defer ba.mutex.Unlock()

	ba.values[ba.nextValueIdx] = value
	ba.nextValueIdx = ba.nextValueIdx + 1

	if ba.nextValueIdx >= ba.windowSize {
		return true, ba.average()
	}

	return false, 0
}

func (ba *BlockAverage) average() float64 {

	var total = float64(0)

	for i := 0; i < ba.windowSize; i++ {
		total += ba.values[i]
	}
	ba.nextValueIdx = 0

	return total / float64(ba.windowSize)
}

func NewBlockAverage(windowSize int) Aggregator {
	return &BlockAverage{
		windowSize: windowSize,
		values:     make([]float64, windowSize),
		mutex:      &sync.Mutex{},
	}
}
