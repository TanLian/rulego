package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Index struct {
	Data Expression // 表达式的求值结果是 map 或 slice
	Key  Expression
	End  Expression
}

func (idx *Index) Eval(env *environment.Environment) object.Object {
	data := idx.Data.Eval(env)
	key := idx.Key.Eval(env).GetValue()
	if m, ok := data.(*object.Map); ok {
		return object.New(m.Val[key])
	}

	if m, ok := data.(*object.MapPointer); ok {
		return object.New(m.GetField(key))
	}

	if s, ok := data.(*object.Slice); ok {
		if idx.End != nil {
			startObj := idx.Key.Eval(env)
			endObj := idx.End.Eval(env)
			var start, end int
			switch obj := startObj.(type) {
			case *object.Int:
				start = int(obj.Val)
			case *object.Float:
				start = int(obj.Val)
			default:
				panic("invalid index expression")
			}

			switch obj := endObj.(type) {
			case *object.Int:
				end = int(obj.Val)
			case *object.Float:
				end = int(obj.Val)
			default:
				panic("invalid index expression")
			}
			return &object.Slice{Val: s.Val[start:end]}
		}

		if index, ok := key.(int64); !ok {
			panic("invalid index expression")
		} else {
			return object.New(s.Val[int(index)])
		}
	}

	if s, ok := data.(*object.String); ok {
		if idx.End != nil {
			start, okStart := idx.Key.Eval(env).GetValue().(float64)
			end, okEnd := idx.End.Eval(env).GetValue().(float64)
			if !okStart || !okEnd {
				panic("invalid index expression")
			}
			return &object.String{Val: s.Val[int(start):int(end)]}
		}

		if index, ok := key.(float64); !ok {
			panic("invalid index expression")
		} else {
			str := []rune(s.Val)
			return &object.Rune{Val: str[int(index)]}
		}
	}
	panic("invalid index expression")
}

func (idx *Index) String() string {
	return fmt.Sprintf("%s[%s]", idx.Data.String(), idx.Key.String())
}

func (idx *Index) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Index {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Data: ")
	s.WriteString(idx.Data.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Key: ")
	s.WriteString(idx.Key.AST(num + 1))
	if idx.End != nil {
		s.WriteString(strings.Repeat(". ", num+1) + " End: ")
		s.WriteString(idx.End.AST(num + 1))
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (idx *Index) expressionNode() {}
