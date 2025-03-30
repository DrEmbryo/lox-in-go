package main

import (
	"github.com/DrEmbryo/clox/src/utils"
	"github.com/DrEmbryo/clox/src/vm"
)

func main() {
	disassembler := utils.Disassembler{}
	chunk := vm.Chunk{Code: make([]byte, 0), Constants: vm.ValuePool{Value: make([]vm.Value, 0)}}

	constant := chunk.Constants.AddConstant(1.2)
	chunk.WriteChunk(byte(vm.OP_CONSTANT), 123)
	chunk.WriteChunk(byte(constant), 123)
	chunk.WriteChunk(byte(vm.OP_RETURN), 0)
	disassembler.DisassembleChunk(&chunk, "test chunk")
}
