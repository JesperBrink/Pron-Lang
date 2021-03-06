package parser

import (
	"Pron-Lang/ast"
	"Pron-Lang/lexer"
	"Pron-Lang/token"
	"fmt"
	"strconv"
	"strings"
)

const (
	_ int = iota // Gives priority to the operators
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:   EQUALS,
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.MODULO:   PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

// addPeekError adds an error to the parser about
// the expected token and the token it got
func (p *Parser) addPeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken makes the parser's lexer focus on the next token
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram parses the program that comes from the given lexer.
// Returns an *ast.Program node
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseDirectFunctionStatement()
	case token.CLASS:
		return p.parseClassStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseClassStatement() *ast.ClassStatement {
	stmt := &ast.ClassStatement{Token: p.curToken}

	p.nextToken()

	stmt.Name = p.parseIdentifier().(*ast.Identifier)
	fields := []*ast.VarStatement{}
	functions := []*ast.DirectFunctionStatement{}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) {
		switch p.curToken.Type {
		case token.VAR:
			fields = append(fields, p.parseVarStatement().(*ast.VarStatement))
		case token.FUNCTION:
			functions = append(functions, p.parseDirectFunctionStatement())
		case token.INIT:
			initParams, initBody := p.parseInitFunction()
			stmt.InitParams = initParams
			stmt.InitBody = initBody
		}

		p.nextToken()
	}

	stmt.Fields = fields
	stmt.Functions = functions

	return stmt
}

func (p *Parser) parseInitFunction() ([]*ast.InitParam, *ast.BlockStatement) {
	initParams := []*ast.InitParam{}

	if !p.expectPeek(token.LPAREN) {
		return nil, nil
	}

	if !p.peekTokenIs(token.RPAREN) {
		p.nextToken()

		if p.curTokenIs(token.THIS) {
			p.nextToken()
			p.nextToken()
			ident := p.parseIdentifier().(*ast.Identifier)
			param := &ast.InitParam{Token: p.curToken, Parameter: ident, IsThisParam: true}
			initParams = append(initParams, param)
		} else if p.curTokenIs(token.IDENT) {
			ident := p.parseIdentifier().(*ast.Identifier)
			param := &ast.InitParam{Token: p.curToken, Parameter: ident, IsThisParam: false}
			initParams = append(initParams, param)
		} else {
			return nil, nil
		}

		for p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.nextToken()
			if p.curTokenIs(token.THIS) {
				p.nextToken()
				p.nextToken()
				ident := p.parseIdentifier().(*ast.Identifier)
				param := &ast.InitParam{Token: p.curToken, Parameter: ident, IsThisParam: true}
				initParams = append(initParams, param)
			} else if p.curTokenIs(token.IDENT) {
				ident := p.parseIdentifier().(*ast.Identifier)
				param := &ast.InitParam{Token: p.curToken, Parameter: ident, IsThisParam: false}
				initParams = append(initParams, param)
			} else {
				return nil, nil
			}
		}

		if !p.expectPeek(token.RPAREN) {
			return nil, nil
		}
	} else {
		p.nextToken()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil, nil
	}

	body := p.parseBlockStatement()

	return initParams, body
}

func (p *Parser) parseDirectFunctionStatement() *ast.DirectFunctionStatement {
	stmt := &ast.DirectFunctionStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// If first letter is lower case: private else: public
	firstLetter := stmt.Name.Value[:1]
	if firstLetter == strings.ToLower(firstLetter) {
		stmt.IsPublic = false
	} else {
		stmt.IsPublic = true
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Function.Token = token.Token{Type: token.FUNCTION, Literal: token.FUNCTION}
	stmt.Function.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Function.Body = p.parseBlockStatement()

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)
	// the following makes sure that you can type
	// expressions like 5 + 5 into the REPL (without ; after)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseVarStatement() ast.Statement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	} else {
		stmt.Value = &ast.Null{}
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.addPeekError(t)
	return false
}

func (p *Parser) parseIdentifier() ast.Expression {
	// Check for suffixes that change the context
	if p.peekTokenIs(token.DOT) {
		return p.parseCallObjectFunction()
	} else if p.peekTokenIs(token.INCREMENT) {
		return p.parseIncrement()
	} else if p.peekTokenIs(token.DECREMENT) {
		return p.parseDecrement()
	}

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIncrement() ast.Expression {
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	return &ast.Increment{Token: p.curToken, Name: *name}
}

func (p *Parser) parseDecrement() ast.Expression {
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken()
	return &ast.Decrement{Token: p.curToken, Name: *name}
}

func (p *Parser) parseThisPrefixedIdentifier() ast.Expression {
	p.nextToken()
	p.nextToken()
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal, HasThisPrefix: true}
}

