package program

import (
	"log"
	"strings"

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

func (p *Program) Run(input string) error {
	log.SetFlags(0)
	l := lexer.New(input)
	ps := parser.NewParser(l)
	states, err := ps.Parse()
	if err != nil {
		log.Println(err)
		return err
	}
	for _, v := range states {
		obj, _ := v.Exec(p.env)
		if p.repl && obj != object.Null {
			log.Println(obj.GetValue())
		}
	}
	return nil
}

func (p *Program) SetRepl(repl bool) {
	p.repl = repl
}

func (p *Program) AST(input string) string {
	log.SetFlags(0)
	l := lexer.New(input)
	ps := parser.NewParser(l)
	states, err := ps.Parse()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	var s strings.Builder
	for _, v := range states {
		s.WriteString(v.AST(0))
	}
	return s.String()
}
