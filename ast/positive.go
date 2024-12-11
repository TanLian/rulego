package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Positive struct {
	Expr Expression
}

func (b *Positive) Eval(env *environment.Environment) object.Object {
	numObj := b.Expr.Eval(env)
	if intObj, ok := numObj.(*object.Int); ok {
		return &object.Int{Val: intObj.Val}
	}
	if floatObj, ok := numObj.(*object.Float); ok {
		return &object.Float{Val: floatObj.Val}
	}
	panic("invalid + expression")
}

func (b *Positive) String() string {
	return fmt.Sprintf("+%s", b.Expr.String())
}

func (b *Positive) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Positive {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Expr: ")
	s.WriteString(b.Expr.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (b *Positive) expressionNode() {}
