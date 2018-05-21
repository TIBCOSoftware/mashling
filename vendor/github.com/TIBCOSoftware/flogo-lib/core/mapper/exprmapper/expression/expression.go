package expression

import (
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/function"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/gocc/lexer"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/gocc/parser"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/funcexprtype"
)

var log = logger.GetLogger("expression")

const (
	EXPRESSION = iota
	TERNARY_EXPRESSION
	FUNCTION
	STRING
)

type functionExpression struct {
	expr string
}

type weExpr struct {
	expr string
}

func NewFunctionExpression(expression string) *functionExpression {
	return &functionExpression{expr: expression}
}

func NewExpression(expression string) *weExpr {
	return &weExpr{expr: expression}
}

func (f *functionExpression) Eval() ([]interface{}, error) {
	function, err := f.GetFunction()
	if err != nil {
		log.Errorf("Get function error %+v", err.Error())
		return nil, err
	}

	return function.Eval()
}

func (f *functionExpression) EvalWithScope(inputScope data.Scope, resolver data.Resolver) ([]interface{}, error) {
	function, err := f.GetFunction()
	if err != nil {
		log.Errorf("Get function error %+v", err.Error())
		return nil, err
	}

	return function.EvalWithScope(inputScope, resolver)
}

func (f *functionExpression) GetFunction() (*function.FunctionExp, error) {
	st, err := GetParser(f.expr)
	if err != nil {
		log.Errorf("Error to parser functions, %+v ", err.Error())
		return nil, err
	}
	function := st.(*function.FunctionExp)

	return function, nil
}

func (e *weExpr) Eval() (interface{}, error) {
	expression, err := e.GetExpression()
	if err != nil {
		log.Errorf("Get expression error %+v", err.Error())
		return nil, err
	}

	return expression.Eval()
}

func (e *weExpr) EvalWithScope(inputScope data.Scope, resolver data.Resolver) (interface{}, error) {
	expression, err := e.GetExpression()
	if err != nil {
		log.Errorf("Get expression error %+v", err.Error())
		return nil, err
	}

	return expression.EvalWithScope(inputScope, resolver)
}

func (e *weExpr) GetExpression() (*expr.Expression, error) {
	st, err := GetParser(e.expr)
	if err != nil {
		log.Errorf("Error to parser functions %+v", err.Error())
		return nil, err
	}
	expression := st.(*expr.Expression)

	return expression, nil
}

func (e *weExpr) GetTernaryExpression() (*expr.TernaryExpressio, error) {
	st, err := GetParser(e.expr)
	if err != nil {
		log.Errorf("Error to parser functions %+v", err.Error())
		return nil, err
	}
	expression := st.(*expr.TernaryExpressio)

	return expression, nil
}

func GetParser(exampleStr string) (interface{}, error) {
	lex := lexer.NewLexer([]byte(exampleStr))
	p := parser.NewParser()
	st, err := p.Parse(lex)
	return st, err
}

func IsExpression(mapValue string) bool {
	return GetExpressionType(mapValue) == EXPRESSION
}

func GetExpressionType(mapValue string) int {
	st, err := GetParser(mapValue)
	if err != nil {
		return STRING
	}

	switch t := st.(type) {
	case *expr.TernaryExpressio:
		return TERNARY_EXPRESSION
	case *expr.Expression:
		if t.Type == funcexprtype.FUNCTION || t.Type == funcexprtype.EXPRESSION {
			return EXPRESSION
		} else {
			return STRING
		}
	case *function.FunctionExp:
		return FUNCTION
	}

	return STRING
}

func IsTernaryExpression(mapValue string) bool {
	return GetExpressionType(mapValue) == TERNARY_EXPRESSION
}

func IsFunction(mapValue string) bool {
	return GetExpressionType(mapValue) == FUNCTION
}
