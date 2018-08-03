package trigger

import (
	"context"

	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

type Handler struct {
	runner action.Runner
	act    action.Action

	outputMd map[string]*data.Attribute
	replyMd  map[string]*data.Attribute

	config *HandlerConfig

	actionInputMapper  data.Mapper
	actionOutputMapper data.Mapper
}

func NewHandler(config *HandlerConfig, act action.Action, outputMd map[string]*data.Attribute, replyMd map[string]*data.Attribute, runner action.Runner) *Handler {
	handler := &Handler{config: config, act: act, outputMd: outputMd, replyMd: replyMd, runner: runner}

	if config != nil {
		if config.Action.Mappings != nil {
			if len(config.Action.Mappings.Input) > 0 {
				handler.actionInputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: config.Action.Mappings.Input}, nil)
			}
			if len(config.Action.Mappings.Output) > 0 {
				handler.actionOutputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: config.Action.Mappings.Output}, nil)
			}
		} else if config.ActionMappings != nil {
			// temporary for backwards compatibility
			if len(config.ActionMappings.Input) > 0 {
				handler.actionInputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: config.ActionMappings.Input}, nil)
			}
			if len(config.ActionMappings.Output) > 0 {
				handler.actionOutputMapper = mapper.GetFactory().NewMapper(&data.MapperDef{Mappings: config.ActionMappings.Output}, nil)
			}
		}
	}

	return handler
}

func (h *Handler) GetSetting(setting string) (interface{}, bool) {

	if h.config == nil {
		return nil, false
	}

	val, exists := data.GetValueWithResolver(h.config.Settings, setting)

	if !exists {
		val, exists = data.GetValueWithResolver(h.config.parent.Settings, setting)
	}

	return val, exists
}

func (h *Handler) GetOutput() (map[string]interface{}) {

	if h.config == nil {
		return nil
	}
	return h.config.Output
}

func (h *Handler) GetStringSetting(setting string) string {
	val, exists := h.GetSetting(setting)

	if !exists {
		return ""
	}

	strVal, err := data.CoerceToString(val)

	if err != nil {
		return ""
	}

	return strVal
}

func (h *Handler) Handle(ctx context.Context, triggerData map[string]interface{}) (map[string]*data.Attribute, error) {

	inputs, err := h.generateInputs(triggerData)

	if err != nil {
		return nil, err
	}

	results, err := h.runner.Execute(ctx, h.act, inputs)

	if err != nil {
		return nil, err
	}

	retValue, err := h.generateOutputs(results)

	return retValue, err
}

func (h *Handler) dataToAttrs(triggerData map[string]interface{}) ([]*data.Attribute, error) {
	attrs := make([]*data.Attribute, 0, len(h.outputMd))

	for k, a := range h.outputMd {
		v, _ := triggerData[k]

		switch t := v.(type) {
		case *data.Attribute:
			attr, err := data.NewAttribute(t.Name(), t.Type(), t.Value())
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, attr)
		default:
			attr, err := data.NewAttribute(a.Name(), a.Type(), v)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, attr)
		}
	}

	return attrs, nil
}

func (h *Handler) generateInputs(triggerData map[string]interface{}) (map[string]*data.Attribute, error) {

	triggerAttrs, err := h.dataToAttrs(triggerData)

	if err != nil {
		logger.Errorf("Failed parsing attrs: %s, Error: %s", triggerData, err)
		return nil, err
	}

	var inputs map[string]*data.Attribute

	//if h.act.IOMetadata() == nil {
	//	inputs = make(map[string]*data.Attribute, len(triggerAttrs))
	//
	//	for _, attr := range triggerAttrs {
	//		inputs[attr.Name()] = attr
	//	}
	//
	//	return inputs, nil
	//}
	//inputMetadata := h.act.IOMetadata().Input

	logger.Debugf("iomd %#v", h.act.IOMetadata())

	//todo verify this behavior
	if h.actionInputMapper != nil && h.act.IOMetadata() != nil && h.act.IOMetadata().Input != nil {

		inputMetadata := h.act.IOMetadata().Input

		inScope := data.NewSimpleScope(triggerAttrs, nil)
		outScope := data.NewFixedScope(inputMetadata)

		err := h.actionInputMapper.Apply(inScope, outScope)
		if err != nil {
			return nil, err
		}

		attrs := outScope.GetAttrs()

		inputs = make(map[string]*data.Attribute, len(inputMetadata))

		for _, attr := range attrs {
			inputs[attr.Name()] = attr
		}
	} else {
		// for backwards compatibility make trigger outputs map directly to action inputs

		logger.Debug("No mapping specified, adding trigger outputs as inputs to action")

		inputs = make(map[string]*data.Attribute, len(triggerAttrs))

		for _, attr := range triggerAttrs {

			logger.Debugf(" Attr: %s, Type: %s, Value: %v", attr.Name(), attr.Type().String(), attr.Value())
			//inputs = append(inputs, data.NewAttribute( attr.Name, attr.Type, attr.Value))
			inputs[attr.Name()] = attr

			attrName := "_T." + attr.Name()
			inputs[attrName] = data.CloneAttribute(attrName, attr)
		}

		//Add action metadata into flow
		if h.act.IOMetadata() != nil && h.act.IOMetadata().Input != nil {
			//Adding action metadat into inputs
			for _, attr := range h.act.IOMetadata().Input {
				inputs[attr.Name()] = attr
			}
		}
	}

	return inputs, nil
}

func (h *Handler) generateOutputs(actionResults map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	if len(actionResults) == 0 {
		return nil, nil
	}

	if h.actionOutputMapper == nil {
		//for backwards compatibility
		return actionResults, nil
	}

	outputMetadata := h.act.IOMetadata().Output

	if outputMetadata != nil {

		outScope := data.NewFixedScopeFromMap(h.replyMd)
		inScope := data.NewSimpleScopeFromMap(actionResults, nil)

		err := h.actionOutputMapper.Apply(inScope, outScope)
		if err != nil {
			return nil, err
		}

		return outScope.GetAttrs(), nil
	}

	return actionResults, nil
}

func (h *Handler) String() string {
	return fmt.Sprintf("Handler[action:%s]", h.config.Action.Ref)
}
