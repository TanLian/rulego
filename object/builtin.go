package object

import "fmt"

type BuiltinFn struct {
	Name string
}

func (fl *BuiltinFn) object() {}

func (fl *BuiltinFn) Type() Type {
	return TypeFnLiteral
}

func (fl *BuiltinFn) GetValue() any {
	return fl
}

func (fl *BuiltinFn) Call(args []Object) Object {
	if fl.Name == "println" {
		var as []any
		for _, v := range args {
			as = append(as, v.GetValue())
		}
		fmt.Println(as...)
		return Null
	}
	return Null
}
