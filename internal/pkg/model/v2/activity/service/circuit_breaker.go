package service

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	// CircuitBreakerModeA triggers the circuit breaker when there are contiguous errors
	CircuitBreakerModeA = "a"
	// CircuitBreakerModeB triggers the circuit breaker when there are errors over time
	CircuitBreakerModeB = "b"
	// CircuitBreakerModeC triggers the circuit breaker when there are contiguous errors over time
	CircuitBreakerModeC = "c"
	// CircuitBreakerModeD is a probabilistic smart circuit breaker
	CircuitBreakerModeD = "d"
	// CircuitBreakerFailure is a failure
	CircuitBreakerFailure = -1.0
	// CircuitBreakerUnknown is an onknown status
	CircuitBreakerUnknown = 0.0
	// CircuitBreakerSuccess isa success
	CircuitBreakerSuccess = 1.0
)

// ErrorCircuitBreakerTripped happens when the circuit breaker has tripped
var ErrorCircuitBreakerTripped = errors.New("circuit breaker tripped")

var now = time.Now

// CircuitBreaker is a circuit breaker service
type CircuitBreaker struct {
	operation, context, mode string
	threshold                int
	period, timeout          time.Duration
	Tripped                  bool `json:"tripped"`
}

// InitializeCircuitBreaker creates a circuit breaker service
func InitializeCircuitBreaker(settings map[string]interface{}) (service *CircuitBreaker, err error) {
	circuit := &CircuitBreaker{
		mode:      CircuitBreakerModeA,
		threshold: 5,
		period:    60 * time.Second,
		timeout:   60 * time.Second,
	}
	circuit.UpdateRequest(settings)
	return circuit, nil
}

// Record is a record of a request
type Record struct {
	Weight float64
	Stamp  time.Time
}

// CircuitBreakerContext is a circuit breaker context
type CircuitBreakerContext struct {
	counter   int
	processed uint64
	timeout   time.Time
	index     int
	buffer    []Record
	tripped   bool
	sync.RWMutex
}

// Trip trips the circuit breaker
func (c *CircuitBreakerContext) Trip(now time.Time, timeout time.Duration) {
	c.timeout = now.Add(timeout)
	c.counter = 0
	c.tripped = true
}

func (c *CircuitBreakerContext) AddRecord(weight float64, now time.Time) {
	c.processed++
	c.buffer[c.index].Weight = weight
	c.buffer[c.index].Stamp = now
	c.index = (c.index + 1) % len(c.buffer)
}

// Probability computes the probability for mode d
func (c *CircuitBreakerContext) Probability(now time.Time) float64 {
	records, factor, sum := c.buffer, 0.0, 0.0
	max := float64(now.Sub(records[c.index].Stamp))
	for _, record := range records {
		a := math.Exp(-float64(now.Sub(record.Stamp)) / max)
		factor += a
		sum += record.Weight * a
	}
	sum /= factor
	return 1 / (1 + math.Exp(8*sum))
}

// CircuitBreakerContexts holds a bunch of circuit breaker contexts
type CircuitBreakerContexts struct {
	contexts map[string]*CircuitBreakerContext
	sync.RWMutex
}

var circuitBreakerContexts = CircuitBreakerContexts{
	contexts: make(map[string]*CircuitBreakerContext),
}

// GetContext gets a circuit breaker context
func (c *CircuitBreakerContexts) GetContext(context string, threshold int) *CircuitBreakerContext {
	context = fmt.Sprintf("%s-%d", context, threshold)
	c.RLock()
	cbContext := c.contexts[context]
	c.RUnlock()

	if cbContext != nil {
		return cbContext
	}

	buffer := make([]Record, threshold)
	cbContext = &CircuitBreakerContext{
		buffer: buffer,
	}
	for i := range buffer {
		buffer[i].Weight = CircuitBreakerSuccess
	}

	c.Lock()
	c.contexts[context] = cbContext
	c.Unlock()

	return cbContext
}

// Execute executes the circuit breaker service
func (c *CircuitBreaker) Execute() (err error) {
	if c.context == "" {
		return errors.New("invalid context")
	}
	if c.threshold <= 0 {
		return errors.New("invalid threshold")
	}

	context, now := circuitBreakerContexts.GetContext(c.context, c.threshold), now()
	switch c.operation {
	case "counter":
		context.Lock()
		if context.timeout.Sub(now) > 0 {
			context.Unlock()
			break
		}
		context.counter++
		context.AddRecord(CircuitBreakerFailure, now)
		if context.tripped {
			context.Trip(now, c.timeout)
			context.Unlock()
			break
		}
		switch c.mode {
		case CircuitBreakerModeA:
			if context.counter >= c.threshold {
				context.Trip(now, c.timeout)
			}
		case CircuitBreakerModeB:
			if context.processed < uint64(c.threshold) {
				break
			}
			if now.Sub(context.buffer[context.index].Stamp) < c.period {
				context.Trip(now, c.timeout)
			}
		case CircuitBreakerModeC:
			if context.processed < uint64(c.threshold) {
				break
			}
			if context.counter >= c.threshold &&
				now.Sub(context.buffer[context.index].Stamp) < c.period {
				context.Trip(now, c.timeout)
			}
		}
		context.Unlock()
	case "reset":
		context.Lock()
		switch c.mode {
		case CircuitBreakerModeA, CircuitBreakerModeB, CircuitBreakerModeC:
			if context.timeout.Sub(now) <= 0 {
				context.counter = 0
				context.tripped = false
			}
		case CircuitBreakerModeD:
			context.AddRecord(CircuitBreakerSuccess, now)
		}
		context.Unlock()
	default:
		switch c.mode {
		case CircuitBreakerModeA, CircuitBreakerModeB, CircuitBreakerModeC:
			context.RLock()
			timeout := context.timeout
			context.RUnlock()
			if timeout.Sub(now) > 0 {
				c.Tripped = true
				return ErrorCircuitBreakerTripped
			}
		case CircuitBreakerModeD:
			context.RLock()
			p := context.Probability(now)
			context.RUnlock()
			if rand.Float64()*1000 < math.Floor(p*1000) {
				context.Lock()
				context.AddRecord(CircuitBreakerUnknown, now)
				context.Unlock()
				c.Tripped = true
				return ErrorCircuitBreakerTripped
			}
		}
	}
	return nil
}

// UpdateRequest updates the circuit breaker service
func (c *CircuitBreaker) UpdateRequest(values map[string]interface{}) (err error) {
	for k, v := range values {
		switch k {
		case "mode":
			mode, ok := v.(string)
			if !ok {
				return errors.New("mode is not a string")
			}
			switch mode {
			case CircuitBreakerModeA:
			case CircuitBreakerModeB:
			case CircuitBreakerModeC:
			case CircuitBreakerModeD:
			default:
				return errors.New("invalid mode")
			}
			c.mode = mode
		case "operation":
			operation, ok := v.(string)
			if !ok {
				return errors.New("operation is not a string")
			}
			c.operation = operation
		case "context":
			context, ok := v.(string)
			if !ok {
				return errors.New("context is not a string")
			}
			c.context = context
		case "threshold":
			threshold, ok := v.(float64)
			if !ok {
				return errors.New("threshold is not a number")
			}
			c.threshold = int(threshold)
		case "timeout":
			timeout, ok := v.(float64)
			if !ok {
				return errors.New("timeout is not a number")
			}
			c.timeout = time.Duration(timeout) * time.Second
		case "period":
			period, ok := v.(float64)
			if !ok {
				return errors.New("period is not a number")
			}
			c.period = time.Duration(period) * time.Second
		}
	}
	return nil
}
