package program

import (
	"log"

	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/object"
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
	log.SetFlags(0)
	l := lexer.New(input)
	ps := parser.NewParser(l, p.env)
	states, err := ps.Parse()
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range states {
		obj, _ := v.Exec(p.env)
		if p.repl && obj != object.Null {
			log.Println(obj.GetValue())
		}
	}
}

func (p *Program) SetRepl(repl bool) {
	p.repl = repl
}
