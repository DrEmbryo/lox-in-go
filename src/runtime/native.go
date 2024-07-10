package runtime

type LoxCallable interface {
	Call(interpreter Interpreter, arguments []any) any
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

func (native *NativeCall) Call(interpreter Interpreter, arguments []any) any {
	return native.NativeCallFunc()
}

func (native *NativeCall) ToString() string {
	return "<native func>"
}
