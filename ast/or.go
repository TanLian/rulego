package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Or struct {
	Left  Expression
	Right Expression
}

func (a *Or) Eval(env *environment.Environment) object.Object {
	left, ok := a.Left.Eval(env).(*object.Int)
	if !ok {
		panic(fmt.Sprintf("type error: %s is not int", a.Left.String()))
	}

	right, ok := a.Right.Eval(env).(*object.Int)
	if !ok {
		panic(fmt.Sprintf("type error: %s is not int", a.Right.String()))
	}
	return &object.Int{Val: left.Val | right.Val}
}

func (a *Or) String() string {
	return fmt.Sprintf("%s | %s", a.Left.String(), a.Right.String())
}

func (a *Or) expressionNode() {}
