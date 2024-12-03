package ast

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type ExecFlag uint8

const (
	RETURN ExecFlag = 1 << iota
	BREAK
	CONTINUE
)

type Statement interface {
	statementNode()
	Exec(env *environment.Environment) (object.Object, ExecFlag)
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

func (es *ExpressionStatement) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	if state, ok := es.Expr.(Statement); ok {
		return state.Exec(env)
	}
	return es.Expr.Eval(env), 0
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) String() string {
	return es.Expr.String()
}
