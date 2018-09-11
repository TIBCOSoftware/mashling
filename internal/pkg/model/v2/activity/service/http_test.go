package service

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	"github.com/TIBCOSoftware/mashling/lib/util"
)

func TestHTTP(t *testing.T) {
	const (
		jsonPayload = "{\n \"name\": \"sally\"\n}"
		xmlPayload  = `<pet>sally</pet>`
	)
	server := &http.Server{Addr: ":8181"}
	http.HandleFunc("/pet/json", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Add("Content-Type", "application/json")
			io.WriteString(w, jsonPayload)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != jsonPayload {
				t.Fatalf("payload is %s and should be %s", string(data), jsonPayload)
			}
		}
	})
	http.HandleFunc("/pet/xml", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Add("Content-Type", "text/xml")
			io.WriteString(w, xmlPayload)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != xmlPayload {
				t.Fatalf("payload is %s and should be %s", string(data), jsonPayload)
			}
		}
	})
	http.HandleFunc("/path/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/path/1?a=b" {
			t.Fatalf("url is %s but should be /path/1?a=b", r.URL.String())
		}
		if r.Header.Get("TEST") != "TEST" {
			t.Fatal("there should be a TEST header")
		}
	})
	done := make(chan bool, 1)
	go func() {
		server.ListenAndServe()
		done <- true
	}()
	_, err := http.Get("http://localhost:8181/pet/json")
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
		Type:     "http",
		Settings: map[string]interface{}{},
	}

	test := func(values map[string]interface{}, mime, input, output string) {
		instance, err := Initialize(service)
		if err != nil {
			t.Fatal(err)
		}
		if input != "" {
			var body interface{}
			err = util.Unmarshal(mime, []byte(input), &body)
			if err != nil {
				t.Fatal(err)
			}
			values["body"] = body
		}
		err = instance.UpdateRequest(values)
		if err != nil {
			t.Fatal(err)
		}
		err = instance.Execute()
		if err != nil {
			t.Fatal(err)
		}
		data, err := util.Marshal(instance.(*HTTP).Response.Body)
		if err != nil {
			panic(err)
		}
		if string(data) != output {
			t.Fatalf("payload is %s and should be %s", string(data), output)
		}
	}

	// json GET test
	test(map[string]interface{}{
		"method": methodGET,
		"url":    "http://localhost:8181/pet/json",
	}, "", "", jsonPayload)

	// xml GET test
	test(map[string]interface{}{
		"method": methodGET,
		"url":    "http://localhost:8181/pet/xml",
	}, "", "", xmlPayload)

	// json POST test
	test(map[string]interface{}{
		"method": methodPOST,
		"url":    "http://localhost:8181/pet/json",
	}, "application/json", jsonPayload, "")

	// xml POST test
	test(map[string]interface{}{
		"method": methodPOST,
		"url":    "http://localhost:8181/pet/xml",
	}, "text/xml", xmlPayload, "")

	// json PUT test
	test(map[string]interface{}{
		"method": methodPUT,
		"url":    "http://localhost:8181/pet/json",
	}, "application/json", jsonPayload, "")

	// xml PUT test
	test(map[string]interface{}{
		"method": methodPUT,
		"url":    "http://localhost:8181/pet/xml",
	}, "text/xml", xmlPayload, "")

	// json PATCH test
	test(map[string]interface{}{
		"method": methodPATCH,
		"url":    "http://localhost:8181/pet/json",
	}, "application/json", jsonPayload, "")

	// xml PATCH test
	test(map[string]interface{}{
		"method": methodPATCH,
		"url":    "http://localhost:8181/pet/xml",
	}, "text/xml", xmlPayload, "")

	// path test
	test(map[string]interface{}{
		"method":     methodGET,
		"url":        "http://localhost:8181/",
		"path":       "path/:id",
		"headers":    map[string]interface{}{"TEST": "TEST"},
		"query":      map[string]string{"a": "b"},
		"pathParams": map[string]interface{}{"id": 1},
	}, "", "", "")

	// DELETE test
	test(map[string]interface{}{
		"method":     methodDELETE,
		"url":        "http://localhost:8181/",
		"path":       "path/:id",
		"headers":    map[string]interface{}{"TEST": "TEST"},
		"query":      map[string]string{"a": "b"},
		"pathParams": map[string]interface{}{"id": 1},
	}, "", "", "")
}
