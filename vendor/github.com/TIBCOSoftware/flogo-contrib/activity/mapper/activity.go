package mapper

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-flogo-mapper")

const (
	ivMappings = "mappings"
)

// MapperActivity is an Activity that is used to reply/return via the trigger
// inputs : {method,uri,params}
// outputs: {result}
type MapperActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new MapperActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MapperActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *MapperActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Invokes a REST Operation
func (a *MapperActivity) Eval(context activity.Context) (done bool, err error) {

	mappings := context.GetInput(ivMappings).([]interface{})

	log.Debugf("Mappings: %+v", mappings)

	mapperDef, err := mapper.NewMapperDefFromAnyArray(mappings)

	//todo move this to a action instance level initialization, need the notion of static inputs or config
	actionMapper := mapper.NewBasicMapper(mapperDef, context.ActivityHost().GetResolver())

	if err != nil {
		return false, err
	}

	activityHost := context.ActivityHost()
	actionScope := activityHost.WorkingData() // action/flow data

	err = actionMapper.Apply(actionScope, actionScope)

	if err != nil {
		return false, err
	}

	return true, nil
}
