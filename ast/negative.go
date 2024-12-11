package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Negative struct {
	Expr Expression
}

func (b *Negative) Eval(env *environment.Environment) object.Object {
	numObj := b.Expr.Eval(env)
	if intObj, ok := numObj.(*object.Int); ok {
		return &object.Int{Val: -intObj.Val}
	}
	if floatObj, ok := numObj.(*object.Float); ok {
		return &object.Float{Val: -floatObj.Val}
	}
	panic("invalid - expression")
}

func (b *Negative) String() string {
	return fmt.Sprintf("-%s", b.Expr.String())
}

func (b *Negative) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Negative {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Expr: ")
	s.WriteString(b.Expr.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (b *Negative) expressionNode() {}
