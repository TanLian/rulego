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

func (fl *FnLiteralObj) Call(env *environment.Environment) object.Object {
	for _, v := range fl.Statements {
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
