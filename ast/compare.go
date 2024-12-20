package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type CompareFlag uint8

const (
	_                   CompareFlag = iota
	CompareGreaterThan              // >
	CompareGreaterEqual             // >=
	CompareEqual                    // ==
	CompareNotEqual                 // ==
	CompareLessThan                 // <
	CompareLessEqual                // <=
)

type Compare struct {
	Flag  CompareFlag
	Left  Expression
	Right Expression
}

func (c *Compare) Eval(env *environment.Environment) object.Object {
	leftObj := c.Left.Eval(env)
	rightObj := c.Right.Eval(env)
	if leftInt, ok := leftObj.(*object.Int); ok {
		if rightInt, ok := rightObj.(*object.Int); ok {
			return &object.Bool{Val: compare(leftInt.Val, rightInt.Val, c.Flag)}
		}
		if rightFloat, ok := rightObj.(*object.Float); ok {
			return &object.Bool{Val: compare(float64(leftInt.Val), rightFloat.Val, c.Flag)}
		}
	}

	if leftFloat, ok := leftObj.(*object.Float); ok {
		if rightInt, ok := rightObj.(*object.Int); ok {
			return &object.Bool{Val: compare(leftFloat.Val, float64(rightInt.Val), c.Flag)}
		}
		if rightFloat, ok := rightObj.(*object.Float); ok {
			return &object.Bool{Val: compare(leftFloat.Val, rightFloat.Val, c.Flag)}
		}
	}

	if leftStr, ok := leftObj.(*object.String); ok {
		if rightStr, ok := rightObj.(*object.String); ok {
			return &object.Bool{Val: compare(string(leftStr.Val), string(rightStr.Val), c.Flag)}
		}
		if rightRune, ok := rightObj.(*object.Rune); ok {
			return &object.Bool{Val: compare(string(leftStr.Val), string(rightRune.Val), c.Flag)}
		}
	}

	if leftRune, ok := leftObj.(*object.Rune); ok {
		if rightRune, ok := rightObj.(*object.Rune); ok {
			return &object.Bool{Val: compare(leftRune.Val, rightRune.Val, c.Flag)}
		}
		if rightStr, ok := rightObj.(*object.String); ok {
			return &object.Bool{Val: compare(string(leftRune.Val), string(rightStr.Val), c.Flag)}
		}
	}
	panic(fmt.Sprintf("TypeError: unsupported operand type(s) for %s: '%s' and '%s'", c.fmtCompareFlag(), leftObj.Type(), rightObj.Type()))
}

func (c *Compare) String() string {
	switch c.Flag {
	case CompareGreaterThan:
		return fmt.Sprintf("%s > %s", c.Left.String(), c.Right.String())
	case CompareGreaterEqual:
		return fmt.Sprintf("%s >= %s", c.Left.String(), c.Right.String())
	case CompareEqual:
		return fmt.Sprintf("%s == %s", c.Left.String(), c.Right.String())
	case CompareLessThan:
		return fmt.Sprintf("%s < %s", c.Left.String(), c.Right.String())
	case CompareLessEqual:
		return fmt.Sprintf("%s <= %s", c.Left.String(), c.Right.String())
	case CompareNotEqual:
		return fmt.Sprintf("%s != %s", c.Left.String(), c.Right.String())
	default:
		return ""
	}
}

func (c *Compare) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Compare {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(c.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(c.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Flag: " + c.fmtCompareFlag() + "\n")
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (c *Compare) fmtCompareFlag() string {
	switch c.Flag {
	case CompareGreaterThan:
		return ">"
	case CompareGreaterEqual:
		return ">="
	case CompareEqual:
		return "=="
	case CompareLessThan:
		return "<"
	case CompareLessEqual:
		return "<="
	case CompareNotEqual:
		return "!="
	default:
		return ""
	}
}

func (c *Compare) expressionNode() {}

// Ordered 定义一个类型约束，要求类型必须实现有序比较
type Ordered interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 | string
}

func compare[T Ordered](a, b T, flg CompareFlag) bool {
	switch flg {
	case CompareGreaterThan:
		return a > b
	case CompareGreaterEqual:
		return a >= b
	case CompareEqual:
		return a == b
	case CompareLessThan:
		return a < b
	case CompareLessEqual:
		return a <= b
	case CompareNotEqual:
		return a != b
	default:
		return false
	}
}
