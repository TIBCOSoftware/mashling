#goScript

Golang like scripting language

Uses go/ast and reflect to allow arbitrary expression evaluation
in a given context.

#Highlights

* Simple and hackable code.

* Fully reflect function call including variable number of args.

* Automatic casting
  * Numeric casting  always to the bigger more generic representation.
  * True && 0 = false
  * Left operand governs the type 1 + "1" = 2,  "1" + 1 = "11"

* Map, slice, array and string indexer access

* Full control over evaluation context.
  * Generic map[string]interface{} / *map[string]interface{} contexts
  * Any struct with its fields as variables
  * Custom Context for implentation/user defined context.

* Mostly tested ;)

#Roadmap

* 1) Rock solid expression evaluation.
* 2) Enrich evaluation context
* 3) Full script program implementation 
* 4) Performance optimizations

#Examples

```go

import (
	"fmt"
	"time"

	"github.com/japm/goScript"
)

//Example type
type ab struct {
	A int
}

//Function as identifier
func T() int {
	return 0
}

func (a ab) Test(x int, z ...interface{}) []interface{} {
	return z
}


func (a ab) Test2(x int, z ...int) []int {
	return z
}

func (a ab) Test3(x float64, y float64) float64 {
	return x + y
}

//Custom dummy context, the ident value is its own name
func (a ab) GetIdent(name string) (val interface{}, err error) {
	return name, nil
}

func examples() {

	d := make(map[string]interface{})
	d["a"] = 1
	ctxt := make(map[string]interface{})
	ctxt["a"] = time.Now()
	ctxt["b"], _ = time.ParseDuration("3h")
	ctxt["c"] = []int{0, 1, 2, 3}
	ctxt["d"] = d
	ctxt["e"] = &ab{45}
	ctxt["f"] = ab{45}
	ctxt["g"] = []interface{}{nil, nil}

	i := new(int)
	*i = 3
	ctxt["h"] = i
	ctxt["i"] = T //A Function

	exp := &goScript.Expr{}
	err := exp.Prepare("(*e).A")
	//err := exp.Prepare("f.A")
	//err := exp.Prepare("a.Day() + 1.4")
	//err := exp.Prepare("f.Test(1,2)")
	//err := exp.Prepare("f.Test2(3,4,c[3])")
	//err := exp.Prepare("*e")
	//err := exp.Prepare("g[0]")
	//err := exp.Prepare("f.Test3(2,3)")
	//err := exp.Prepare("(1 + c[2]) * 3")
	//err := exp.Prepare("3.5 + len(c)")
	//err := exp.Prepare("i()")
	//err := exp.Prepare("d[\"a\"]")
	//err := exp.Prepare("*h + 2")

	if err != nil {
		fmt.Println(err)
		return
	}


  val, err := exp.Eval(ctxt)
  //val, err := exp.Eval(&ctxt)
  
  //Custom context examples
	//err := exp.Prepare("Name1")
	//val, err := exp.Eval(ab{1})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Result...", val)
}
```
