package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
struct person {
	age,
	name,
}
*/

type StructLiteral struct {
	*object.Empty
	Name    string
	Fields  []string
	Methods map[string]*FnLiteral
}

func (rs *StructLiteral) Eval(env *environment.Environment) object.Object {
	env.Set(rs.Name, rs)
	return rs
}

func (rs *StructLiteral) String() string {
	return ""
}

func (rs *StructLiteral) Type() object.Type {
	return object.TypeStruct
}

func (rs *StructLiteral) expressionNode() {}

func (rs *StructLiteral) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.StructLiteral {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Name: " + rs.Name + "\n")
	for i, v := range rs.Fields {
		s.WriteString(strings.Repeat(". ", num+1) + fmt.Sprintf(" Field %d: ", i) + v + "\n")
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (rs *StructLiteral) CheckFieldExist(field string) bool {
	for _, v := range rs.Fields {
		if field == v {
			return true
		}
	}
	return false
}

/*
p1 = person{1,"leo"};
p2 = person{age: 2};
*/

type RgStructInstantiate struct {
	Ident  Expression
	Values []Expression
	KV     map[Expression]Expression
}

func (rsi *RgStructInstantiate) Eval(env *environment.Environment) object.Object {
	obj := rsi.Ident.Eval(env)
	rs, ok := obj.(*StructLiteral)
	if !ok {
		panic("not a struct")
	}

	res := &object.Struct{
		Name:   rs.Name,
		Values: make(map[string]object.Object),
	}

	if len(rsi.Values) > 0 {
		if len(rs.Fields) != len(rsi.Values) {
			panic("struct fields count not match")
		}
		for i, v := range rs.Fields {
			res.SetFieldValue(v, rsi.Values[i].Eval(env))
		}
	}

	if len(rsi.KV) > 0 {
		for k, v := range rsi.KV {
			ident, ok := k.(*Ident)
			if !ok {
				panic("key must be ident")
			}
			field := ident.Token.Value

			if !rs.CheckFieldExist(field) {
				panic(fmt.Sprintf("field %s not exist in struct %s", field, rs.Name))
			}
			res.SetFieldValue(field, v.Eval(env))
		}
	}
	return res
}

func (rsi *RgStructInstantiate) String() string {
	return ""
}

func (rsi *RgStructInstantiate) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.RgStructInstantiate {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Ident: " + rsi.Ident.AST(num+1))
	if len(rsi.Values) > 0 {
		for i, v := range rsi.Values {
			s.WriteString(strings.Repeat(". ", num+1) + fmt.Sprintf(" Values %d: ", i) + v.AST(num+1))
		}
	} else {
		for k, v := range rsi.KV {
			s.WriteString(strings.Repeat(". ", num+1) + fmt.Sprintf(" Key: %s", k.AST(num+1)))
			s.WriteString(strings.Repeat(". ", num+1) + fmt.Sprintf(" Value: %s", v.AST(num+1)))
		}
	}
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (rsi *RgStructInstantiate) expressionNode() {}
