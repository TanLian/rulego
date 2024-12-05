package ast

import (
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Block struct {
	States []Statement
}

func (b *Block) statementNode() {}

func (b *Block) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	childEnv := environment.New(env)
	for _, v := range b.States {
		if obj, flg := v.Exec(childEnv); flg&RETURN != 0 {
			return obj, flg
		}
	}
	return object.Null, 0
}

func (b *Block) String() string {
	var s strings.Builder
	s.WriteString("{")
	for _, v := range b.States {
		s.WriteString(v.String() + "\r\n")
	}
	s.WriteString("}")
	return s.String()
}
