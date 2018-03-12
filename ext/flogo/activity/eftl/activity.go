/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */

package eftl

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling/lib/eftl"
	"github.com/TIBCOSoftware/mashling/lib/util"

	opentracing "github.com/opentracing/opentracing-go"
	ctx "golang.org/x/net/context"
)

const (
	ivContent  = "content"
	ivDest     = "dest"
	ivURL      = "url"
	ivID       = "id"
	ivUser     = "user"
	ivPassword = "password"
	ivCA       = "ca"
	ivTracing  = "tracing"

	ovTracing = "tracing"
)

var log = logger.GetLogger("activity-tibco-rest")

// Activity is a EFTL client
type Activity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new EFTL activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &Activity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	var span opentracing.Span
	if tracing := context.GetInput(ivTracing); tracing != nil {
		span = opentracing.SpanFromContext(tracing.(ctx.Context))
	}

	if span != nil {
		span = opentracing.StartSpan(
			context.TaskName(),
			opentracing.ChildOf(span.Context()))
		context.SetOutput(ovTracing, opentracing.ContextWithSpan(ctx.Background(), span))
		defer span.Finish()
	}

	setTag := func(key string, value interface{}) {
		if span != nil {
			span.SetTag(key, value)
		}
	}

	logError := func(format string, a ...interface{}) {
		str := fmt.Sprintf(format, a...)
		setTag("error", str)
		log.Error(str)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if ca, ok := context.GetInput(ivCA).(string); ok && ca != "" {
		certificate, err := ioutil.ReadFile(ca)
		if err != nil {
			logError("can't open certificate", err)
			return false, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(certificate)
		tlsConfig = &tls.Config{
			RootCAs: pool,
		}
	}

	options := &eftl.Options{
		TLSConfig: tlsConfig,
	}
	if id, ok := context.GetInput(ivID).(string); ok {
		options.ClientID = id
	}
	if user, ok := context.GetInput(ivUser).(string); ok {
		options.Username = user
	}
	if password, ok := context.GetInput(ivPassword).(string); ok {
		options.Password = password
	}

	url := context.GetInput(ivURL).(string)
	errorsChannel := make(chan error, 1)
	connection, err := eftl.Connect(url, options, errorsChannel)
	if err != nil {
		logError("connection failed: %s", err)
		return false, err
	}
	defer connection.Disconnect()

	content := context.GetInput(ivContent)
	data, err := util.Marshal(content)
	if err != nil {
		logError("failed to marshal: %s", err)
		return false, err
	}

	if dest, ok := context.GetInput(ivDest).(string); ok {
		err = connection.Publish(eftl.Message{
			"_dest":   dest,
			"content": data,
		})
		if err != nil {
			logError("failed to publish: %s", err)
			return false, err
		}
	}

	return true, nil
}
