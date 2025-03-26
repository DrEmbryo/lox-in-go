package main

import (
	"github.com/DrEmbryo/clox/src/utils"
	"github.com/DrEmbryo/clox/src/vm"
)

func main() {
	disassembler := utils.Disassembler{}
	chunk := vm.Chunk{Code: make([]byte, 0), Constants: vm.ValuePool{Value: make([]vm.Value, 0)}}

	constant := chunk.Constants.AddConstant(1.2)
	chunk.WriteChunk(byte(vm.OP_CONSTANT))
	chunk.WriteChunk(byte(constant))
	chunk.WriteChunk(byte(vm.OP_RETURN))
	disassembler.DisassembleChunk(&chunk, "test chunk")
}
