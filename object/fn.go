package object

import (
	"fmt"
	"reflect"
)

type Fn struct {
	Value reflect.Value
}

func (fn *Fn) object() {}

func (fn *Fn) Type() Type {
	return TypeFn
}

func (fn *Fn) GetValue() any {
	return fn.Value.Interface()
}

func (fn *Fn) Call(params []Object) (Object, error) {
	// 获取函数的参数数量
	numIn := fn.Value.Type().NumIn()
	if len(params) != numIn {
		return Null, fmt.Errorf("func expects %d parameters, but %d were provided", numIn, len(params))
	}

	var in []reflect.Value
	for i, v := range params {
		in = append(in, convert(v.GetValue(), fn.Value.Type().In(i)))
	}
	out := fn.Value.Call(in)
	if len(out) == 0 {
		return Null, nil
	}
	if len(out) == 1 {
		return ToObject(out[0]), nil
	}
	res := &Slice{}
	for i := 0; i < len(out); i++ {
		res.Val = append(res.Val, ToObject(out[i]))
	}
	return res, nil
}
