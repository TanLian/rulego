package ast

import (
	"fmt"
	"math"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Power struct {
	Left  Expression
	Right Expression
}

func (m *Power) Eval(env *environment.Environment) object.Object {
	left := m.Left.Eval(env)
	right := m.Right.Eval(env)
	if leftObj, ok := left.(*object.Int); ok {
		if rightObj, ok := right.(*object.Int); ok {
			if rightObj.Val > 0 {
				return &object.Int{Val: int64(math.Pow(float64(leftObj.Val), float64(rightObj.Val)))}
			}
			return &object.Float{Val: math.Pow(float64(leftObj.Val), float64(rightObj.Val))}
		}

		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: math.Pow(float64(leftObj.Val), rightObj.Val)}
		}
	}

	if leftObj, ok := left.(*object.Float); ok {
		if rightObj, ok := right.(*object.Int); ok {
			return &object.Float{Val: math.Pow(leftObj.Val, float64(rightObj.Val))}
		}

		if rightObj, ok := right.(*object.Float); ok {
			return &object.Float{Val: math.Pow(leftObj.Val, rightObj.Val)}
		}
	}
	panic("invalid power expression")
}

func (m *Power) String() string {
	return fmt.Sprintf("(%s ** %s)", m.Left.String(), m.Right.String())
}

func (m *Power) expressionNode() {}
