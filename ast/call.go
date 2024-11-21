package ast

import (
	"fmt"
	"reflect"

	"github.com/tanlian/rulego/builtin"

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
		panic("expect Slice expression")
	}
	var in []reflect.Value
	var inObj []object.Object
	for _, vv := range args.Data {
		in = append(in, reflect.ValueOf(vv.Eval(env).GetValue()))
		inObj = append(inObj, vv.Eval(env))
	}

	if fn, ok := c.Left.(*Ident); ok {
		//fmt.Println("env: ", env)
		// 优先使用用户自定义的函数
		if obj, ok := env.Get(fn.Token.Value); ok {
			if fnLiteral, ok := obj.(*FnLiteralObj); ok {
				if len(fnLiteral.Args) != len(args.Data) {
					panic("the length of parameters is not equal")
				}

				childEnv := environment.New(env)
				for i := 0; i < len(fnLiteral.Args); i++ {
					childEnv.SetCurrent(fnLiteral.Args[i], args.Data[i].Eval(env))
				}
				return fnLiteral.Call(childEnv)
			}

			if fnObj, ok := obj.(*object.Fn); ok {
				res, err := fnObj.Call(inObj)
				if err != nil {
					panic(err)
				}
				return res
			}

			//if rule, ok := obj.(*Rule); ok {
			//	return rule.Call(env)
			//}
		}

		// 看看是否是内建函数
		if fnIn, ok := builtin.BuiltInFn[fn.Token.Value]; ok {
			var objs []object.Object
			for i := 0; i < len(args.Data); i++ {
				objs = append(objs, args.Data[i].Eval(env))
			}
			return fnIn.Call(objs...)
		}

		panic(fmt.Sprintf("%s is not a func", fn.Token.Value))
	}

	if dot, ok := c.Left.(*Dot); ok {
		dotLeft := dot.Left.Eval(env)
		dotRight, ok := dot.Right.(*Ident)
		if !ok {
			panic("invalid Call expression")
		}

		// 调用结构体的成员方法
		if obj, ok := dotLeft.(*object.Struct); ok {
			if out, ok := obj.Call(dotRight.Token.Value, inObj); ok {
				return out
			}
			return object.Null
		}
		v := reflect.ValueOf(dotLeft)
		method := v.MethodByName(dotRight.Token.Value)
		out, _ := c.MethodCall(method, in)
		return out
	}
	panic("invalid Call expression")
}

func (c *Call) String() string {
	return fmt.Sprintf("%s(%s)", c.Left.String(), c.Arguments.String())
}

func (c *Call) expressionNode() {}

func (c *Call) MethodCall(method reflect.Value, in []reflect.Value) (object.Object, bool) {
	if !method.IsValid() {
		return object.Null, false
	}

	out := method.Call(in)
	if len(out) == 0 {
		return object.Null, true
	}
	if len(out) == 1 {
		return object.ToObject(out[0]), true
	}
	res := &object.Slice{}
	for i := 0; i < len(out); i++ {
		res.Val = append(res.Val, object.ToObject(out[i]))
	}
	return res, true
}
