package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("activity-tibco-app")

const (
	ivAttrName = "attribute"
	ivOp       = "operation"
	ivType     = "type"
	ivValue    = "value"

	ovValue = "value"
)

// AppActivity is a App Activity implementation
type AppActivity struct {
	sync.Mutex
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &AppActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *AppActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *AppActivity) Eval(context activity.Context) (done bool, err error) {

	attrName := context.GetInput(ivAttrName).(string)
	op := strings.ToUpper(context.GetInput(ivOp).(string)) //ADD,UPDATE,GET

	switch op {
	case "ADD":
		log.Debug("In ADD operation")
		dt, ok := data.ToTypeEnum(strings.ToLower(context.GetInput(ivType).(string)))

		if !ok {
			errorMsg := fmt.Sprintf("Unsupported type '%s'", context.GetInput(ivType).(string))
			log.Error(errorMsg)
			return false, activity.NewError(errorMsg, "", nil)
		}

		val := context.GetInput(ivValue)
		//data.CoerceToValue(val, dt)

		data.GetGlobalScope().AddAttr(attrName, dt, val)
		context.SetOutput(ovValue, val)
	case "GET":
		log.Debug("In GET operation")
		typedVal, ok := data.GetGlobalScope().GetAttr(attrName)

		if !ok {
			errorMsg := fmt.Sprintf("Attribute not defined: '%s'", attrName)
			log.Error(errorMsg)
			return false, activity.NewError(errorMsg, "", nil)
		}

		context.SetOutput(ovValue, typedVal.Value())
	case "UPDATE":
		log.Debug("In UPDATE operation")
		val := context.GetInput(ivValue)
		//data.CoerceToValue(val, dt)

		data.GetGlobalScope().SetAttrValue(attrName, val)
		context.SetOutput(ovValue, val)
	default:
		errorMsg := fmt.Sprintf("Unsupported Op:'%s' ", op)
		log.Error(errorMsg)
		return false, activity.NewError(errorMsg, "", nil)
	}

	return true, nil
}
