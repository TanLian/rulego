package object

import (
	"fmt"
	"reflect"
)

type InjectFn struct {
	Value reflect.Value
}

func (fn *InjectFn) object() {}

func (fn *InjectFn) Type() Type {
	return TypeFn
}

func (fn *InjectFn) GetValue() any {
	return fn.Value.Interface()
}

func (fn *InjectFn) Call(params []Object) Object {
	// 获取函数的参数数量
	numIn := fn.Value.Type().NumIn()
	if len(params) != numIn {
		panic(fmt.Sprintf("fn.Call: 参数数量不匹配, 需要参数数量: %d, 实际参数数量: %d", numIn, len(params)))
	}

	var in []reflect.Value
	for i, v := range params {
		in = append(in, convert(v.GetValue(), fn.Value.Type().In(i)))
	}
	out := fn.Value.Call(in)
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
