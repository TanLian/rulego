package ast

import (
	"fmt"

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
	panic("invalid ^ expression")
}

func (b *BitwiseNot) String() string {
	return fmt.Sprintf("^%s", b.Expr.String())
}

func (b *BitwiseNot) expressionNode() {}
