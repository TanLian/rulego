package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type String struct {
	Token token.Token
	Value *object.String
}

func (be *String) Eval(env *environment.Environment) object.Object {
	return be.Value
}

func (be *String) AST(num int) string {
	return fmt.Sprintf("*ast.String { %s }\n", string(be.Value.Val))
}

func (be *String) String() string {
	return be.Token.String()
}

func (be *String) expressionNode() {}
