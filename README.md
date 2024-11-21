# 表达式
## 解析
支持 
- ==、!=、> 、>=、<、<=、==、!=
- &&、||、!
- +、-、*、/
- myFunc(X)、a.b()
- array[index]、a.b.c

## 注入数据
- 我们强调代码（包括表达式、规则等）与数据**分离**的思想，规则是预先定义好的，而数据是动态注入的
```go
env := environment.New(environment.Root)
env.Inject("User", &User{Name:"leo", Age: 18}) // 注入一个结构体对象
env.Inject("Nums", []uint32{0,1,2,3}) // 注入一个slice对象
env.Inject("MapInfo", map[string]interface{}{"name":"leo"}) // 注入一个map对象
env.Inject("println",fmt.Println) // 注入一个函数
```

## 执行
result := NewExpression("your expression").Eval(env)

## 示例
```go
// 数字运算
NewExpression("(12+8)*5-(36/6)+(4*7)-10").Eval(nil) // result is 112

// 比较运算符
NewExpression("1 == 1").Eval(nil) // result is true
NewExpression("1 != 1").Eval(nil) // result is false
NewExpression("1 > 1").Eval(nil)  // result is false
NewExpression("1 >= 1").Eval(nil) // result is true
NewExpression("1 < 1").Eval(nil)  // result is false
NewExpression("1 <= 1").Eval(nil) // result is true

// 与或非
NewExpression("!true && false").Eval(nil)                                                               // result is false
NewExpression("(true || !false) && (false && !true)").Eval(nil)                                         // result is false
NewExpression("(!true && false) || (true || !false) && (false && !true) || (!false && true)").Eval(nil) // result is true

// 注入普通变量
env := environment.New(environment.Root)
env.Inject("a", 1)
env.Inject("b", 3)
NewExpression("a+b").Eval(env) // result is 4

// 注入结构体
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

user := &TestUser{Name: "leo", Age: 20}
env.Inject("User", user)
NewExpression("User.Name").Eval(env)        // result is "leo"
NewExpression("User.Age").Eval(env)         // result is 20
NewExpression("User.Info(1)").Eval(env)     // result is "name: leo, age: 20, sex: 1"
NewExpression("User.AddAge(10)").Eval(env)  // result is 30

// 注入slice
env.Inject("Users", []TestUser{{Name: "a1", Age: 10}, {Name: "a2", Age: 11}})
NewExpression("Users").Eval(env)         // result is [{a1 10} {a2 11}]
NewExpression("Users[1].Name").Eval(env) // result is "a2"
NewExpression("Users[1]").Eval(env)      // result is {a2 11}

// 注入函数
env.Inject("join", strings.Join)
env.Inject("elems", []string{"hello", "my", "name", "is", "leo"})
env.Inject("sep", " ")
NewExpression("join(elems, sep)").Eval(env) // result is "hello my name is leo"
```

# 规则
## 语法格式
```go
rule rule_name
{
    your_statement_1;
    your_statement_2;
    your_statement_3;
    return expression; // optional
}
```

## 支持的语句
### for语句
```go
rule rule_for
{
    ans = 0;
    for i = 0; i < 10; i++ {
        ans += i;
    }
    return ans; // ans is 45
}
```

### if语句
```go
rule rule_if
{
    if age >= 18 {
        println("You are an adult.");	
    } else if age >= 10 {
        println("You are a teenager.");
    } else {
        println("You are a child.");
    }
}
```
## 解析

```go
rule := NewRule(input string)
```

1. 规则可以有返回值，也可以没有
2. 规则内支持行注释，用 //
3. 支持 if、else、for、break、return等关键字
4. 规则之间**不能相互调用**，规则也不能递归调用自己
5. 限制：一次性只能解析**单个**规则

## 注入数据
同表达式的数据部分

## 执行
```go
result := rule.Execute(env)
```

## 示例
```go
env := environment.New(environment.Root)
input := `rule rule_if
	{
		if age >= 18 {
			println("You are an adult.");	
		} else if age >= 10 {
			println("You are a teenager.");
		} else {
			println("You are a child.");
		}
	}`
env.Inject("age", 18)
env.Inject("println", fmt.Println)
fmt.Println(NewRule(input).Execute(env)) // it will print "You are an adult."
```

# 多规则
上面的规则一次性只能解析单个规则，本小节就是突破这个限制，使得一次性可以解析多个规则。

## 定义多个规则
```go
fn fib(n) {
	if n < 2 {
		return 1;
	}
	return fib(n-1) + fib(n-2);
}

rule rule1
{
	if n < 2 {
		return n;
	}
	return n+2;
}

rule rule2
{
	println(n);   // 这是通过environment注入的函数
	return fib(n);   // 这是我们自定义的函数
}
```
- 上面定义了一个函数和两个规则
- 规则内**可以调用函数**，不管是用户自定义的函数还是通过environment注入的函数
- **优先使用**通过environment注入的函数

## 规则的执行、更新、删除

```go
mr := NewMultiRule(input)   // input就是上面定义的多规则
env := mr.GetEnv()    // 获取环境变量
env.Inject("n", 10)   // 注入变量n
env.Inject("println", fmt.Println) // 注入函数
fmt.Println(mr.ExecuteOne("rule1", env))    // it will print 12
fmt.Println(mr.ExecuteOne("rule2", env))    // 执行规则rule2 

newRule1 := `rule rule1
		{
			if n < 10 {
				return n;
			}
			return n+2;
		}`
mr.Upsert(newRule1) // 动态更新rule1
mr.Remove("rule2")  // 删除rule2
```

# 规则的编排
TODO