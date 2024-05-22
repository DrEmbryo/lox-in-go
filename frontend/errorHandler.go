package Lox

import "fmt"

type LoxError struct {
	line uint32
	position uint32
	message string 
}

func (e LoxError) Print() {
	fmt.Printf("[%d:%d]: Error: %s \n", e.line, e.position, e.message)
}

func (e LoxError) Error() string {
	return fmt.Sprintf("[%d:%d]: Error: %s \n", e.line, e.position, e.message)
}