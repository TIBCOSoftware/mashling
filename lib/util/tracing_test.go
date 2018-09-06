/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"testing"
)

func TestTracer(t *testing.T) {
	tracer, settings := Tracer{}, make(map[string]interface{})
	err := tracer.ConfigureTracer(settings, "localhost", "test")
	if err != nil {
		t.Fatal(err)
	}
	err = tracer.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestTracerInvalidTracer(t *testing.T) {
	tracer, settings := Tracer{}, make(map[string]interface{})
	settings[settingTracer] = "invalid"
	err := tracer.ConfigureTracer(settings, "localhost", "test")
	if err != ErrorInvalidTracer {
		t.Fatal("error should be ErrorInvalidTracer")
	}
}

func TestTracerNoEndpoint(t *testing.T) {
	for _, name := range []string{tracerZipKin, tracerAPPDash, tracerJaeger} {
		t.Log("tracer", name)
		tracer, settings := Tracer{}, make(map[string]interface{})
		settings[settingTracer] = name
		err := tracer.ConfigureTracer(settings, "localhost", "test")
		if err != ErrorTracerEndpointRequired {
			t.Fatal("error should be ErrorTracerEndpointRequired")
		}
	}
}

func TestTracerNoToken(t *testing.T) {
	tracer, settings := Tracer{}, make(map[string]interface{})
	settings[settingTracer] = tracerLightStep
	err := tracer.ConfigureTracer(settings, "localhost", "test")
	if err != ErrorTracerTokenRequired {
		t.Fatal("error should be ErrorTracerTokenRequired")
	}
}
