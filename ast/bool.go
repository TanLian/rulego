package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type Bool struct {
	Token token.Token
	Value *object.Bool
}

func (be *Bool) Eval(env *environment.Environment) object.Object {
	return be.Value
}

func (be *Bool) String() string {
	return be.Token.String()
}

func (be *Bool) AST(num int) string {
	return fmt.Sprintf("*ast.Bool { %s }\n", be.Token.Value)
}

func (be *Bool) expressionNode() {}
