package lexer

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Lexer struct {
	Source  []rune
	current int
	line 	int
}

func (lexer *Lexer) consume() rune {
	char := lexer.Source[lexer.current]
	lexer.current++
	return char
}

func (lexer *Lexer) lookahead() rune {
	char := lexer.Source[lexer.current]
	return char
}

func (lexer Lexer) Tokenize() ([]grammar.Token, []LexerError) {
	tokens := make([]grammar.Token, 0)
	lexErrors := make([]LexerError, 0)

	if len(lexer.Source) == 0 {
		lexErrors = append(lexErrors, LexerError{Message: "source contains 0 characters"})
		return tokens, lexErrors
	}

	for lexer.current <= len(lexer.Source) - 1 {
		char := lexer.consume()
		switch token := lexer.parseSingleCharToken(&char).(type) {
		case grammar.Token:
			tokens = append(tokens, token)
		case LexerError:
			lexErrors = append(lexErrors, token)
		}

	}

	tokens = append(tokens, grammar.Token{TokenType: grammar.EOF, Lexeme: "EOF"})

	return tokens, lexErrors
}

func (lexer *Lexer) parseSingleCharToken(char *rune) any {
	fmt.Println(string(*char))
	switch *char {
	case '(':
		return grammar.Token{TokenType: grammar.LEFT_PAREN, Lexeme: string(*char)}
	case ')':
		return grammar.Token{TokenType: grammar.RIGHT_PAREN, Lexeme: string(*char)}
	case '{':
		return grammar.Token{TokenType: grammar.LEFT_BRACE, Lexeme: string(*char)}
	case '}':
		return grammar.Token{TokenType: grammar.RIGHT_BRACE, Lexeme: string(*char)}
	case ',':
		return grammar.Token{TokenType: grammar.COMMA, Lexeme: string(*char)}
	case '.':
		return grammar.Token{TokenType: grammar.DOT, Lexeme: string(*char)}
	case '-':
		return grammar.Token{TokenType: grammar.MINUS, Lexeme: string(*char)}
	case '+':
		return grammar.Token{TokenType: grammar.PLUS, Lexeme: string(*char)}
	case ';':
		return grammar.Token{TokenType: grammar.SEMICOLON, Lexeme: string(*char)}
	case '*':
		return grammar.Token{TokenType: grammar.STAR, Lexeme: string(*char)}
	case '/':
		switch {
		case lexer.lookahead() == '/':
			lexer.parseSingleLineComments(char)
		case lexer.lookahead() == '*':
			lexer.parseMultilineLineComments(char)
		default:
			return grammar.Token{TokenType: grammar.SLASH, Lexeme: string(*char)}
		}
	}
	
	return lexer.parseMultiCahrToken(char)
}

func (lexer *Lexer) parseMultiCahrToken(char *rune) any {
	switch *char {
	case '=':
		if lexer.lookahead() == '=' {
			return grammar.Token{TokenType: grammar.EQUAL_EQUAL, Lexeme: fmt.Sprintf("%s%s", string(*char), string(lexer.consume()))}
		} else {
			return grammar.Token{TokenType: grammar.EQUAL, Lexeme: string(*char)}
		}
	case '!':
		if lexer.lookahead() == '=' {
			return grammar.Token{TokenType: grammar.BANG_EQUAL, Lexeme: fmt.Sprintf("%s%s", string(*char), string(lexer.consume()))}
		} else {
			return grammar.Token{TokenType: grammar.BANG, Lexeme: string(*char)}
		}
	case '<':
		if lexer.lookahead() == '=' {
			return grammar.Token{TokenType: grammar.LESS_EQUAL, Lexeme: fmt.Sprintf("%s%s", string(*char), string(lexer.consume()))}
		} else {
			return grammar.Token{TokenType: grammar.LESS, Lexeme: string(*char)}
		}
	case '>':
		if lexer.lookahead() == '=' {
			return grammar.Token{TokenType: grammar.GREATER_EQUAL, Lexeme: fmt.Sprintf("%s%s", string(*char), string(lexer.consume()))}
		} else {
			return grammar.Token{TokenType: grammar.GREATER, Lexeme: string(*char)}
		}
	case '"':
		token, err := lexer.parseString(char)
		if err != nil {
			return err
		}
		return token
	case '\n':
		lexer.line++
	default:
		switch {
		case parseDigit(char):
			return lexer.parseNumerics(char)
		case parseChar(char):
			return lexer.parseIdentifiers(char)
		case parseSkippable(char):
			return nil
		default:
			return LexerError{Line: lexer.line, Position: lexer.current, Message: fmt.Sprintf("Unknown token: %c", *char)}
			}	
		} 
		return nil
}

