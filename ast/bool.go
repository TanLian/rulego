package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Bool struct {
	Value *object.Bool
}

func (b *Bool) Eval(env *environment.Environment) object.Object {
	return b.Value
}

func (b *Bool) String() string {
	if b.Value == nil {
		return ""
	}
	if b.Value.Val {
		return "true"
	}
	return "false"
}

func (b *Bool) AST(num int) string {
	return fmt.Sprintf("*ast.Bool { %s }\n", b)
}

func (b *Bool) expressionNode() {}
