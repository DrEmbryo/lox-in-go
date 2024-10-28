package lexer

import "fmt"

type LexerError struct {
	Line     int
	Position int
	Message  string
}

func (e LexerError) Print() {
	fmt.Printf("[%d:%d]: lexer error: %s \n", e.Line, e.Position, e.Message)
}

func (e LexerError) Error() string {
	return fmt.Sprintf("[%d:%d]: lexer error: %s \n", e.Line, e.Position, e.Message)
}
