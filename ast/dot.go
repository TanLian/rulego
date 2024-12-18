package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
a = [1,2,3];
a.Push(4);

type user struct {
	name string
}
u := &user{name: "name"}
u.name
*/

type Dot struct {
	Left  Expression // slice、struct
	Right Expression // ident
}

func (p *Dot) Eval(env *environment.Environment) object.Object {
	right, ok := p.Right.(*Ident)
	if !ok {
		return object.Null
	}

	left := p.Left.Eval(env)
	// 注入的结构体
	if structObj, ok := left.(*object.InjectStruct); ok {
		if res := structObj.GetFieldValue(right.Token.Value); res != object.Null {
			return res
		}

		return &object.InjectStructCall{
			MethodName: right.Token.Value,
			Struct:     structObj,
		}
	}

	// 内置结构体
	if structObj, ok := left.(*object.Struct); ok {
		if res := structObj.GetFieldValue(right.Token.Value); res != object.Null {
			return res
		}

		// 检查是不是结构体的方法
		if sObj, ok := env.Get(structObj.Name); ok {
			if obj, ok := sObj.(*StructLiteral); ok {
				if fn, ok := obj.Methods[right.Token.Value]; ok {
					return &StructCall{Self: structObj, Method: fn}
				}
			}
		}
	}

	v := reflect.ValueOf(left)
	if method := v.MethodByName(right.Token.Value); method.IsValid() {
		return &object.InnerStructMethodCall{
			Name:   right.Token.Value,
			Method: method,
		}
	}
	return object.Null
}

func (p *Dot) String() string {
	return fmt.Sprintf("(%s.%s)", p.Left.String(), p.Right.String())
}

func (p *Dot) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Dot {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(p.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Right: ")
	s.WriteString(p.Right.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (p *Dot) expressionNode() {}
