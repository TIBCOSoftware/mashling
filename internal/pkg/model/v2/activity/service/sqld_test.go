package service

import (
	"testing"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

func TestSQLD(t *testing.T) {
	service := types.Service{
		Type:     "sqld",
		Settings: map[string]interface{}{},
	}

	test := func(a string, attack bool) {
		instance, err := Initialize(service)
		if err != nil {
			t.Fatal(err)
		}
		var payload interface{} = map[string]interface{}{
			"content": map[string]interface{}{
				"test": a,
			},
		}
		err = instance.UpdateRequest(map[string]interface{}{
			"payload": &payload,
		})
		if err != nil {
			t.Fatal(err)
		}
		err = instance.Execute()
		if err != nil {
			t.Fatal(err)
		}
		if attack {
			if instance.(*SQLD).Attack < 50 {
				t.Fatal("should be an attack", a, instance.(*SQLD).Attack)
			}
			if instance.(*SQLD).AttackValues["content"].(map[string]interface{})["test"].(float64) < 50 {
				t.Fatal("should be an attack", a, instance.(*SQLD).Attack)
			}
		} else {
			if instance.(*SQLD).Attack > 50 {
				t.Fatal("should not be an attack", a, instance.(*SQLD).Attack)
			}
			if instance.(*SQLD).AttackValues["content"].(map[string]interface{})["test"].(float64) > 50 {
				t.Fatal("should not be an attack", a, instance.(*SQLD).Attack)
			}
		}
	}
	test("test or 1337=1337 --\"", true)
	test(" or 1=1 ", true)
	test("/**/or/**/1337=1337", true)
	test("abc123", false)
	test("abc123 123abc", false)
	test("123", false)
	test("abcorabc", false)
	test("available", false)
	test("orcat1", false)
	test("cat1or", false)
	test("cat1orcat1", false)
}
