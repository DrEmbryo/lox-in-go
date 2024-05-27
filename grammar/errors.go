package Lox

import "fmt"

type LoxError interface {
	Error () string
	Print ()
}

type LexerError struct {
	line uint32
	position uint32
	message	string
	stage string 
}

func (e LexerError) Print() {
	fmt.Printf("[%d:%d]: %s error: %s \n", e.line, e.position, e.stage, e.message)
}

func (e LexerError) Error() string {
	return fmt.Sprintf("[%d:%d]: %s error: %s \n", e.line, e.position, e.stage, e.message)
}

type ParserError struct {
	position uint32
	token Token
	message	string
}

func (e ParserError) Print() {
	if (e.token.tokenType == EOF) {
		fmt.Printf("[%d]: EOF at %s, %s", e.position, e.token.lexeme, e.message)
	} else {
		fmt.Printf("[%d]: Token %d: failed at %s with value %s: %s", e.position, e.token.tokenType, e.token.lexeme, e.token.literal, e.message)
	}
}

func (e ParserError) Error() string {
	return fmt.Sprintf("[%d]: Token %d: failed at %s with value %s: %s", e.position, e.token.tokenType, e.token.lexeme, e.token.literal, e.message)
}