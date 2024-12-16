package object

import (
	"reflect"
)

/*
TODO: 有以下几类函数调用
1. 自己定义的函数 Fn
fn aaa(a) {}
aaa([1,2,3]);

2. 内置函数 BuiltinFn
println();

3. 注入的函数 InjectFn

4. 自己定义的结构体方法 Caller
struct person {
	age,
	name,
}

impl person {
	fn get_name(self) {
		self.name
	}

	fn set_name(self, name) {
		self.name = name;
	}
}

p1 = person{1,"leo"};
println(p1.get_name());

5. object的内置方法 BuiltinMethod
a = [1,2,3];
a.Push(4);

6. 注入的结构体的方法 InjectMethod
*/

type StructCall struct {
	Self   Object
	Method Caller
}

func (fn *StructCall) object() {}

func (fn *StructCall) Type() Type {
	return TypeCall
}

func (fn *StructCall) GetValue() any {
	return Null
}

func (fn *StructCall) Call(args []Object) Object {
	return fn.Method.Call(args)
}

type InjectStructCall struct {
	MethodName string
	Struct     *InjectStruct
}

func (fn *InjectStructCall) object() {}

func (fn *InjectStructCall) Type() Type {
	return TypeCall
}

func (fn *InjectStructCall) GetValue() any {
	return Null
}

func (fn *InjectStructCall) Call(args []Object) Object {
	m := fn.Struct.Value.MethodByName(fn.MethodName)
	if !m.IsValid() {
		return Null
	}

	types, ok := fn.Struct.lazyLoadMethodType(fn.MethodName)
	if !ok {
		return Null
	}

	// 参数数量不匹配
	if len(types) != len(args) {
		return Null
	}

	var in []reflect.Value
	for i, v := range args {
		in = append(in, convert(v.GetValue(), types[i]))
	}

	out := m.Call(in)
	if len(out) == 0 {
		return Null
	}
	if len(out) == 1 {
		return ToObject(out[0])
	}
	res := &Slice{}
	for i := 0; i < len(out); i++ {
		res.Val = append(res.Val, ToObject(out[i]))
	}
	return res
}

type InnerStructMethodCall struct {
	Name   string
	Method reflect.Value
}

func (fn *InnerStructMethodCall) object() {}

func (fn *InnerStructMethodCall) Type() Type {
	return TypeCall
}

func (fn *InnerStructMethodCall) GetValue() any {
	return Null
}

func (fn *InnerStructMethodCall) Call(args []Object) Object {
	var in []reflect.Value
	for _, v := range args {
		in = append(in, reflect.ValueOf(v.GetValue()))
	}
	out := fn.Method.Call(in)
	if len(out) == 0 {
		return Null
	}
	if len(out) == 1 {
		return ToObject(out[0])
	}
	res := &Slice{}
	for i := 0; i < len(out); i++ {
		res.Val = append(res.Val, ToObject(out[i]))
	}
	return res
}
