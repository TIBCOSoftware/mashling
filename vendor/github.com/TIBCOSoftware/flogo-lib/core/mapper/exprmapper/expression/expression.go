package expression

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/gocc/lexer"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/gocc/parser"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("expression")

func ParseExpression(exprString string) (expr.Expr, error) {
	st, err := getParser(exprString)
	if err != nil {
		return nil, err
	}

	ex, ok := st.(expr.Expr)
	if ok {
		return ex, nil
	}
	return nil, fmt.Errorf("Not a valid expression str")
}

func getParser(exampleStr string) (interface{}, error) {
	lex := lexer.NewLexer([]byte(exampleStr))
	p := parser.NewParser()
	st, err := p.Parse(lex)
	return st, err
}

func IsExpression(mapValue string) bool {
	_, err := ParseExpression(mapValue)
	return err == nil
}
