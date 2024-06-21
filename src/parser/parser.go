package parser

import (
	"slices"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Parser struct {
	Tokens  []grammar.Token
	current int
}

func (parser *Parser) consume() grammar.Token {
	token := parser.Tokens[parser.current]
	parser.current++
	return token
}

func (parser *Parser) lookahead() grammar.Token {
	token := parser.Tokens[parser.current]
	return token
}

func (parser *Parser) lookbehind() grammar.Token {
	token := parser.Tokens[parser.current - 1]
	return token
}

func (parser *Parser) expect(tokenType int, message string) grammar.LoxError {
	if tokenType != grammar.EOF && parser.lookahead().TokenType == tokenType{
		parser.current++
		return nil
	}
	return ParserError{Token: parser.lookahead(), Message: message, Position: parser.current}
}

func (parser *Parser) matchToken(tokenTypes ...int) bool {
	for _, tokenType := range tokenTypes {
		if tokenType != grammar.EOF && parser.lookahead().TokenType == tokenType {
			parser.consume()
			return true
		} else {
			return false
		}
	}
	return false
}

func (parser Parser) Parse() ([]grammar.Statement, grammar.LoxError) {
	statements := make([]grammar.Statement, 0)

	if len(parser.Tokens) == 0 {
		return statements, ParserError{Position: 0, Message: "source contains 0 tokens"}
	}

	current := parser.consume()
	for  current.TokenType != grammar.EOF {
		stmt, err := parser.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
		current = parser.consume()
	}

	return statements, nil
}

func (parser *Parser) statement() (grammar.Statement, grammar.LoxError) {
	switch {
	case parser.matchToken(grammar.PRINT):
		return parser.printStatment()
	default:
		return parser.expressionStatement()
	}
}

func (parser *Parser) printStatment() (grammar.Statement, grammar.LoxError) {
	value, err := parser.expression()
	if err != nil {
		return nil, err
	}

	return grammar.PrintStatment{Value: value}, nil
}

func (parser *Parser) expressionStatement() (grammar.Statement, grammar.LoxError) {
	expr, err := parser.expression()
	if err != nil {
		return nil, err
	}

	return grammar.ExpressionStatement{Expression: expr}, nil
 }

func (parser *Parser) expression() (grammar.Expression, grammar.LoxError) {
	return parser.equality()
}

func (parser *Parser) equality() (grammar.Expression, grammar.LoxError) {
	leftExpr, err := parser.comparison()

	for parser.matchToken(grammar.BANG, grammar.EQUAL_EQUAL) {
		operator := parser.lookbehind()
		rightExpr, err := parser.comparison()
		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}
	return leftExpr, err
}

func (parser *Parser) comparison() (grammar.Expression, grammar.LoxError) {
	leftExpr, err := parser.term()

	for parser.matchToken(grammar.GREATER, grammar.GREATER_EQUAL, grammar.LESS, grammar.LESS_EQUAL) {
		operator := parser.lookbehind()
		rightExpr, err := parser.term()
		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}

	return leftExpr, err
}

func (parser *Parser) term() (grammar.Expression, grammar.LoxError) {
	leftExpr, err := parser.factor()

	for parser.matchToken(grammar.MINUS, grammar.PLUS) {
		operator := parser.lookbehind()
		rightExpr, err := parser.factor()
		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}

	return leftExpr, err
}

func (parser *Parser) factor() (grammar.Expression, grammar.LoxError) {
	leftExpr, err := parser.unary()

	for parser.matchToken(grammar.SLASH, grammar.STAR) {
		operator := parser.lookbehind()
		rightExpr, err := parser.unary()
		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
		return leftExpr, err
	}

	return leftExpr, err
}

func (parser *Parser) unary() (grammar.Expression, grammar.LoxError) {
	if parser.matchToken(grammar.BANG, grammar.MINUS) {
		operator := parser.lookbehind()
		rightExpr, err := parser.unary()
		return grammar.UnaryExpression{Right: rightExpr, Operator: operator}, err
	}
	return parser.primary()
}

func (parser *Parser) primary() (grammar.Expression, grammar.LoxError) {

	switch {
	case parser.matchToken(grammar.FALSE):
		return grammar.LiteralExpression{Literal: false}, nil
	case parser.matchToken(grammar.TRUE):
		return grammar.LiteralExpression{Literal: true}, nil
	case parser.matchToken(grammar.NULL):
		return grammar.LiteralExpression{Literal: nil}, nil
	case parser.matchToken(grammar.NUMBER, grammar.STRING):
		return grammar.LiteralExpression{Literal: parser.lookbehind().Lexeme}, nil
	case parser.matchToken(grammar.LEFT_PAREN):
		expr, _ := parser.expression()
		return grammar.GroupingExpression{Expression: expr}, parser.expect(grammar.RIGHT_PAREN, "Expect ')' after expression.")
	}
	parser.sync()
	return nil, ParserError{Position: parser.current, Message: "Expect expression.", Token: parser.Tokens[parser.current]}
}

func (parser *Parser) sync() {
	for {
		switch {
		case parser.lookbehind().TokenType == grammar.SEMICOLON:
			fallthrough
		case slices.Contains(grammar.SYNC_TOKENS, parser.lookahead().TokenType):
			fallthrough
		case parser.lookahead().TokenType == grammar.EOF:
			return
		}

		parser.consume()
	}
}





// type Parser struct {
// 	Tokens  []grammar.Token
// 	current int
// }

// func (parser Parser) prev() grammar.Token {
// 	if parser.current == 0 {
// 		return parser.Tokens[parser.current]
// 	}
// 	if parser.current < (len(parser.Tokens)) {
// 		return parser.Tokens[parser.current-1]
// 	}
// 	return grammar.Token{TokenType: grammar.EOF}
// }

// func (parser *Parser) next() grammar.Token {
// 	if parser.current < (len(parser.Tokens)) {
// 		nextToken := parser.Tokens[parser.current]
// 		parser.current++
// 		return nextToken
// 	}
// 	return grammar.Token{TokenType: grammar.EOF}
// }

// func (parser Parser) lookahead() grammar.Token {
// 	if parser.current < (len(parser.Tokens)) {
// 		return parser.Tokens[parser.current]
// 	}
// 	return grammar.Token{TokenType: grammar.EOF}
// }

// func (parser *Parser) matchToken(tokenTypes ...int) bool {
// 	for _, tokenType := range tokenTypes {
// 		if tokenType != grammar.EOF && parser.lookahead().TokenType == tokenType {
// 			parser.next()
// 			return true
// 		}
// 	}
// 	return false
// }

// func (parser *Parser) consume(tokenType int, message string) grammar.LoxError {
// 	if tokenType != grammar.EOF && parser.lookahead().TokenType == tokenType {
// 		parser.next()
// 		return nil
// 	}
// 	return ParserError{Token: parser.lookahead(), Message: message, Position: parser.current}
// }

// func (parser Parser) Parse() ([]grammar.Statement, grammar.LoxError) {
// 	statements := make([]grammar.Statement, 0)
// 	current := parser.prev().TokenType
// 	for current != grammar.EOF {
// 		stmt, err := parser.statement()
// 		if err != nil {
// 			return nil, err
// 		}
// 		statements = append(statements, stmt)
// 		current = parser.next().TokenType
// 	}
// 	return statements, nil
// }

// func (parser *Parser) statement() (grammar.Statement, grammar.LoxError) {
// 	switch {
// 	case parser.matchToken(grammar.PRINT):
// 		return parser.printStatment()
// 	default:
// 		return parser.expressionStatement()
// 	}
// }

// func (parser *Parser) printStatment() (grammar.Statement, grammar.LoxError) {
// 	value, err := parser.expression()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return grammar.PrintStatment{Value: value}, nil
// }

// func (parser *Parser) expressionStatement() (grammar.Statement, grammar.LoxError) {
// 	expr, err := parser.expression()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return grammar.ExpressionStatement{Expression: expr}, nil
// }

// func (parser *Parser) expression() (grammar.Expression, grammar.LoxError) {
// 	return parser.equality()
// }

// func (parser *Parser) equality() (grammar.Expression, grammar.LoxError) {
// 	leftExpr, err := parser.comparison()

// 	for parser.matchToken(grammar.BANG, grammar.EQUAL_EQUAL) {
// 		operator := parser.prev()
// 		rightExpr, err := parser.comparison()
// 		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
// 		return leftExpr, err
// 	}
// 	return leftExpr, err
// }

// func (parser *Parser) comparison() (grammar.Expression, grammar.LoxError) {
// 	leftExpr, err := parser.term()

// 	for parser.matchToken(grammar.GREATER, grammar.GREATER_EQUAL, grammar.LESS, grammar.LESS_EQUAL) {
// 		operator := parser.prev()
// 		rightExpr, err := parser.term()
// 		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
// 		return leftExpr, err
// 	}

// 	return leftExpr, err
// }

// func (parser *Parser) term() (grammar.Expression, grammar.LoxError) {
// 	leftExpr, err := parser.factor()

// 	for parser.matchToken(grammar.MINUS, grammar.PLUS) {
// 		operator := parser.prev()
// 		rightExpr, err := parser.factor()
// 		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
// 		return leftExpr, err
// 	}

// 	return leftExpr, err
// }

// func (parser *Parser) factor() (grammar.Expression, grammar.LoxError) {
// 	leftExpr, err := parser.unary()

// 	for parser.matchToken(grammar.SLASH, grammar.STAR) {
// 		operator := parser.prev()
// 		rightExpr, err := parser.unary()
// 		leftExpr = grammar.BinaryExpression{Left: leftExpr, Right: rightExpr, Operator: operator}
// 		return leftExpr, err
// 	}

// 	return leftExpr, err
// }

// func (parser *Parser) unary() (grammar.Expression, grammar.LoxError) {
// 	if parser.matchToken(grammar.BANG, grammar.MINUS) {
// 		operator := parser.prev()
// 		rightExpr, err := parser.unary()
// 		return grammar.UnaryExpression{Right: rightExpr, Operator: operator}, err
// 	}
// 	return parser.primary()
// }

// func (parser *Parser) primary() (grammar.Expression, grammar.LoxError) {

// 	switch {
// 	case parser.matchToken(grammar.FALSE):
// 		return grammar.LiteralExpression{Literal: false}, nil
// 	case parser.matchToken(grammar.TRUE):
// 		return grammar.LiteralExpression{Literal: true}, nil
// 	case parser.matchToken(grammar.NULL):
// 		return grammar.LiteralExpression{Literal: nil}, nil
// 	case parser.matchToken(grammar.NUMBER, grammar.STRING):
// 		return grammar.LiteralExpression{Literal: parser.prev().Literal}, nil
// 	case parser.matchToken(grammar.LEFT_PAREN):
// 		expr, _ := parser.expression()
// 		return grammar.GroupingExpression{Expression: expr}, parser.consume(grammar.RIGHT_PAREN, "Expect ')' after expression.")
// 	}
// 	parser.sync()
// 	return nil, ParserError{Position: parser.current, Message: "Expect expression."}
// }

// func (parser *Parser) sync() {

// 	prev := parser.prev().TokenType
// 	if prev != grammar.EOF {
// 		parser.current++
// 		prev = parser.prev().TokenType
// 	}
// 	for prev != grammar.EOF {
// 		if prev == grammar.SEMICOLON {
// 			return
// 		}
// 		switch parser.lookahead().TokenType {
// 		case grammar.CLASS:
// 			fallthrough
// 		case grammar.FUNC:
// 			fallthrough
// 		case grammar.VAR:
// 			fallthrough
// 		case grammar.FOR:
// 			fallthrough
// 		case grammar.IF:
// 			fallthrough
// 		case grammar.WHILE:
// 			fallthrough
// 		case grammar.PRINT:
// 			fallthrough
// 		case grammar.RETURN:
// 			return
// 		}

// 		parser.current++
// 	}
// }