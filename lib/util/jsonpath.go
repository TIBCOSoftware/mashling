/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package util

import (
	"encoding/json"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/oliveagle/jsonpath"
	"github.com/pkg/errors"
)

var newLogger = logger.GetLogger("jsonpath-eval")

//JsonPathEval evaluates value from the content passed based on expression
func JsonPathEval(jsonData string, expression string) (*string, error) {
	var jsondata interface{}
	json.Unmarshal([]byte(jsonData), &jsondata)

	res, err := jsonpath.JsonPathLookup(jsondata, expression)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error evaluating jsonpath expression[%v] on content [%v]", expression, jsonData))
	}
	strVar := res.(string)
	newLogger.Debugf("jsonpath [%v] evaluated to value [%v]", expression, strVar)
	return &strVar, nil
}
