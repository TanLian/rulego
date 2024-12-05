package ast

import (
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
	for _, exprStates := range f.Ifs {
		if object.TransToBool(exprStates.Expr.Eval(env)) {
			childEnv := environment.New(env)
			for i, v := range exprStates.States {
				// the last statement, if it is an expression
				if i == len(exprStates.States)-1 {
					if expr, ok := v.(*ExpressionStatement); ok {
						return expr.Expr.Eval(childEnv)
					}
				}

				if obj, flg := v.Exec(childEnv); flg&RETURN != 0 {
					return obj
				}
			}
		}
	}

	childEnv := environment.New(env)
	for i, v := range f.Else {
		if i == len(f.Else)-1 {
			if expr, ok := v.(*ExpressionStatement); ok {
				return expr.Expr.Eval(childEnv)
			}
		}

		if obj, flg := v.Exec(childEnv); flg&RETURN != 0 {
			return obj
		}
	}
	return object.Null
}

func (f *If) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	for _, exprStates := range f.Ifs {
		if object.TransToBool(exprStates.Expr.Eval(env)) {
			childEnv := environment.New(env)
			for _, v := range exprStates.States {
				obj, flg := v.Exec(childEnv)
				if flg&RETURN != 0 || flg&BREAK != 0 || flg&CONTINUE != 0 {
					return obj, flg
				}
			}
			return object.Null, 0
		}
	}

	childEnv := environment.New(env)
	for _, v := range f.Else {
		obj, flg := v.Exec(childEnv)
		if flg&RETURN != 0 || flg&BREAK != 0 || flg&CONTINUE != 0 {
			return obj, flg
		}
	}
	return object.Null, 0
}

func (f *If) String() string {
	return ""
}
