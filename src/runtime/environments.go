package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type Environment struct {
	Values map[string]any
	Parent *Environment
}

func (env *Environment) defineEnvValue(name grammar.Token, value any) {
	field := fmt.Sprintf("%s", name.Lexeme)
	env.Values[field] = value
}

func (env *Environment) getEnvValue(name grammar.Token) (any, grammar.LoxError) {
	lookup := fmt.Sprintf("%s", name.Lexeme)
	val, ok := env.Values[lookup]
	if ok {
		return val, nil
	}
	if env.Parent != nil {
		return env.Parent.getEnvValue(name)
	}
	return nil, RuntimeError{Token: name, Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)}
}

func (env *Environment) assignEnvValue(name grammar.Token, value any) grammar.LoxError {
	lookup := fmt.Sprintf("%s", name.Lexeme)
	_, ok := env.Values[lookup]
	if ok {
		env.Values[lookup] = value
		return nil
	}
	if env.Parent != nil {
		return env.Parent.assignEnvValue(name, value)
	}
	return RuntimeError{Token: name, Message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme)}
}
