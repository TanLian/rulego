package ast

import (
	"fmt"

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
	Name   string
	Fields []string
}

func (rs *StructLiteral) Eval(env *environment.Environment) object.Object {
	env.Set(rs.Name, &object.RgStruct{
		Name:    rs.Name,
		Fields:  rs.Fields,
		Methods: make(map[string]object.Method),
		Values:  make(map[string]object.Object),
	})
	return object.Null
}

func (rs *StructLiteral) String() string {
	return ""
}

func (rs *StructLiteral) expressionNode() {}

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
	rs, ok := obj.(*object.RgStruct)
	if !ok {
		panic("not a struct")
	}
	rs = rs.Clone()

	if len(rsi.Values) > 0 {
		if len(rs.Fields) != len(rsi.Values) {
			panic("struct fields count not match")
		}
		for i, v := range rs.Fields {
			rs.SetFieldValue(v, rsi.Values[i].Eval(env))
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
			rs.SetFieldValue(field, v.Eval(env))
		}
	}
	return rs
}

func (rsi *RgStructInstantiate) String() string {
	return ""
}

func (rsi *RgStructInstantiate) expressionNode() {}
