package condition

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/mashling-lib/util"
)

var fLogger = logger.GetLogger("event-link-operator")

// ExpressionType for expression type
type ExpressionType int

const (
	EXPR_TYPE_NOT_VALID ExpressionType = iota - 1
	EXPR_TYPE_CONTENT
	EXPR_TYPE_ENV
	EXPR_TYPE_HEADER
)

var (
	OperatorRegistry = NewOperatorRegistry()
)

type Condition struct {
	Operator
	LHS string
	RHS string
}

type Operator interface {
	HasOperatorInfo
	Eval(lhs string, rhs string) bool
}

// OperRegistry is a registry for condition operators
type OperRegistry struct {
	operatorsMu sync.Mutex
	operators   map[string]Operator
}

// NewOperatorRegistry creates a new operator registry
func NewOperatorRegistry() *OperRegistry {
	return &OperRegistry{operators: make(map[string]Operator)}
}

// RegisterOperator registers an operator
func (r *OperRegistry) RegisterOperator(operator Operator) {
	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	if operator == nil {
		panic("OperatorRegistry: operator cannot be nil")
	}

	operatorNames := operator.OperatorInfo().Names

	for _, name := range operatorNames {
		if _, exists := r.operators[name]; exists {
			panic("OperatorRegistry: operator [" + name + "] already registered")
		}

		r.operators[name] = operator
	}

}

// Operator gets the specified operator
func (r *OperRegistry) Operator(operatorName string) (o Operator, exists bool) {

	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	operator, exists := r.operators[operatorName]
	return operator, exists
}

// Operators gets all the registered operators
func (r *OperRegistry) Operators() ([]string, []Operator) {

	r.operatorsMu.Lock()
	defer r.operatorsMu.Unlock()

	var names []string
	var opers []Operator
	for k, v := range r.operators {
		names = append(names, k)
		opers = append(opers, v)
	}

	return names, opers
}

// EvaluateExpression evaluates the specified expression
func EvaluateExpression(expression string, content string) bool {
	condition, err := getCondition(expression)
	if err != nil {
		fLogger.Debugf("Error getting the condition from expression [%v], [%v]", expression, err)
		return false
	}
	operator := condition.Operator
	lhsExpression := condition.LHS

	//evaluate the lhs against the content
	lhs, err := util.JsonPathEval(content, lhsExpression)
	if err != nil {
		return false
	}

	return operator.Eval(*lhs, condition.RHS)
}

// Evaluate evaluates the specified operator
func EvaluateCondition(condition Condition, content string) (bool, error) {
	//check if the message content is JSON first. mashling only supports JSON payloads for condition/content evaluation
	if !IsJSON(content) {
		return false, errors.New(fmt.Sprintf("Content is not a valid JSON payload [%v]", content))
	}

	operator := condition.Operator
	lhsExpression := condition.LHS

	//evaluate the lhs against the content
	lhs, err := util.JsonPathEval(content, lhsExpression)
	if err != nil {
		fLogger.Debugf("Error evaluating lhs jsonpath [%v] on the content [%v], [%v]", lhsExpression, content, err)
		return false, nil
	}

	return operator.Eval(*lhs, condition.RHS), nil
}

func getCondition(conditionExpr string) (*Condition, error) {
	oper, name, err := GetOperatorInExpression(conditionExpr)
	if err != nil {
		fLogger.Debugf("Error getting the operator from expression [%v], [%v]", conditionExpr, err)
		return nil, err
	}
	//found the operation!
	index := strings.Index(conditionExpr, *name)
	// find the LHS
	// Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
	lhs := strings.TrimSpace(conditionExpr[:index]) + "+"
	//get the value for LHS
	fLogger.Debugf("condition: left hand side found to be [%v", lhs)

	//find the RHS
	rhs := strings.TrimSpace(conditionExpr[index+len(*name):])
	fLogger.Debugf("condition: right hand side found to be [%v]", rhs)

	//create the condition
	condition := Condition{*oper, lhs, rhs}
	return &condition, nil

}

