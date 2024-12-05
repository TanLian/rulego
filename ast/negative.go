package ast

import (
	"fmt"

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

func (b *Negative) expressionNode() {}
