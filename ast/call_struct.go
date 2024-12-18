package ast

import (
	"github.com/tanlian/rulego/object"
)

type StructCall struct {
	*object.Empty
	Self   *object.Struct // 指向的 struct 对象
	Method *FnLiteral     // 调用的方法
}

func (sc *StructCall) Call(args []object.Object) object.Object {
	if sc.Self != nil {
		args = append([]object.Object{sc.Self}, args...)
	}
	return sc.Method.Call(args)
}
