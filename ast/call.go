package ast

import (
	"fmt"
	"strings"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Call struct {
	Left      Expression // Ident or Dot expression
	Arguments Expression // Slice expression
}

func (c *Call) Eval(env *environment.Environment) object.Object {
	args, ok := c.Arguments.(*Slice)
	if !ok {
		panic("TypeError: expect slice expression")
	}

	var inObj []object.Object
	for _, v := range args.Data {
		inObj = append(inObj, v.Eval(env))
	}

	leftObj := c.Left.Eval(env)
	if leftObj == nil || leftObj.Type() == object.TypeUndefined {
		return object.Null
	}

	// 自定义的函数
	if fn, ok := leftObj.(*object.LiteralFn); ok {
		return fn.Call(inObj)
	}

	// 注入的函数
	if fnObj, ok := leftObj.(*object.InjectFn); ok {
		return fnObj.Call(inObj)
	}

	// 调用内部结构体（原生支持的结构体）的函数
	if callObj, ok := leftObj.(*object.StructCall); ok {
		if callObj.Self != nil && callObj.Self != object.Null {
			inObj = append([]object.Object{callObj.Self}, inObj...)
		}
		return callObj.Call(inObj)
	}

	// 调用注入的结构体的函数
	if fnObj, ok := leftObj.(*object.InjectStructCall); ok {
		return fnObj.Call(inObj)
	}

	// 调用object的扩展函数
	if fnObj, ok := leftObj.(*object.InnerStructMethodCall); ok {
		return fnObj.Call(inObj)
	}

	// 内置函数
	if builtInFn, ok := leftObj.(*object.BuiltinFn); ok {
		return builtInFn.Call(inObj)
	}

	panic(fmt.Sprintf("TypeError: expect function, got %s, %v", leftObj.Type(), leftObj.GetValue()))
}

func (c *Call) String() string {
	return fmt.Sprintf("%s(%s)", c.Left.String(), c.Arguments.String())
}

func (c *Call) AST(num int) string {
	var s strings.Builder
	s.WriteString("*ast.Call {\n")
	s.WriteString(strings.Repeat(". ", num+1) + " Left: ")
	s.WriteString(c.Left.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num+1) + " Arguments: ")
	s.WriteString(c.Arguments.AST(num + 1))
	s.WriteString(strings.Repeat(". ", num) + " }\n")
	return s.String()
}

func (c *Call) expressionNode() {}
