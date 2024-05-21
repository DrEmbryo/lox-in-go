package Lox

import "fmt"

type Error struct {
	line uint32
	position uint32
	message string 
}

func (e Error) ThrowLexerError () {
	fmt.Printf("[%d:%d]: Error: %s \n", e.line, e.position, e.message)
}