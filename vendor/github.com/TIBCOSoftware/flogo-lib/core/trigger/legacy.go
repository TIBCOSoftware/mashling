package trigger

import (
	"context"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Deprecated: Temporarily generated for backwards compatibility support
type LegacyRunner struct {
	currentRunner   action.Runner
	triggerMetadata *Metadata
	handlers        map[*HandlerConfig]*Handler
	actToHandlers   map[action.Action]*Handler
}

func NewLegacyRunner(runner action.Runner, metadata *Metadata) *LegacyRunner {
	return &LegacyRunner{currentRunner: runner, triggerMetadata: metadata, handlers: make(map[*HandlerConfig]*Handler), actToHandlers: make(map[action.Action]*Handler)}
}

func (lr *LegacyRunner) RegisterHandler(h *Handler) {
	lr.handlers[h.config] = h
	lr.actToHandlers[h.act] = h
}

func (lr *LegacyRunner) Run(ctx context.Context, act action.Action, uri string, options interface{}) (code int, data interface{}, err error) {

	newOptions := make(map[string]interface{})
	newOptions["deprecated_options"] = options

	results, err := lr.RunAction(ctx, act, newOptions)

	if len(results) != 0 {
		defData, ok := results["data"]
		if ok {
			data = defData.Value()
		}
		defCode, ok := results["code"]
		if ok {
			code = defCode.Value().(int)
		}
	}

	return code, data, err
}

func (lr *LegacyRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {

	trgHandler, trgData := lr.getHandler(ctx, act)
	return trgHandler.Handle(ctx, trgData)
}

func (*LegacyRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	//only called by handler so not needed

	return nil, errors.New("not supported")
}

func (lr *LegacyRunner) getHandler(ctx context.Context, act action.Action) (*Handler, map[string]interface{}) {
	var values map[string]interface{}
	var handler *Handler

	if ctx != nil {
		var exists bool
		ctxData, exists := ExtractContextData(ctx)

		if exists {
			values = attrsToData(ctxData.Attrs)
			if ctxData.HandlerCfg != nil {
				handler = lr.handlers[ctxData.HandlerCfg]
			}

		}
	}

	if handler == nil {
		handler = lr.actToHandlers[act]
	}

	if handler == nil {
		logger.Warn("Unable to determine handler, creating generic one")
		handler = NewHandler(nil, act, lr.triggerMetadata.Output, lr.triggerMetadata.Reply, lr.currentRunner)
	}

	return handler, values
}

func attrsToData(attrs []*data.Attribute) map[string]interface{} {

	if attrs == nil {
		return nil
	}

	values := make(map[string]interface{}, len(attrs))

	for _, attr := range attrs {
		values[attr.Name()] = attr
	}

	return values
}
