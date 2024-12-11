package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type Ident struct {
	Token token.Token
}

func (ie *Ident) Eval(env *environment.Environment) object.Object {
	//fmt.Println("ident eval: ", ie.Token.Value, " env: ", env)
	if obj, has := env.Get(ie.Token.Value); has {
		return obj
	}
	panic(fmt.Sprintf("NameError: name '%s' is not defined", ie.Token.Value))
	//return object.Null
}

func (ie *Ident) String() string {
	return ie.Token.String()
}

func (ie *Ident) AST(num int) string {
	return fmt.Sprintf("*ast.Ident { %s }\n", ie.Token.Value)
}

func (ie *Ident) expressionNode() {}
