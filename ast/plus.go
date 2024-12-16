package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Plus struct {
	Left  Expression // int, float, string
	Right Expression
}

func (p *Plus) Eval(env *environment.Environment) object.Object {
	left := p.Left.Eval(env)
	right := p.Right.Eval(env)
	if leftObj, ok := left.(*object.Int); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Int{Val: leftObj.Val + rightObj.Val}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: float64(leftObj.Val) + rightObj.Val}
		}
		if right.Type() == object.TypeUndefined {
			return &object.Int{Val: leftObj.Val}
		}
		panic(fmt.Sprintf("TypeError: unsuported operand type for +: %s and %s", left.Type(), right.Type()))
	}

	if leftObj, ok := left.(*object.Float); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Float{Val: leftObj.Val + float64(rightObj.Val)}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: leftObj.Val + rightObj.Val}
		}
		if right.Type() == object.TypeUndefined {
			return &object.Float{Val: leftObj.Val}
		}
		panic(fmt.Sprintf("TypeError: unsuported operand type for +: %s and %s", left.Type(), right.Type()))
	}

	if leftObj, ok := left.(*object.String); ok {
		if rightObj, ok := right.(*object.String); ok {
			return &object.String{Val: append(leftObj.Val, rightObj.Val...)}
		}
		if right.Type() == object.TypeUndefined {
			return &object.String{Val: append([]rune(nil), leftObj.Val...)}
		}
		panic(fmt.Sprintf("TypeError: unsuported operand type for +: %s and %s", left.Type(), right.Type()))
	}

	// slice + slice
	if leftObj, ok := left.(*object.Slice); ok {
		if rightObj, ok := right.(*object.Slice); ok {
			return &object.Slice{Val: append(leftObj.Val, rightObj.Val...)}
		}
	}

	if left.Type() == object.TypeUndefined {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Int{Val: rightObj.Val}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: rightObj.Val}
		}
		if rightObj, ok := right.(*object.String); ok {
			return &object.String{Val: append([]rune(nil), rightObj.Val...)}
		}
	}

	panic(fmt.Sprintf("TypeError: unsuported operand type for +: %s and %s", left.Type(), right.Type()))
}

func (p *Plus) String() string {
	return fmt.Sprintf("(%s + %s)", p.Left.String(), p.Right.String())
}

func (p *Plus) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Plus {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(p.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(p.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (p *Plus) expressionNode() {}
