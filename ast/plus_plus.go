package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type PlusPlus struct {
	Left Expression
}

func (p *PlusPlus) Eval(env *environment.Environment) object.Object {
	as := &Assign{Left: p.Left, Right: &Plus{Left: p.Left, Right: &Number{Value: &object.Float{Val: 1}}}}
	as.Exec(env)
	return object.Null
}

func (p *PlusPlus) Exec(env *environment.Environment) (object.Object, bool, bool) {
	as := &Assign{Left: p.Left, Right: &Plus{Left: p.Left, Right: &Number{Value: &object.Float{Val: 1}}}}
	as.Exec(env)
	return object.Null, false, false
}

func (p *PlusPlus) String() string {
	return fmt.Sprintf("%s++", p.Left.String())
}

func (p *PlusPlus) expressionNode() {}

func (p *PlusPlus) statementNode() {}
