package fggos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/japm/goScript"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/definition"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// GosLinkExprManager is the Lua Implementation of a Link Expression Manager
type GosLinkExprManager struct {
	values map[int][]string
	exprs  map[int]*goScript.Expr
}

type varInfo struct {
	name string
	isd  int
}

type GosLinkExprManagerFactory struct {
}

// NewGosLinkExprManager creates a new LuaLinkExprManager
func (f *GosLinkExprManagerFactory) NewLinkExprManager(def *definition.Definition) definition.LinkExprManager {

	mgr := &GosLinkExprManager{}
	mgr.values = make(map[int][]string)
	mgr.exprs = make(map[int]*goScript.Expr)

	links := definition.GetExpressionLinks(def)

	for _, link := range links {

		if len(strings.TrimSpace(link.Value())) > 0 {

			fixedExpr := fixupExpression(link.Value())
			vars, exprStr := transExpr(fixedExpr)

			mgr.values[link.ID()] = vars

			logger.Debugf("expr: %v\n", exprStr)

			expr := &goScript.Expr{}
			err := expr.Prepare(exprStr)

			if err == nil {
				mgr.exprs[link.ID()] = expr
			} else {
				logger.Errorf("Error preparing expression: %s - %v", link.Value(), err)
			}
		}
	}

	return mgr
}

func fixupExpression(expr string) string {

	fixed := strings.Replace(expr, "${T.", "${trigger.", -1)
	fixed = strings.Replace(fixed, "${TriggerData.", "${trigger.", -1)
	fixed = strings.Replace(fixed, "${A", "${activity.", 1)

	return fixed
}

func transExpr(s string) ([]string, string) {

	var vars []string

	strLen := len(s)

	var buffer bytes.Buffer

	//todo cleanup
	for i := 0; i < strLen; i++ {

		if s[i] == '"' {
			buffer.WriteByte('"')

			var j int
			for j = i + 1; j < strLen; j++ {
				buffer.WriteByte(s[j])
				if s[j] == '"' {

					break
				}
			}
			i = j
		} else if s[i] == '$' && s[i+1] == '{' {
			//variable old

			buffer.WriteString("v[\"${")

			for j := i + 2; j < strLen; j++ {

				if !isPartOfNameOld(s[j]) {
					//fmt.Printf("\n")
					vars = append(vars, s[i:j])
					buffer.WriteString("\"]")
					buffer.WriteByte(s[j])
					i = j
					break
				} else if j == strLen-1 {
					//last char
					vars = append(vars, s[i:j+1])
					buffer.WriteByte(s[j])
					buffer.WriteString("\"]")
					i = j + 1
					break
				}
				buffer.WriteByte(s[j])
			}
		} else if s[i] == '$' {
			//variable

			buffer.WriteString("v[\"$")

			for j := i + 1; j < strLen; j++ {

				if !isPartOfName(s[j]) {
					//fmt.Printf("\n")
					vars = append(vars, s[i:j])
					buffer.WriteString("\"]")
					buffer.WriteByte(s[j])
					i = j
					break
				} else if j == strLen-1 {
					//last char
					vars = append(vars, s[i:j+1])
					buffer.WriteByte(s[j])
					buffer.WriteString("\"]")
					i = j + 1
					break
				}
				buffer.WriteByte(s[j])
			}
		} else {
			buffer.WriteByte(s[i])
		}
	}

	return vars, buffer.String()
}

func isPartOfName(char byte) bool {

	if (char < '0' || char > '9') && (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != '.' && char != '_' && char != '[' && char != ']' {
		return false
	}

	return true
}

func isPartOfNameOld(char byte) bool {

	if (char < '0' || char > '9') && (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != '.' && char != '_' && char != '}' && char != '[' && char != ']' {
		return false
	}

	return true
}

// EvalLinkExpr implements LinkExprManager.EvalLinkExpr
func (em *GosLinkExprManager) EvalLinkExpr(link *definition.Link, scope data.Scope) (ret bool, err error) {

	defer func() {
		if r := recover(); r != nil {
			ret = false
			err = definition.NewLinkExprError(fmt.Sprintf("Error evaluting expression: %s", r))
		}
	}()

	if link.Type() == definition.LtDependency {
		// dependency links are always true
		return true, nil
	}

	vars, attrsOK := em.values[link.ID()]
	expr, exprOK := em.exprs[link.ID()]

	if !attrsOK || !exprOK {

		return false, fmt.Errorf("Unable to evaluate expression '%s', did not compile properly\n", link.Value())
	}

	ctxt := make(map[string]interface{})
	vals := make(map[string]interface{})

	for _, varRep := range vars {

		resolver := definition.GetDataResolver()

		val, err := resolver.Resolve(varRep, scope)

		if err == nil {
			//	return false, err
			vals[varRep] = FixUpValue(val)
		}

		//var attrValue interface{}
		//
		//
		//
		//attrName, attrPath, _ := data.PathDeconstruct(varRep)
		//attr, exists := scope.GetAttr(attrName)
		//
		//attrValue = attr.Value
		//
		//if exists && len(attrPath) > 0 {
		//
		//	valMap, ok := attrValue.(map[string]interface{})
		//
		//	var val interface{}
		//
		//	if ok && len(valMap) > 0 {
		//		val = nil //data.GetMapValue(valMap, attrPath)
		//	} else {
		//		//assume its a map[string]string
		//		strMap, ok := attrValue.(map[string]string)
		//		if ok {
		//			val = strMap[attrPath]
		//		}
		//	}
		//
		//	attrValue = FixUpValue(val)
		//
		//	vals[varRep] = attrValue
		//}
	}

	ctxt["v"] = vals

	f := isDefinedFunc{scope: scope}
	ctxt["isDefined"] = f.isDefined

	logger.Debugf("Vals: %v", vals)

	val, err := expr.Eval(ctxt)

	if err != nil {
		return false, definition.NewLinkExprError(fmt.Sprintf("Error evaluting expression: %s", err.Error()))
	}

	return val.(bool), nil
}

// FixUpValue fixes json numbers
func FixUpValue(val interface{}) interface{} {

	ret := val
	var err error

	switch t := val.(type) {
	case json.Number:
		if strings.Index(t.String(), ".") > -1 {
			ret, err = t.Float64()
		} else {
			ret, err = t.Int64()
		}
	}

	if err != nil {
		ret = val
	}

	return ret
}

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

type isDefinedFunc struct {
	scope data.Scope
}

func (f *isDefinedFunc) isDefined(value string) bool {

	resolver := definition.GetDataResolver()
	_, err := resolver.Resolve(value, f.scope)

	return err == nil
}
