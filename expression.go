package rulego

import (
	"github.com/tanlian/rulego/ast"
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/parser"
	"github.com/tanlian/rulego/token"
)

type Expression struct {
	exp ast.Expression
}

func NewExpression(input string) *Expression {
	l := lexer.New(input)
	p := parser.NewParser(l, nil)
	expr, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		panic(err)
	}
	return &Expression{exp: expr}
}

func (exp *Expression) Eval(data any) any {
	if env, ok := data.(*environment.Environment); ok {
		return exp.exp.Eval(env).GetValue()
	}

	env := environment.New(environment.Root)
	if kv, ok := data.(map[string]any); ok {
		for k, v := range kv {
			env.Inject(k, v)
		}
	}
	return exp.exp.Eval(env).GetValue()
}

func (exp *Expression) AST() string {
	return exp.exp.AST(0)
}
