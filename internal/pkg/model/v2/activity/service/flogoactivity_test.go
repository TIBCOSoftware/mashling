package service

import (
	"io"
	"net/http"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/activity/rest"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/mashling/internal/app/gateway/flogo/registry/activities"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

func TestFlogoActivity(t *testing.T) {
	act, err := activities.Asset("vendor/github.com/TIBCOSoftware/flogo-contrib/activity/rest/activity.json")
	if err != nil {
		t.Fatal(err)
	}
	actmd := activity.NewMetadata(string(act))
	activity.Register(rest.NewActivity(actmd))

	const (
		jsonPayload = "{\n \"name\": \"sally\"\n}"
	)
	server := &http.Server{Addr: ":8181"}
	http.HandleFunc("/pet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, jsonPayload)
	})
	done := make(chan bool, 1)
	go func() {
		server.ListenAndServe()
		done <- true
	}()
	_, err = http.Get("http://localhost:8181/pet/json")
	for err != nil {
		_, err = http.Get("http://localhost:8181/pet/json")
	}
	defer func() {
		err := server.Shutdown(nil)
		if err != nil {
			t.Fatal(err)
		}
		<-done
	}()

	service := types.Service{
		Type: "flogoActivity",
		Settings: map[string]interface{}{
			"ref": "github.com/TIBCOSoftware/flogo-contrib/activity/rest",
			"inputs": map[string]interface{}{
				"uri":    "http://localhost:8181/pet",
				"method": "GET",
			},
		},
	}

	instance, err := Initialize(service)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{
		"method": "GET",
		"url":    "http://localhost:8181/pet/json",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}
	data, err := util.Marshal(instance.(*FlogoActivity).Response.Outputs["result"])
	if err != nil {
		panic(err)
	}
	if string(data) != jsonPayload {
		t.Fatalf("payload is %s and should be %s", string(data), jsonPayload)
	}
}
