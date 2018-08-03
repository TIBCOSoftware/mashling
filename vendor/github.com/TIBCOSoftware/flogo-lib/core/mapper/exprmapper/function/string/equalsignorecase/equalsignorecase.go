package equalsignorecase

import (
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("equalsIgnoreCase-function")

type EqualsIgnoreCase struct {
}

func init() {
	function.Registry(&EqualsIgnoreCase{})
}

func (s *EqualsIgnoreCase) GetName() string {
	return "equalsIgnoreCase"
}

func (s *EqualsIgnoreCase) GetCategory() string {
	return "string"
}
func (s *EqualsIgnoreCase) Eval(str, str2 string) bool {
	log.Debugf(`Reports whether "%s" equels "%s" with ignore case`, str, str2)
	return strings.EqualFold(str, str2)
}
