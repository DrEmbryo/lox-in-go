package lexer

import (
	"fmt"
	"regexp"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Lexer struct {
	Source  []rune
	current int
	line 	int
}

func (lexer *Lexer) consume() rune {
	char := lexer.Source[lexer.current]
	fmt.Println(char)
	lexer.current++
	return char
}

func (lexer Lexer) Tokenize() ([]grammar.Token, []LexerError) {
	tokens := make([]grammar.Token, 0)
	lexErrors := make([]LexerError, 0)
	fmt.Println(lexer.Source)

	char := lexer.consume()
	for char != 0 && lexer.current < len(lexer.Source) {
		switch token := lexer.parseSingleCharToken(char).(type) {
		case grammar.Token:
			tokens = append(tokens, token)
		}
		if char == 0 {
			tokens = append(tokens, grammar.Token{TokenType: grammar.EOF, Lexeme: "EOF"})
			lexer.consume()
		}

		char = lexer.consume()
	}
	fmt.Println(tokens)
	return tokens, lexErrors
}

func (lexer *Lexer) parseSingleCharToken(char rune) any {
	switch {
	//handle single character tokens
	case char == '(':
		return grammar.Token{TokenType: grammar.LEFT_PAREN, Lexeme: char}
	case char == ')':
		return grammar.Token{TokenType: grammar.RIGHT_PAREN, Lexeme: char}
	case char == '{':
		return grammar.Token{TokenType: grammar.LEFT_BRACE, Lexeme: char}
	case char == '}':
		return grammar.Token{TokenType: grammar.RIGHT_BRACE, Lexeme: char}
	case char == ',':
		return grammar.Token{TokenType: grammar.COMMA, Lexeme: char}
	case char == '.':
		return grammar.Token{TokenType: grammar.DOT, Lexeme: char}
	case char == '-':
		return grammar.Token{TokenType: grammar.MINUS, Lexeme: char}
	case char == '+':
		return grammar.Token{TokenType: grammar.PLUS, Lexeme: char}
	case char == ';':
		return grammar.Token{TokenType: grammar.SEMICOLON, Lexeme: char}
	case char == '*':
		return grammar.Token{TokenType: grammar.STAR, Lexeme: char}
	case char == '/':
		// handle comments
		switch {
		// case lexer.lookahead() == '/':
		// 	lexer.handleSingleLineComments(&char)
		// case lexer.lookahead() == '*':
		// 	lexer.handleMultilineLineComments(&char)
		// default:
		// 	tokenType = grammar.SLASH
		// 	tokenValue = char
		}
	}
	return nil
}


// func (lexer *Lexer) next() rune {
// 	if lexer.current < uint32(len(lexer.source)) {
// 		nextChar := lexer.source[lexer.current]
// 		lexer.current++
// 		return nextChar
// 	}
// 	return 0
// }

// func (lexer *Lexer) lookahead() rune {
// 	if lexer.current < uint32(len(lexer.source)) {
// 		return lexer.source[lexer.current]
// 	}
// 	return 0
// }

// func (lexer Lexer) Tokenize() ([]grammar.Token, []LexerError) {
// 	tokens := make([]grammar.Token, 0)
// 	lexErrors := make([]LexerError, 0)

// 	for {
// 		var err error
// 		tokenType := -1
// 		var tokenValue any
// 		char := lexer.next()
// 		if char == 0 {
// 			break
// 		}

// 		switch {
// 		//handle single character tokens
// 		case char == '(':
// 			tokenType = grammar.LEFT_PAREN
// 			tokenValue = char
// 		case char == ')':
// 			tokenType = grammar.RIGHT_PAREN
// 			tokenValue = char
// 		case char == '{':
// 			tokenType = grammar.LEFT_BRACE
// 			tokenValue = char
// 		case char == '}':
// 			tokenType = grammar.RIGHT_BRACE
// 			tokenValue = char
// 		case char == ',':
// 			tokenType = grammar.COMMA
// 			tokenValue = char
// 		case char == '.':
// 			tokenType = grammar.DOT
// 			tokenValue = char
// 		case char == '-':
// 			tokenType = grammar.MINUS
// 			tokenValue = char
// 		case char == '+':
// 			tokenType = grammar.PLUS
// 			tokenValue = char
// 		case char == ';':
// 			tokenType = grammar.SEMICOLON
// 			tokenValue = char
// 		case char == '*':
// 			tokenType = grammar.STAR
// 			tokenValue = char
// 		case char == '/':
// 			// handle comments
// 			switch {
// 			case lexer.lookahead() == '/':
// 				lexer.handleSingleLineComments(&char)
// 			case lexer.lookahead() == '*':
// 				lexer.handleMultilineLineComments(&char)
// 			default:
// 				tokenType = grammar.SLASH
// 				tokenValue = char
// 			}
// 		//handle multi character tokens
// 		case char == '=':
// 			if lexer.lookahead() == '=' {
// 				tokenType = grammar.EQUAL_EQUAL
// 				tokenValue = "=="
// 				lexer.current++
// 			} else {
// 				tokenType = grammar.EQUAL
// 				tokenValue = char
// 			}
// 		case char == '!':
// 			if lexer.lookahead() == '=' {
// 				tokenType = grammar.BANG_EQUAL
// 				tokenValue = "!="
// 				lexer.current++
// 			} else {
// 				tokenType = grammar.BANG
// 				tokenValue = char
// 			}
// 		case char == '<':
// 			if lexer.lookahead() == '=' {
// 				tokenType = grammar.LESS_EQUAL
// 				tokenValue = "<="
// 				lexer.current++
// 			} else {
// 				tokenType = grammar.LESS
// 				tokenValue = char
// 			}
// 		case char == '>':
// 			if lexer.lookahead() == '=' {
// 				tokenType = grammar.GREATER_EQUAL
// 				tokenValue = ">="
// 				lexer.current++
// 			} else {
// 				tokenType = grammar.GREATER
// 				tokenValue = char
// 			}
// 		// handle strings
// 		case char == '"':
// 			var lerr LexerError
// 			tokenType, tokenValue, err = lexer.handleStrings(&char)
// 			if err != nil {
// 				lexErrors = append(lexErrors, lerr)
// 			}
// 		// handle numeric values
// 		case parseDigit(char):
// 			tokenType, tokenValue = lexer.handleNumerics(&char)
// 		// handle keywords
// 		case parseChar(char):
// 			tokenType, tokenValue = lexer.handleIdentifiers(&char)
// 		//  handle skippable characters 
// 		case parseSkippable(char):
// 			tokenType = -1
// 		case char == '\n':
// 			lexer.line++
// 			tokenType = -1
// 		default:
// 			lexErrors = append(lexErrors, LexerError{Line: lexer.line, Position: lexer.current, Stage: "lexer",  Message: fmt.Sprintf("Unknown token: %c", char)})
// 			tokenType = -1
// 		}
// 		if tokenType != -1 {
// 			var literal any
// 			var lexeme string
// 			switch tokenType {
// 			case grammar.NUMBER:
// 				literal = tokenValue.(float64)
// 			case grammar.STRING:
// 				literal = tokenValue.(string)
// 			case grammar.IDENTIFIER:
// 				literal = tokenValue.(string)
// 			default:
// 				literal = nil
// 			} 
			
// 			switch val := tokenValue.(type) {
// 			case string:
// 				lexeme = val
// 			case rune: 
// 				lexeme = string(val)
// 			default: 
// 				lexeme = fmt.Sprintf("%v", val)
// 			}
			
// 			tokens = append(tokens, grammar.Token{TokenType: tokenType, Lexeme: lexeme, Literal: literal})
// 		}
// 	}

// 	tokens = append(tokens, grammar.Token{TokenType: grammar.EOF})
// 	return tokens, lexErrors
// }

// func (lexer *Lexer) handleSingleLineComments(char *rune) {
// 	lexer.current++
// 	for *char != '\n' {
// 		*char = lexer.next()
// 	}
// 	lexer.line++
// }

// func (lexer *Lexer) handleMultilineLineComments(char *rune) {
// 	lexer.current++
// 	for *char != 0 {
// 		*char = lexer.next()
// 		if *char == '\n' {
// 			lexer.line++
// 		}
// 		if *char == '*' && lexer.lookahead() == '/' {
// 			*char = lexer.next()
// 			break
// 		}
// 	}
// }

// func (lexer *Lexer) handleStrings(char *rune) (int, string, grammar.LoxError) {
// 	var tokenType int
// 	var tokenValue string
// 	buff := bytes.NewBufferString("")
// 	for {
// 		*char = lexer.next()
// 		if *char == 0 {
// 			tokenType = -1
// 			return tokenType, tokenValue, LexerError{Line: lexer.line, Position: lexer.current, Stage: "lexer", Message: "Unterminated string"}
// 		}
// 		if *char != '"' {
// 			if *char == '\n' {
// 				lexer.line++
// 			}
// 			buff.WriteRune(*char)
// 		} else {
// 			tokenType = grammar.STRING
// 			tokenValue = buff.String()
// 			break
// 		}
// 	}
// 	return tokenType, tokenValue, nil
// }

// func (lexer *Lexer) handleNumerics(char *rune) (int, float64) {
// 	buff := bytes.NewBufferString("")
// 	for parseDigit(*char) {
// 		buff.WriteRune(*char)
// 		*char = lexer.next()
// 	}
// 	if *char == '.' && parseDigit(lexer.lookahead()) {
// 		buff.WriteRune(*char)
// 		*char = lexer.next()
// 		for parseDigit(*char) {
// 			buff.WriteRune(*char)
// 			*char = lexer.next()
// 		}
// 	}
// 	lexer.current--
// 	value, _ := strconv.ParseFloat(buff.String(), 64)
// 	return grammar.NUMBER, value
// }

// func (lexer *Lexer) handleIdentifiers(char *rune) (int, string) {
// 	buff := bytes.NewBufferString("")
// 	for parseChar(*char) || parseDigit(*char){
// 		buff.WriteRune(*char)
// 		*char = lexer.next()
// 	}
// 	tokenValue := buff.String()
// 	keyword, ok := grammar.KEYWORDS[tokenValue]
// 	if ok {
// 		return keyword, tokenValue
// 	} 
// 	return grammar.IDENTIFIER, tokenValue
// }

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