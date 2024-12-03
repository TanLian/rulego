package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type LogicOr struct {
	Left  Expression
	Right Expression
}

func (a *LogicOr) Eval(env *environment.Environment) object.Object {
	return &object.Bool{Val: object.TransToBool(a.Left.Eval(env)) || object.TransToBool(a.Right.Eval(env))}
}

func (a *LogicOr) String() string {
	return fmt.Sprintf("%s || %s", a.Left.String(), a.Right.String())
}

func (a *LogicOr) expressionNode() {}
