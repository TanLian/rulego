package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type DivideAssign struct {
	Left  Expression
	Right Expression
}

func (p *DivideAssign) Eval(env *environment.Environment) object.Object {
	as := &Assign{Left: p.Left, Right: &Divide{Left: p.Left, Right: p.Right}}
	as.Exec(env)
	return object.Null
}

func (p *DivideAssign) String() string {
	return fmt.Sprintf("%s /= %s", p.Left.String(), p.Right.String())
}

func (p *DivideAssign) expressionNode() {}
