package service

import (
	"testing"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
)

func TestJS(t *testing.T) {
	service := types.Service{
		Type:     "js",
		Settings: map[string]interface{}{"script": "result.sum = parameters.a + parameters.b"},
	}
	instance, err := Initialize(service)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{
		"parameters": map[string]interface{}{"a": 1.0, "b": 2.0},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}
	if instance.(*JS).Response.Result["sum"].(float64) != 3.0 {
		t.Fatal("sum should be 3.0")
	}
}
