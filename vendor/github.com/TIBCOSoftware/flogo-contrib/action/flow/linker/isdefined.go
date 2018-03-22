package linker

import (
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
)

type IsDefine struct {
}

func init() {
	function.Registry(&IsDefine{})
}

func (s *IsDefine) GetName() string {
	return "isDefined"
}

func (s *IsDefine) GetCategory() string {
	return ""
}

func (s *IsDefine) Eval(in interface{}) bool {
	return in != nil
}
