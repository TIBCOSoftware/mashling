package window

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Settings struct {
	Size          int
	Resolution    int
	ExternalTimer bool

	TotalCountModifier int
}

func (s *Settings) SetAdditionalSettings(as map[string]string) error {

	for key, value := range as {
		if strings.ToLower(key) == "totalcountmodifier" {
			//todo should we return an error?
			s.TotalCountModifier, _ = strconv.Atoi(value)
		}
	}

	return nil
}

///////////////////
// Tumbling Window

func NewTumblingWindow(addFunc AddSampleFunc, aggFunc AggregateSingleFunc, settings *Settings) Window {

	return &TumblingWindow{addFunc: addFunc, aggFunc: aggFunc, settings: settings, mutex: &sync.Mutex{}}
}

//note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
type TumblingWindow struct {
	addFunc  AddSampleFunc
	aggFunc  AggregateSingleFunc
	settings *Settings

	data       interface{}
	numSamples int

	mutex *sync.Mutex
}

// AddSample implements window.Window.AddSample
func (w *TumblingWindow) AddSample(sample interface{}) (bool, interface{}) {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	//sample size should match data size
	w.data = w.addFunc(w.data, sample)
	w.numSamples++

	if w.numSamples == w.settings.Size {
		// aggregate and emit
		val := w.aggFunc(w.data, w.settings.Size)

		w.numSamples = 0
		w.data, _ = zero(w.data)

		return true, val
	}

	return false, nil
}

///////////////////////
// Tumbling Time Window

func NewTumblingTimeWindow(addFunc AddSampleFunc, aggFunc AggregateSingleFunc, settings *Settings) TimeWindow {
	return &TumblingTimeWindow{addFunc: addFunc, aggFunc: aggFunc, settings: settings, mutex: &sync.Mutex{}}
}

// TumblingTimeWindow - A tumbling window based on time. Relies on external entity moving window along
// by calling NextBlock at the appropriate time.
//note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
type TumblingTimeWindow struct {
	addFunc  AddSampleFunc
	aggFunc  AggregateSingleFunc
	settings *Settings

	data       interface{}
	maxSamples int
	numSamples int

	nextEmit int
	lastAdd  int

	mutex *sync.Mutex
}

func (w *TumblingTimeWindow) AddSample(sample interface{}) (bool, interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.data = w.addFunc(w.data, sample)
	w.numSamples++

	if w.numSamples > w.maxSamples {
		w.maxSamples = w.numSamples
	}

	if !w.settings.ExternalTimer {
		currentTime := getTimeMillis()

		//todo what do we do if this greatly exceeds the nextEmit time?
		if currentTime >= w.nextEmit {
			w.nextEmit = +w.settings.Size // size == time in millis
			return w.nextBlock()
		}
	}

	return false, nil
}

func (w *TumblingTimeWindow) NextBlock() (bool, interface{}) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.nextBlock()
}

func (w *TumblingTimeWindow) nextBlock() (bool, interface{}) {

	// aggregate and emit
	val := w.aggFunc(w.data, w.maxSamples) //num samples or max samples?

	w.numSamples = 0
	w.data, _ = zero(w.data)

	if w.settings.TotalCountModifier > 0 {
		//local, so reset max samples
		//todo in the future use average of last N 'numSamples' to calculate max
		w.maxSamples = 0
	}

	return true, val
}

///////////////////
// Sliding Window

func NewSlidingWindow(aggFunc AggregateBlocksFunc, settings *Settings) Window {

	w := &SlidingWindow{aggFunc: aggFunc, settings: settings}
	w.blocks = make([]interface{}, settings.Size)
	w.mutex = &sync.Mutex{}

	return w
}

//note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
// todo split external vs on-add timer
type SlidingWindow struct {
	aggFunc  AggregateBlocksFunc
	settings *Settings

	blocks       []interface{}
	numSamples   int
	currentBlock int
	canEmit      bool

	mutex *sync.Mutex
}

