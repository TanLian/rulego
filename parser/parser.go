package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tanlian/rulego/ast"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(exp ast.Expression) (ast.Expression, error)
)

type Parser struct {
	l *lexer.Lexer

	prefixFn     map[token.TokenType]prefixParseFn
	infixParseFn map[token.TokenType]infixParseFn

	// token 相关
	currentToken token.Token
	peekToken    token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:            l,
		prefixFn:     make(map[token.TokenType]prefixParseFn),
		infixParseFn: make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefixFn(token.IDENTIFIER, p.parseIdent)
	p.registerPrefixFn(token.NUMBER, p.parseNumber)
	p.registerPrefixFn(token.STRING, p.parseString)
	p.registerPrefixFn(token.BOOL, p.parseBool)
	p.registerPrefixFn(token.LPAREN, p.parseGroup)
	p.registerPrefixFn(token.LBRACKET, p.parseSlice)
	p.registerPrefixFn(token.LBRACE, p.parseHashMap)
	p.registerPrefixFn(token.FUNC, p.parseFnLiteral)
	p.registerPrefixFn(token.RULE, p.parseRuleLiteral)
	p.registerPrefixFn(token.BANG, p.parseBang)
	p.registerPrefixFn(token.PLUS, p.parsePositive)
	p.registerPrefixFn(token.MINUS, p.parseNegative)
	p.registerPrefixFn(token.XOR, p.parseBitwiseNot)
	p.registerPrefixFn(token.IF, p.parseIf)
	p.registerPrefixFn(token.FOR, p.parseFor)
	p.registerPrefixFn(token.SWITCH, p.parseSwitch)
	p.registerPrefixFn(token.RETURN, p.parseReturn)
	p.registerPrefixFn(token.BREAK, p.parseBreak)
	p.registerPrefixFn(token.CONTINUE, p.parseContinue)
	p.registerPrefixFn(token.STRUCT, p.parseStructLiteral)
	p.registerPrefixFn(token.IMPL, p.parseImpl)

	p.registerInfixFn(token.PLUS, p.parsePlus)
	p.registerInfixFn(token.MINUS, p.parseMinus)
	p.registerInfixFn(token.TIMES, p.parseTimes)
	p.registerInfixFn(token.DIVIDE, p.parseDivide)
	p.registerInfixFn(token.LPAREN, p.parseCall)
	p.registerInfixFn(token.LBRACKET, p.parseIndex)
	p.registerInfixFn(token.GREATER, p.parseGreater)
	p.registerInfixFn(token.LESS, p.parseLesser)
	p.registerInfixFn(token.EQUAL, p.parseEqual)
	p.registerInfixFn(token.NOT_EQUAL, p.parseNotEqual)
	p.registerInfixFn(token.AND, p.parseAnd)
	p.registerInfixFn(token.LOGIC_AND, p.parseLogicAnd)
	p.registerInfixFn(token.OR, p.parseOr)
	p.registerInfixFn(token.XOR, p.parseXor)
	p.registerInfixFn(token.LOGIC_OR, p.parseLogicOr)
	p.registerInfixFn(token.LEFT_SHIFT, p.parseLeftShift)
	p.registerInfixFn(token.RIGHT_SHIFT, p.parseRightShift)
	p.registerInfixFn(token.MOD, p.parseMod)
	p.registerInfixFn(token.DOT, p.parseDot)
	p.registerInfixFn(token.ASSIGN, p.parseAssign)
	p.registerInfixFn(token.POWER, p.parsePower)
	p.registerInfixFn(token.LBRACE, p.parseStructInstantiate)

	p.forward()
	p.forward()

	return p
}

func (p *Parser) Parse() ([]ast.Statement, error) {
	var statements []ast.Statement
	for p.currentToken.Type != token.EOF {
		state, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, state)
		p.forward()
	}
	return statements, nil
}

