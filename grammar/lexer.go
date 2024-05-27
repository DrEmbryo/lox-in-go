package Lox

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

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

func (lexer *Lexer) lookahead() rune {
	if lexer.current < uint32(len(lexer.source)) {
		return lexer.source[lexer.current]
	}
	return 0
}

func (lexer Lexer) Tokenize(source string) ([]Token, []LexerError) {
	tokens := make([]Token, 0)
	lexer.init([]rune(source), 1, 0, 1)
	lexErrors := make([]LexerError, 0)

	for {
		var err error
		tokenType := -1
		var tokenValue any
		char := lexer.next()
		if char == 0 {
			break
		}

		switch {
		//handle single character tokens
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
			// handle comments
			switch {
			case lexer.lookahead() == '/':
				lexer.handleSingleLineComments(&char)
			case lexer.lookahead() == '*':
				lexer.handleMultilineLineComments(&char)
			default:
				tokenType = SLASH
			}
		//handle multi character tokens
		case char == '=':
			if lexer.lookahead() == '=' {
				tokenType = EQUAL_EQUAL
				lexer.current++
			} else {
				tokenType = EQUAL
			}
		case char == '!':
			if lexer.lookahead() == '=' {
				tokenType = BANG_EQUAL
				lexer.current++
			} else {
				tokenType = BANG
			}
		case char == '<':
			if lexer.lookahead() == '=' {
				tokenType = LESS_EQUAL
				lexer.current++
			} else {
				tokenType = EQUAL
			}
		case char == '>':
			if lexer.lookahead() == '=' {
				tokenType = GREATER_EQUAL
				lexer.current++
			} else {
				tokenType = EQUAL
			}
		// handle strings
		case char == '"':
			var lerr LexerError
			tokenType, tokenValue, err = lexer.handleStrings(&char)
			if err != nil {
				lexErrors = append(lexErrors, lerr)
			}
		// handle numeric values
		case parseDigit(char):
			tokenType, tokenValue = lexer.handleNumerics(&char)
		// handle keywords
		case parseChar(char):
			tokenType, tokenValue = lexer.handleIdentifiers(&char)
		//  handle skippable characters 
		case parseSkippable(char):
			tokenType = -1
		case char == '\n':
			lexer.line++
			tokenType = -1
		default:
			lexErrors = append(lexErrors, LexerError{line: lexer.line, position: lexer.current, stage: "lexer", message: fmt.Sprintf("Unknown token: %c", char)})
			tokenType = -1
		}
		if tokenType != -1 {
			tokens = append(tokens, Token{tokenType: int8(tokenType), lexeme: tokenValue})
		}
	}

	tokens = append(tokens, Token{tokenType: EOF})
	return tokens, lexErrors
}

func (lexer *Lexer) handleSingleLineComments(char *rune) {
	lexer.current++
	for {
		*char = lexer.next()
		if *char == 0 {
			break
		}
		if *char == '\n' {
			lexer.line++
			break
		}
	}
}

func (lexer *Lexer) handleMultilineLineComments(char *rune) {
	lexer.current++
	for {
		*char = lexer.next()
		if *char == 0 {
			break
		}
		if *char == '\n' {
			lexer.line++
		}
		if *char == '*' && lexer.lookahead() == '/' {
			*char = lexer.next()
			break
		}
	}
}

func (lexer *Lexer) handleStrings(char *rune) (int, string, LoxError) {
	var tokenType int
	var tokenValue string
	buff := bytes.NewBufferString("")
	for {
		*char = lexer.next()
		if *char == 0 {
			tokenType = -1
			return tokenType, tokenValue, LexerError{line: lexer.line, position: lexer.current, stage: "lexer", message: "Unterminated string"}
		}
		if *char != '"' {
			if *char == '\n' {
				lexer.line++
			}
			buff.WriteRune(*char)
		} else {
			tokenType = STRING
			tokenValue = buff.String()
			break
		}
	}
	return tokenType, tokenValue, nil
}

func (lexer *Lexer) handleNumerics(char *rune) (int, float64) {
	buff := bytes.NewBufferString("")
	for {
		if !parseDigit(*char) {
			break
		}
		buff.WriteRune(*char)
		*char = lexer.next()
	}
	if *char == '.' && parseDigit(lexer.lookahead()) {
		buff.WriteRune(*char)
		*char = lexer.next()
		for {
			if !parseDigit(*char) {
				break
			}
			buff.WriteRune(*char)
			*char = lexer.next()
		}
	}
	value, _ := strconv.ParseFloat(buff.String(), 32)
	return NUMBER, value
}

func (lexer *Lexer) handleIdentifiers(char *rune) (int, string) {
	buff := bytes.NewBufferString("")
	for {
		if !(parseChar(*char) || parseDigit(*char)) {
			break
		}
		buff.WriteRune(*char)
		*char = lexer.next()
	}
	tokenValue := buff.String()
	keyword, ok := KEYWORDS[tokenValue]
	if ok {
		return keyword, tokenValue
	} 
	return IDENTIFIER, tokenValue
}

func parseDigit(char rune) bool {
	expr, _ := regexp.Compile("[0-9]")
	return expr.MatchString(string(char))
}

func parseChar(char rune) bool {
	expr, _ := regexp.Compile("[A-Za-z_]")
	return expr.MatchString(string(char))
}

func parseSkippable(char rune) bool {
	expr, _ := regexp.Compile("[\t\r ]*")
	return expr.MatchString(string(char))
}