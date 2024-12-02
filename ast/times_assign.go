package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type TimesAssign struct {
	Left  Expression
	Right Expression
}

func (p *TimesAssign) Eval(env *environment.Environment) object.Object {
	as := &Assign{Left: p.Left, Right: &Times{Left: p.Left, Right: p.Right}}
	as.Exec(env)
	return object.Null
}

func (p *TimesAssign) String() string {
	return fmt.Sprintf("%s *= %s", p.Left.String(), p.Right.String())
}

func (p *TimesAssign) expressionNode() {}