// parseStatement 解析语句
// 函数执行前 currentToken 指向语句的第一个 token，函数执行后 currentToken 指向语句的最后一个token（一般是;）
func (p *Parser) parseStatement(precedence ...int) (ast.Statement, error) {
	defer func() {
		if p.expectPeekToken(token.SEMICOLON) {
			p.forward()
		}
	}()

	// block statement
	if p.expectToken(token.LBRACE) {
		p.forward() // 跳过 {
		var states []ast.Statement
		for !p.expectToken(token.RBRACE) {
			state, err := p.parseStatement(precedence...)
			if err != nil {
				return nil, err
			}
			states = append(states, state)
			p.forward()
		}
		if !p.expectToken(token.RBRACE) {
			return nil, fmt.Errorf("expect }, but got %s on line %d, col: %d", p.currentToken.Value, p.currentToken.Row, p.currentToken.Col)
		}
		return &ast.Block{States: states}, nil
	}

	pre := token.PrecedenceLowest
	if len(precedence) > 0 {
		pre = precedence[0]
	}
	expr, err := p.ParseExpression(pre)
	if err != nil {
		return nil, err
	}
	return &ast.ExpressionStatement{Expr: expr}, nil
}

func (p *Parser) ParseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixFn[p.currentToken.Type]
	if prefix == nil {
		return nil, fmt.Errorf("unrecognized token: %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for !p.expectPeekToken(token.SEMICOLON) && precedence < token.GetPrecedence(p.peekToken.Type) {
		infix := p.infixParseFn[p.peekToken.Type]
		if infix == nil {
			return leftExp, nil
		}

		p.forward()
		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}
	return leftExp, nil
}

func (ep *Parser) registerPrefixFn(t token.TokenType, fn prefixParseFn) {
	ep.prefixFn[t] = fn
}

func (ep *Parser) registerInfixFn(t token.TokenType, fn infixParseFn) {
	ep.infixParseFn[t] = fn
}

// 以下是前缀表达式的解析

func (p *Parser) parseBool() (ast.Expression, error) {
	return &ast.Bool{Token: p.currentToken, Value: &object.Bool{Val: p.currentToken.Value == "true"}}, nil
}

func (p *Parser) parseIdent() (ast.Expression, error) {
	return &ast.Ident{Token: p.currentToken}, nil
}

func (p *Parser) parseNumber() (ast.Expression, error) {
	if strings.Contains(p.currentToken.Value, ".") {
		val, err := strconv.ParseFloat(p.currentToken.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("parseNumber err: %v, value: %v on line %d, col: %d", err, p.currentToken.Value, p.currentToken.Row, p.currentToken.Col)
		}
		return &ast.Number{Token: p.currentToken, Value: &object.Float{Val: val}}, nil
	}

	val, err := strconv.ParseInt(p.currentToken.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parseNumber err: %v, value: %v on line %d, col: %d", err, p.currentToken.Value, p.currentToken.Row, p.currentToken.Col)
	}
	return &ast.Number{Token: p.currentToken, Value: &object.Int{Val: val}}, nil
}

func (p *Parser) parseString() (ast.Expression, error) {
	return &ast.String{Token: p.currentToken, Value: &object.String{Val: []rune(p.currentToken.Value)}}, nil
}

