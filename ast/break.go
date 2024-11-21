package ast

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Break struct{}

func (b *Break) statementNode() {}

func (b *Break) Exec(env *environment.Environment) (object.Object, bool, bool) {
	return object.Null, false, true
}

func (b *Break) String() string {
	return "break"
}
