package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

type Impl struct {
	Name    string
	Methods map[string]*FnLiteralObj
}

func (im *Impl) Eval(env *environment.Environment) object.Object {
	obj, has := env.Get(im.Name)
	if !has {
		panic(fmt.Sprintf("no such impl %s", im.Name))
	}

	stu, ok := obj.(*object.RgStruct)
	if !ok {
		panic(fmt.Sprintf("impl %s is not a struct", im.Name))
	}

	for k, v := range im.Methods {
		stu.Methods[k] = v
	}
	return object.Null
}

func (im *Impl) String() string {
	return ""
}

func (im *Impl) expressionNode() {}
