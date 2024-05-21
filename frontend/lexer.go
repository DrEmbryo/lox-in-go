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
	tokenType int8
	lexeme    string
	literal   string
}

func (token Token) Create(tokenType int8, lexeme string, literal string) Token {
	token.tokenType = tokenType
	token.lexeme = lexeme
	token.literal = literal
	return token
}

type Lexer struct {
	start   uint32
	current uint32
	line    uint32
	source  []rune
}

func (lexer *Lexer) init(source []rune, start uint32, current uint32, line uint32) {
	lexer.start = start
	lexer.current = current
	lexer.line = line
	lexer.source = source
}

func (lexer *Lexer) next() rune {
	if lexer.current < uint32(len(lexer.source)) {
		nextChar := lexer.source[lexer.current]
		lexer.current++
		return nextChar
	}
	return 0
}

// TODO: seems like its causing issues, should be runeCount?
func (lexer Lexer) lookahead() rune {
	if lexer.current < uint32(len(lexer.source)) {
		return lexer.source[lexer.current]
	}
	return 0
}

func (lexer Lexer) Tokenize(source string) ([]Token, []Error) {
	tokens := make([]Token, 10)
	lexer.init([]rune(source), 1, 0, 1)
	lexErrors := make([]Error, 0)

	for {
		var tokenType int8
		char := lexer.next()
		if char == 0 {
			break
		}

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
			// multiline comment 
			if lexer.lookahead() == '/' {
				lexer.current++
				for  {
					char = lexer.next();
					if char == 0 {
						break
					}
					if char == '\n' {
						lexer.line++
					}
					if char == '/' && lexer.lookahead() == '*' {
						lexer.next()
						break
					}
				}
			} else {
				tokenType = STAR
			}
		case '/':
			if lexer.lookahead() == '/' {
				for {
					if lexer.next() == '\n' || lexer.next() == 0 {
						lexer.line++
						break
					}
				}
			} else {
				tokenType = SLASH
			}
		case '=':
			if lexer.lookahead() == '=' {
				tokenType = EQUAL_EQUAL
				lexer.next()
			} else {
				tokenType = EQUAL
			}
		case '!':
			if lexer.lookahead() == '=' {
				tokenType = BANG_EQUAL
				lexer.next()
			} else {
				tokenType = BANG
			}
		case '<':
			if lexer.lookahead() == '=' {
				tokenType = LESS_EQUAL
				lexer.next()
			} else {
				tokenType = EQUAL
			}
		case '>':
			if lexer.lookahead() == '=' {
				tokenType = GREATER_EQUAL
				lexer.next()
			} else {
				tokenType = EQUAL
			}
		case ' ':
		case '\r':
		case '\t':
		case '\n':
			lexer.line++
		default:
			lexErrors = append(lexErrors, Error{line: lexer.line, position: lexer.current, message: fmt.Sprintf("Unknown token: %c", char)})
			tokenType = -1
		}

		if tokenType != -1 {
			tokens = append(tokens, Token{}.Create(tokenType, "", ""))
		}
	}

	tokens = append(tokens, Token{}.Create(EOF, "", ""))
	return tokens, lexErrors
}
