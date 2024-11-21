package ast

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Statement interface {
	statementNode()
	Exec(env *environment.Environment) (object.Object, bool, bool)
	String() string
}

type Expression interface {
	expressionNode()
	Eval(env *environment.Environment) object.Object
	String() string
}

type ExpressionStatement struct {
	Expr Expression
}

func (es *ExpressionStatement) Exec(env *environment.Environment) (object.Object, bool, bool) {
	if state, ok := es.Expr.(Statement); ok {
		return state.Exec(env)
	}
	return es.Expr.Eval(env), false, false
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) String() string {
	return es.Expr.String()
}
