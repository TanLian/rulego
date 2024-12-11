package rulego

import (
	"sync"

	"github.com/tanlian/rulego/ast"
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/parser"
	"github.com/tanlian/rulego/token"
)

func NewRule(input string) *Rule {
	l := lexer.New(input)
	p := parser.NewParser(l)
	exp, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		panic(err)
	}
	r, ok := exp.(*ast.Rule)
	if !ok {
		panic("invalid rule")
	}
	return &Rule{r: r}
}

type Rule struct {
	r *ast.Rule
}

func (r *Rule) Execute(env *environment.Environment) any {
	return r.r.Call(env).GetValue()
}

type RuleBatch struct {
	Rules []*Rule // 一批规则
	Type  uint8   // 0-串行 1-并行执行
}

func (rb *RuleBatch) Execute(env *environment.Environment) []any {
	if rb.Type == 1 {
		return rb.executeConcurrently(env)
	}

	result := make([]any, len(rb.Rules))
	for i, rule := range rb.Rules {
		result[i] = rule.Execute(env)
	}
	return result
}

func (rb *RuleBatch) executeConcurrently(env *environment.Environment) []any {
	var wg sync.WaitGroup
	result := make([]any, len(rb.Rules))
	ch := make(chan ruleBatchResult, len(rb.Rules))
	for i, rule := range rb.Rules {
		wg.Add(1)
		go func(i int, rule *Rule) {
			defer wg.Done()
			ch <- ruleBatchResult{i, rule.Execute(env)}
		}(i, rule)
	}
	wg.Wait()
	close(ch)

	for v := range ch {
		result[v.idx] = v.res
	}
	return result
}

type ruleBatchResult struct {
	idx int
	res any
}

type RuleChain []RuleBatch

func (rc *RuleChain) Add(rb RuleBatch) {
	*rc = append(*rc, rb)
}

func (rc *RuleChain) Execute(env *environment.Environment) []any {
	var result []any
	for _, v := range *rc {
		result = append(result, v.Execute(env))
	}
	return result
}
