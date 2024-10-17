package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type LoxFunction struct {
	Declaration grammar.FunctionDeclarationStatement
	Closure     *Environment
	Initializer bool
}

func (function *LoxFunction) Call(interpreter Interpreter, arguments []any) (any, grammar.LoxError) {
	env := Environment{Parent: function.Closure, Values: function.Closure.Values}
	for i := 0; i < len(function.Declaration.Params); i++ {
		env.defineEnvValue(function.Declaration.Params[i], arguments[i])
	}

	val, err := interpreter.executeBlock(function.Declaration.Body.Statements, env)

	if function.Initializer {
		return env.Parent.getEnvValueAt(0, grammar.Token{TokenType: grammar.THIS, Lexeme: "this"})
	}

	return val, err
}

func (function *LoxFunction) GetAirity() int {
	return len(function.Declaration.Params)
}

func (function *LoxFunction) Bind(instance LoxClassInstance) LoxFunction {
	function.Closure.defineEnvValue(grammar.Token{TokenType: grammar.THIS, Lexeme: "this"}, instance)
	return *function
}

func (function *LoxFunction) ToString() string {
	return fmt.Sprintf("<fn %v>", function.Declaration.Name.Lexeme)
}