func (lexer *Lexer) parseString(char *rune) (grammar.Token, grammar.LoxError) {
	buff := bytes.NewBufferString("")

	for lexer.current <= len(lexer.Source) - 1 {
		*char = lexer.consume()
		switch *char {
		case '"':
			return grammar.Token{TokenType: grammar.STRING, Lexeme: buff.String()}, nil
		case '\n':
			lexer.line++
		default :
			buff.WriteRune(*char)
		}
	}
	return grammar.Token{}, LexerError{Line: lexer.line, Position: lexer.current,  Message: "Unterminated string"}
}

func (lexer *Lexer) parseNumerics(char *rune) grammar.Token {
	buff := bytes.NewBufferString("")
	buff.WriteRune(*char)
	for lexer.current <= len(lexer.Source) - 1 {
		*char = lexer.consume()
		if (parseDigit(char)) {
			buff.WriteRune(*char)
		} else {
			break
		}
	}

	if *char == '.' {
		buff.WriteRune(*char)
		for lexer.current <= len(lexer.Source) - 1 {
			*char = lexer.consume()
			if (parseDigit(char)) {
				buff.WriteRune(*char)
			} else {
				break
			}
		}
	} else {
		lexer.current--
	}

	value, _ := strconv.ParseFloat(buff.String(), 64)
 	return grammar.Token{TokenType: grammar.NUMBER, Lexeme: value}
}

func (lexer *Lexer) parseIdentifiers(char *rune) grammar.Token {
	buff := bytes.NewBufferString("")
	buff.WriteRune(*char)
	for lexer.current <= len(lexer.Source) - 1 {
		*char = lexer.consume()
		if (parseDigit(char) || parseChar(char)) {
			buff.WriteRune(*char)
		} else {
			break
		}
	}
	lexer.current--
	tokenValue := buff.String()
	keyword, ok := grammar.KEYWORDS[tokenValue]
	if ok {
		return grammar.Token{TokenType: keyword, Lexeme: tokenValue}
	} 
	return grammar.Token{TokenType: grammar.IDENTIFIER, Lexeme: tokenValue}
}

func (lexer *Lexer) parseSingleLineComments(char *rune) {
	for lexer.current <= len(lexer.Source) - 1 {
		*char = lexer.consume()
		if *char == '\n' {
			lexer.line++
			return
		}
	}
}

func (lexer *Lexer) parseMultilineLineComments(char *rune) {
	for lexer.current <= len(lexer.Source) - 1 {
		*char = lexer.consume()
		if *char == '\n' {
			lexer.line++
		}
		if *char == '/' && lexer.lookahead() == '*' {
			*char = lexer.consume()
			lexer.parseMultilineLineComments(char)
		}
		if *char == '*' && lexer.lookahead() == '/' {
			*char = lexer.consume()
			return 
		}
	}
}

func parseDigit(char *rune) bool {
	expr, _ := regexp.Compile("[0-9]")
	return expr.MatchString(string(*char))
}

func parseChar(char *rune) bool {
	expr, _ := regexp.Compile("[A-Za-z_]")
	return expr.MatchString(string(*char))
}

func parseSkippable(char *rune) bool {
	expr, _ := regexp.Compile("[\t\r ]")
	return expr.MatchString(string(*char))
}