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
	return &object.Bool{Val: object.TransToBool(a.Left.Eval(env)) || object.TransToBool(a.Right.Eval(env))}
}

func (a *Or) String() string {
	return fmt.Sprintf("%s || %s", a.Left.String(), a.Right.String())
}

func (a *Or) expressionNode() {}
