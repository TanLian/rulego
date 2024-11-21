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
