package parser

import (
	"fmt"

	"github.com/DrEmbryo/jlox/src/grammar"
)

type ParserError struct {
	Position int
	Token    grammar.Token
	Message  string
}

func (e ParserError) Print() {
	fmt.Printf("[%d]: Token %d: failed at %s with value %s: %s\n", e.Position, e.Token.TokenType, e.Token.Lexeme, e.Token.Literal, e.Message)
}

func (e ParserError) Error() string {
	return fmt.Sprintf("[%d]: Token %d: failed at %s with value %s: %s", e.Position, e.Token.TokenType, e.Token.Lexeme, e.Token.Literal, e.Message)
}
