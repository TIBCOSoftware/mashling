package runner

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)


func generateInputs(act action.Action, ctxData *trigger.ContextData) ([]*data.Attribute) {

	if ctxData == nil || ctxData.Attrs == nil {
		return nil
	}

	inputMetadata := action.GetConfigInputMetadata(act)

	if ctxData.HandlerCfg != nil && inputMetadata != nil {

		outputMapper := ctxData.HandlerCfg.GetActionInputMapper()

		outScope := data.NewFixedScope(inputMetadata)
		inScope := data.NewSimpleScope(ctxData.Attrs, nil)

		outputMapper.Apply(inScope, outScope)

		attrs := outScope.GetAttrs()

		inputs := make([]*data.Attribute, 0, len(inputMetadata))

		for _, attr := range attrs {
			inputs = append(inputs, attr)
		}

		return inputs

	} else {
		// for backwards compatibility make trigger outputs map directly to action inputs

		if len(ctxData.Attrs) > 0 {
			logger.Debug("No mapping specified, adding trigger outputs as inputs to action")

			inputs := make([]*data.Attribute, 0, len(ctxData.Attrs))

			for _, attr := range ctxData.Attrs {

				logger.Debugf(" Attr: %s, Type: %s, Value: %v", attr.Name(), attr.Type().String(), attr.Value())
				//inputs = append(inputs, data.NewAttribute( attr.Name, attr.Type, attr.Value))
				inputs = append(inputs, attr)

				attrName := "_T." + attr.Name()

				inputs = append(inputs, data.CloneAttribute(attrName, attr))
			}

			return inputs
		}
	}

	return nil
}

func generateOutputs(act action.Action, ctxData *trigger.ContextData, actionResults map[string]*data.Attribute) (map [string]*data.Attribute) {

	if len(actionResults) == 0 {
		return nil
	}

	if ctxData == nil {
		//for backwards compatibility
		return actionResults
	}

	outputMetadata := action.GetConfigOutputMetadata(act)

	if ctxData.HandlerCfg != nil && outputMetadata != nil {

		outputMapper := ctxData.HandlerCfg.GetActionOutputMapper()

		triggerId := ctxData.HandlerCfg.GetTriggerConfig().Id
		triggerMd := trigger.Instance(triggerId).Interf.Metadata()
		outScope := data.NewFixedScopeFromMap(triggerMd.Reply)
		inScope := data.NewSimpleScopeFromMap(actionResults, nil)

		outputMapper.Apply(inScope, outScope)

		return outScope.GetAttrs()
	}

	return actionResults
}

