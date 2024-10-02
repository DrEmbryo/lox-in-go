package runtime

import (
	"github.com/DrEmbryo/lox/src/grammar"
)

type LoxCallable interface {
	Call(interpreter Interpreter, arguments []any) (any, grammar.LoxError)
	GetAirity() int
}

type NativeCallFunc func(...any) any

type NativeCall struct {
	Airity         int
	NativeCallFunc NativeCallFunc
}

func (native *NativeCall) GetAirity() int {
	return native.Airity
}

func (native *NativeCall) Call(interpreter Interpreter, arguments []any) (any, grammar.LoxError) {
	return native.NativeCallFunc(arguments), nil
}

func (native *NativeCall) ToString() string {
	return "<native func>"
}
