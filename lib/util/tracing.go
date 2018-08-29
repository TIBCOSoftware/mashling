/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"errors"
	"fmt"

	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"sourcegraph.com/sourcegraph/appdash"
	appdashtracing "sourcegraph.com/sourcegraph/appdash/opentracing"
)

const (
	tracerNoOP            = "noop"
	tracerZipKin          = "zipkin"
	tracerAPPDash         = "appdash"
	tracerLightStep       = "lightstep"
	settingTracer         = "tracer"
	settingTracerEndpoint = "tracerEndpoint"
	settingTracerToken    = "tracerToken"
	settingTracerDebug    = "tracerDebug"
	settingTracerSameSpan = "tracerSameSpan"
	settingTracerID128Bit = "tracerID128Bit"
)

var (
	// ErrorTracerEndpointRequired a tracer endpoint is required
	ErrorTracerEndpointRequired = errors.New("tracer endpoint required")
	// ErrorInvalidTracer invalid tracer
	ErrorInvalidTracer = errors.New("invalid tracer")
	// ErrorTracerTokenRequired tracer token required
	ErrorTracerTokenRequired = errors.New("tracer token required")
)

// ConfigureTracer configures the distributed tracer
func ConfigureTracer(settings map[string]interface{}, address, name string) error {
	tracer := tracerNoOP
	if setting, ok := settings[settingTracer]; ok {
		tracer = setting.(string)
	}
	tracerEndpoint := ""
	if setting, ok := settings[settingTracerEndpoint]; ok {
		tracerEndpoint = setting.(string)
	}
	tracerToken := ""
	if setting, ok := settings[settingTracerToken]; ok {
		tracerToken = setting.(string)
	}
	tracerDebug := false
	if setting, ok := settings[settingTracerDebug]; ok {
		tracerDebug = setting.(bool)
	}
	tracerSameSpan := false
	if setting, ok := settings[settingTracerSameSpan]; ok {
		tracerSameSpan = setting.(bool)
	}
	tracerID128Bit := true
	if setting, ok := settings[settingTracerID128Bit]; ok {
		tracerID128Bit = setting.(bool)
	}

	switch tracer {
	case tracerNoOP:
		opentracing.SetGlobalTracer(&opentracing.NoopTracer{})
	case tracerZipKin:
		if tracerEndpoint == "" {
			return ErrorTracerEndpointRequired
		}

		collector, err := zipkin.NewHTTPCollector(tracerEndpoint)
		if err != nil {
			return fmt.Errorf("unable to create Zipkin HTTP collector: %+v\n", err)
		}

		recorder := zipkin.NewRecorder(collector, tracerDebug, address, name)

		tracer, err := zipkin.NewTracer(
			recorder,
			zipkin.ClientServerSameSpan(tracerSameSpan),
			zipkin.TraceID128Bit(tracerID128Bit),
		)
		if err != nil {
			return fmt.Errorf("unable to create Zipkin tracer: %+v\n", err)
		}

		opentracing.SetGlobalTracer(tracer)
	case tracerAPPDash:
		if tracerEndpoint == "" {
			return ErrorTracerEndpointRequired
		}

		collector := appdash.NewRemoteCollector(tracerEndpoint)
		chunkedCollector := appdash.NewChunkedCollector(collector)
		tracer := appdashtracing.NewTracer(chunkedCollector)
		opentracing.SetGlobalTracer(tracer)
	case tracerLightStep:
		if tracerToken == "" {
			return ErrorTracerTokenRequired
		}

		lightstepTracer := lightstep.NewTracer(lightstep.Options{
			AccessToken: tracerToken,
		})

		opentracing.SetGlobalTracer(lightstepTracer)
	default:
		return ErrorInvalidTracer
	}

	return nil
}
