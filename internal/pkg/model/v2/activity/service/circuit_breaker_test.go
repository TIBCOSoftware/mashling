package service

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

func TestCircuitBreakerModeA(t *testing.T) {
	rand.Seed(1)
	clock := time.Unix(1533930608, 0)
	now = func() time.Time {
		now := clock
		clock = clock.Add(time.Duration(rand.Intn(2)+1) * time.Second)
		return now
	}
	defer func() {
		now = time.Now
	}()

	service := types.Service{
		Type: "circuitBreaker",
		Settings: map[string]interface{}{
			"context": "testA",
		},
	}
	execute := func(values map[string]interface{}, should error) {
		breaker, err := Initialize(service)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.UpdateRequest(values)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.Execute()
		if err != should {
			t.Fatalf("error should be %v but is %v", should, err)
		}
	}

	for i := 0; i < 4; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	execute(nil, nil)
	execute(map[string]interface{}{"operation": "reset"}, nil)

	for i := 0; i < 5; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	execute(nil, ErrorCircuitBreakerTripped)

	clock = clock.Add(60 * time.Second)

	execute(nil, nil)
	execute(map[string]interface{}{"operation": "counter"}, nil)

	execute(nil, ErrorCircuitBreakerTripped)

	clock = clock.Add(60 * time.Second)

	execute(nil, nil)
}

func TestCircuitBreakerModeB(t *testing.T) {
	rand.Seed(1)
	clock := time.Unix(1533930608, 0)
	now = func() time.Time {
		now := clock
		clock = clock.Add(time.Duration(rand.Intn(2)+1) * time.Second)
		return now
	}
	defer func() {
		now = time.Now
	}()

	service := types.Service{
		Type: "circuitBreaker",
		Settings: map[string]interface{}{
			"mode":    CircuitBreakerModeB,
			"context": "testB",
		},
	}
	execute := func(values map[string]interface{}, should error) {
		breaker, err := Initialize(service)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.UpdateRequest(values)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.Execute()
		if err != should {
			t.Fatalf("error should be %v but is %v", should, err)
		}
	}

	for i := 0; i < 4; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	clock = clock.Add(60 * time.Second)

	for i := 0; i < 5; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	execute(nil, ErrorCircuitBreakerTripped)

	clock = clock.Add(60 * time.Second)

	execute(nil, nil)
	execute(map[string]interface{}{"operation": "counter"}, nil)

	execute(nil, ErrorCircuitBreakerTripped)

	clock = clock.Add(60 * time.Second)

	execute(nil, nil)
}

func TestCircuitBreakerModeC(t *testing.T) {
	rand.Seed(1)
	clock := time.Unix(1533930608, 0)
	now = func() time.Time {
		now := clock
		clock = clock.Add(time.Duration(rand.Intn(2)+1) * time.Second)
		return now
	}
	defer func() {
		now = time.Now
	}()

	service := types.Service{
		Type: "circuitBreaker",
		Settings: map[string]interface{}{
			"mode":    CircuitBreakerModeC,
			"context": "testC",
		},
	}
	execute := func(values map[string]interface{}, should error) {
		breaker, err := Initialize(service)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.UpdateRequest(values)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.Execute()
		if err != should {
			t.Fatalf("error should be %v but is %v", should, err)
		}
	}

	for i := 0; i < 4; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	clock = clock.Add(60 * time.Second)

	for i := 0; i < 4; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	execute(nil, nil)
	execute(map[string]interface{}{"operation": "reset"}, nil)

	for i := 0; i < 5; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "counter"}, nil)
	}

	execute(nil, ErrorCircuitBreakerTripped)

	clock = clock.Add(60 * time.Second)

	execute(nil, nil)
	execute(map[string]interface{}{"operation": "counter"}, nil)

	execute(nil, ErrorCircuitBreakerTripped)

	clock = clock.Add(60 * time.Second)

	execute(nil, nil)
}

func TestCircuitBreakerModeD(t *testing.T) {
	rand.Seed(1)
	clock := time.Unix(1533930608, 0)
	now = func() time.Time {
		now := clock
		clock = clock.Add(time.Duration(rand.Intn(2)+1) * time.Second)
		return now
	}
	defer func() {
		now = time.Now
	}()

	service := types.Service{
		Type: "circuitBreaker",
		Settings: map[string]interface{}{
			"mode":    CircuitBreakerModeD,
			"context": "testD",
		},
	}
	execute := func(values map[string]interface{}, should error) error {
		breaker, err := Initialize(service)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.UpdateRequest(values)
		if err != nil {
			t.Fatal(err)
		}
		err = breaker.Execute()
		if err != should {
			t.Fatalf("error should be %v but is %v", should, err)
		}
		return err
	}

	for i := 0; i < 100; i++ {
		execute(nil, nil)
		execute(map[string]interface{}{"operation": "reset"}, nil)
	}
	p := circuitBreakerContexts.GetContext("testD", 5).Probability(now())
	if math.Floor(p*100) != 0.0 {
		t.Fatalf("probability should be zero but is %v", math.Floor(p*100))
	}

	type Test struct {
		a, b error
	}
	tests := []Test{
		{nil, nil},
		{nil, nil},
		{ErrorCircuitBreakerTripped, nil},
		{ErrorCircuitBreakerTripped, nil},
		{nil, nil},
		{ErrorCircuitBreakerTripped, nil},
		{ErrorCircuitBreakerTripped, nil},
		{ErrorCircuitBreakerTripped, nil},
	}
	for _, test := range tests {
		err := execute(nil, test.a)
		if err != nil {
			continue
		}
		execute(map[string]interface{}{"operation": "counter"}, test.b)
	}

	tests = []Test{
		{nil, nil},
		{nil, nil},
		{nil, nil},
		{nil, nil},
		{nil, nil},
	}
	for _, test := range tests {
		err := execute(nil, test.a)
		if err != nil {
			continue
		}
		execute(map[string]interface{}{"operation": "reset"}, test.b)
	}
	p = circuitBreakerContexts.GetContext("testD", 5).Probability(now())
	if math.Floor(p*100) != 0.0 {
		t.Fatalf("probability should be zero but is %v", math.Floor(p*100))
	}
}
