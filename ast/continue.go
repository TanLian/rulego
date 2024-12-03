package ast

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Continue struct{}

func (b *Continue) statementNode() {}

func (b *Continue) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	return object.Null, CONTINUE
}

func (b *Continue) Eval(env *environment.Environment) object.Object {
	return object.Null
}

func (b *Continue) String() string {
	return "continue"
}

func (b *Continue) expressionNode() {}
