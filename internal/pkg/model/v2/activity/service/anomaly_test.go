package service

import (
	"encoding/json"
	"math"
	"math/rand"
	"testing"
)

var complexityTests = []string{`{
 "alfa": [
  {"alfa": "1"},
	{"bravo": "2"}
 ],
 "bravo": [
  {"alfa": "3"},
	{"bravo": "4"}
 ]
}`, `{
 "a": [
  {"a": "aa"},
	{"b": "bb"}
 ],
 "b": [
  {"a": "aa"},
	{"b": "bb"}
 ]
}`}

func generateRandomJSON(rnd *rand.Rand) map[string]interface{} {
	sample := func(stddev float64) int {
		return int(math.Abs(rnd.NormFloat64()) * stddev)
	}
	sampleCount := func() int {
		return sample(1) + 1
	}
	sampleName := func() string {
		const symbols = 'z' - 'a'
		s := sample(8)
		if s > symbols {
			s = symbols
		}
		return string('a' + s)
	}
	sampleValue := func() string {
		value := sampleName()
		return value + value
	}
	sampleDepth := func() int {
		return sample(3)
	}
	var generate func(hash map[string]interface{}, depth int)
	generate = func(hash map[string]interface{}, depth int) {
		count := sampleCount()
		if depth > sampleDepth() {
			for i := 0; i < count; i++ {
				hash[sampleName()] = sampleValue()
			}
			return
		}
		for i := 0; i < count; i++ {
			array := make([]interface{}, sampleCount())
			for j := range array {
				sub := make(map[string]interface{})
				generate(sub, depth+1)
				array[j] = sub
			}
			hash[sampleName()] = array
		}
	}
	object := make(map[string]interface{})
	generate(object, 0)
	return object
}

func TestComplexity(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	complexity := NewComplexity(2)
	for i := 0; i < 1024; i++ {
		data, err := json.Marshal(generateRandomJSON(rnd))
		if err != nil {
			t.Fatal(err)
		}
		complexity.Complexity(data)
	}
	a, _ := complexity.Complexity([]byte(complexityTests[0]))
	b, _ := complexity.Complexity([]byte(complexityTests[1]))
	if a < b {
		t.Fatal("complexity sanity check failed", a, b)
	}
}
