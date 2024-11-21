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

func (r *Return) Exec(env *environment.Environment) (object.Object, bool, bool) {
	return r.Expr.Eval(env), true, false
}

func (r *Return) String() string {
	return fmt.Sprintf("return %s;", r.Expr.String())
}
