package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Slice struct {
	Data     []Expression
	InitExpr Expression
	LenExpr  Expression
}

func (se *Slice) Eval(env *environment.Environment) object.Object {
	var res []any
	if se.InitExpr != nil && se.LenExpr != nil {
		lenObj := se.LenExpr.Eval(env)
		if lenFloat, ok := lenObj.(*object.Float); ok {
			for i := int64(0); i < int64(lenFloat.Val); i++ {
				res = append(res, se.InitExpr.Eval(env).GetValue())
			}
			return &object.Slice{Val: res}
		}
		if lenInt, ok := lenObj.(*object.Int); ok {
			for i := int64(0); i < lenInt.Val; i++ {
				res = append(res, se.InitExpr.Eval(env).GetValue())
			}
			return &object.Slice{Val: res}
		}
		panic("invalid slice")
	}
	for _, v := range se.Data {
		res = append(res, v.Eval(env).GetValue())
	}
	return &object.Slice{Val: res}
}

func (se *Slice) String() string {
	var s strings.Builder
	s.WriteString("[")
	var str []string
	for _, v := range se.Data {
		str = append(str, v.String())
	}
	s.WriteString(strings.Join(str, ","))
	s.WriteString("]")
	return s.String()
}

func (se *Slice) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Slice {\n")
	if len(se.Data) > 0 {
		for i, v := range se.Data {
			s.WriteString(strings.Repeat(". ", num+1) + fmt.Sprint(i) + ": " + v.AST(num+1))
		}
	} else if se.InitExpr != nil && se.LenExpr != nil {
		s.WriteString(strings.Repeat(". ", num+1) + " InitExpr: ")
		s.WriteString(se.InitExpr.AST(num + 1))
		s.WriteString(strings.Repeat(". ", num+1) + " LenExpr: ")
		s.WriteString(se.LenExpr.AST(num + 1))
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (se *Slice) expressionNode() {}
