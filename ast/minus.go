package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Minus struct {
	Left  Expression
	Right Expression
}

func (p *Minus) Eval(env *environment.Environment) object.Object {
	left := p.Left.Eval(env)
	right := p.Right.Eval(env)
	if leftObj, ok := left.(*object.Int); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Int{Val: leftObj.Val - rightObj.Val}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: float64(leftObj.Val) - rightObj.Val}
		}
		if right.Type() == object.TypeUndefined {
			return &object.Int{Val: leftObj.Val}
		}
		panic("invalid minus expression")
	}

	if leftObj, ok := left.(*object.Float); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Float{Val: leftObj.Val - float64(rightObj.Val)}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: leftObj.Val - rightObj.Val}
		}
		if right.Type() == object.TypeUndefined {
			return &object.Float{Val: leftObj.Val}
		}
		panic("invalid minus expression")
	}

	if left.Type() == object.TypeUndefined {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Int{Val: -rightObj.Val}
		}
		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: -rightObj.Val}
		}
	}
	panic("invalid minus expression")
}

func (p *Minus) String() string {
	return fmt.Sprintf("(%s - %s)", p.Left.String(), p.Right.String())
}

func (p *Minus) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Minus {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(p.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(p.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (p *Minus) expressionNode() {}
