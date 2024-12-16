package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type BitwiseNot struct {
	Expr Expression
}

func (b *BitwiseNot) Eval(env *environment.Environment) object.Object {
	numObj := b.Expr.Eval(env)
	if intObj, ok := numObj.(*object.Int); ok {
		return &object.Int{Val: ^intObj.Val}
	}
	panic(fmt.Sprintf("TypeError: invalid operand for bitwise not(^): %s", numObj.Type()))
}

func (b *BitwiseNot) String() string {
	return fmt.Sprintf("^%s", b.Expr.String())
}

func (b *BitwiseNot) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.BitwiseNot {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Expr: ")
	s.WriteString(b.Expr.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (b *BitwiseNot) expressionNode() {}
