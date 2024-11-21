package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type Compare struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (c *Compare) Eval(env *environment.Environment) object.Object {
	if leftInt, ok := c.Left.Eval(env).(*object.Int); ok {
		if rightInt, ok := c.Right.Eval(env).(*object.Int); ok {
			return &object.Bool{Val: compare(leftInt.Val, rightInt.Val, c.Token.Type)}
		}
		if rightFloat, ok := c.Right.Eval(env).(*object.Float); ok {
			return &object.Bool{Val: compare(float64(leftInt.Val), rightFloat.Val, c.Token.Type)}
		}
		panic("invalid compare expression")
	}

	if leftFloat, ok := c.Left.Eval(env).(*object.Float); ok {
		if rightInt, ok := c.Right.Eval(env).(*object.Int); ok {
			return &object.Bool{Val: compare(leftFloat.Val, float64(rightInt.Val), c.Token.Type)}
		}
		if rightFloat, ok := c.Right.Eval(env).(*object.Float); ok {
			return &object.Bool{Val: compare(leftFloat.Val, rightFloat.Val, c.Token.Type)}
		}
		panic("invalid compare expression")
	}

	if leftStr, ok := c.Left.Eval(env).(*object.String); ok {
		if rightStr, ok := c.Right.Eval(env).(*object.String); ok {
			return &object.Bool{Val: compare(string(leftStr.Val), string(rightStr.Val), c.Token.Type)}
		}
		if rightRune, ok := c.Right.Eval(env).(*object.Rune); ok {
			return &object.Bool{Val: compare(string(leftStr.Val), string(rightRune.Val), c.Token.Type)}
		}
		panic("invalid compare expression")
	}

	if leftRune, ok := c.Left.Eval(env).(*object.Rune); ok {
		if rightRune, ok := c.Right.Eval(env).(*object.Rune); ok {
			return &object.Bool{Val: compare(leftRune.Val, rightRune.Val, c.Token.Type)}
		}
		if rightStr, ok := c.Right.Eval(env).(*object.String); ok {
			return &object.Bool{Val: compare(string(leftRune.Val), string(rightStr.Val), c.Token.Type)}
		}
		panic("invalid compare expression")
	}
	panic("invalid compare expression")
}

func (c *Compare) String() string {
	switch c.Token.Type {
	case token.GREATER:
		return fmt.Sprintf("%s > %s", c.Left.String(), c.Right.String())
	case token.GREATER_EQUAL:
		return fmt.Sprintf("%s >= %s", c.Left.String(), c.Right.String())
	case token.EQUAL:
		return fmt.Sprintf("%s == %s", c.Left.String(), c.Right.String())
	case token.LESS:
		return fmt.Sprintf("%s < %s", c.Left.String(), c.Right.String())
	case token.LESS_EQUAL:
		return fmt.Sprintf("%s <= %s", c.Left.String(), c.Right.String())
	case token.NOT_EQUAL:
		return fmt.Sprintf("%s != %s", c.Left.String(), c.Right.String())
	default:
		return ""
	}
}

func (c *Compare) expressionNode() {}

// 定义一个类型约束，要求类型必须实现有序比较
type Ordered interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | string
}

func compare[T Ordered](a, b T, tokenType token.TokenType) bool {
	switch tokenType {
	case token.GREATER:
		return a > b
	case token.GREATER_EQUAL:
		return a >= b
	case token.EQUAL:
		return a == b
	case token.LESS:
		return a < b
	case token.LESS_EQUAL:
		return a <= b
	case token.NOT_EQUAL:
		return a != b
	default:
		return false
	}
}
