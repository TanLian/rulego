package ast

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type Number struct {
	Token token.Token
	Value *object.Float
}

func (ne *Number) Eval(env *environment.Environment) object.Object {
	return ne.Value
}

func (ne *Number) String() string {
	return ne.Token.String()
}

func (ne *Number) expressionNode() {}