// AddSample implements window.Window.AddSample
func (w *SlidingWindow) AddSample(sample interface{}) (bool, interface{}) {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	//sample size should match data size
	w.blocks[w.currentBlock] = sample //no addSampleFunc required, just tracking all values

	if !w.canEmit {
		if w.currentBlock == w.settings.Size-1 {
			w.canEmit = true
		}
	}

	w.numSamples++

	if w.canEmit && w.numSamples >= w.settings.Resolution {

		// aggregate and emit
		val := w.aggFunc(w.blocks, w.currentBlock, 1)

		w.numSamples = 0
		w.currentBlock++

		w.currentBlock = w.currentBlock % w.settings.Size

		return true, val
	}

	w.currentBlock++

	return false, nil
}

//////////////////////
// Sliding Time Window

func NewSlidingTimeWindow(addFunc AddSampleFunc, aggFunc AggregateBlocksFunc, settings *Settings) TimeWindow {

	numBlocks := settings.Size / settings.Resolution

	w := &SlidingTimeWindow{addFunc: addFunc, aggFunc: aggFunc, numBlocks: numBlocks, settings: settings}

	w.blocks = make([]interface{}, numBlocks)
	w.mutex = &sync.Mutex{}

	return w
}

// SlidingTimeWindow - A sliding window based on time. Relies on external entity moving window along
// by calling NextBlock at the appropriate time.
// note:  using interface{} 4x slower than using specific types, starting with interface{} for expediency
type SlidingTimeWindow struct {
	addFunc  AddSampleFunc
	aggFunc  AggregateBlocksFunc
	settings *Settings

	numBlocks    int
	blocks       []interface{}
	maxSamples   int
	numSamples   int
	currentBlock int
	canEmit      bool

	nextBlockTime int
	lastAdd       int

	mutex *sync.Mutex
}

// AddSample implements window.Window.AddSample
func (w *SlidingTimeWindow) AddSample(sample interface{}) (bool, interface{}) {

	w.mutex.Lock()
	defer w.mutex.Lock()

	//sample size should match data size
	w.blocks[w.currentBlock] = w.addFunc(w.blocks[w.currentBlock], sample)

	w.numSamples++

	if w.numSamples > w.maxSamples {
		w.maxSamples = w.numSamples
	}

	if !w.settings.ExternalTimer {
		currentTime := getTimeMillis()

		if currentTime > w.nextBlockTime {
			w.nextBlockTime += w.settings.Resolution
			return w.nextBlock()
		}

		return false, nil
	}

	return false, nil
}

func (w *SlidingTimeWindow) NextBlock() (bool, interface{}) {

	w.mutex.Lock()
	defer w.mutex.Lock()

	return w.nextBlock()
}

func (w *SlidingTimeWindow) nextBlock() (bool, interface{}) {

	if !w.canEmit {
		if w.currentBlock == w.numBlocks-1 {
			w.canEmit = true
		}
	}

	w.numSamples = 0
	w.currentBlock++

	if w.canEmit {

		// aggregate and emit
		val := w.aggFunc(w.blocks, w.currentBlock, w.maxSamples)

		w.currentBlock = w.currentBlock % w.numBlocks
		w.blocks[w.currentBlock], _ = zero(w.blocks[w.currentBlock])
		return true, val
	}

	return false, nil
}

///////////////////
// utils

func zero(a interface{}) (interface{}, error) {
	switch x := a.(type) {
	case int:
		return 0, nil
	case float64:
		return 0.0, nil
	case []int:
		for idx := range x {
			x[idx] = 0
		}
		return x, nil
	case []float64:
		for idx := range x {
			x[idx] = 0.0
		}
		return x, nil
	}

	return nil, fmt.Errorf("unsupported type")
}

func getTimeMillis() int {
	now := time.Now()
	nano := now.Nanosecond()
	return nano / 1000000
}
