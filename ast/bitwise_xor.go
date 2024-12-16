package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type BitwiseXOR struct {
	Left  Expression
	Right Expression
}

func (a *BitwiseXOR) Eval(env *environment.Environment) object.Object {
	leftObj := a.Left.Eval(env)
	rightObj := a.Right.Eval(env)
	left, ok := leftObj.(*object.Int)
	if !ok {
		panic(fmt.Sprintf("TypeError: unsuported operand type for ^: %s and %s", leftObj.Type(), rightObj.Type()))
	}

	right, ok := rightObj.(*object.Int)
	if !ok {
		panic(fmt.Sprintf("TypeError: unsuported operand type for ^: %s and %s", leftObj.Type(), rightObj.Type()))
	}
	return &object.Int{Val: left.Val ^ right.Val}
}

func (a *BitwiseXOR) String() string {
	return fmt.Sprintf("%s ^ %s", a.Left.String(), a.Right.String())
}

func (a *BitwiseXOR) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.BitwiseXOR {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(a.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(a.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (a *BitwiseXOR) expressionNode() {}
