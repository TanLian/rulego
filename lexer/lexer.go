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
	row      int
	col      int
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
		if l.currentChar() == '\r' || l.currentChar() == '\n' {
			l.row++
			l.col = 0
		}
		l.position++
		l.col++
	}()
	ch := l.currentChar()
	switch ch {
	case '+':
		return token.Token{
			Type:  token.PLUS,
			Value: "+",
			Row:   l.row,
			Col:   l.col,
		}
	case '-':
		return token.Token{
			Type:  token.MINUS,
			Value: "-",
			Row:   l.row,
			Col:   l.col,
		}
	case '*':
		if l.peekChar() == '*' {
			l.position++
			l.col++
			return token.Token{Type: token.POWER, Value: "**", Row: l.row, Col: l.col}
		}
		return token.Token{
			Type:  token.TIMES,
			Value: "*",
			Row:   l.row,
			Col:   l.col,
		}
	case '/':
		if l.peekChar() == '/' { // inline comments
			l.position += 2
			l.col += 2
			for l.currentChar() != '\n' {
				l.position++
				l.col++
			}
			return token.Token{Type: token.INLINE_COMMENTS, Value: "//", Row: l.row, Col: l.col}
		}

		return token.Token{
			Type:  token.DIVIDE,
			Value: "/",
			Row:   l.row,
			Col:   l.col,
		}
	case '!':
		if l.peekChar() == '=' {
			l.position++
			l.col++
			return token.Token{Type: token.NOT_EQUAL, Value: "!=", Row: l.row, Col: l.col}
		}
		return token.Token{
			Type:  token.BANG,
			Value: "!",
			Row:   l.row,
			Col:   l.col,
		}
	case '%':
		return token.Token{
			Type:  token.MOD,
			Value: "%",
			Row:   l.row,
			Col:   l.col,
		}
	case '&':
		if l.peekChar() == '&' {
			l.position++
			l.col++
			return token.Token{Type: token.LOGIC_AND, Value: "&&", Row: l.row, Col: l.col}
		}
		return token.Token{Type: token.AND, Value: "&", Row: l.row, Col: l.col}
	case '|':
		if l.peekChar() == '|' {
			l.position++
			l.col++
			return token.Token{Type: token.LOGIC_OR, Value: "||", Row: l.row, Col: l.col}
		}
		return token.Token{Type: token.OR, Value: "|", Row: l.row, Col: l.col}
	case '^':
		return token.Token{Type: token.XOR, Value: "^", Row: l.row, Col: l.col}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return l.eatNumber()
	case '(':
		return token.Token{Type: token.LPAREN, Value: "(", Row: l.row, Col: l.col}
	case ')':
		return token.Token{Type: token.RPAREN, Value: ")", Row: l.row, Col: l.col}
	case '[':
		return token.Token{Type: token.LBRACKET, Value: "[", Row: l.row, Col: l.col}
	case ']':
		return token.Token{Type: token.RBRACKET, Value: "]", Row: l.row, Col: l.col}
	case '{':
		return token.Token{Type: token.LBRACE, Value: "{", Row: l.row, Col: l.col}
	case '}':
		return token.Token{Type: token.RBRACE, Value: "}", Row: l.row, Col: l.col}
	case '=':
		if l.peekChar() == '=' {
			l.position++
			l.col++
			return token.Token{Type: token.EQUAL, Value: "==", Row: l.row, Col: l.col}
		}
		return token.Token{Type: token.ASSIGN, Value: "=", Row: l.row, Col: l.col}
	case ';':
		return token.Token{Type: token.SEMICOLON, Value: ";", Row: l.row, Col: l.col}
	case ',':
		return token.Token{Type: token.COMMA, Value: ",", Row: l.row, Col: l.col}
	case ':':
		return token.Token{Type: token.COLON, Value: ":", Row: l.row, Col: l.col}
	case '.':
		return token.Token{Type: token.DOT, Value: ".", Row: l.row, Col: l.col}
	case '"':
		return l.eatString('"')
	case '\'':
		return l.eatString('\'')
	case '>':
		if l.peekChar() == '>' {
			l.position++
			l.col++
			return token.Token{Type: token.RIGHT_SHIFT, Value: ">>", Row: l.row, Col: l.col}
		}
		return token.Token{Type: token.GREATER, Value: ">", Row: l.row, Col: l.col}
	case '<':
		if l.peekChar() == '<' {
			l.position++
			l.col++
			return token.Token{Type: token.LEFT_SHIFT, Value: "<<", Row: l.row, Col: l.col}
		}
		return token.Token{Type: token.LESS, Value: "<", Row: l.row, Col: l.col}
	default:
		if util.IsAlphabet(ch) {
			return l.eatIdent()
		}
		return token.InvalidToken
	}
}

func (l *Lexer) skipBlank() {
	for l.position < len(l.input) && (l.currentChar() == ' ' || l.currentChar() == '\n' || l.currentChar() == '\r' || l.currentChar() == '\t') {
		if l.currentChar() == '\n' || l.currentChar() == '\r' {
			l.row++
			l.col = 0
		}
		l.position++
		l.col++
	}
}

func (l *Lexer) eatNumber() token.Token {
	var s bytes.Buffer
	for l.position < len(l.input) && (l.currentChar() == '.' || util.IsDigit(l.currentChar())) {
		s.WriteByte(l.currentChar())
		l.position++
		l.col++
	}
	l.position--
	l.col--
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
		Row:   l.row,
		Col:   l.col,
	}
}

func (l *Lexer) eatString(end byte) token.Token {
	var s bytes.Buffer
	l.position++
	l.col++
	for l.position < len(l.input) && l.currentChar() != end {
		s.WriteByte(l.currentChar())
		l.position++
		l.col++
	}
	str := s.String()
	return token.Token{
		Type:  token.STRING,
		Value: str,
		Row:   l.row,
		Col:   l.col,
	}
}

func (l *Lexer) eatIdent() token.Token {
	var s bytes.Buffer
	for l.position < len(l.input) && (util.IsAlphabet(l.input[l.position]) || util.IsDigit(l.input[l.position]) || l.input[l.position] == '_') {
		s.WriteByte(l.input[l.position])
		l.position++
		l.col++
	}
	l.position--
	l.col--
	str := s.String()
	if typ, ok := token.Keyword[str]; ok {
		return token.Token{Type: typ, Value: str, Row: l.row, Col: l.col}
	}
	return token.Token{Type: token.IDENTIFIER, Value: str, Row: l.row, Col: l.col}
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
