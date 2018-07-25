package service

import (
	"errors"
	"sync"
	"time"
)

// CircuitBreaker is a circuit breaker service
type CircuitBreaker struct {
	operation, context string
	threshold, timeout int
	Tripped            bool `json:"tripped"`
}

// InitializeCircuitBreaker creates a circuit breaker service
func InitializeCircuitBreaker(settings map[string]interface{}) (service *CircuitBreaker, err error) {
	circuit := &CircuitBreaker{
		threshold: 5,
		timeout:   60,
	}
	circuit.UpdateRequest(settings)
	return circuit, nil
}

// CircuitBreakerContext is a circuit breaker context
type CircuitBreakerContext struct {
	counter int
	timeout time.Time
}

// CircuitBreakerContexts holds a bunch of circuit breaker contexts
type CircuitBreakerContexts struct {
	contexts map[string]CircuitBreakerContext
	sync.RWMutex
}

var circuitBreakerContexts = CircuitBreakerContexts{
	contexts: make(map[string]CircuitBreakerContext),
}

// Execute executes the circuit breaker service
func (c *CircuitBreaker) Execute() (err error) {
	if c.context == "" {
		return errors.New("invalid context")
	}

	switch c.operation {
	case "counter":
		circuitBreakerContexts.Lock()
		context := circuitBreakerContexts.contexts[c.context]
		context.counter++
		if context.counter >= c.threshold {
			context.timeout = time.Now().Add(time.Duration(c.timeout) * time.Second)
			context.counter = 0
		}
		circuitBreakerContexts.contexts[c.context] = context
		circuitBreakerContexts.Unlock()
	case "reset":
		circuitBreakerContexts.Lock()
		context := circuitBreakerContexts.contexts[c.context]
		context.counter = 0
		circuitBreakerContexts.contexts[c.context] = context
		circuitBreakerContexts.Unlock()
	default:
		circuitBreakerContexts.RLock()
		context := circuitBreakerContexts.contexts[c.context]
		circuitBreakerContexts.RUnlock()
		if context.timeout.Sub(time.Now()) > 0 {
			c.Tripped = true
			return errors.New("circuit breaker tripped")
		}
	}
	return nil
}

// UpdateRequest updates the circuit breaker service
func (c *CircuitBreaker) UpdateRequest(values map[string]interface{}) (err error) {
	for k, v := range values {
		switch k {
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
			c.timeout = int(timeout)
		}
	}
	return nil
}
