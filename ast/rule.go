package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Rule struct {
	object.Empty
	Name       string
	Statements []Statement
}

func (rl *Rule) expressionNode() {}

func (rl *Rule) Eval(env *environment.Environment) object.Object {
	env.Set(rl.Name, rl)
	return object.Null
}

func (rl *Rule) Call(env *environment.Environment) object.Object {
	for _, v := range rl.Statements {
		if obj, flg := v.Exec(env); flg&RETURN != 0 {
			return obj
		}
	}
	return object.Null
}

func (rl *Rule) AST(num int) string { // TODO
	return ""
}

func (rl *Rule) String() string {
	var s strings.Builder
	s.WriteString("rule ")
	s.WriteString(fmt.Sprintf("%s {", rl.Name))
	for _, v := range rl.Statements {
		s.WriteString("	" + v.String() + ";")
	}
	s.WriteString("}")
	return s.String()
}
