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
	obj := b.Expr.Eval(env)
	boolObj, ok := obj.(*object.Bool)
	if !ok {
		panic(fmt.Sprintf("TypeError: ! operator does not support converting %s to bool", obj.Type()))
	}
	return &object.Bool{Val: !boolObj.Val}
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
