package substring

import (
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("substring-function")

type Substring struct {
}

func init() {
	function.Registry(&Substring{})
}

func (s *Substring) GetName() string {
	return "substring"
}

func (s *Substring) GetCategory() string {
	return "string"
}

func (s *Substring) Eval(str string, position int, length int) string {
	log.Debugf("Return the substring of string \"%s\"", str)
	runes := []rune(str)

	if length > 0 {
		l := position + length
		if l > len(runes) {
			l = len(runes)
		}
		return string(runes[position:l])
	} else if length < 0 {
		l := len(str) + length
		if l < 0 {
			return string(runes[0:0])
		}
		return string(runes[len(str)+length:])
	}
	return string(runes[0:0])
}
