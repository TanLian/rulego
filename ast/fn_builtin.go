package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

var BuiltInFn = map[string]*FnBuiltin{
	"println": {Name: "println", Fn: func(args []object.Object) object.Object {
		var as []any
		for _, v := range args {
			as = append(as, v.GetValue())
		}
		fmt.Println(as...)
		return object.Null
	}},
	"assert_eq": {Name: "assert_eq", Fn: func(args []object.Object) object.Object {
		if len(args) != 2 {
			panic(fmt.Sprintf("expect 2 args, but got %d", len(args)))
		}
		v1 := args[0].GetValue()
		v2 := args[1].GetValue()
		if !reflect.DeepEqual(v1, v2) {
			panic(fmt.Sprintf("assert_eq failed: %v != %v", v1, v2))
		}
		return object.Null
	}},
}

type FnBuiltin struct {
	*object.Empty
	Name string
	Fn   func(args []object.Object) object.Object
}

func (fb *FnBuiltin) Eval(env *environment.Environment) object.Object {
	return fb
}

func (fb *FnBuiltin) String() string {
	return fmt.Sprintf("<builtin function: %s>", fb.Name)
}

func (fb *FnBuiltin) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.FnBuiltin {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Name: " + fb.Name + "\n")
	s.WriteString(strings.Repeat(". ", num) + "}")
	return s.String()
}

func (fb *FnBuiltin) expressionNode() {}

func (fb *FnBuiltin) Call(args []object.Object) object.Object {
	return fb.Fn(args)
}

func (fb *FnBuiltin) Type() object.Type {
	return object.TypeFn
}

func (fb *FnBuiltin) GetValue() any {
	return fb
}
