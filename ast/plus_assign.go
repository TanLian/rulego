package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type PlusAssign struct {
	Left  Expression
	Right Expression
}

func (p *PlusAssign) Eval(env *environment.Environment) object.Object {
	as := &Assign{Left: p.Left, Right: &Plus{Left: p.Left, Right: p.Right}}
	as.Exec(env)
	return object.Null
}

func (p *PlusAssign) String() string {
	return fmt.Sprintf("%s += %s", p.Left.String(), p.Right.String())
}

func (p *PlusAssign) expressionNode() {}
