package vm

import "fmt"

const (
	INTERPRET_OK = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

const DEBUG_TRACE_EXECUTION = true

type VM struct {
	Chunk *Chunk
	Ip    int
}

func (vm *VM) Interpret(chunk *Chunk) int {
	vm.Chunk = chunk
	vm.Ip = 0
	return vm.Run()
}

func (vm *VM) readByte() byte {
	return vm.Chunk.Code[vm.Ip]
}

func (vm *VM) readConstant() Value {
	return vm.Chunk.Constants.Value[vm.readByte()]
}

func (vm *VM) Run() int {
	for {
		instruction := vm.readByte()
		switch instruction {
		case OP_CONSTANT:
			vm.handleConstantOp()
		case OP_RETURN:
			return vm.handleReturnOp()
		}
		vm.Ip++
	}
}

func (vm *VM) handleConstantOp() {
	fmt.Printf("%v \n", vm.readConstant())
}

func (vm *VM) handleReturnOp() int {
	return INTERPRET_OK
}
