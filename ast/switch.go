package ast

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
switch expr {
case expr1:
	state1;
case expr2:
	state2;
default:
	state3;
}
*/

type Switch struct {
	Expr    Expression
	Cases   []ExprStates
	Default []Statement
}

func (s *Switch) expressionNode() {}

func (s *Switch) statementNode() {}

func (s *Switch) Eval(env *environment.Environment) object.Object {
	res, _ := s.Exec(env)
	return res
}

func (s *Switch) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	val := s.Expr.Eval(env).GetValue()
	for _, v := range s.Cases {
		if v.Expr.Eval(env).GetValue() == val {
			for _, vv := range v.States {
				if res, flg := vv.Exec(env); flg&RETURN != 0 {
					return res, RETURN
				}
			}
			return object.Null, 0
		}
	}

	// 执行default
	for _, v := range s.Default {
		if res, flg := v.Exec(env); flg&RETURN != 0 {
			return res, RETURN
		}
	}
	return object.Null, 0
}

func (s *Switch) String() string {
	return ""
}
