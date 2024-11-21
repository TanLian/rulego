package lexer

import (
	"bytes"
	"strings"

	"github.com/tanlian/rulego/token"
	"github.com/tanlian/rulego/util"
)

type Lexer struct {
	input    string
	position int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) ReadNextToken() token.Token {
	l.skipBlank()
	if l.position >= len(l.input) {
		return token.EofToken
	}
	defer func() {
		l.position++
	}()
	ch := l.currentChar()
	switch ch {
	case '+':
		if l.peekChar() == '=' {
			l.position++
			return token.Token{Type: token.PLUS_ASSIGN, Value: "+="}
		}
		if l.peekChar() == '+' {
			l.position++
			return token.Token{Type: token.PLUS_PLUS, Value: "++"}
		}
		return token.Token{
			Type:  token.PLUS,
			Value: "+",
		}
	case '-':
		return token.Token{
			Type:  token.MINUS,
			Value: "-",
		}
	case '*':
		return token.Token{
			Type:  token.TIMES,
			Value: "*",
		}
	case '/':
		if l.peekChar() == '/' { // inline comments
			l.position += 2
			for l.currentChar() != '\n' {
				l.position++
			}
			return l.ReadNextToken()
		}
		return token.Token{
			Type:  token.DIVIDE,
			Value: "/",
		}
	case '!':
		if l.peekChar() == '=' {
			l.position++
			return token.Token{Type: token.NOT_EQUAL, Value: "!="}
		}
		return token.Token{
			Type:  token.BANG,
			Value: "!",
		}
	case '%':
		return token.Token{
			Type:  token.MOD,
			Value: "%",
		}
	case '&':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '&' {
			l.position++
			return token.Token{Type: token.AND, Value: "&&"}
		}
		return token.InvalidToken // TODO
	case '|':
		if l.position+1 < len(l.input) && l.input[l.position+1] == '|' {
			l.position++
			return token.Token{Type: token.OR, Value: "||"}
		}
		return token.InvalidToken // TODO
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.eatNumber()
	case '(':
		return token.Token{Type: token.LPAREN, Value: "("}
	case ')':
		return token.Token{Type: token.RPAREN, Value: ")"}
	case '[':
		return token.Token{Type: token.LBRACKET, Value: "["}
	case ']':
		return token.Token{Type: token.RBRACKET, Value: "]"}
	case '{':
		return token.Token{Type: token.LBRACE, Value: "{"}
	case '}':
		return token.Token{Type: token.RBRACE, Value: "}"}
	case '=':
		if l.peekChar() == '=' {
			l.position++
			return token.Token{Type: token.EQUAL, Value: "=="}
		}
		return token.Token{Type: token.ASSIGN, Value: "="}
	case ';':
		return token.Token{Type: token.SEMICOLON, Value: ";"}
	case ',':
		return token.Token{Type: token.COMMA, Value: ","}
	case ':':
		return token.Token{Type: token.COLON, Value: ":"}
	case '.':
		return token.Token{Type: token.DOT, Value: "."}
	case '"':
		return l.eatString('"')
	case '\'':
		return l.eatString('\'')
	case '>':
		if l.peekChar() == '=' {
			l.position++
			return token.Token{Type: token.GREATER_EQUAL, Value: ">="}
		}
		return token.Token{Type: token.GREATER, Value: ">"}
	case '<':
		if l.peekChar() == '=' {
			l.position++
			return token.Token{Type: token.LESS_EQUAL, Value: "<="}
		}
		return token.Token{Type: token.LESS, Value: "<"}
	default:
		if util.IsAlphabet(ch) {
			return l.eatIdent()
		}
		return token.InvalidToken
	}
}

func (l *Lexer) skipBlank() {
	for l.position < len(l.input) && (l.currentChar() == ' ' || l.currentChar() == '\n' || l.currentChar() == '\r' || l.currentChar() == '\t') {
		l.position++
	}
}

func (l *Lexer) eatNumber() token.Token {
	var s bytes.Buffer
	for l.position < len(l.input) && (l.currentChar() == '.' || util.IsDigit(l.currentChar())) {
		s.WriteByte(l.currentChar())
		l.position++
	}
	l.position--
	str := s.String()
	if str[len(str)-1] == '.' {
		return token.InvalidToken
	}
	if strings.Count(str, ".") > 1 {
		return token.InvalidToken
	}
	return token.Token{
		Type:  token.NUMBER,
		Value: str,
	}
}

func (l *Lexer) eatString(end byte) token.Token {
	var s bytes.Buffer
	l.position++
	for l.position < len(l.input) && l.currentChar() != end {
		s.WriteByte(l.currentChar())
		l.position++
	}
	str := s.String()
	return token.Token{
		Type:  token.STRING,
		Value: str,
	}
}

func (l *Lexer) eatIdent() token.Token {
	var s bytes.Buffer
	for l.position < len(l.input) && (util.IsAlphabet(l.input[l.position]) || util.IsDigit(l.input[l.position]) || l.input[l.position] == '_') {
		s.WriteByte(l.input[l.position])
		l.position++
	}
	l.position--
	str := s.String()
	switch str {
	case "true", "false":
		return token.Token{Type: token.BOOL, Value: str}
	case "if":
		return token.Token{Type: token.IF, Value: str}
	case "else":
		return token.Token{Type: token.ELSE, Value: str}
	case "fn":
		return token.Token{Type: token.FUNC, Value: str}
	case "rule":
		return token.Token{Type: token.RULE, Value: str}
	case "return":
		return token.Token{Type: token.RETURN, Value: str}
	case "for":
		return token.Token{Type: token.FOR, Value: str}
	case "break":
		return token.Token{Type: token.BREAK, Value: str}
	case "switch":
		return token.Token{Type: token.SWITCH, Value: str}
	case "case":
		return token.Token{Type: token.CASE, Value: str}
	case "default":
		return token.Token{Type: token.DEFAULT, Value: str}
	default:
		return token.Token{Type: token.IDENTIFIER, Value: s.String()}
	}
}

func (l *Lexer) currentChar() byte {
	return l.input[l.position]
}

func (l *Lexer) peekChar() byte {
	if l.position+1 < len(l.input) {
		return l.input[l.position+1]
	}
	return 0
}
