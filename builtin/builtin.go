package builtin

/*import (
	"fmt"

	"github.com/tanlian/rulego/object"
)

var BuiltInFn = map[string]*FuncBuiltin{
	"println": {fn: func(a ...object.Object) object.Object {
		var args []any
		for _, v := range a {
			args = append(args, v.GetValue())
		}
		fmt.Println(args...)
		return object.Null
	}},

	"len": {fn: func(a ...object.Object) object.Object {
		if len(a) == 0 {
			return object.Null
		}

		if arrObj, ok := a[0].(*object.Slice); ok {
			return &object.Int{Val: int64(len(arrObj.Val))}
		}
		if mapObj, ok := a[0].(*object.Map); ok {
			return &object.Int{Val: int64(len(mapObj.Val))}
		}
		if strObj, ok := a[0].(*object.String); ok {
			return &object.Int{Val: int64(len(strObj.Val))}
		}
		return object.Null
	}},

	"append": {fn: func(a ...object.Object) object.Object {
		if len(a) == 0 {
			return object.Null
		}

		arr, ok := a[0].(*object.Slice)
		if !ok {
			return object.Null
		}
		res := &object.Slice{Val: append([]any(nil), arr.Val...)}
		for i := 1; i < len(a); i++ {
			res.Val = append(res.Val, a[i])
		}
		return res
	}},
}

type Fn func(...object.Object) object.Object

type FuncBuiltin struct {
	fn Fn
}

func (f FuncBuiltin) Call(args ...object.Object) object.Object {
	return f.fn(args...)
}*/
