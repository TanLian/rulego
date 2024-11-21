package rulego

import (
	"github.com/tanlian/rulego/ast"
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/parser"
)

type MultiRule struct {
	parentEnv *environment.Environment
}

func NewMultiRule(input string) *MultiRule {
	l := lexer.New(input)
	env := environment.New(environment.Root)
	p := parser.NewParser(l, env)
	for _, v := range p.Parse() {
		v.Exec(env)
	}
	return &MultiRule{parentEnv: env}
}

func (mr *MultiRule) ExecuteOne(name string, env *environment.Environment) any {
	obj, ok := mr.parentEnv.Get(name)
	if !ok {
		return nil
	}
	r, ok := obj.(*ast.Rule)
	if !ok {
		return nil
	}
	return r.Call(env).GetValue()
}

func (mr *MultiRule) ExecuteOneByOne(names []string, env *environment.Environment) []any {
	var res []any
	for _, v := range names {
		res = append(res, mr.ExecuteOne(v, env))
	}
	return res
}

func (mr *MultiRule) GetEnv() *environment.Environment {
	return mr.parentEnv
}

func (mr *MultiRule) Upsert(content string) {
	l := lexer.New(content)
	p := parser.NewParser(l, mr.parentEnv)
	for _, v := range p.Parse() {
		v.Exec(mr.parentEnv)
	}
}

func (mr *MultiRule) Remove(name string) {
	mr.parentEnv.Remove(name)
}
