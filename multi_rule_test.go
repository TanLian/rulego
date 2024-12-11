package rulego

import (
	"fmt"
	"testing"
)

// TODO
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
			return fib(n);
		}
		`
	mr := NewMultiRule(input)
	env := mr.GetEnv()
	fmt.Println("1111 env: ", env)
	env.Inject("n", 10)
	fmt.Println("2222 env: ", env)
	//env.Inject("println", fmt.Println)
	fmt.Println("33333 env: ", env)
	fmt.Println(mr.ExecuteOne("rule1"))
	fmt.Println(mr.ExecuteOne("rule2"))

	//newRule1 := `rule rule1
	//	{
	//		if n < 10 {
	//			return n;
	//		}
	//		return n+2;
	//	}`
	//mr.Upsert(newRule1) // 动态更新rule1
	//mr.Remove("rule2")  // 删除rule2
}
