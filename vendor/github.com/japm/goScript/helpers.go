//Package goScript, javascript Eval() for go
//The MIT License (MIT)
//Copyright (c) 2016 Juan Pascual

package goScript

import (
	"fmt"
	"go/ast"
	//"reflect"
	//"time"
)

func isEmpty(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, c := range s {
		if c != ' ' && c != '\t' {
			return false
		}
	}
	return true
}


//AST structires  with values resolved for constants
//specialized by each expression
type constNodeBasicLit struct {
	ast.BasicLit
	value interface{}
}

type constNodeIdent struct {
	ast.Ident
	value interface{}
}

type constNodeUnaryExpr struct {
	ast.UnaryExpr
	value interface{}
}

type constNodeBinaryExpr struct {
	ast.BinaryExpr
	value interface{}
}

type constNodeParenExpr struct {
	ast.ParenExpr
	value interface{}
}

//Replace constants in expressions for its values precalculated
//using AST nodes with values.
func resolveConstants(expr ast.Node) ast.Node {
	switch expr.(type) {
	case *ast.BasicLit:
		v, e := evalBasicLit(expr.(*ast.BasicLit), nil)
		if e != nil {
			return expr
		}
		return &constNodeBasicLit{*expr.(*ast.BasicLit), v}
	case *ast.BinaryExpr:
		bexp := expr.(*ast.BinaryExpr)
		bexp.X = resolveConstants(bexp.X).(ast.Expr)
		bexp.Y = resolveConstants(bexp.Y).(ast.Expr)
		v, e := evalBinaryExpr(expr.(*ast.BinaryExpr), nilContext{})
		if e != nil {
			return expr
		}
		return &constNodeBinaryExpr{*expr.(*ast.BinaryExpr), v}
	case *ast.ParenExpr:
		pexp := expr.(*ast.ParenExpr)
		pexp.X = resolveConstants(pexp.X).(ast.Expr)
		v, e := eval(pexp.X, nilContext{})
		if e != nil {
			return expr
		}
		return &constNodeParenExpr{*expr.(*ast.ParenExpr), v}

	case *ast.Ident:
		v, e := evalIdent(expr.(*ast.Ident), nilContext{})
		if e != nil {
			return expr
		}
		return &constNodeIdent{*expr.(*ast.Ident), v}
	case *ast.UnaryExpr:
		v, e := evalUnaryExpr(expr.(*ast.UnaryExpr), nilContext{})
		if e != nil {
			fmt.Println(e)
			return expr
		}
		return &constNodeUnaryExpr{*expr.(*ast.UnaryExpr), v}
	case *ast.CallExpr:
		callexp := expr.(*ast.CallExpr)
		for key, value := range callexp.Args {
			callexp.Args[key] = resolveConstants(value).(ast.Expr)
		}
		return expr
	case *ast.IndexExpr:
		idexpr := expr.(*ast.IndexExpr)
		idexpr.Index = resolveConstants(idexpr.Index).(ast.Expr)
		return expr
	case *ast.SliceExpr:
		slexpr := expr.(*ast.SliceExpr)
		slexpr.Low = resolveConstants(slexpr.Low).(ast.Expr)
		slexpr.High = resolveConstants(slexpr.High).(ast.Expr)
		slexpr.X = resolveConstants(slexpr.X).(ast.Expr)
		return expr
	default:
		return expr
	}
}
