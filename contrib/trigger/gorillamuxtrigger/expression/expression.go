package expression

import (
	"strings"

	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/elgs/gojq"
)

//MashlingExpr mashling expression
type MashlingExpr struct {
	lhs      string
	rhs      string
	operator string
}

//package logger
var log = logger.GetLogger("triggerhttpnew-expression")

//EvalMashlingExpr evaluates mashling expression
func EvalMashlingExpr(exprStr string, content string) (bool, error) {
	result := false

	log.SetLogLevel(logger.DebugLevel)

	//tokenize expression string
	expr, err := tokenizeExpr(exprStr)
	if err != nil {
		log.Errorf("not a valid expression")
		return false, err
	}

	//parse content
	parser, err := gojq.NewStringQuery(content)
	if err != nil {
		log.Errorf("error while parsing content - %v", err)
		return false, err
	}

	//resolve expression-lhs token value
	exprLhs := expr.lhs
	if strings.HasPrefix(exprLhs, "trigger.content.") {
		exprLhs = exprLhs[16:len(exprLhs)]
	}
	val, err := parser.Query(exprLhs)
	if err != nil {
		log.Errorf("no element found for the expression '%v'", exprLhs)
		return false, fmt.Errorf("no element found for the expression '%v'", exprLhs)
	}
	expr.lhs = val.(string)

	//evaluate expression
	result, err = eval(expr)
	if err != nil {
		log.Errorf("error while evaluating expressio")
		return false, fmt.Errorf("not able to evaluate expression - %v", err)
	}

	return result, nil
}

func eval(expr MashlingExpr) (bool, error) {
	result := false
	switch expr.operator {
	case "==":
		result = strings.Compare(expr.lhs, expr.rhs) == 0
	default:
		return false, fmt.Errorf("operator '%v' not supported", expr.operator)
	}
	return result, nil
}

func tokenizeExpr(exprStr string) (expr MashlingExpr, err error) {
	if strings.HasPrefix(exprStr, "${") {
		exprStr = exprStr[2:len(exprStr)]
	}
	if strings.HasSuffix(exprStr, "}") {
		exprStr = exprStr[0 : len(exprStr)-1]
	}

	exprTokens := strings.Split(exprStr, " ")
	expr.lhs = exprTokens[0]
	expr.operator = exprTokens[1]
	expr.rhs = exprTokens[2]

	if len(exprTokens) != 3 {
		return expr, fmt.Errorf("invalid expression '%v'", exprStr)
	}

	return expr, nil
}