func (p *Parser) parseCallObjectFunction() ast.Expression {
	callObjectFunction := &ast.CallObjectFunction{}

	// UPS: Cannot use the parseIdentifier,
	// because it would end in a infinite loop if we call it here
	callObjectFunction.ObjectName = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken()

	callObjectFunction.Token = p.curToken

	p.nextToken()

	callObjectFunction.FunctionName = p.parseIdentifier().(*ast.Identifier)

	p.nextToken()

	callObjectFunction.Arguments = p.parseExpressionList(token.RPAREN)

	return callObjectFunction
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseRealLiteral() ast.Expression {
	lit := &ast.RealLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if value, ok := precedences[p.peekToken.Type]; ok {
		return value
	}

	return LOWEST
}

func (p *Parser) curPrecendence() int {
	if value, ok := precedences[p.curToken.Type]; ok {
		return value
	}

	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifToken := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	firstCondition := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	firstConsequence := p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		// Else statement
		expression := &ast.IfExpression{Token: p.curToken,
			Condition: firstCondition, Consequence: firstConsequence}
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()

		return expression

	} else if p.peekTokenIs(token.ELIF) {
		// Elif statement
		expression := &ast.ElseIfExpression{Token: p.curToken}
		conditionAndBlockstatements := []*ast.ConditionAndBlockstatementExpression{}

		p.nextToken()

		// add the 'if'
		firstConditionAndBlockstatement := &ast.ConditionAndBlockstatementExpression{Token: ifToken,
			Condition: firstCondition, Consequence: firstConsequence}
		conditionAndBlockstatements = append(conditionAndBlockstatements, firstConditionAndBlockstatement)

		// add all the 'elif'
		for {
			conditionAndBlockstatement := &ast.ConditionAndBlockstatementExpression{Token: p.curToken}

			if !p.expectPeek(token.LPAREN) {
				return nil
			}

			p.nextToken()

			conditionAndBlockstatement.Condition = p.parseExpression(LOWEST)

			if !p.expectPeek(token.RPAREN) {
				return nil
			}

			if !p.expectPeek(token.LBRACE) {
				return nil
			}

			conditionAndBlockstatement.Consequence = p.parseBlockStatement()

			// add the conditionAndBlockstatement to the list
			conditionAndBlockstatements = append(conditionAndBlockstatements, conditionAndBlockstatement)

			if p.peekTokenIs(token.ELSE) {
				// This was the last 'elif', but there is an Else left
				p.nextToken()

				if !p.expectPeek(token.LBRACE) {
					return nil
				}

				expression.Alternative = p.parseBlockStatement()
				expression.ConditionAndBlockstatementList = conditionAndBlockstatements
				return expression
			} else if p.peekTokenIs(token.ELIF) {
				p.nextToken()
			} else {
				break
			}
		}
		expression.ConditionAndBlockstatementList = conditionAndBlockstatements
		return expression
	} else {
		// return simple 'if' without 'else if' or 'else'
		return &ast.IfExpression{Token: p.curToken, Condition: firstCondition, Consequence: firstConsequence}
	}
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

func (p *Parser) parseForloopExpression() ast.Expression {
	curToken := p.curToken

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	localVar := p.parseIdentifier()

	if p.peekTokenIs(token.FROM) {
		// increment forloop
		expression := &ast.IncrementForloopExpression{Token: curToken, LocalVar: localVar}

		p.nextToken()
		p.nextToken()
		expression.From = p.parseExpression(LOWEST)

		if !p.expectPeek(token.TO) {
			return nil
		}

		p.nextToken()
		expression.To = p.parseExpression(LOWEST)

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Body = p.parseBlockStatement()

		return expression

	} else if p.peekTokenIs(token.IN) {
		// array forloop
		expression := &ast.ArrayForloopExpression{Token: curToken, LocalVar: localVar}

		p.nextToken()
		p.nextToken()
		expression.ArrayName = p.parseIdentifier()

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Body = p.parseBlockStatement()

		return expression
	} else {
		return nil
	}
}

func (p *Parser) parseObjectInitialization() ast.Expression {
	objectInitialiation := &ast.ObjectInitialization{Token: p.curToken}

	p.nextToken()

	objectInitialiation.Name = p.parseIdentifier().(*ast.Identifier)

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	objectInitialiation.Arguments = p.parseExpressionList(token.RPAREN)

	return objectInitialiation
}

func (p *Parser) parseBlockComment() ast.Expression {
	for !p.peekTokenIs(token.ENDBLOCKCOMMENT) {
		p.nextToken()
	}

	p.nextToken()

	return &ast.Null{}
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Init both curToken and peekToken
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.REAL, p.parseRealLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)
	p.registerPrefix(token.FOR, p.parseForloopExpression)
	p.registerPrefix(token.NEW, p.parseObjectInitialization)
	p.registerPrefix(token.THIS, p.parseThisPrefixedIdentifier)
	p.registerPrefix(token.STARTBLOCKCOMMENT, p.parseBlockComment)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseInfixExpression)

	return p
}
