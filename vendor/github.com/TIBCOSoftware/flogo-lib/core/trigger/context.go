package trigger

import (
	"context"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type key int
var handlerKey key

type HandlerInfo struct {
	Name string
}

// NewHandlerContext add the handler info to a new child context
func NewHandlerContext(parentCtx context.Context, config *HandlerConfig) context.Context {
	if config != nil && config.Name != "" {
		return context.WithValue(parentCtx, handlerKey, &HandlerInfo{Name:config.Name})
	}
	return parentCtx
}

// HandlerFromContext returns the handler info stored in the context, if any.
func HandlerFromContext(ctx context.Context) (*HandlerInfo, bool) {
	u, ok := ctx.Value(handlerKey).(*HandlerInfo)
	return u, ok
}

// DEPRECATED
var ctxDataKey key

type ContextData struct {
	Attrs      []*data.Attribute
	HandlerCfg *HandlerConfig
}

// NewContext returns a new Context that carries the trigger data.
// DEPRECATED
func NewContext(parentCtx context.Context, attrs []*data.Attribute) context.Context {
	ctxData := &ContextData{Attrs: attrs}
	return context.WithValue(parentCtx, ctxDataKey, ctxData)
}

// DEPRECATED
func NewInitialContext(attrs []*data.Attribute, config *HandlerConfig) context.Context {
	return context.WithValue(context.Background(), ctxDataKey, &ContextData{Attrs: attrs, HandlerCfg: config})
}

// NewContext returns a new Context that carries the trigger data.
// DEPRECATED
func NewContextWithData(parentCtx context.Context, contextData *ContextData) context.Context {
	return context.WithValue(parentCtx, ctxDataKey, contextData)
}

// DEPRECATED
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
