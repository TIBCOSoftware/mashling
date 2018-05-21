package linker

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper"
)

var log = logger.GetLogger("linker")

type linkerManager struct {
}

type linkerFactory struct {
}

func NewDefaultLinkerFactory() definition.LinkExprManagerFactory {
	return &linkerFactory{}
}

func (factory *linkerFactory) NewLinkExprManager() definition.LinkExprManager {
	return &linkerManager{}
}

func (em *linkerManager) EvalLinkExpr(link *definition.Link, scope data.Scope) (bool, error) {
	value := link.Value()
	if value == "" {
		return true, nil
	}

	log.Debugf("WI link expression value [%s]", value)
	funcValue, err := exprmapper.GetMappingValue(value, scope, definition.GetDataResolver())
	if err != nil {
		log.Error("Get value from link value %+v, error %s", value, err.Error())
		return false, fmt.Errorf("Get value from link value %+v, error %s", value, err.Error())
	}

	b, err := data.CoerceToBoolean(funcValue)
	if err != nil {
		log.Error("Parser [%+v] to boolean error [%s]", value, err.Error())
		return false, fmt.Errorf("Parser [%+v] to boolean error [%s]", value, err.Error())
	}
	log.Debugf("Linking %s result %b", link.Value(), b)
	return b, nil
}
