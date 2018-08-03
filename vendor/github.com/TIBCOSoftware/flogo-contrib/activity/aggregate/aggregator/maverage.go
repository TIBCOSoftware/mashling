package aggregator

import "sync"

type MovingAverage struct {
	windowSize   int
	values       []float64
	nextValueIdx int
	full         bool
	mutex        *sync.Mutex
}

func init() {
	RegisterFactory("moving_avg", NewMovingAverage)
}

func (ma *MovingAverage) Add(value float64) (bool, float64) {

	ma.mutex.Lock()
	defer ma.mutex.Unlock()

	//if ma.full && ma.nextValueIdx == 0 && ma.autoReset {
	//	ma.full = false
	//}

	ma.values[ma.nextValueIdx] = value

	ma.nextValueIdx = (ma.nextValueIdx + 1) % ma.windowSize

	if !ma.full && ma.nextValueIdx == 0 {
		ma.full = true
	}

	if ma.full {
		return true, ma.result()
	}

	return false, 0
}

func (ma *MovingAverage) result() float64 {

	var total = float64(0)

	var count = ma.windowSize

	if !ma.full {
		if ma.nextValueIdx == 0 {
			return 0
		}

		count = ma.nextValueIdx
	}

	for i := 0; i < count; i++ {
		total += ma.values[i]
	}

	return total / float64(count)
}

func NewMovingAverage(windowSize int) Aggregator {
	return &MovingAverage{
		windowSize: windowSize,
		values:     make([]float64, windowSize),
		mutex:      &sync.Mutex{},
	}
}
