package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Bang struct {
	Expr Expression
}

func (b *Bang) Eval(env *environment.Environment) object.Object {
	expr, ok := b.Expr.Eval(env).(*object.Bool)
	if !ok {
		panic(fmt.Sprintf("type error: %s is not bool", b.Expr.String()))
	}
	return &object.Bool{Val: !expr.Val}
}

func (b *Bang) String() string {
	return fmt.Sprintf("!%s", b.Expr.String())
}

func (b *Bang) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Bang {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Expr: ")
	s.WriteString(b.Expr.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (b *Bang) expressionNode() {}
