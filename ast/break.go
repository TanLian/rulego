package ast

import (
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Break struct{}

func (b *Break) statementNode() {}

func (b *Break) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	return object.Null, BREAK
}

func (b *Break) Eval(env *environment.Environment) object.Object {
	return object.Null
}

func (b *Break) String() string {
	return "break"
}

func (b *Break) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Break {\n")
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (b *Break) expressionNode() {}
