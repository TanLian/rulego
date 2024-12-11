package rulego

import (
	"fmt"
	"log"

	"github.com/tanlian/rulego/ast"
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/parser"
)

type MultiRule struct {
	env *environment.Environment
}

func NewMultiRule(input string) *MultiRule {
	l := lexer.New(input)
	env := environment.New(environment.Root)
	p := parser.NewParser(l, env)
	states, err := p.Parse()
	if err != nil {
		log.Println(err)
		return nil
	}
	for _, v := range states {
		v.Exec(env)
	}
	fmt.Println("... env: ", env)
	return &MultiRule{env: env}
}

func (mr *MultiRule) ExecuteOne(name string) any {
	obj, ok := mr.env.Get(name)
	if !ok {
		return nil
	}
	r, ok := obj.(*ast.Rule)
	if !ok {
		return nil
	}
	return r.Call(mr.env).GetValue()
}

func (mr *MultiRule) ExecuteOneByOne(names []string) []any {
	var res []any
	for _, v := range names {
		res = append(res, mr.ExecuteOne(v))
	}
	return res
}

func (mr *MultiRule) GetEnv() *environment.Environment {
	return mr.env
}

func (mr *MultiRule) Upsert(content string) {
	l := lexer.New(content)
	p := parser.NewParser(l, mr.env)
	states, err := p.Parse()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range states {
		v.Exec(mr.env)
	}
}

func (mr *MultiRule) Remove(name string) {
	mr.env.Remove(name)
}
