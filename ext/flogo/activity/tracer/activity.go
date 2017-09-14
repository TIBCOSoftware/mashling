package tracer

import (
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	opentracing "github.com/opentracing/opentracing-go"
	ctx "golang.org/x/net/context"
)

var (
	ErrorTracingSpanRequired = errors.New("either tracing or span is required")
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	tracing, span := context.GetInput("tracing"), context.GetInput("span")
	if span == nil && tracing == nil {
		return false, ErrorTracingSpanRequired
	}

	if span != nil {
		span.(opentracing.Span).Finish()
	}

	if tracing != nil {
		if span := opentracing.SpanFromContext(tracing.(ctx.Context)); span != nil {
			span = opentracing.StartSpan(
				context.TaskName(),
				opentracing.ChildOf(span.Context()))
			context.SetOutput("span", span)
			context.SetOutput("tracing", opentracing.ContextWithSpan(ctx.Background(), span))
		}
	}

	return true, nil
}
