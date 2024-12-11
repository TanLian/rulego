package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type FnLiteralObj struct {
	*object.Empty
	Name       string
	Args       []string
	Statements []Statement
}

func (fl *FnLiteralObj) Eval(env *environment.Environment) object.Object {
	env.Set(fl.Name, fl)
	return object.Null
}

func (fl *FnLiteralObj) Call(args []object.Object) object.Object {
	if len(fl.Args) != len(args) {
		panic("the length of parameters is not equal")
	}

	env := environment.New(environment.Root) // TODO: 这个env是从新定义还是从参数传递进来？
	for i := 0; i < len(fl.Args); i++ {
		env.SetCurrent(fl.Args[i], args[i])
	}
	for i, v := range fl.Statements {
		if i == len(fl.Statements)-1 {
			obj, _ := v.Exec(env)
			return obj
		}

		if obj, flg := v.Exec(env); flg&RETURN != 0 {
			return obj
		}
	}
	return object.Null
}

func (fl *FnLiteralObj) String() string {
	var s strings.Builder
	s.WriteString("fn ")
	s.WriteString(fl.Name)
	s.WriteString(fmt.Sprintf("(%s) {", strings.Join(fl.Args, ",")))
	for _, v := range fl.Statements {
		s.WriteString(v.String() + ";")
	}
	s.WriteString("}")
	return s.String()
}

func (fl *FnLiteralObj) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.FnLiteralObj {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Name: " + fl.Name + "\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Args: (" + strings.Join(fl.Args, ",") + ")\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Statements: {\n")
	for i, v := range fl.Statements {
		s.WriteString(strings.Repeat(". ", num+2) + strconv.Itoa(i) + ": " + v.AST(num+2))
	}
	s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (fl *FnLiteralObj) expressionNode() {}
