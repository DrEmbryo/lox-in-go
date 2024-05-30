package grammar

import "fmt"

type Parser struct {
	tokens  []Token
	current uint32
}

func (parser Parser) prev() Token {
	if parser.current > 0 {
		prevToken := parser.tokens[parser.current-1]
		return prevToken
	}
	return Token{}
}

func (parser *Parser) next() Token {
	if parser.current < uint32(len(parser.tokens)) {
		nextToken := parser.tokens[parser.current]
		parser.current++
		return nextToken
	}
	return Token{}
}

func (parser Parser) lookahead() Token {
	if parser.current < uint32(len(parser.tokens)) {
		return parser.tokens[parser.current]
	}
	return Token{}
}

func (parser *Parser) init(tokens []Token, current uint32) {
	parser.tokens = tokens
	parser.current = current
}

func (parser *Parser) matchToken(tokenTypes ...int8) bool {
	for _, tokenType := range tokenTypes {
		if tokenType != EOF && parser.lookahead().TokenType == tokenType {
			parser.next()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(tokenType int8, message string) (LoxError) {
	fmt.Println(tokenType)
	fmt.Println(parser.lookahead())
	if parser.lookahead().TokenType == tokenType {
		parser.current++
		return nil
	}
	return  ParserError{Token: parser.lookahead(), Message: message, Position: parser.current}
}

func (parser Parser) Parse(tokens []Token) (Expression, LoxError) {
	parser.init(tokens, 0)
	return parser.expression()
}

func (parser *Parser) expression() (Expression, LoxError) {
	return parser.equality()
}

func (parser *Parser) equality() (Expression, LoxError) {
	leftExpr, err := parser.comparison()

	for parser.matchToken(BANG, EQUAL_EQUAL) {
		operator := parser.prev()
		rightExpr, err := parser.comparison()
		leftExpr = BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	} 
	return leftExpr, err
}

func (parser *Parser) comparison() (Expression, LoxError) {
	leftExpr, err := parser.term()

	for  parser.matchToken(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := parser.prev()
		rightExpr, err := parser.term()
		leftExpr = BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}

	return leftExpr, err
}

func (parser *Parser) term() (Expression, LoxError) {
	leftExpr, err := parser.factor()

	for parser.matchToken(MINUS, PLUS) {
		operator := parser.prev()
		rightExpr, err := parser.factor()
		leftExpr = BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}

	return leftExpr, err
}

func (parser *Parser) factor() (Expression, LoxError) {
	leftExpr, err := parser.unary()

	for parser.matchToken(SLASH, STAR) {
		operator := parser.prev()
		rightExpr, err := parser.unary()
		leftExpr = BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}

	return leftExpr, err
}

func (parser *Parser) unary() (Expression, LoxError) {
	if parser.matchToken(BANG, MINUS) {
		operator := parser.prev()
		rightExpr, err := parser.unary()
		return UnaryExpression{Right: rightExpr, Operator: operator}, err
	}
	return parser.primary()
}

func (parser *Parser) primary() (Expression, LoxError) {

	switch {
	case parser.matchToken(FALSE):
		return LiteralExpression{Literal: false}, nil
	case parser.matchToken(TRUE):
		return LiteralExpression{Literal: true}, nil
	case parser.matchToken(NULL):
		return LiteralExpression{Literal: nil}, nil
	case parser.matchToken(NUMBER, STRING):
		return LiteralExpression{Literal: parser.prev().Literal}, nil
	case parser.matchToken(LEFT_PAREN):
		expr, _ := parser.expression()
		return GroupingExpression{Expression: expr}, parser.consume(RIGHT_PAREN, "Expect ')' after expression.")
	}
	parser.sync()
	return nil, ParserError{Position: parser.current, Message: "Expect expression."}
}

func (parser *Parser) sync() {
	parser.current++

	for {
		switch parser.prev().TokenType {
		case EOF:
		case SEMICOLON:
			return
		}

		switch parser.lookahead().TokenType {
		case CLASS:
			fallthrough
		case FUNC:
			fallthrough
		case VAR:
			fallthrough
		case FOR:
			fallthrough
		case IF:
			fallthrough
		case WHILE:
			fallthrough
		case PRINT:
			fallthrough
		case RETURN:
			return
		}

		parser.current++
	}
}