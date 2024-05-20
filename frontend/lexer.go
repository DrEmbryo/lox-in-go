package Lox

import "fmt"

const (
	// single char tokens
	LEFT_PAREN = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// multi char tokens
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// literals
	IDENTIFIER
	STRING
	NUMBER

	// keywords
	AND
	CLASS
	ELSE
	FALSE
	FUNC
	FOR
	IF
	NULL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE
	EOF
)

type Token struct {
	tokenType uint8
	lexeme    string
	literal   string
}

func (token Token) Create(tokenType uint8, lexeme string, literal string) Token {
	token.tokenType = tokenType
	token.lexeme = lexeme
	token.literal = literal
	return token
}

type Lexer struct {
	start   uint32
	current uint32
	line    uint32
}

func (lexer *Lexer) Next(source []rune) rune {
	nextChar := source[lexer.current]
	lexer.current++
	return nextChar
}

func (lexer Lexer) Tokenize(source string) ([]Token, []Error) {
	lexErrors := make([]Error, 0)
	tokens := make([]Token, 10)
	sourceStr := []rune(source)
	lexer.start = 1
	lexer.current = 0
	lexer.line = 0

	for {
		char := lexer.Next(sourceStr)
		if lexer.current >= uint32(len(source)) {
			break
		}
		fmt.Println(string(char))
		var tokenType uint8
		switch char {
		case '(':
			tokenType = LEFT_PAREN
		case ')':
			tokenType = RIGHT_PAREN
		case '{':
			tokenType = LEFT_BRACE
		case '}':
			tokenType = RIGHT_BRACE
		case ',':
			tokenType = COMMA
		case '.':
			tokenType = DOT
		case '-':
			tokenType = MINUS
		case '+':
			tokenType = PLUS
		case ';':
			tokenType = SEMICOLON
		case '*':
			tokenType = STAR
		case '/':
			if lexer.Next(sourceStr) == '/' {
				// double slash comment support

			} else {
				tokenType = SLASH
			}
		case '=':
			
			if lexer.Next(sourceStr) == '=' {
				tokenType = EQUAL_EQUAL
			} else {
				tokenType = EQUAL
			}
		case '!':
			if lexer.Next(sourceStr) == '=' {
				tokenType = BANG_EQUAL
			} else {
				tokenType = BANG
			}
		case '<':
			if lexer.Next(sourceStr) == '=' {
				tokenType = LESS_EQUAL
			} else {
				tokenType = EQUAL
			}
		case '>':
			if lexer.Next(sourceStr) == '=' {
				tokenType = GREATER_EQUAL
			} else {
				tokenType = EQUAL
			}
		default:
			lexErrors = append(lexErrors, Error{line: lexer.line, position: lexer.current, message: "Unidentified token"})
			tokenType = 0
		}

		tokens = append(tokens, Token{}.Create(tokenType, "", ""))
	}

	tokens = append(tokens, Token{
		tokenType: EOF,
		lexeme:    "",
		literal:   "",
	})
	return tokens, lexErrors
}
