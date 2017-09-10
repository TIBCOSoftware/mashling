package util

import (
	"fmt"
	"github.com/NodePrime/jsonpath"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/pkg/errors"
	"strings"
)

var newLogger = logger.GetLogger("jsonpath-eval")

func JsonPathEval(jsonData string, expression string) (*string, error) {
	paths, err := jsonpath.ParsePaths(expression)
	if err != nil {
		return nil, err
	}
	newLogger.Debugf("jsonpath expression is [%v], data is [%v]", expression, jsonData)

	eval, err := jsonpath.EvalPathsInBytes([]byte(jsonData), paths)
	if err != nil {
		newLogger.Infof("unable to evaluate jsonpath expression, error is [%v]", err)
		return nil, err
	}

	for {
		if result, ok := eval.Next(); ok {
			//return after the first match
			value := string(result.Value) //The value obtained will be encased in double quote characters, e.g. "USA" when the value happens to be USA.
			//Trim the double quote suffix and prefix
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, "\"")
			newLogger.Debugf("jsonpath [%v] evaluated to value [%v]", expression, value)
			return &value, nil // true -> show keys in pretty string
		} else {
			return nil, errors.New(fmt.Sprintf("Error evaluating jsonpath expression[%v] on content [%v]", expression, jsonData))
		}
	}
}
