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
	POWER

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
	STRUCT
	IMPL

	// 标点符号
	DOT       // .
	SEMICOLON // ;
	COMMA     // ,
	COLON     // :
	MOD       // %

	GREATER   // >
	LESS      // <
	EQUAL     // ==
	NOT_EQUAL // !=

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

	AND         // &
	OR          // |
	XOR         // ^
	LEFT_SHIFT  // <<
	RIGHT_SHIFT // >>
	LOGIC_AND   // &&
	LOGIC_OR    // ||
	BANG        // !

	INLINE_COMMENTS

	// end of file
	EOF
)

type Token struct {
	Type  TokenType
	Value string
	Row   int
	Col   int
}

func (t Token) String() string {
	m := map[TokenType]string{
		Invalid:         "invalid",
		PLUS:            "plus",
		MINUS:           "minus",
		TIMES:           "times",
		DIVIDE:          "divide",
		POWER:           "power",
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
		STRUCT:          "struct",
		IMPL:            "impl",
		DOT:             "dot",
		SEMICOLON:       "semicolon",
		COMMA:           "comma",
		COLON:           "colon",
		IDENTIFIER:      "ident",
		ASSIGN:          "assign",
		GREATER:         "greater",
		LESS:            "less",
		EQUAL:           "equal",
		NOT_EQUAL:       "not_equal",
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
		XOR:             "^",
		LEFT_SHIFT:      "<<",
		RIGHT_SHIFT:     ">>",
		LOGIC_AND:       "&&",
		LOGIC_OR:        "||",
		INLINE_COMMENTS: "//",
		EOF:             "eof",
	}
	return fmt.Sprintf("(%s,%s)", m[t.Type], t.Value)
}

var (
	InvalidToken = Token{Type: Invalid}
	EofToken     = Token{Type: EOF}
)

var (
	Keyword = map[string]TokenType{
		"true":     BOOL,
		"false":    BOOL,
		"if":       IF,
		"else":     ELSE,
		"fn":       FUNC,
		"rule":     RULE,
		"return":   RETURN,
		"for":      FOR,
		"break":    BREAK,
		"continue": CONTINUE,
		"switch":   SWITCH,
		"case":     CASE,
		"default":  DEFAULT,
		"struct":   STRUCT,
		"impl":     IMPL,
	}
)
