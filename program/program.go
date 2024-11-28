package program

import (
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/parser"
)

type Program struct {
	env *environment.Environment
}

func New() *Program {
	return &Program{env: environment.Root}
}

func (p *Program) Run(input string) {
	l := lexer.New(input)
	ps := parser.NewParser(l, p.env)
	states := ps.Parse()
	for _, v := range states {
		v.Exec(p.env)
	}
}
