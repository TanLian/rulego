package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Return struct {
	Expr Expression
}

func (r *Return) statementNode() {}

func (r *Return) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	return r.Expr.Eval(env), RETURN
}

func (r *Return) Eval(env *environment.Environment) object.Object {
	return r.Expr.Eval(env)
}

func (r *Return) String() string {
	return fmt.Sprintf("return %s;", r.Expr.String())
}

func (r *Return) expressionNode() {}
