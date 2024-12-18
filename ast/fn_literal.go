package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type FnLiteral struct {
	*object.Empty
	Name   string
	Args   []string
	States []Statement
}

func (fl *FnLiteral) Eval(env *environment.Environment) object.Object {
	env.SetCurrent(fl.Name, fl)
	return fl
}

func (fl *FnLiteral) Call(args []object.Object) object.Object {
	env := environment.New(environment.Root)
	for i := 0; i < len(fl.Args); i++ {
		env.SetCurrent(fl.Args[i], args[i])
	}
	for i, v := range fl.States {
		if i == len(fl.States)-1 {
			obj, _ := v.Exec(env)
			return obj
		}

		if obj, flg := v.Exec(env); flg&RETURN != 0 {
			return obj
		}
	}
	return object.Null
}

func (fl *FnLiteral) Type() object.Type {
	return object.TypeFnLiteral
}

func (fl *FnLiteral) String() string {
	if fl == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString("fn ")
	s.WriteString(fl.Name)
	s.WriteString(fmt.Sprintf("(%s) {", strings.Join(fl.Args, ",")))
	for _, v := range fl.States {
		s.WriteString(v.String() + ";")
	}
	s.WriteString("}")
	return s.String()
}

func (fl *FnLiteral) AST(num int) string {
	if fl == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString("*ast.LiteralFn {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Name: " + fl.Name + "\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Args: (" + strings.Join(fl.Args, ",") + ")\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Statements: {\n")
	for i, v := range fl.States {
		s.WriteString(strings.Repeat(". ", num+2) + strconv.Itoa(i) + ": " + v.AST(num+2))
	}
	s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (fl *FnLiteral) expressionNode() {}
