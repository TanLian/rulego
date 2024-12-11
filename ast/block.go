package ast

import (
	"strconv"
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

func (b *Block) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Block {\n")
	for i, v := range b.States {
		s.WriteString(strings.Repeat(". ", num+1) + " " + strconv.Itoa(i) + ": " + v.AST(num+1) + "\n")
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}
