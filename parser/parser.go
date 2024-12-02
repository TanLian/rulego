package parser

import (
	"fmt"
	"strconv"

	"github.com/tanlian/rulego/ast"
	"github.com/tanlian/rulego/environment"
	"github.com/tanlian/rulego/lexer"
	"github.com/tanlian/rulego/object"
	"github.com/tanlian/rulego/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(exp ast.Expression) ast.Expression
)

type Parser struct {
	l   *lexer.Lexer
	env *environment.Environment

	prefixFn     map[token.TokenType]prefixParseFn
	infixParseFn map[token.TokenType]infixParseFn

	// token 相关
	currentToken token.Token
	peekToken    token.Token
}

func NewParser(l *lexer.Lexer, env *environment.Environment) *Parser {
	p := &Parser{
		l:            l,
		env:          env,
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
	p.registerPrefixFn(token.MINUS, p.parseNegative)
	p.registerPrefixFn(token.IF, p.parseIf)
	p.registerPrefixFn(token.FOR, p.parseFor)
	p.registerPrefixFn(token.SWITCH, p.parseSwitch)

	p.registerInfixFn(token.PLUS, p.parsePlus)
	p.registerInfixFn(token.MINUS, p.parseMinus)
	p.registerInfixFn(token.TIMES, p.parseTimes)
	p.registerInfixFn(token.DIVIDE, p.parseDivide)
	p.registerInfixFn(token.LPAREN, p.parseCall)
	p.registerInfixFn(token.LBRACKET, p.parseIndex)
	p.registerInfixFn(token.GREATER, p.parseCompare)
	p.registerInfixFn(token.GREATER_EQUAL, p.parseCompare)
	p.registerInfixFn(token.EQUAL, p.parseCompare)
	p.registerInfixFn(token.LESS, p.parseCompare)
	p.registerInfixFn(token.LESS_EQUAL, p.parseCompare)
	p.registerInfixFn(token.NOT_EQUAL, p.parseCompare)
	p.registerInfixFn(token.PLUS_ASSIGN, p.parsePlusAssign)
	p.registerInfixFn(token.MINUS_ASSIGN, p.parseMinusAssign)
	p.registerInfixFn(token.TIMES_ASSIGN, p.parseTimesAssign)
	p.registerInfixFn(token.DIVIDE_ASSIGN, p.parseDivideAssign)
	p.registerInfixFn(token.PLUS_PLUS, p.parsePlusPlus)
	p.registerInfixFn(token.AND, p.parseAnd)
	p.registerInfixFn(token.OR, p.parseOr)
	p.registerInfixFn(token.MOD, p.parseMod)
	p.registerInfixFn(token.DOT, p.parseDot)

	p.forward()
	p.forward()

	return p
}

func (p *Parser) Parse() []ast.Statement {
	var statements []ast.Statement
	for p.currentToken.Type != token.EOF {
		s, err := p.parseStatement()
		if err != nil {
			panic(err)
		}
		statements = append(statements, s)
		p.forward()
	}
	return statements
}

// parseStatement 解析语句
// 函数执行前 currentToken 指向语句的第一个token，函数执行后 currentToken 指向语句的最后一个token（一般是;）
func (p *Parser) parseStatement() (ast.Statement, error) {
	defer func() {
		if p.expectPeekToken(token.SEMICOLON) {
			p.forward()
		}
	}()

	if p.expectToken(token.RETURN) {
		p.forward() // 跳过 return
		return &ast.Return{Expr: p.ParseExpression(token.PrecedenceLowest)}, nil
	}

	if p.expectToken(token.BREAK) {
		p.forward() // 跳过 break
		return &ast.Break{}, nil
	}

	expr := p.ParseExpression(token.PrecedenceLowest)
	if p.expectPeekToken(token.ASSIGN) {
		p.forward()
		p.forward() // 跳过 =
		return &ast.Assign{Left: expr, Right: p.ParseExpression(token.PrecedenceLowest)}, nil
	} else if p.expectPeekToken(token.PLUS_PLUS) {
		p.forward()
		return &ast.PlusPlus{Left: expr}, nil
	}
	return &ast.ExpressionStatement{Expr: expr}, nil
}

func (p *Parser) ParseExpression(precedence int) ast.Expression {
	prefix := p.prefixFn[p.currentToken.Type]
	if prefix == nil {
		panic(fmt.Errorf("unrecognized token: %s", p.currentToken.String()))
	}
	leftExp := prefix()

	for !p.expectPeekToken(token.SEMICOLON) && precedence < token.GetPrecedence(p.peekToken.Type) {
		infix := p.infixParseFn[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.forward()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (ep *Parser) registerPrefixFn(t token.TokenType, fn prefixParseFn) {
	ep.prefixFn[t] = fn
}

func (ep *Parser) registerInfixFn(t token.TokenType, fn infixParseFn) {
	ep.infixParseFn[t] = fn
}

// 以下是前缀表达式的解析

func (p *Parser) parseBool() ast.Expression {
	return &ast.Bool{Token: p.currentToken, Value: &object.Bool{Val: p.currentToken.Value == "true"}}
}

func (p *Parser) parseIdent() ast.Expression {
	return &ast.Ident{Token: p.currentToken}
}

func (p *Parser) parseNumber() ast.Expression {
	val, err := strconv.ParseFloat(p.currentToken.Value, 64)
	if err != nil {
		panic(fmt.Sprintf("parseNumber err: %v, value: %v", err, p.currentToken.Value))
	}
	return &ast.Number{Token: p.currentToken, Value: &object.Float{Val: val}}
}

func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.currentToken, Value: &object.String{Val: []rune(p.currentToken.Value)}}
}

func (p *Parser) parseSlice() ast.Expression {
	res := &ast.Slice{}
	p.forward() // 跳过 [
	for !p.expectToken(token.RBRACKET) {
		expr := p.ParseExpression(token.PrecedenceLowest)
		p.forward()
		if p.expectToken(token.COMMA) {
			p.forward() // 跳过,
		}
		res.Data = append(res.Data, expr)
	}
	if !p.expectToken(token.RBRACKET) {
		panic(fmt.Sprintf("expect ], but got %s", p.currentToken.String()))
	}
	return res
}

func (p *Parser) parseHashMap() ast.Expression {
	res := &ast.Map{KV: make(map[ast.Expression]ast.Expression)}
	p.forward() // 跳过 {
	for !p.expectToken(token.RBRACE) {
		k := p.ParseExpression(token.PrecedenceLowest)
		p.forward()
		if !p.expectToken(token.COLON) {
			panic(fmt.Errorf("expect :, but got %s", p.currentToken.String()))
		}
		p.forward() // 跳过 :
		v := p.ParseExpression(token.PrecedenceLowest)
		res.KV[k] = v
		p.forward()

		if p.expectToken(token.COMMA) {
			p.forward()
		}
	}

	if !p.expectToken(token.RBRACE) {
		panic(fmt.Errorf("expect }, but got %s", p.currentToken.String()))
	}
	return res
}

func (p *Parser) parseFnLiteral() ast.Expression {
	res := &ast.FnLiteralObj{}
	p.forward() // 跳过 fn
	if !p.expectToken(token.IDENTIFIER) {
		panic(fmt.Sprintf("exepect identifier, found %s", p.currentToken.String()))
	}
	res.Name = p.currentToken.Value

	p.forward() // 跳过函数名
	if !p.expectToken(token.LPAREN) {
		panic(fmt.Sprintf("exepect (, but got %s", p.currentToken.String()))
	}
	p.forward() // 跳过 (

	// 解析参数列表
	for !p.expectToken(token.RPAREN) {
		if !p.expectToken(token.IDENTIFIER) {
			panic(fmt.Sprintf("exepect identifier, found %s", p.currentToken.String()))
		}
		res.Args = append(res.Args, p.currentToken.Value)
		p.forward() // 跳过形参
		if p.expectToken(token.COMMA) {
			p.forward() // 跳过 ,
		}
	}

	if !p.expectToken(token.RPAREN) {
		panic(fmt.Sprintf("exepect ), found %s", p.currentToken.String()))
	}
	p.forward() // 跳过 )

	// 解析函数体
	if !p.expectToken(token.LBRACE) {
		panic(fmt.Sprintf("exepect {, found %s", p.currentToken.String()))
	}
	p.forward() // 跳过 {

	for !p.expectToken(token.RBRACE) {
		state, err := p.parseStatement()
		if err != nil {
			panic(err)
		}
		res.Statements = append(res.Statements, state)
		p.forward()
	}

	if !p.expectToken(token.RBRACE) {
		panic(fmt.Sprintf("exepect }, found %s", p.currentToken.String()))
	}
	return res
}

func (p *Parser) parseRuleLiteral() ast.Expression {
	res := &ast.Rule{}
	p.forward() // 跳过 rule
	if !p.expectToken(token.IDENTIFIER) {
		panic(fmt.Sprintf("exepect identifier, found %s", p.currentToken.String()))
	}
	res.Name = p.currentToken.Value

	p.forward() // 跳过名字
	// 解析函数体
	if !p.expectToken(token.LBRACE) {
		panic(fmt.Sprintf("exepect {, found %s", p.currentToken.String()))
	}
	p.forward() // 跳过 {

	for !p.expectToken(token.RBRACE) {
		state, err := p.parseStatement()
		if err != nil {
			panic(err)
		}
		res.Statements = append(res.Statements, state)
		p.forward()
	}

	if !p.expectToken(token.RBRACE) {
		panic(fmt.Sprintf("exepect }, found %s", p.currentToken.String()))
	}
	return res
}

func (p *Parser) parseBang() ast.Expression {
	p.forward() // 跳过 !
	return &ast.Bang{Expr: p.ParseExpression(token.PrecedencePrefix)}
}

func (p *Parser) parseNegative() ast.Expression {
	p.forward() // 跳过 -
	return &ast.Negative{Expr: p.ParseExpression(token.PrecedencePrefix)}
}

func (p *Parser) parseIf() ast.Expression {
	ifs, elseStates := p.parseIfStatement()
	return &ast.If{
		Ifs:  ifs,
		Else: elseStates,
	}
}

func (p *Parser) parseIfStatement() ([]ast.ExprStates, []ast.Statement) {
	p.forward() // 跳过 if
	expr := p.ParseExpression(token.PrecedenceLowest)
	p.forward()

	states := p.parseBlockStatement()

	ifs := []ast.ExprStates{{Expr: expr, States: states}}
	var elseStates []ast.Statement
	if p.expectPeekToken(token.ELSE) {
		p.forward() // 跳过 }
		p.forward() // 跳过 else
		if p.expectToken(token.LBRACE) {
			elseStates = p.parseBlockStatement()
		} else if p.expectToken(token.IF) {
			tmpIfs, tmpElseStates := p.parseIfStatement()
			ifs = append(ifs, tmpIfs...)
			elseStates = tmpElseStates
		} else {
			panic(fmt.Sprintf("expect if or }, but got %s", p.currentToken.String()))
		}
	}
	return ifs, elseStates
}

func (p *Parser) parseFor() ast.Expression {
	p.forward() // 跳过 for
	res := &ast.For{}

	initial, err := p.parseStatement()
	if err != nil {
		panic(err)
	}
	res.Initial = initial

	if !p.expectToken(token.SEMICOLON) {
		panic(fmt.Sprintf("expect ;, but got %s", p.currentToken.String()))
	}
	p.forward() // 跳过 ;

	res.Condition = p.ParseExpression(token.PrecedenceLowest)
	p.forward()
	if !p.expectToken(token.SEMICOLON) {
		panic(fmt.Sprintf("expect ;, but got %s", p.currentToken.String()))
	}
	p.forward() // 跳过 ;

	post, err := p.parseStatement()
	if err != nil {
		panic(err)
	}
	res.Post = post
	p.forward()

	res.Statements = p.parseBlockStatement()
	return res
}

func (p *Parser) parseSwitch() ast.Expression {
	p.forward() // 跳过 switch
	res := &ast.Switch{
		Expr: p.ParseExpression(token.PrecedenceLowest),
	}
	p.forward()

	if !p.expectToken(token.LBRACE) {
		panic(fmt.Sprintf("expect {, but got %s", p.currentToken.String()))
	}
	p.forward() // 跳过 {

	for p.expectToken(token.CASE) {
		p.forward() // 跳过 case
		expr := p.ParseExpression(token.PrecedenceLowest)
		p.forward()

		var states []ast.Statement
		if !p.expectToken(token.COLON) {
			panic(fmt.Sprintf("expect :, but got %s", p.currentToken.String()))
		}

		p.forward() // 跳过 :
		for !p.expectToken(token.CASE) && !p.expectToken(token.DEFAULT) {
			state, err := p.parseStatement()
			if err != nil {
				panic(err)
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
			panic(fmt.Sprintf("expect :, but got %s", p.currentToken.String()))
		}
		p.forward() // 跳过 :

		for !p.expectToken(token.RBRACE) {
			state, err := p.parseStatement()
			if err != nil {
				panic(err)
			}
			res.Default = append(res.Default, state)
			p.forward()
		}
	}

	if !p.expectToken(token.RBRACE) {
		panic(fmt.Sprintf("expect }, but got %s", p.currentToken.String()))
	}
	return res
}

func (p *Parser) parseGroup() ast.Expression {
	p.forward() // 跳过 (
	res := &ast.Group{Expr: p.ParseExpression(token.PrecedenceLowest)}
	p.forward()
	if !p.expectToken(token.RPAREN) {
		panic(fmt.Sprintf("expect ), but got %s", p.currentToken.String()))
	}
	return res
}

// parseBlockStatement 解析块语句
// 函数执行前 currentToken 指向{，函数执行后 currentToken 指向}
func (p *Parser) parseBlockStatement() []ast.Statement {
	if !p.expectToken(token.LBRACE) {
		panic(fmt.Sprintf("expect {, but got %s", p.currentToken.String()))
	}
	p.forward() // 跳过 {

	var res []ast.Statement
	for !p.expectToken(token.RBRACE) {
		state, err := p.parseStatement()
		if err != nil {
			panic(err)
		}
		res = append(res, state)
		p.forward()
	}

	if !p.expectToken(token.RBRACE) {
		panic(fmt.Sprintf("expect }, but got %s", p.currentToken.String()))
	}
	return res
}

// 以下是中缀表达式的解析

func (p *Parser) parsePlus(left ast.Expression) ast.Expression {
	res := &ast.Plus{Left: left}
	p.forward() // 跳过 +
	res.Right = p.ParseExpression(token.PrecedenceAddMinus)
	return res
}

func (p *Parser) parseMinus(left ast.Expression) ast.Expression {
	res := &ast.Minus{Left: left}
	p.forward() // 跳过 -
	res.Right = p.ParseExpression(token.PrecedenceAddMinus)
	return res
}

func (p *Parser) parseTimes(left ast.Expression) ast.Expression {
	res := &ast.Times{Left: left}
	p.forward() // 跳过 *
	res.Right = p.ParseExpression(token.PrecedenceMultiplyDivide)
	return res
}

func (p *Parser) parseDivide(left ast.Expression) ast.Expression {
	res := &ast.Divide{Left: left}
	p.forward() // 跳过 /
	res.Right = p.ParseExpression(token.PrecedenceMultiplyDivide)
	return res
}

func (p *Parser) parseCall(left ast.Expression) ast.Expression {
	args := &ast.Slice{}
	res := &ast.Call{Left: left, Arguments: args}

	p.forward() // skip (
	// parse arguments
	for !p.expectToken(token.RPAREN) {
		expr := p.ParseExpression(token.PrecedenceLowest)
		args.Data = append(args.Data, expr)
		p.forward() // skip the last character of the expression.
		if p.expectToken(token.COMMA) {
			p.forward() // skip ,
		}
	}

	if !p.expectToken(token.RPAREN) {
		panic(fmt.Sprintf("expect ), but got %s", p.currentToken.Value))
	}
	return res
}

func (p *Parser) parseIndex(left ast.Expression) ast.Expression {
	res := &ast.Index{Data: left}
	p.forward() // 跳过 [

	// 解析 key
	res.Key = p.ParseExpression(token.PrecedenceLowest)
	p.forward() // 跳过 key

	if p.expectToken(token.COLON) {
		p.forward() // 跳过 :
		res.End = p.ParseExpression(token.PrecedenceLowest)
		p.forward() // 跳过 end
	}

	if !p.expectToken(token.RBRACKET) {
		panic(fmt.Sprintf("expect ], but got %s", p.currentToken.String()))
	}
	return res
}

func (p *Parser) parseCompare(left ast.Expression) ast.Expression {
	res := &ast.Compare{
		Left:  left,
		Token: p.currentToken,
	}
	p.forward() // 跳过比较运算符
	res.Right = p.ParseExpression(token.PrecedenceCompare)
	return res
}

func (p *Parser) parsePlusAssign(left ast.Expression) ast.Expression {
	res := &ast.PlusAssign{Left: left}
	p.forward() // 跳过 +=
	res.Right = p.ParseExpression(token.PrecedenceLowest)
	return res
}

func (p *Parser) parseMinusAssign(left ast.Expression) ast.Expression {
	res := &ast.MinusAssign{Left: left}
	p.forward() // 跳过 -=
	res.Right = p.ParseExpression(token.PrecedenceLowest)
	return res
}

func (p *Parser) parseTimesAssign(left ast.Expression) ast.Expression {
	res := &ast.TimesAssign{Left: left}
	p.forward() // 跳过 *=
	res.Right = p.ParseExpression(token.PrecedenceLowest)
	return res
}

func (p *Parser) parseDivideAssign(left ast.Expression) ast.Expression {
	res := &ast.DivideAssign{Left: left}
	p.forward() // 跳过 /=
	res.Right = p.ParseExpression(token.PrecedenceLowest)
	return res
}

func (p *Parser) parsePlusPlus(left ast.Expression) ast.Expression {
	res := &ast.PlusPlus{Left: left}
	p.forward() // 跳过 ++
	return res
}

func (p *Parser) parseAnd(left ast.Expression) ast.Expression {
	res := &ast.And{Left: left}
	p.forward() // 跳过 &&
	res.Right = p.ParseExpression(token.PrecedenceAnd)
	return res
}

func (p *Parser) parseOr(left ast.Expression) ast.Expression {
	res := &ast.Or{Left: left}
	p.forward() // 跳过 ||
	res.Right = p.ParseExpression(token.PrecedenceOr)
	return res
}

func (p *Parser) parseMod(left ast.Expression) ast.Expression {
	res := &ast.Mod{Left: left}
	p.forward() // 跳过 %
	res.Right = p.ParseExpression(token.GetPrecedence(token.MOD))
	return res
}

func (p *Parser) parseDot(left ast.Expression) ast.Expression {
	res := &ast.Dot{Left: left}
	p.forward() // 跳过 .
	res.Right = p.ParseExpression(token.GetPrecedence(token.DOT))
	return res
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
