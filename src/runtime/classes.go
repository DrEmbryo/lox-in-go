package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type LoxClass struct {
	Name grammar.Token
}

func (class *LoxClass) Call(interpreter Interpreter, args []any) (any, grammar.LoxError) {
	instance := LoxClassInstance{Class: *class}
	return instance, nil
}

func (function *LoxClass) GetAirity() int {
	return 0
}

func (class *LoxClass) ToString() string {
	return fmt.Sprintf("<class %v>", class.Name.Lexeme)
}

type LoxClassInstance struct {
	Class LoxClass
}

func (instance *LoxClassInstance) ToString() string {
	return fmt.Sprintf("<instance of %v class>", instance.Class.Name.Lexeme)
}
