package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type FnLiteral struct {
	Obj *object.LiteralFn
}

func (fl *FnLiteral) Eval(env *environment.Environment) object.Object {
	if fl.Obj != nil {
		env.Set(fl.Obj.Name, fl.Obj)
	}
	return fl.Obj
}

func (fl *FnLiteral) String() string {
	var s strings.Builder
	s.WriteString("fn ")
	s.WriteString(fl.Obj.Name)
	s.WriteString(fmt.Sprintf("(%s) {", strings.Join(fl.Obj.Args, ",")))
	states := fl.Obj.Block.(*FnLiteralBlock).States
	for _, v := range states {
		s.WriteString(v.String() + ";")
	}
	s.WriteString("}")
	return s.String()
}

func (fl *FnLiteral) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.LiteralFn {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Name: " + fl.Obj.Name + "\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Args: (" + strings.Join(fl.Obj.Args, ",") + ")\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Statements: {\n")
	states := fl.Obj.Block.(*FnLiteralBlock).States
	for i, v := range states {
		s.WriteString(strings.Repeat(". ", num+2) + strconv.Itoa(i) + ": " + v.AST(num+2))
	}
	s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (fl *FnLiteral) expressionNode() {}

type FnLiteralBlock struct {
	Args   []string
	States []Statement
}

func (fl *FnLiteralBlock) Call(args []object.Object) object.Object {
	if len(fl.Args) != len(args) {
		panic("the length of parameters is not equal")
	}

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
