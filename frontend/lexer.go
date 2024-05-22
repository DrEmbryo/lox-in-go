package Lox

import (
	"bytes"
	"fmt"
)

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

var KEYWORDS = map[string]int{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"func":   FUNC,
	"if":     IF,
	"null":   NULL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

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
		tokenType := -1
		tokenValue := ""
		char := lexer.next()
		if char == 0 {
			break
		}

		switch {
		case char == '(':
			tokenType = LEFT_PAREN
		case char == ')':
			tokenType = RIGHT_PAREN
		case char == '{':
			tokenType = LEFT_BRACE
		case char == '}':
			tokenType = RIGHT_BRACE
		case char == ',':
			tokenType = COMMA
		case char == '.':
			tokenType = DOT
		case char == '-':
			tokenType = MINUS
		case char == '+':
			tokenType = PLUS
		case char == ';':
			tokenType = SEMICOLON
		case char == '*':
			tokenType = STAR
		case char == '/':
			if lexer.lookahead() == '/' {
				lexer.current++
				for {
					char = lexer.next()
					if char == 0 {
						break
					}
					if char == '\n' {
						lexer.line++
						break
					}
				}
			}  else if lexer.lookahead() == '/' {
				lexer.current++
				for {
					char = lexer.next()
					if char == 0 {
						break
					}
					if char == '\n' {
						lexer.line++
					}
					if char == '*' && lexer.lookahead() == '/' {
						lexer.next()
						break
					}
				}
			} else {
				tokenType = SLASH
			}
		case char == '=':
			if lexer.lookahead() == '=' {
				tokenType = EQUAL_EQUAL
				lexer.next()
			} else {
				tokenType = EQUAL
			}
		case char == '!':
			if lexer.lookahead() == '=' {
				tokenType = BANG_EQUAL
				lexer.next()
			} else {
				tokenType = BANG
			}
		case char == '<':
			if lexer.lookahead() == '=' {
				tokenType = LESS_EQUAL
				lexer.next()
			} else {
				tokenType = EQUAL
			}
		case char == '>':
			if lexer.lookahead() == '=' {
				tokenType = GREATER_EQUAL
				lexer.next()
			} else {
				tokenType = EQUAL
			}
		case char == '"':
			startLine := lexer.current
			buff := bytes.NewBufferString("")
			for {
				char = lexer.next()
				if char == 0 {
					tokenType = -1
					lexErrors = append(lexErrors, Error{line: lexer.line, position: lexer.current, message: fmt.Sprintf("Unterminated string at : %c", startLine)})
					break
				}
				if char != '"' {
					if char == '\n' {
						lexer.line++
					}
					buff.WriteRune(char)
				} else {
					tokenType = STRING
					tokenValue = buff.String()
					break
				}
			}
		case parseDigit(char):
			tokenType = NUMBER
			buff := bytes.NewBufferString("")
			for {
				if !parseDigit(char) {
					break
				}
				buff.WriteRune(char)
				char = lexer.next()
			}
			if char == '.' && parseDigit(lexer.lookahead()) {
				buff.WriteRune(char)
				char = lexer.next()
				for {
					if !parseDigit(char) {
						break
					}
					buff.WriteRune(char)
					char = lexer.next()
				}
			}
			tokenValue = buff.String()
		case parseChar(char):
			buff := bytes.NewBufferString("")
			for {
				if !(parseChar(char) || parseDigit(char)) {
					break
				}
				buff.WriteRune(char)
				char = lexer.next()
			}
			tokenValue = buff.String()
			keyword, ok := KEYWORDS[tokenValue]
			if ok {
				tokenType = keyword
			} else {
				tokenType = IDENTIFIER
			}
		case char == ' ':
		case char == '\r':
		case char == '\t':
		case char == '\n':
			lexer.line++
		default:
			lexErrors = append(lexErrors, Error{line: lexer.line, position: lexer.current, message: fmt.Sprintf("Unknown token: %c", char)})
			tokenType = -1
		}
		if tokenType != -1 {
			tokens = append(tokens, Token{}.Create(int8(tokenType), tokenValue, ""))
		}
	}

	tokens = append(tokens, Token{}.Create(EOF, "", ""))
	return tokens, lexErrors
}

func parseDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func parseChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}
