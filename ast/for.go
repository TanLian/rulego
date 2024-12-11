package ast

import (
	"strconv"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
for i = 0; i < 10; i+=1 {
	state1;
	state2;
}
*/

type For struct {
	Initial    Statement
	Condition  Expression
	Post       Statement
	Statements []Statement
}

func (f *For) statementNode() {}

func (f *For) expressionNode() {}

func (f *For) Eval(env *environment.Environment) object.Object {
	res, _ := f.Exec(env)
	return res
}

func (f *For) Exec(env *environment.Environment) (object.Object, ExecFlag) {
	childEnv := environment.New(env)
	f.Initial.Exec(childEnv)
	for object.TransToBool(f.Condition.Eval(childEnv)) {
		for _, v := range f.Statements {
			obj, flg := v.Exec(childEnv)
			if flg&RETURN != 0 {
				return obj, flg
			}
			if flg&BREAK != 0 {
				goto end
			}
			if flg&CONTINUE != 0 {
				break
			}
		}
		f.Post.Exec(childEnv)
	}
end:
	return object.Null, 0
}

func (f *For) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.For {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Initial: ")
	s.WriteString(f.Initial.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Condition: ")
	s.WriteString(f.Condition.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Post: ")
	s.WriteString(f.Post.AST(num + 1))
	for i, v := range f.Statements {
		s.WriteString(strings.Repeat(". ", num+1) + " Statements[" + strconv.Itoa(i) + "]: " + v.AST(0))
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (f *For) String() string {
	return ""
}
