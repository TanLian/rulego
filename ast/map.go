package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Map struct {
	KV map[Expression]Expression
}

func (m *Map) Eval(env *environment.Environment) object.Object {
	res := make(map[any]object.Object)
	for k, v := range m.KV {
		res[k.Eval(env).GetValue()] = v.Eval(env)
	}
	return &object.Map{Val: res}
}

func (m *Map) String() string {
	var s strings.Builder
	s.WriteString("{")
	for k, v := range m.KV {
		s.WriteString(fmt.Sprintf("%s:%s,", k.String(), v.String()))
	}
	s.WriteString("}")
	return s.String()
}

func (m *Map) expressionNode() {}
