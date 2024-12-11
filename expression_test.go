package rulego

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/tanlian/rulego/environment"
)

func TestNewExpression1(t *testing.T) {
	env := environment.New(environment.Root)        // 定义一个 environment 对象
	env.Inject("m", map[string]int{"a": 1, "b": 2}) // 注入map
	t.Log(NewExpression(`m["a"]+m["b"]`).Eval(env)) // result is 3
}

func TestNewExpression(t *testing.T) {
	env := environment.New(environment.Root)

	// 数字运算
	t.Log(NewExpression("(12+8)*5-(36/6)+(4*7)-10").Eval(nil)) // result is 112

	// 比较运算符
	t.Log(NewExpression("1 == 1").Eval(nil)) // result is true
	t.Log(NewExpression("1 != 1").Eval(nil)) // result is false
	t.Log(NewExpression("1 > 1").Eval(nil))  // result is false
	t.Log(NewExpression("1 >= 1").Eval(nil)) // result is true
	t.Log(NewExpression("1 < 1").Eval(nil))  // result is false
	t.Log(NewExpression("1 <= 1").Eval(nil)) // result is true

	// 与或非
	t.Log(NewExpression("!true && false").Eval(nil))                                                               // result is false
	t.Log(NewExpression("(true || !false) && (false && !true)").Eval(nil))                                         // result is false
	t.Log(NewExpression("(!true && false) || (true || !false) && (false && !true) || (!false && true)").Eval(nil)) // result is true

	// 注入普通变量
	env.Inject("a", 1)
	env.Inject("b", 3)
	t.Log("a+b=", NewExpression("a+b").Eval(env)) // result is 4

	// 注入结构体
	user := &TestUser{Name: "leo", Age: 20}
	env.Inject("User", user)
	t.Log("User.Name: ", NewExpression("User.Name").Eval(env))      // result is "leo"
	t.Log("User.Age: ", NewExpression("User.Age").Eval(env))        // result is 20
	t.Log("User.Info(): ", NewExpression("User.Info(1)").Eval(env)) // result is "name: leo, age: 20, sex: 1"
	t.Log(NewExpression("User.AddAge(10)").Eval(env))               // result is 30
	t.Logf("user: %v", user)

	// 注入slice
	env.Inject("Users", []TestUser{{Name: "a1", Age: 10}, {Name: "a2", Age: 11}})
	t.Log(NewExpression("Users").Eval(env))         // result is [{a1 10} {a2 11}]
	t.Log(NewExpression("Users[1].Name").Eval(env)) // result is "a2"
	t.Log(NewExpression("Users[1]").Eval(env))      // result is {a2 11}

	// 注入函数
	env.Inject("println", fmt.Println)
	env.Inject("join", strings.Join)
	env.Inject("split", strings.Split)
	env.Inject("fnInjectMap", fnInjectMap)
	env.Inject("fnInjectSlice", fnInjectSlice)
	env.Inject("fnInjectStruct", fnInjectStruct)
	env.Inject("elems", []string{"hello", "my", "name", "is", "leo"})
	env.Inject("sep", " ")
	env.Inject("m", map[string]string{"name": "leo", "age": "19"})
	env.Inject("u", TestUser{Name: "name", Age: 21})
	env.Inject("Fib", Fib)
	NewExpression("println(100+5*7)").Eval(env)
	t.Log(NewExpression("join(elems, sep)").Eval(env))
	t.Log(NewExpression(`split("aaa,bbb,ccc", ",")`).Eval(env))
	t.Log(NewExpression(`fnInjectMap(m)`).Eval(env))
	t.Log(NewExpression(`fnInjectSlice(elems)`).Eval(env))
	t.Log(NewExpression(`fnInjectStruct(u)`).Eval(env))
	t.Log(NewExpression(`Fib(20)`).Eval(env))
}

type TestUser struct {
	Name string
	Age  int
}

func (u TestUser) Info(sex int) string {
	return fmt.Sprintf("name: %s, age: %d, sex: %d", u.Name, u.Age, sex)
}

func (u *TestUser) AddAge(val int) int {
	u.Age += val
	return u.Age
}

func fnInjectMap(m map[string]string) {
	for k, v := range m {
		fmt.Println("key: ", k, " value: ", v)
	}
}

func fnInjectSlice(s []string) {
	for _, v := range s {
		fmt.Println(v)
	}
}

func fnInjectStruct(u TestUser) {
	fmt.Println("user: ", u)
}

func Fib(n int) int {
	if n < 2 {
		return 1
	}
	return Fib(n-1) + Fib(n-2)
}

func TestAnd(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{input: "true && true"},
			want: true,
		},
		{
			args: args{input: "true && false"},
			want: false,
		},
		{
			args: args{input: "false && false"},
			want: false,
		},
		{
			args: args{input: "false && true"},
			want: false,
		},
		{
			args: args{input: "2 > 1 && true"},
			want: true,
		},
		{
			args: args{input: "2 == 1 && true"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExpression(tt.args.input).Eval(nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpressionAST(t *testing.T) {
	input := `
a = 5;
`
	t.Log(NewExpression(input).AST())
}
