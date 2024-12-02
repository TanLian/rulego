package token

import "fmt"

type TokenType uint32

// 一系列 token 定义
const (
	Invalid TokenType = iota

	// 运算符
	PLUS
	MINUS
	TIMES
	DIVIDE

	// 括号
	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]
	LBRACE   // {
	RBRACE   // }

	NUMBER // 数字
	BOOL
	STRING

	// 标点符号
	DOT       // .
	SEMICOLON // ;
	COMMA     // ,
	COLON     // :
	MOD       // %

	GREATER // >
	LESS    // <

	IDENTIFIER
	ASSIGN // =
	IF
	ELSE
	FUNC
	RETURN
	FOR
	BREAK
	CONTINUE
	SWITCH
	CASE
	DEFAULT
	RULE

	AND  // &
	OR   // |
	BANG // !

	INLINE_COMMENTS

	// end of file
	EOF
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) String() string {
	m := map[TokenType]string{
		Invalid:         "invalid",
		PLUS:            "plus",
		MINUS:           "minus",
		TIMES:           "times",
		DIVIDE:          "divide",
		BANG:            "bang",
		MOD:             "mod",
		LPAREN:          "lparen",
		RPAREN:          "rparen",
		LBRACKET:        "lbracket",
		RBRACKET:        "rbracket",
		LBRACE:          "lbrace",
		RBRACE:          "rbrace",
		NUMBER:          "number",
		BOOL:            "bool",
		STRING:          "string",
		DOT:             "dot",
		SEMICOLON:       "semicolon",
		COMMA:           "comma",
		COLON:           "colon",
		IDENTIFIER:      "ident",
		ASSIGN:          "assign",
		GREATER:         "greater",
		LESS:            "less",
		IF:              "if",
		ELSE:            "else",
		SWITCH:          "switch",
		CASE:            "case",
		DEFAULT:         "default",
		FUNC:            "fn",
		RULE:            "rule",
		RETURN:          "return",
		BREAK:           "break",
		CONTINUE:        "continue",
		AND:             "&",
		OR:              "|",
		INLINE_COMMENTS: "//",
		EOF:             "eof",
	}
	return fmt.Sprintf("(%s,%s)", m[t.Type], t.Value)
}

var (
	InvalidToken = Token{Type: Invalid}
	EofToken     = Token{Type: EOF}
)
