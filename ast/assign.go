package ast

import (
	"fmt"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/object"
)

/*
赋值语句的左边是一个ident、index表达式 或 dot表达式
赋值语句的右边是一个表达式

// ident
a = 10;
a = expression;

// index
nums = [1,2,3];
nums[1] = 4; // 这个1有可能是别的函数返回的

person = {"name": "leo"};
person["name"] = "tan";


// dot
type User struct {
	Name string
}
user := new(User)
user.Name = "leo"
*/

type Assign struct {
	Left  Expression
	Right Expression
}

func (as *Assign) Exec(env *environment.Environment) (object.Object, bool, bool) {
	if ident, ok := as.Left.(*Ident); ok {
		// fmt.Println("assign to ", ident.Token.Value, " val: ", as.Right.String())
		env.Set(ident.Token.Value, as.Right.Eval(env))
		return object.Null, false, false
	}

	if idx, ok := as.Left.(*Index); ok {
		data := idx.Data.Eval(env)
		key := idx.Key.Eval(env).GetValue()
		if m, ok := data.(*object.Map); ok {
			m.Val[key] = as.Right.Eval(env)
		}
		if s, ok := data.(*object.Slice); ok {
			s.Val[int(key.(float64))] = as.Right.Eval(env)
		}
		return object.Null, false, false
	}

	if dot, ok := as.Left.(*Dot); ok {
		dotLeft := dot.Left.Eval(env)
		if s, ok := dotLeft.(*object.Struct); ok {
			if dotRight, ok := dot.Right.(*Ident); ok {
				s.SetField(dotRight.Token.Value, as.Right.Eval(env).GetValue())
			}
		}
		return object.Null, false, false
	}
	panic("invalid assign statement")
}

func (as *Assign) String() string {
	return fmt.Sprintf("%s = %s", as.Left.String(), as.Right.String())
}

func (as *Assign) statementNode() {}