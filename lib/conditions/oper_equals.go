package condition

import (
	"strings"
)

/**
specify the exact name of the operator such that the operator can be
used directly in an expression in mashling event-links. the operator
must be preceded by a space (' ') and succeeded by a space (' ') when
used in an expression.

an operator can have more than one names. make sure that the names do not
collide with other operators.

this is a string (alphanumeric) equality operator.

e.g. ${trigger.content.country == USA}
*/
var infoEquals = &OperatorInfo{
	Names:       []string{"==", "equals"},
	Description: `Support for equals operation to be used in the conditions`,
}

func init() {
	OperatorRegistry.RegisterOperator(&Equals{info: infoEquals})
}

type Equals struct {
	info *OperatorInfo
}

func (o *Equals) OperatorInfo() *OperatorInfo {
	return o.info
}

// Eval implementation of condition.Operator.Eval
func (o *Equals) Eval(lhs string, rhs string) bool {
	if strings.Compare(lhs, rhs) == 0 {
		return true
	}
	return false
}
