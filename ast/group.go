package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Group struct {
	Expr Expression
}

func (g *Group) Eval(env *environment.Environment) object.Object {
	return g.Expr.Eval(env)
}

func (g *Group) String() string {
	return fmt.Sprintf("(%s)", g.Expr.String())
}

func (g *Group) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Group {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Expr: ")
	s.WriteString(g.Expr.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (g *Group) expressionNode() {}
