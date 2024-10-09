package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type LoxFunction struct {
	Declaration grammar.FunctionDeclarationStatement
	Closure     Environment
}

func (function *LoxFunction) Call(interpreter Interpreter, arguments []any) (any, grammar.LoxError) {
	env := function.Closure

	for i := 0; i < len(function.Declaration.Params); i++ {
		env.defineEnvValue(function.Declaration.Params[i], arguments[i])
	}

	return interpreter.executeBlock(function.Declaration.Body.Statements, env)
}

func (function *LoxFunction) GetAirity() int {
	return len(function.Declaration.Params)
}

func (function *LoxFunction) Bind(instance *LoxClassInstance) LoxFunction {
	env := function.Closure

	env.defineEnvValue(grammar.Token{TokenType: grammar.THIS, Lexeme: "this"}, instance)

	return LoxFunction{Declaration: function.Declaration, Closure: env}
}

func (function *LoxFunction) ToString() string {
	return fmt.Sprintf("<fn %v>", function.Declaration.Name.Lexeme)
}
