package runtime

import (
	"fmt"

	"github.com/DrEmbryo/jlox/src/grammar"
)

const CONSTRUCTOR string = "constructor"

type LoxClass struct {
	Name    grammar.Token
	Fields  map[any]any
	Methods map[string]LoxFunction
	Super   any
}

func (class *LoxClass) Call(interpreter Interpreter, args []any) (any, grammar.LoxError) {
	instance := LoxClassInstance{Class: class}
	initMethod := instance.Class.FindMethod(CONSTRUCTOR)
	if init, ok := initMethod.(LoxFunction); ok {
		init.Bind(instance)
		init.Call(interpreter, args)
	}
	return instance, nil
}

func (class *LoxClass) GetAirity() int {
	initMethod := class.FindMethod(CONSTRUCTOR)
	if init, ok := initMethod.(LoxFunction); ok {
		return init.GetAirity()
	}
	return 0
}

func (class *LoxClass) FindMethod(name string) any {
	if method, ok := class.Methods[name]; ok {
		return method
	}

	super, ok := class.Super.(LoxClass)
	if ok {
		return super.FindMethod(name)
	}
	return nil
}

func (class *LoxClass) ToString() string {
	return fmt.Sprintf("<class %v>", class.Name.Lexeme)
}

type LoxClassInstance struct {
	Class *LoxClass
}

func (instance *LoxClassInstance) GetProperty(name grammar.Token) (any, grammar.LoxError) {
	lookup := fmt.Sprintf("%s", name.Lexeme)
	if _, ok := instance.Class.Fields[lookup]; ok {
		return instance.Class.Fields[lookup], nil
	}

	if method := instance.Class.FindMethod(lookup); method != nil {
		m, _ := method.(LoxFunction)
		return m.Bind(*instance), nil
	}

	return nil, RuntimeError{Token: name, Message: fmt.Sprintf("Undefined property '%v'.", name.Lexeme)}
}

func (instance *LoxClassInstance) SetProperty(name grammar.Token, value any) grammar.LoxError {
	lookup := fmt.Sprintf("%s", name.Lexeme)
	instance.Class.Fields[lookup] = value
	return nil
}

func (instance *LoxClassInstance) ToString() string {
	return fmt.Sprintf("<instance of %v class>", instance.Class.Name.Lexeme)
}
