package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type LoxFunction struct {
	Declaration grammar.FunctionDeclarationStatement
}

func (function *LoxFunction) Call(interpreter Interpreter, arguments []any) any {
	env := Environment{Values: interpreter.globalEnv.Values}

	for i := 0; i < len(function.Declaration.Params); i++ {
		env.defineEnvValue(function.Declaration.Params[i], arguments[i])
	}

	interpreter.executeBlock(function.Declaration.Body.Statements, env)
	return nil
}

func (function *LoxFunction) GetAirity() int {
	return len(function.Declaration.Params)
}

func (function *LoxFunction) ToString() string {
	return fmt.Sprintf("<fn  %v>", function.Declaration.Name.Lexeme)
}
