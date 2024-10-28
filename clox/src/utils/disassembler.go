package utils

import (
	"fmt"

	"github.com/DrEmbryo/clox/src/vm"
)

type Disassembler struct{}

func (disassembler *Disassembler) DisassembleChunk(chunk *vm.Chunk, name string) {
	fmt.Printf("=== %s ==\n", name)
	for offset := 0; offset < len(chunk.Code); {
		offset = disassembler.disassembleInstruction(chunk, offset)
	}
}

func (disassembler *Disassembler) disassembleInstruction(chunk *vm.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	instruction := chunk.Code[offset]
	switch instruction {
	case vm.OP_RETURN:
		return disassembler.simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown ocode %d\n", instruction)
		return offset + 1
	}
}

func (disassembler *Disassembler) simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
