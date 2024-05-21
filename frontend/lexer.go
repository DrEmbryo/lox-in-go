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
	lexer.init([]rune(source), 1, 0, 0)
	lexErrors := make([]Error, 0)

	for {
		char := lexer.next()
		if lexer.current >= uint32(len(source)) {
			break
		}

		var tokenType int8
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
				for {
					// TODO: something is wrong with \n count inside this loop
					next := lexer.next()
					fmt.Println(string(next))
					lookup := lexer.lookahead()
					fmt.Println(string(lookup))
					if lookup == '\n' {
						lexer.line++
					}
					if next == 0 || lookup == 0 {
						break
					}
					if next == '/' && lookup == '*' {
						lexer.next() 
						break
					}
				}
			} else {
				tokenType = STAR
			}
		case '/':
			if lexer.lookahead() == '/' {
				// comment line skip
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
			if lexer.next() == '=' {
				tokenType = EQUAL_EQUAL
			} else {
				tokenType = EQUAL
			}
		case '!':
			if lexer.next() == '=' {
				tokenType = BANG_EQUAL
			} else {
				tokenType = BANG
			}
		case '<':
			if lexer.next() == '=' {
				tokenType = LESS_EQUAL
			} else {
				tokenType = EQUAL
			}
		case '>':
			if lexer.next() == '=' {
				tokenType = GREATER_EQUAL
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

		tokens = append(tokens, Token{}.Create(tokenType, "", ""))
	}

	tokens = append(tokens, Token{}.Create(EOF, "", ""))
	return tokens, lexErrors
}
