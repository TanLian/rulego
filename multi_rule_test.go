package rulego

import (
	"fmt"
	"testing"
)

func TestMultiRule(t *testing.T) {
	input := `
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
			println(n);
			return fib(n);
		}
		`
	mr := NewMultiRule(input)
	env := mr.GetEnv()
	env.Inject("n", 10)
	env.Inject("println", fmt.Println)
	fmt.Println(mr.ExecuteOne("rule1", env))
	fmt.Println(mr.ExecuteOne("rule2", env))

	newRule1 := `rule rule1
		{
			if n < 10 {
				return n;
			}
			return n+2;
		}`
	mr.Upsert(newRule1) // 动态更新rule1
	mr.Remove("rule2")  // 删除rule2
}
