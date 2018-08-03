package aggregator

import (
	"sync"
	"time"
)

type TimeBlockAverage struct {
	windowSize   time.Duration
	values       []float64
	windowMtx    *sync.Mutex
	startMtx     *sync.RWMutex
	windowActive bool
}

func init() {
	RegisterFactory("timeblockavg", NewTimeBlockAverage)
}

func (ta *TimeBlockAverage) Add(value float64) (bool, float64) {

	ta.windowMtx.Lock()
	ta.values = append(ta.values, value)
	ta.windowMtx.Unlock()

	if ta.startWindow() {
		time.Sleep(ta.windowSize * time.Millisecond)
		return true, ta.average()
	} else {
		return false, 0
	}
}

func (ta *TimeBlockAverage) average() float64 {

	ta.windowMtx.Lock()

	var total = float64(0)

	count := len(ta.values)

	for i := 0; i < count; i++ {
		total += ta.values[i]
	}

	ta.resetWindow()

	return total / float64(count)
}

func (ta *TimeBlockAverage) startWindow() bool {

	ta.startMtx.RLock()

	if ta.windowActive {
		ta.startMtx.RUnlock()
		return false
	}
	ta.startMtx.RUnlock()

	ta.startMtx.Lock()
	defer ta.startMtx.Unlock()

	if !ta.windowActive {
		ta.windowActive = true
		return true
	}

	return false
}

func (ta *TimeBlockAverage) resetWindow() {
	ta.values = nil
	ta.windowMtx.Unlock()

	ta.startMtx.Lock()
	ta.windowActive = false
	ta.startMtx.Unlock()
}

func NewTimeBlockAverage(windowSize int) Aggregator {
	return &TimeBlockAverage{
		windowSize: time.Duration(windowSize),
		windowMtx:  &sync.Mutex{},
		startMtx:   &sync.RWMutex{},
	}
}
