package rulego

import (
	"fmt"
	"testing"

	"github.com/tanlian/rulego/environment"
)

func TestNewRuleFor(t *testing.T) {
	env := environment.New(environment.Root)
	input := `rule rule_for
			{
				ans = 0;
				for i = 0; i < 10; i++ {
					ans += i;
				}
				return ans; // ans is 45
			}`
	fmt.Println(NewRule(input).Execute(env))
}

func TestNewRuleIf(t *testing.T) {
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
	fmt.Println(NewRule(input).Execute(env))
}

func TestNewRule(t *testing.T) {
	env := environment.New(environment.Root)

	var input string
	input = `rule rule1
		{
			// 判断用户的年龄
			if 7 == User.GetAge() { // 判断用户的年龄是否等于7
				User.Age = User.GetAge() + 10000000;
				User.Print("6666");
				return 1;
			} else {
				User.Name = "yyyy";
			}
			return User;
		}`
	user := &TestRuleUser{
		Name: "leo",
		Age:  7,
	}
	env.Inject("User", user)
	fmt.Println(NewRule(input).Execute(env))
	fmt.Println(user)
}

func TestRuleMap(t *testing.T) {
	input := `
rule rule_map
{
	// 过滤掉age小于18 或者 gender不是男的数据
	if data["age"] < 18 || data["gender"] != "男" {
		println("age less than 18 or gender not 男");
		return 0;
	}

	// 对手机号进行脱敏处理
	data["tel"] = fuzzyTel(data["tel"]);

	// 补全id
	data["id"] = 10000;
}
`
	data := map[string]any{
		"age":    18,
		"gender": "男",
		"tel":    "13011110000",
	}
	var fuzzyTel = func(tel string) string {
		if len(tel) != 11 {
			return tel
		}
		return tel[:7] + "****"
	}
	// 注入数据
	env := environment.New(environment.Root)
	env.Inject("data", &data)
	env.Inject("fuzzyTel", fuzzyTel)
	// 生成并执行规则
	NewRule(input).Execute(env)
	fmt.Println("data: ", data)
}

func TestRuleBatch(t *testing.T) {
	env := environment.New(environment.Root)

	input1 := `rule rule1
	{
		print("in rule1");
		return 1;
	}
	`

	input2 := `rule rule2
	{
		print("in rule2");
		return 2;
	}
	`

	input3 := `rule rule3
	{
		print("in rule3");
		return 3;
	}
	`

	env.Inject("print", fmt.Println)
	rb := &RuleBatch{
		Rules: []*Rule{NewRule(input1), NewRule(input2), NewRule(input3)},
		Type:  1,
	}
	fmt.Println(rb.Execute(env))
}

type TestRuleUser struct {
	Name string
	Age  int
}

func (u *TestRuleUser) GetAge() int {
	return u.Age
}

func (u *TestRuleUser) Print(s string) {
	fmt.Println(s)
}
