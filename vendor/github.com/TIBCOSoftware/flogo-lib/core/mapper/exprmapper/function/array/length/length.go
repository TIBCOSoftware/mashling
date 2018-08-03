package length

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("length-function")

type Length struct {
}

func init() {
	function.Registry(&Length{})
}

func (s *Length) GetName() string {
	return "length"
}

func (s *Length) GetCategory() string {
	return "array"
}

func (s *Length) Eval(arr interface{}) (int, error) {
	log.Debugf("Return the length of array \"%v\"", arr)
	myArr, err := data.CoerceToArray(arr)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}
	return len(myArr), nil
}
