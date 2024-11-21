package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Group struct {
	Expr Expression
}

func (g *Group) Eval(env *environment.Environment) object.Object {
	return g.Expr.Eval(env)
}

func (g *Group) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

func (g *Group) expressionNode() {}
