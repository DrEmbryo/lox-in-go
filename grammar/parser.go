package Lox

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
		if tokenType != EOF && parser.lookahead().tokenType == tokenType {
			parser.next()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(tokenType int8, message string) (LoxError) {
	fmt.Println(tokenType)
	fmt.Println(parser.lookahead())
	if parser.lookahead().tokenType == tokenType {
		parser.current++
		return nil
	}
	return  ParserError{token: parser.lookahead(), message: message, position: parser.current}
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

	for {
		if parser.matchToken(BANG, EQUAL_EQUAL) {
			operator := parser.prev()
			rightExpr, err := parser.comparison()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr, err
		} else {
			break
		}
	}

	return leftExpr, err
}

func (parser *Parser) comparison() (Expression, LoxError) {
	leftExpr, err := parser.term()

	for {
		if parser.matchToken(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
			operator := parser.prev()
			rightExpr, err := parser.term()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr, err
		} else {
			break
		}
	}

	return leftExpr, err
}

func (parser *Parser) term() (Expression, LoxError) {
	leftExpr, err := parser.factor()

	for {
		if parser.matchToken(MINUS, PLUS) {
			operator := parser.prev()
			rightExpr, err := parser.factor()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr, err
		} else {
			break
		}
	}

	return leftExpr, err
}

func (parser *Parser) factor() (Expression, LoxError) {
	leftExpr, err := parser.unary()

	for {
		if parser.matchToken(SLASH, STAR) {
			operator := parser.prev()
			rightExpr, err := parser.unary()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr, err
		} else {
			break
		}
	}

	return leftExpr, err
}

func (parser *Parser) unary() (Expression, LoxError) {
	if parser.matchToken(BANG, MINUS) {
		operator := parser.prev()
		rightExpr, err := parser.unary()
		return UnaryExpression{right: rightExpr, operator: operator}, err
	}
	return parser.primary()
}

func (parser *Parser) primary() (Expression, LoxError) {

	switch {
	case parser.matchToken(FALSE):
		return LiteralExpression{literal: false}, nil
	case parser.matchToken(TRUE):
		return LiteralExpression{literal: true}, nil
	case parser.matchToken(NULL):
		return LiteralExpression{literal: nil}, nil
	case parser.matchToken(NUMBER, STRING):
		return LiteralExpression{literal: parser.prev().literal}, nil
	case parser.matchToken(LEFT_PAREN):
		expr, _ := parser.expression()
		return GroupingExpression{expression: expr}, parser.consume(RIGHT_PAREN, "Expect ')' after expression.")
	}

	return nil, ParserError{position: parser.current, message: "Expect expression."}
}

func (parser *Parser) sync() {
	parser.current++

	for {
		switch parser.prev().tokenType {
		case EOF:
		case SEMICOLON:
			return
		}

		switch parser.lookahead().tokenType {
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