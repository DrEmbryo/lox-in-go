package main

import (
	"github.com/DrEmbryo/clox/src/utils"
	"github.com/DrEmbryo/clox/src/vm"
)

func main() {
	disassembler := utils.Disassembler{}
	chunk := vm.Chunk{Code: make([]byte, 0)}
	chunk.WriteChunk(byte(vm.OP_RETURN))
	disassembler.DisassembleChunk(&chunk, "test chunk")
}