func (p *Parser) parseSlice() (ast.Expression, error) {
	res := &ast.Slice{}
	p.forward() // 跳过 [
	for !p.expectToken(token.RBRACKET) {
		expr, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		p.forward()
		if p.expectToken(token.SEMICOLON) {
			p.forward() // 跳过;
			if res.InitExpr != nil {
				return nil, fmt.Errorf("invalid slice on line %d, col: %d", p.currentToken.Row, p.currentToken.Col)
			}
			res.InitExpr = expr
			continue
		}
		if res.InitExpr != nil {
			if res.LenExpr != nil {
				return nil, fmt.Errorf("invalid slice on line %d, col: %d", p.currentToken.Row, p.currentToken.Col)
			}
			res.LenExpr = expr
			continue
		}
		if p.expectToken(token.COMMA) {
			p.forward() // 跳过,
		}
		res.Data = append(res.Data, expr)
	}
	if !p.expectToken(token.RBRACKET) {
		return nil, fmt.Errorf("expect ], but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseHashMap() (ast.Expression, error) {
	res := &ast.Map{KV: make(map[ast.Expression]ast.Expression)}
	p.forward() // 跳过 {
	for !p.expectToken(token.RBRACE) {
		k, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		p.forward()
		if !p.expectToken(token.COLON) {
			return nil, fmt.Errorf("expect :, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}
		p.forward() // 跳过 :
		v, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		res.KV[k] = v
		p.forward()

		if p.expectToken(token.COMMA) {
			p.forward()
		}
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("expect }, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseFnLiteral() (ast.Expression, error) {
	res := &ast.FnLiteralObj{}
	p.forward() // 跳过 fn
	if !p.expectToken(token.IDENTIFIER) {
		return nil, fmt.Errorf("exepect identifier, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	res.Name = p.currentToken.Value

	p.forward() // 跳过函数名
	if !p.expectToken(token.LPAREN) {
		return nil, fmt.Errorf("exepect (, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 (

	// 解析参数列表
	for !p.expectToken(token.RPAREN) {
		if !p.expectToken(token.IDENTIFIER) {
			return nil, fmt.Errorf("exepect identifier, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}
		res.Args = append(res.Args, p.currentToken.Value)
		p.forward() // 跳过形参
		if p.expectToken(token.COMMA) {
			p.forward() // 跳过 ,
		}
	}

	if !p.expectToken(token.RPAREN) {
		return nil, fmt.Errorf("exepect ), found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 )

	// 解析函数体
	if !p.expectToken(token.LBRACE) {
		return nil, fmt.Errorf("exepect {, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 {

	for !p.expectToken(token.RBRACE) {
		state, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		res.Statements = append(res.Statements, state)
		p.forward()
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("exepect }, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseRuleLiteral() (ast.Expression, error) {
	res := &ast.Rule{}
	p.forward() // 跳过 rule
	if !p.expectToken(token.IDENTIFIER) {
		return nil, fmt.Errorf("exepect identifier, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	res.Name = p.currentToken.Value

	p.forward() // 跳过名字
	// 解析函数体
	if !p.expectToken(token.LBRACE) {
		return nil, fmt.Errorf("exepect {, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 {

	for !p.expectToken(token.RBRACE) {
		state, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		res.Statements = append(res.Statements, state)
		p.forward()
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("exepect }, found %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseBang() (ast.Expression, error) {
	p.forward() // 跳过 !
	expr, err := p.ParseExpression(token.PrecedencePrefix)
	if err != nil {
		return nil, err
	}
	return &ast.Bang{Expr: expr}, nil
}

func (p *Parser) parseNegative() (ast.Expression, error) {
	p.forward() // 跳过 -
	expr, err := p.ParseExpression(token.PrecedencePrefix)
	if err != nil {
		return nil, err
	}
	return &ast.Negative{Expr: expr}, nil
}

func (p *Parser) parsePositive() (ast.Expression, error) {
	p.forward() // 跳过 +
	expr, err := p.ParseExpression(token.PrecedencePrefix)
	if err != nil {
		return nil, err
	}
	return &ast.Positive{Expr: expr}, nil
}

func (p *Parser) parseBitwiseNot() (ast.Expression, error) {
	p.forward() // 跳过 ^
	expr, err := p.ParseExpression(token.PrecedencePrefix)
	if err != nil {
		return nil, err
	}
	return &ast.BitwiseNot{Expr: expr}, nil
}

func (p *Parser) parseIf() (ast.Expression, error) {
	ifs, elseStates, err := p.parseIfStatement()
	if err != nil {
		return nil, err
	}
	return &ast.If{
		Ifs:  ifs,
		Else: elseStates,
	}, nil
}

func (p *Parser) parseIfStatement() ([]ast.ExprStates, []ast.Statement, error) {
	p.forward() // 跳过 if
	expr, err := p.ParseExpression(token.PrecedencePlaceholder)
	if err != nil {
		return nil, nil, err
	}
	p.forward()

	states, err := p.parseBlockStatement()
	if err != nil {
		return nil, nil, err
	}

	ifs := []ast.ExprStates{{Expr: expr, States: states}}
	var elseStates []ast.Statement
	if p.expectPeekToken(token.ELSE) {
		p.forward() // 跳过 }
		p.forward() // 跳过 else
		if p.expectToken(token.LBRACE) {
			elseStates, err = p.parseBlockStatement()
			if err != nil {
				return nil, nil, err
			}
		} else if p.expectToken(token.IF) {
			tmpIfs, tmpElseStates, err := p.parseIfStatement()
			if err != nil {
				return nil, nil, err
			}
			ifs = append(ifs, tmpIfs...)
			elseStates = tmpElseStates
		} else {
			return nil, nil, fmt.Errorf("expect if or }, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}
	}
	return ifs, elseStates, nil
}

func (p *Parser) parseReturn() (ast.Expression, error) {
	p.forward() // 跳过 return
	if p.expectToken(token.SEMICOLON) {
		return nil, nil
	}
	expr, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	return &ast.Return{Expr: expr}, nil
}

func (p *Parser) parseBreak() (ast.Expression, error) {
	p.forward() // 跳过 break
	return &ast.Break{}, nil
}

func (p *Parser) parseContinue() (ast.Expression, error) {
	p.forward() // 跳过 continue
	return &ast.Continue{}, nil
}

func (p *Parser) parseFor() (ast.Expression, error) {
	p.forward() // 跳过 for
	initial, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	res := &ast.For{Initial: initial}

	if !p.expectToken(token.SEMICOLON) {
		return nil, fmt.Errorf("expect ;, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 ;

	condition, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	res.Condition = condition

	p.forward()
	if !p.expectToken(token.SEMICOLON) {
		return nil, fmt.Errorf("expect ;, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 ;

	post, err := p.parseStatement(token.PrecedencePlaceholder)
	if err != nil {
		return nil, err
	}
	res.Post = post
	p.forward()

	states, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	res.Statements = states
	return res, nil
}

func (p *Parser) parseSwitch() (ast.Expression, error) {
	p.forward() // 跳过 switch
	expr, err := p.ParseExpression(token.PrecedencePlaceholder)
	if err != nil {
		return nil, err
	}
	res := &ast.Switch{
		Expr: expr,
	}
	p.forward()

	if !p.expectToken(token.LBRACE) {
		return nil, fmt.Errorf("expect {, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 {

	for p.expectToken(token.CASE) {
		p.forward() // 跳过 case
		expr, err = p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		p.forward()

		var states []ast.Statement
		if !p.expectToken(token.COLON) {
			return nil, fmt.Errorf("expect :, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}

		p.forward() // 跳过 :
		for !p.expectToken(token.CASE) && !p.expectToken(token.DEFAULT) {
			state, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			states = append(states, state)
			p.forward()
		}

		res.Cases = append(res.Cases, ast.ExprStates{
			Expr:   expr,
			States: states,
		})
	}

	if p.expectToken(token.DEFAULT) {
		p.forward() // 跳过 default
		if !p.expectToken(token.COLON) {
			return nil, fmt.Errorf("expect :, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}
		p.forward() // 跳过 :

		for !p.expectToken(token.RBRACE) {
			state, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			res.Default = append(res.Default, state)
			p.forward()
		}
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("expect }, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseGroup() (ast.Expression, error) {
	p.forward() // 跳过 (
	expr, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	res := &ast.Group{Expr: expr}
	p.forward()
	if !p.expectToken(token.RPAREN) {
		return nil, fmt.Errorf("expect ), but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

// parseBlockStatement 解析块语句
// 函数执行前 currentToken 指向{，函数执行后 currentToken 指向}
func (p *Parser) parseBlockStatement() ([]ast.Statement, error) {
	if !p.expectToken(token.LBRACE) {
		return nil, fmt.Errorf("expect {, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 {

	var res []ast.Statement
	for !p.expectToken(token.RBRACE) {
		state, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		res = append(res, state)
		p.forward()
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("expect }, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseStructLiteral() (ast.Expression, error) {
	p.forward() // 跳过 struct
	if !p.expectToken(token.IDENTIFIER) {
		return nil, fmt.Errorf("expect identifier, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}

	stu := &ast.StructLiteral{Name: p.currentToken.Value}
	p.forward() // 跳过 identifier

	if !p.expectToken(token.LBRACE) {
		return nil, fmt.Errorf("expect {, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 {

	for !p.expectToken(token.RBRACE) {
		if !p.expectToken(token.IDENTIFIER) {
			return nil, fmt.Errorf("expect identifier, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}
		stu.Fields = append(stu.Fields, p.currentToken.Value)
		p.forward() // 跳过 identifier

		if p.expectToken(token.COMMA) {
			p.forward() // 跳过 ,
		}
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("expect }, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return stu, nil
}

/*
	impl person {
		fn new(name, age) {
			Person{age, name}
		}
	}
*/
func (p *Parser) parseImpl() (ast.Expression, error) {
	p.forward() // 跳过 impl
	if !p.expectToken(token.IDENTIFIER) {
		return nil, fmt.Errorf("expect identifier, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}

	name := p.currentToken.Value
	p.forward() // 跳过 identifier

	if !p.expectToken(token.LBRACE) {
		return nil, fmt.Errorf("expect {, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	p.forward() // 跳过 {

	res := &ast.Impl{Name: name, Methods: make(map[string]*ast.FnLiteralObj)}
	for !p.expectToken(token.RBRACE) {
		if !p.expectToken(token.FUNC) {
			return nil, fmt.Errorf("expect func, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
		}
		fnLiteral, err := p.parseFnLiteral()
		if err != nil {
			return nil, err
		}
		fn := fnLiteral.(*ast.FnLiteralObj)
		res.Methods[fn.Name] = fn
		p.forward() // 跳过 }
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("expect }, but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

// 以下是中缀表达式的解析

func (p *Parser) parsePlus(left ast.Expression) (ast.Expression, error) {
	p.forward()                                                                                                 // 跳过 +
	if p.expectToken(token.PLUS) && (p.peekToken.Type == token.SEMICOLON || p.peekToken.Type == token.LBRACE) { // ++
		return &ast.Assign{
			Left: left,
			Right: &ast.Plus{
				Left:  left,
				Right: &ast.Number{Value: &object.Int{Val: 1}},
			},
		}, nil
	}

	if p.expectToken(token.ASSIGN) { // +=
		p.forward() // 跳过 =
		right, err := p.ParseExpression(token.PrecedencePlaceholder)
		if err != nil {
			return nil, err
		}
		return &ast.Assign{
			Left: left,
			Right: &ast.Plus{
				Left:  left,
				Right: right,
			},
		}, nil
	}

	right, err := p.ParseExpression(token.PrecedenceAddMinus)
	if err != nil {
		return nil, err
	}
	return &ast.Plus{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseMinus(left ast.Expression) (ast.Expression, error) {
	p.forward()                                                                                                  // 跳过 -
	if p.expectToken(token.MINUS) && (p.peekToken.Type == token.SEMICOLON || p.peekToken.Type == token.LBRACE) { // --
		return &ast.Assign{
			Left: left,
			Right: &ast.Minus{
				Left:  left,
				Right: &ast.Number{Value: &object.Int{Val: 1}},
			},
		}, nil
	}

	if p.expectToken(token.ASSIGN) { // -=
		p.forward() // 跳过 =
		right, err := p.ParseExpression(token.PrecedencePlaceholder)
		if err != nil {
			return nil, err
		}
		return &ast.Assign{
			Left: left,
			Right: &ast.Minus{
				Left:  left,
				Right: right,
			},
		}, nil
	}

	right, err := p.ParseExpression(token.PrecedenceAddMinus)
	if err != nil {
		return nil, err
	}
	return &ast.Minus{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseTimes(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 *

	if p.expectToken(token.ASSIGN) { // *=
		p.forward() // 跳过 =
		right, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		return &ast.Assign{
			Left: left,
			Right: &ast.Times{
				Left:  left,
				Right: right,
			},
		}, nil
	}

	right, err := p.ParseExpression(token.PrecedenceMultiplyDivide)
	if err != nil {
		return nil, err
	}
	return &ast.Times{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseDivide(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 /

	if p.expectToken(token.ASSIGN) { // /=
		p.forward() // 跳过 =
		right, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		return &ast.Assign{
			Left: left,
			Right: &ast.Divide{
				Left:  left,
				Right: right,
			},
		}, nil
	}

	right, err := p.ParseExpression(token.PrecedenceMultiplyDivide)
	if err != nil {
		return nil, err
	}
	return &ast.Divide{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseCall(left ast.Expression) (ast.Expression, error) {
	args := &ast.Slice{}
	res := &ast.Call{Left: left, Arguments: args}

	p.forward() // skip (
	// parse arguments
	for !p.expectToken(token.RPAREN) {
		expr, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		args.Data = append(args.Data, expr)
		p.forward() // skip the last character of the expression.
		if p.expectToken(token.COMMA) {
			p.forward() // skip ,
		}
	}

	if !p.expectToken(token.RPAREN) {
		return nil, fmt.Errorf("expect ), but got %s on line %d, col: %d", p.currentToken.Value, p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseIndex(left ast.Expression) (ast.Expression, error) {
	res := &ast.Index{Data: left}
	p.forward() // 跳过 [

	// 解析 key
	key, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	res.Key = key
	p.forward() // 跳过 key

	if p.expectToken(token.COLON) {
		p.forward() // 跳过 :
		end, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		res.End = end
		p.forward() // 跳过 end
	}

	if !p.expectToken(token.RBRACKET) {
		return nil, fmt.Errorf("expect ], but got %s on line %d, col: %d", p.currentToken.String(), p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) parseGreater(left ast.Expression) (ast.Expression, error) {
	p.forward()                      // 跳过>
	if p.expectToken(token.ASSIGN) { // >=
		p.forward() // 跳过=
		right, err := p.ParseExpression(token.PrecedenceCompare)
		if err != nil {
			return nil, err
		}
		return &ast.Compare{
			Flag:  ast.CompareGreaterEqual,
			Left:  left,
			Right: right,
		}, nil
	}

	right, err := p.ParseExpression(token.PrecedenceCompare)
	if err != nil {
		return nil, err
	}
	return &ast.Compare{
		Flag:  ast.CompareGreaterThan,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseLesser(left ast.Expression) (ast.Expression, error) {
	p.forward()                      // 跳过<
	if p.expectToken(token.ASSIGN) { // <=
		p.forward() // 跳过=
		right, err := p.ParseExpression(token.PrecedenceCompare)
		if err != nil {
			return nil, err
		}
		return &ast.Compare{
			Flag:  ast.CompareLessEqual,
			Left:  left,
			Right: right,
		}, nil
	}

	right, err := p.ParseExpression(token.PrecedenceCompare)
	if err != nil {
		return nil, err
	}
	return &ast.Compare{
		Flag:  ast.CompareLessThan,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseAnd(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 &
	right, err := p.ParseExpression(token.PrecedenceBitwiseAnd)
	if err != nil {
		return nil, err
	}
	return &ast.BitwiseAnd{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseLogicAnd(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 &&
	right, err := p.ParseExpression(token.PrecedenceLogicAnd)
	if err != nil {
		return nil, err
	}
	return &ast.LogicAnd{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseOr(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 |
	right, err := p.ParseExpression(token.PrecedenceBitwiseOr)
	if err != nil {
		return nil, err
	}
	return &ast.Or{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseXor(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 ^
	right, err := p.ParseExpression(token.PrecedenceBitwiseXor)
	if err != nil {
		return nil, err
	}
	return &ast.BitwiseXOR{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseLeftShift(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 <<
	right, err := p.ParseExpression(token.PrecedenceBitShift)
	if err != nil {
		return nil, err
	}
	return &ast.LeftShift{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseRightShift(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 >>
	right, err := p.ParseExpression(token.PrecedenceBitShift)
	if err != nil {
		return nil, err
	}
	return &ast.RightShift{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseLogicOr(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 ||
	right, err := p.ParseExpression(token.PrecedenceLogicOr)
	if err != nil {
		return nil, err
	}
	return &ast.LogicOr{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseMod(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 %
	right, err := p.ParseExpression(token.GetPrecedence(token.MOD))
	if err != nil {
		return nil, err
	}
	return &ast.Mod{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parsePower(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 **
	right, err := p.ParseExpression(token.GetPrecedence(token.POWER))
	if err != nil {
		return nil, err
	}
	return &ast.Power{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseDot(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 .
	right, err := p.ParseExpression(token.GetPrecedence(token.DOT))
	if err != nil {
		return nil, err
	}
	return &ast.Dot{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseAssign(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 =
	right, err := p.ParseExpression(token.PrecedenceLowest)
	if err != nil {
		return nil, err
	}
	return &ast.Assign{
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseEqual(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 ==
	right, err := p.ParseExpression(token.PrecedenceCompare)
	if err != nil {
		return nil, err
	}
	return &ast.Compare{
		Flag:  ast.CompareEqual,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseNotEqual(left ast.Expression) (ast.Expression, error) {
	p.forward() // 跳过 !=
	right, err := p.ParseExpression(token.PrecedenceCompare)
	if err != nil {
		return nil, err
	}
	return &ast.Compare{
		Flag:  ast.CompareNotEqual,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) parseStructInstantiate(left ast.Expression) (ast.Expression, error) {
	_, ok := left.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("not ident on line %d, col: %d", p.currentToken.Row, p.currentToken.Col)
	}

	p.forward() // 跳过 {
	res := &ast.RgStructInstantiate{Ident: left, KV: make(map[ast.Expression]ast.Expression)}
	for !p.expectToken(token.RBRACE) {
		exp, err := p.ParseExpression(token.PrecedenceLowest)
		if err != nil {
			return nil, err
		}
		p.forward()

		if p.expectToken(token.COMMA) {
			p.forward() // 跳过 ,
			res.Values = append(res.Values, exp)
		} else if p.expectToken(token.COLON) {
			p.forward() // 跳过 :
			exp2, err := p.ParseExpression(token.PrecedenceLowest)
			if err != nil {
				return nil, err
			}
			p.forward()
			res.KV[exp] = exp2
		} else if p.expectToken(token.RBRACE) {
			res.Values = append(res.Values, exp)
		}
	}

	if !p.expectToken(token.RBRACE) {
		return nil, fmt.Errorf("expect } but got %s on line %d, col: %d", p.currentToken.Value, p.currentToken.Row, p.currentToken.Col)
	}
	return res, nil
}

func (p *Parser) forward() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.ReadNextToken()
	if p.expectToken(token.INLINE_COMMENTS) {
		p.forward()
		return
	}
	//fmt.Println("curr: ", p.currentToken, " peek: ", p.peekToken)
}

func (p *Parser) expectToken(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) expectPeekToken(t token.TokenType) bool {
	return p.peekToken.Type == t
}
