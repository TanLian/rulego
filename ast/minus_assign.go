package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type MinusAssign struct {
	Left  Expression
	Right Expression
}

func (p *MinusAssign) Eval(env *environment.Environment) object.Object {
	as := &Assign{Left: p.Left, Right: &Minus{Left: p.Left, Right: p.Right}}
	as.Exec(env)
	return object.Null
}

func (p *MinusAssign) String() string {
	return fmt.Sprintf("%s -= %s", p.Left.String(), p.Right.String())
}

func (p *MinusAssign) expressionNode() {}
