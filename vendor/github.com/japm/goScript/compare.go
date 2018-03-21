//Package goScript, javascript Eval() for go
//The MIT License (MIT)
//Copyright (c) 2016 Juan Pascual

package goScript

import (
	"fmt"
	"reflect"
)

func evalBinaryExprNEQ(left interface{}, right interface{}) (interface{}, error) {
	val, err := evalBinaryExprEQL(left, right)
	if err != nil {
		return nilInterf, err
	}
	return !val.(bool), nil
}

func evalBinaryExprEQL(left interface{}, right interface{}) (interface{}, error) {
	tp, e := binaryOperTypeL(left, right)
	if e != nil {
		return nilInterf, e
	}
	if tp.IsNil() {
		return nilInterf, nil
	}
	if tp.Bool() {

		l, err := castBool(left)
		if err != nil {
			return nilInterf, err
		}

		r, err := castBool(right)
		if err != nil {
			return nilInterf, err
		}

		return l == r, nil

	} else if !tp.IsNumeric() {
		switch left.(type) {
		case string:
			val, err := castString(right)
			if err != nil {
				return nilInterf, err
			}
			return left.(string) == val, nil
		case nil:
			switch right.(type) {
			case nil:
				return trueInterf, nil
			}
			return nilInterf, nil
		}
	} else {
		op := opEql{}
		return evalBinary(left, right, tp, op)
	}
	return nil, fmt.Errorf("Unimplemented == for types  %s and %s", reflect.TypeOf(left), reflect.TypeOf(right))
}
