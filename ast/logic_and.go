package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type LogicAnd struct {
	Left  Expression
	Right Expression
}

func (a *LogicAnd) Eval(env *environment.Environment) object.Object {
	left, ok := a.Left.Eval(env).(*object.Bool)
	if !ok {
		panic(fmt.Sprintf("type error: %s is not bool", a.Left.String()))
	}

	right, ok := a.Right.Eval(env).(*object.Bool)
	if !ok {
		panic(fmt.Sprintf("type error: %s is not bool", a.Right.String()))
	}
	return &object.Bool{Val: left.Val && right.Val}
}

func (a *LogicAnd) String() string {
	return fmt.Sprintf("%s && %s", a.Left.String(), a.Right.String())
}

func (a *LogicAnd) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.LogicAnd {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(a.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(a.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (a *LogicAnd) expressionNode() {}
