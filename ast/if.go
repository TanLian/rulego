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
			for _, v := range exprStates.States {
				if obj, ret, _ := v.Exec(childEnv); ret {
					return obj
				}
			}
		}
	}

	childEnv := environment.New(env)
	for _, v := range f.Else {
		if obj, ret, _ := v.Exec(childEnv); ret {
			return obj
		}
	}
	return object.Null
}

func (f *If) Exec(env *environment.Environment) (object.Object, bool, bool) {
	for _, exprStates := range f.Ifs {
		if object.TransToBool(exprStates.Expr.Eval(env)) {
			childEnv := environment.New(env)
			for _, v := range exprStates.States {
				if obj, ret, brk := v.Exec(childEnv); ret || brk {
					return obj, ret, brk
				}
			}
			return object.Null, false, false
		}
	}

	childEnv := environment.New(env)
	for _, v := range f.Else {
		if obj, ret, brk := v.Exec(childEnv); ret || brk {
			return obj, ret, brk
		}
	}
	return object.Null, false, false
}

func (f *If) String() string {
	return ""
}