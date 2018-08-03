package random

import (
	"math/rand"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/data"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("random-function")

type Random struct {
}

func init() {
	function.Registry(&Random{})
}

func (s *Random) GetName() string {
	return "random"
}

func (s *Random) GetCategory() string {
	return "number"
}

func (s *Random) Eval(limitIn interface{}) int {
	limit, err := data.CoerceToInteger(limitIn)
	if err != nil {
		log.Errorf("Convert %+v to int error %s", limitIn, err.Error())
		limit = 10
	}
	log.Debugf("Generate sudo-random integer number within the scope of [0, %d)", limit)
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(limit)
}
