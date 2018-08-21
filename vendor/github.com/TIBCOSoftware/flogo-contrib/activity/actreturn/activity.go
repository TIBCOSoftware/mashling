package actreturn

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-flogo-return")

const (
	ivMappings = "mappings"
)

// ReturnActivity is an Activity that is used to return/return via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type ReturnActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new ReturnActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ReturnActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *ReturnActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *ReturnActivity) Eval(ctx activity.Context) (done bool, err error) {

	actionCtx := ctx.ActivityHost()

	if ctx.GetInput(ivMappings) == nil {
		//No mapping
		actionCtx.Return(nil, nil)
		return true, nil
	}

	mappings, ok := ctx.GetInput(ivMappings).([]interface{})
	if !ok {
		return false, activity.NewError("invalid return mappings, mappings must be array", "", nil)
	}

	log.Debugf("Mappings: %+v", mappings)

	mapperDef, err := mapper.NewMapperDefFromAnyArray(mappings)

	//todo move this to a action instance level initialization, need the notion of static inputs or config
	returnMapper := mapper.NewBasicMapper(mapperDef, ctx.ActivityHost().GetResolver())
	if err != nil {
		return false, activity.NewError(err.Error(), "", nil)
	}

	outputScope := newOutputScope(actionCtx, mapperDef)
	inputScope := actionCtx.WorkingData() //flow data

	err = returnMapper.Apply(inputScope, outputScope)

	if err != nil {
		return false, activity.NewError(err.Error(), "", nil)
	}

	actionCtx.Return(outputScope.GetAttrs(), nil)

	return true, nil
}

func newOutputScope(activityHost activity.Host, mapperDef *data.MapperDef) *data.FixedScope {

	if activityHost.IOMetadata() == nil {
		//todo temporary fix to support tester service
		attrs := make(map[string]*data.Attribute, len(mapperDef.Mappings))

		for _, mappingDef := range mapperDef.Mappings {
			attr, _ := data.NewAttribute(mappingDef.MapTo, data.TypeAny, nil)
			attrs[attr.Name()] = attr
		}

		return data.NewFixedScope(attrs)
	} else {
		outAttrs := activityHost.IOMetadata().Output
		attrs := make(map[string]*data.Attribute, len(outAttrs))

		for _, outAttr := range outAttrs {
			attrs[outAttr.Name()] = outAttr
		}

		//create a fixed scope using the output metadata
		return data.NewFixedScope(attrs)
	}
}
