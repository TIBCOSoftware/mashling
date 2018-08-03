package equals

import (
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("equals-function")

type Equals struct {
}

func init() {
	function.Registry(&Equals{})
}

func (s *Equals) GetName() string {
	return "equals"
}

func (s *Equals) GetCategory() string {
	return "string"
}
func (s *Equals) Eval(str, str2 string) bool {
	log.Debugf(`Reports whether "%s" equals "%s" `, str, str2)
	return str == str2
}
