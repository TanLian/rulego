package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type FnLiteralObj struct {
	*object.Empty
	Name       string
	Args       []string
	Statements []Statement
}

func (fl *FnLiteralObj) Eval(env *environment.Environment) object.Object {
	env.Set(fl.Name, fl)
	return object.Null
}

func (fl *FnLiteralObj) Call(args []object.Object) object.Object {
	if len(fl.Args) != len(args) {
		panic("the length of parameters is not equal")
	}

	env := environment.New(environment.Root)
	for i := 0; i < len(fl.Args); i++ {
		env.SetCurrent(fl.Args[i], args[i])
	}
	for i, v := range fl.Statements {
		if i == len(fl.Statements)-1 {
			obj, _ := v.Exec(env)
			return obj
		}

		if obj, flg := v.Exec(env); flg&RETURN != 0 {
			return obj
		}
	}
	return object.Null
}

func (fl *FnLiteralObj) String() string {
	var s strings.Builder
	s.WriteString("fn ")
	s.WriteString(fl.Name)
	s.WriteString(fmt.Sprintf("(%s) {", strings.Join(fl.Args, ",")))
	for _, v := range fl.Statements {
		s.WriteString(v.String() + ";")
	}
	s.WriteString("}")
	return s.String()
}

func (fl *FnLiteralObj) expressionNode() {}
