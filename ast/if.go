package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
if 既可以作为语句，又可以作为表达式

作为语句时：
if expr {
	state1;
	state2;
} else if expr {
	state1;
	state2;
} else {
	state1;
	state2;
}

作为表达式时：
b = 1;
a = if b > 0 {1} else {-1};
*/

type If struct {
	Ifs  []ExprStates
	Else []Statement
}

type ExprStates struct {
	Expr   Expression
	States []Statement
}

func (f *If) expressionNode() {}

func (f *If) statementNode() {}

func (f *If) Eval(env *environment.Environment) object.Object {
	res, _ := f.Exec(env)
	return res
}

func (f *If) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	for _, exprStates := range f.Ifs {
		if object.TransToBool(exprStates.Expr.Eval(env)) {
			childEnv := environment.New(env)
			for i, v := range exprStates.States {
				if i == len(exprStates.States)-1 {
					return v.Exec(childEnv)
				}

				if obj, flg := v.Exec(childEnv); flg&RETURN != 0 || flg&BREAK != 0 || flg&CONTINUE != 0 {
					return obj, flg
				}
			}
			return object.Null, 0
		}
	}

	childEnv := environment.New(env)
	for i, v := range f.Else {
		if i == len(f.Else)-1 {
			return v.Exec(childEnv)
		}

		if obj, flg := v.Exec(childEnv); flg&RETURN != 0 || flg&BREAK != 0 || flg&CONTINUE != 0 {
			return obj, flg
		}
	}
	return object.Null, 0
}

func (f *If) AST(num int) string {
	var s strings.Builder
	s.WriteString(" *ast.If {\n")
	for _, v := range f.Ifs {
		s.WriteString(strings.Repeat(". ", num+1) + " Condition: " + v.Expr.AST(num+1))
		s.WriteString(strings.Repeat(". ", num+1) + " Statements: {\n")
		for i, vv := range v.States {
			s.WriteString(strings.Repeat(". ", num+2) + fmt.Sprintf(" %d: %s", i, vv.AST(num+2)))
		}
		s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	}
	if len(f.Else) > 0 {
		s.WriteString(strings.Repeat(". ", num+1) + " Else: {\n")
		s.WriteString(strings.Repeat(". ", num+1) + " Statements: \n")
		for i, v := range f.Else {
			s.WriteString(strings.Repeat(". ", num+2) + fmt.Sprintf(" %d: %s", i, v.AST(num+2)))
		}
		s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (f *If) String() string {
	return ""
}
