package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Impl struct {
	Name    string
	Methods []*FnLiteral
}

func (im *Impl) Eval(env *environment.Environment) object.Object {
	obj, has := env.Get(im.Name)
	if !has {
		panic(fmt.Sprintf("NameError: no such impl %s", im.Name))
	}

	stu, ok := obj.(*StructLiteral)
	if !ok {
		panic(fmt.Sprintf("TypeError: impl %s is not a struct", im.Name))
	}

	for _, v := range im.Methods {
		stu.Methods[v.Name] = v
	}
	return object.Null
}

func (im *Impl) String() string {
	return ""
}

func (im *Impl) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Impl {\n")
	s.WriteString(strings.Repeat(". ", num+1) + fmt.Sprintf(" Name: %s\n", im.Name))
	s.WriteString(strings.Repeat(". ", num+1) + " Methods: {\n")
	for i, v := range im.Methods {
		s.WriteString(strings.Repeat(". ", num+2) + fmt.Sprintf(" %d: %s", i, v.AST(num+2)))
	}
	s.WriteString(strings.Repeat(". ", num+1) + " }\n")
	return s.String()
}

func (im *Impl) expressionNode() {}
