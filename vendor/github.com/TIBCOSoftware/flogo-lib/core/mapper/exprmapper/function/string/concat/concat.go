package concat

import (
	"bytes"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("concat-function")

type Concat struct {
}

func init() {
	function.Registry(&Concat{})
}

func (s *Concat) GetName() string {
	return "concat"
}

func (s *Concat) GetCategory() string {
	return "string"
}

func (s *Concat) Eval(strs ...string) (string, error) {
	log.Debugf("Start concat function with parameters %s", strs)
	if len(strs) >= 2 {
		var buffer bytes.Buffer

		for _, v := range strs {
			buffer.WriteString(v)
		}
		log.Debugf("Done concat function with result %s", buffer.String())
		return buffer.String(), nil
	}

	return "", fmt.Errorf("Concat function at least have 2 arguments")
}
