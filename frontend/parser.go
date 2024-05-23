package Lox

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

func (parser Parser) matchToken(tokenTypes ...int8) bool {
	for _, tokenType := range tokenTypes {
		if tokenType != EOF && parser.lookahead().tokenType == tokenType {
			parser.next()
			return true
		}
	}
	return false
}

func (parser *Parser) consume(tokenType int8, err string) {

}

func (parser Parser) Parse(tokens []Token) {
	// parserError := make([]LoxError, 0)
	parser.init(tokens, 0)

}

func (parser *Parser) expression() Expression {
	return parser.equality()
}

func (parser *Parser) equality() Expression {
	leftExpr := parser.comparison()

	for {
		if parser.matchToken(BANG, EQUAL_EQUAL) {
			operator := parser.prev()
			rightExpr := parser.comparison()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
		} else {
			break
		}
	}

	return leftExpr
}

func (parser *Parser) comparison() Expression {
	leftExpr := parser.term()

	for {
		if parser.matchToken(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
			operator := parser.prev()
			rightExpr := parser.term()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr
		} else {
			break
		}
	}

	return leftExpr
}

func (parser *Parser) term() Expression {
	leftExpr := parser.factor()

	for {
		if parser.matchToken(MINUS, PLUS) {
			operator := parser.prev()
			rightExpr := parser.factor()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr
		} else {
			break
		}
	}

	return leftExpr
}

func (parser *Parser) factor() Expression {
	leftExpr := parser.unary()

	for {
		if parser.matchToken(SLASH, STAR) {
			operator := parser.prev()
			rightExpr := parser.unary()
			leftExpr = BinaryExpression{left: leftExpr, right: rightExpr, operator: operator}
			return leftExpr
		} else {
			break
		}
	}

	return leftExpr
}

func (parser *Parser) unary() Expression {
	if parser.matchToken(BANG, MINUS) {
		operator := parser.prev()
		rightExpr := parser.unary()
		return UnaryExpression{right: rightExpr, operator: operator}
	}
	return parser.primary()
}

func (parser *Parser) primary() Expression {
	switch {
	case parser.matchToken(FALSE):
		return LiteralExpression{literal: false}
	case parser.matchToken(TRUE):
		return LiteralExpression{literal: true}
	case parser.matchToken(NULL):
		return LiteralExpression{literal: nil}
	case parser.matchToken(NUMBER, STRING):
		return LiteralExpression{literal: parser.prev().literal}
	case parser.matchToken(LEFT_PAREN):
		expr := parser.expression()
		parser.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return GroupingExpression{expression: expr}
	}
	return nil
}