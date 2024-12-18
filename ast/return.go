package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Return struct {
	Expr Expression
}

func (r *Return) statementNode() {}

func (r *Return) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	if r.Expr == nil {
		return object.Null, RETURN
	}
	return r.Expr.Eval(env), RETURN
}

func (r *Return) Eval(env *environment.Environment) object.Object {
	return r.Expr.Eval(env)
}

func (r *Return) String() string {
	return fmt.Sprintf("return %s;", r.Expr.String())
}

func (r *Return) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Return {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Expr: ")
	s.WriteString(r.Expr.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (r *Return) expressionNode() {}
