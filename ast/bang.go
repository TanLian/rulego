package ast

import (
	"fmt"

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

func (b *Bang) expressionNode() {}
