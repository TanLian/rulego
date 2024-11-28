package ast

import (
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Slice struct {
	Data []Expression
}

func (se *Slice) Eval(env *environment.Environment) object.Object {
	var res []any
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

func (se *Slice) expressionNode() {}
