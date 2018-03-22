package trigger

import (
	"context"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type key int

var ctxDataKey key

type ContextData struct {
	Attrs      []*data.Attribute
	HandlerCfg *HandlerConfig
}

// NewContext returns a new Context that carries the trigger data.
func NewContext(parentCtx context.Context, attrs []*data.Attribute) context.Context {
	ctxData := &ContextData{Attrs: attrs}
	return context.WithValue(parentCtx, ctxDataKey, ctxData)
}

func NewInitialContext(attrs []*data.Attribute, config *HandlerConfig) context.Context {
	return context.WithValue(context.Background(), ctxDataKey, &ContextData{Attrs: attrs, HandlerCfg: config})
}

// NewContext returns a new Context that carries the trigger data.
func NewContextWithData(parentCtx context.Context, contextData *ContextData) context.Context {
	return context.WithValue(parentCtx, ctxDataKey, contextData)
}

func ExtractContextData(ctx context.Context) (*ContextData, bool) {
	if ctx == nil {
		return nil, false
	}
	ctxDataVal := ctx.Value(ctxDataKey)
	if ctxDataVal == nil {
		return nil, false
	}
	ctxData, ok := ctxDataVal.(*ContextData)
	return ctxData, ok
}
