package program

import (
	"fmt"

	"github.com/tanlian/rulego/object"

	"github.com/tanlian/rulego/ast"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/parser"
)

type Program struct {
	env  *environment.Environment
	repl bool
}

func New() *Program {
	return &Program{env: environment.Root}
}

func (p *Program) Run(input string) {
	l := lexer.New(input)
	ps := parser.NewParser(l, p.env)
	states := ps.Parse()
	for _, v := range states {
		if expr, ok := v.(*ast.ExpressionStatement); p.repl && ok {
			if obj := expr.Expr.Eval(p.env); obj != object.Null {
				fmt.Println(obj.GetValue())
			}
			continue
		}
		v.Exec(p.env)
	}
}

func (p *Program) SetRepl(repl bool) {
	p.repl = repl
}
