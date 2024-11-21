package object

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	reflect.DeepEqual(New(10), &Int{Val: 10})
	reflect.DeepEqual(New(10.1), &Float{Val: 10.1})
	reflect.DeepEqual(New("hello"), &String{Val: "hello"})
	reflect.DeepEqual(New(true), &Bool{Val: true})
	reflect.DeepEqual(New([]uint32{1, 2, 3}), &Slice{Val: []Object{&Int{Val: 1}, &Int{Val: 2}, &Int{Val: 3}}})
	user := TestUser{Name: "leo", Age: 20}
	fmt.Println(New(user).GetValue())
	fmt.Println(New(&user).GetValue())
}

type TestUser struct {
	Name string
	Age  int
}