func GetConditionOperation(conditionStr string) (*Condition, error) {
	/**
	Content based conditions rules

	The condition identifier is "${" at the start and "}" at the end.

	If LHS
		If the condition clause starts with "trigger.content" then it refers to the trigger's payload. It maps internally to the "$." JSONPath of the payload.
		The above examples of JSONPath can be expressed as "${trigger.content.phoneNumbers[:1].type" and "${trigger.content.address.city" respectively.
		<<TBD>> If the condition clause does not start with "trigger.content":
		<<TBD>> If it starts with "env" then it is evaluated as an environment variable. So, "${env.PROD_ENV == true}" will be evaluated as a condition based on the environment variable.
	If Operator
		The condition must evaluate to a boolean output. Example operators are "==" and "!=".
	If RHS
		The condition RHS will be interpreted as a string
	*/
	if !strings.HasPrefix(conditionStr, util.Gateway_Link_Condition_LHS_Start_Expr) {
		return nil, errors.New("If does not match expected semantics, missing '${' at the start.")
	}
	if !strings.HasSuffix(conditionStr, util.Gateway_Link_Condition_LHS_End_Expr) {
		return nil, errors.New("condition 'If' does not match expected semantics, missing '}' at the end.")
	}

	condition := conditionStr[len(util.Gateway_Link_Condition_LHS_Start_Expr) : len(conditionStr)-len(util.Gateway_Link_Condition_LHS_End_Expr)]
	contentRoot := GetContentRoot()

	if !strings.HasPrefix(condition, contentRoot) {
		return nil, errors.New(fmt.Sprintf("condition 'If' JSONPath must start with %v", contentRoot))
	}

	condition = strings.Replace(condition, contentRoot, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1)

	condition = strings.TrimSpace(condition)

	condOperation, err := getCondition(condition)
	if err != nil {
		return nil, err
	}
	return condOperation, nil

}

// Insert inserts the value into the slice at the specified index,
// which must be in range.
// The slice must have room for the new element.
func Insert(slice []string, index int, value string) []string {
	// Grow the slice by one element.
	slice = slice[0 : len(slice)+1]
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	return slice
}

// GetConditionOperationAndExpressionType takes expression string as input and
//identifies expression type. type can be condtion/header/environment/etc based expression and
//prepares condtion operation and returs
func GetConditionOperationAndExpressionType(expressionStr string) (*Condition, ExpressionType, error) {
	if !strings.HasPrefix(expressionStr, util.Gateway_Link_Condition_LHS_Start_Expr) {
		return nil, EXPR_TYPE_NOT_VALID, errors.New("expression does not match expected semantics, missing '${' at the start")
	}
	if !strings.HasSuffix(expressionStr, util.Gateway_Link_Condition_LHS_End_Expr) {
		return nil, EXPR_TYPE_NOT_VALID, errors.New("expression does not match expected semantics, missing '}' at the end")
	}

	//trim expression starting and ending semantics i.e ${ and }
	expressionStr = strings.TrimPrefix(expressionStr, util.Gateway_Link_Condition_LHS_Start_Expr)
	expressionStr = strings.TrimSuffix(expressionStr, util.Gateway_Link_Condition_LHS_End_Expr)
	expressionStr = strings.TrimSpace(expressionStr)

	//decode condtion from expression string
	oper, name, err := GetOperatorInExpression(expressionStr)
	if err != nil {
		fLogger.Debugf("error getting the operator from expression [%v], [%v]", expressionStr, err)
		return nil, EXPR_TYPE_NOT_VALID, err
	}
	//found the operation!
	index := strings.Index(expressionStr, *name)
	// find the LHS
	lhs := strings.TrimSpace(expressionStr[:index])
	//find the RHS
	rhs := strings.TrimSpace(expressionStr[index+len(*name):])

	exprType := EXPR_TYPE_NOT_VALID
	contentRoot := GetContentRoot()

	if strings.HasPrefix(expressionStr, contentRoot) {
		//update lhs
		// Important!! The '+' at the end is required to access the value from jsonpath evaluation result!
		lhs = strings.Replace(lhs, contentRoot, util.Gateway_Link_Condition_LHS_JSONPath_Root, -1) + "+"
		//update expression type
		exprType = EXPR_TYPE_CONTENT
	} else if strings.HasPrefix(expressionStr, util.Gateway_Link_Condition_LHS_Header_Prifix) {
		//update lhs
		lhs = strings.TrimPrefix(lhs, util.Gateway_Link_Condition_LHS_Header_Prifix)
		//update expression type
		exprType = EXPR_TYPE_HEADER
	} else if strings.HasPrefix(expressionStr, util.Gateway_Link_Condition_LHS_Environment_Prifix) {
		//update lhs
		lhs = strings.TrimPrefix(lhs, util.Gateway_Link_Condition_LHS_Environment_Prifix)
		//update expression type
		exprType = EXPR_TYPE_ENV
	}

	//prepare condition and return
	condition := Condition{*oper, lhs, rhs}
	return &condition, exprType, nil
}
