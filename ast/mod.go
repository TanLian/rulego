package ast

import (
	"fmt"
	"math"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Mod struct {
	Left  Expression
	Right Expression
}

func (m *Mod) Eval(env *environment.Environment) object.Object {
	left := m.Left.Eval(env)
	right := m.Right.Eval(env)
	if leftObj, ok := left.(*object.Int); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Int{Val: leftObj.Val % rightObj.Val}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: math.Mod(float64(leftObj.Val), rightObj.Val)}
		}
		panic("invalid mod expression")
	}

	if leftObj, ok := left.(*object.Float); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Float{Val: math.Mod(leftObj.Val, float64(rightObj.Val))}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: math.Mod(leftObj.Val, rightObj.Val)}
		}
		panic("invalid mod expression")
	}
	panic("invalid mod expression")
}

func (m *Mod) String() string {
	return fmt.Sprintf("(%s % %s)", m.Left.String(), m.Right.String())
}

func (m *Mod) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Mod {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(m.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(m.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (m *Mod) expressionNode() {}
