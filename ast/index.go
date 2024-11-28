package ast

import (
	"fmt"

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
		return m.Val[key]
	}
	if s, ok := data.(*object.Slice); ok {
		if idx.End != nil {
			start, okStart := idx.Key.Eval(env).GetValue().(float64)
			end, okEnd := idx.End.Eval(env).GetValue().(float64)
			if !okStart || !okEnd {
				panic("invalid index expression")
			}
			return &object.Slice{Val: s.Val[int(start):int(end)]}
		}

		if index, ok := key.(float64); !ok {
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

func (idx *Index) expressionNode() {}
