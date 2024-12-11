package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type LogicOr struct {
	Left  Expression
	Right Expression
}

func (a *LogicOr) Eval(env *environment.Environment) object.Object {
	return &object.Bool{Val: object.TransToBool(a.Left.Eval(env)) || object.TransToBool(a.Right.Eval(env))}
}

func (a *LogicOr) String() string {
	return fmt.Sprintf("%s || %s", a.Left.String(), a.Right.String())
}

func (a *LogicOr) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.LogicOr {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(a.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(a.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (a *LogicOr) expressionNode() {}
