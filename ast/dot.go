package ast

import (
	"fmt"
	"reflect"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
a = [1,2,3];
a.push(4);

type user struct {
	name string
}
u := &user{name: "name"}
u.name
*/

type Dot struct {
	Left  Expression // slice、struct
	Right Expression // ident
}

func (p *Dot) Eval(env *environment.Environment) object.Object {
	right, ok := p.Right.(*Ident)
	if !ok {
		return object.Null
	}
	//fmt.Println(p.Left.String())
	left := p.Left.Eval(env)
	// fmt.Println("leftObj: ", leftObj.GetValue(), " type: ", leftObj.Type())
	if structObj, ok := left.(*object.Struct); ok {
		stt := structObj.Value
		if stt.Kind() == reflect.Pointer {
			stt = stt.Elem()
		}
		val := stt.FieldByName(right.Token.Value)
		if val.IsValid() {
			return object.ToObject(val)
		}
	}

	if structObj, ok := left.(*object.RgStruct); ok {
		return structObj.GetFieldValue(right.Token.Value)
	}
	return object.Null
}

func (p *Dot) String() string {
	return fmt.Sprintf("(%s.%s)", p.Left.String(), p.Right.String())
}

func (p *Dot) expressionNode() {}
