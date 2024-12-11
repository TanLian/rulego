package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type Number struct {
	Token token.Token
	Value object.Object
}

func (ne *Number) Eval(env *environment.Environment) object.Object {
	return ne.Value
}

func (ne *Number) String() string {
	return ne.Token.String()
}

func (ne *Number) AST(num int) string {
	return fmt.Sprintf("*ast.Number { %d }\n", ne.Value.GetValue())
}

func (ne *Number) expressionNode() {}
