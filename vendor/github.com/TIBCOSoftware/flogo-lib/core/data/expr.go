package data

type Expr interface {
	Eval(scope Scope) (interface{}, error)
}
