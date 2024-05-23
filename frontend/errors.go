package Lox

import "fmt"

type LoxError struct {
	line uint32
	position uint32
	message	string
	stage string 
}

func (e LoxError) Print() {
	fmt.Printf("[%d:%d]: %s error: %s \n", e.line, e.position, e.stage, e.message)
}

func (e LoxError) Error() string {
	return fmt.Sprintf("[%d:%d]: %s error: %s \n", e.line, e.position, e.stage, e.message)
}