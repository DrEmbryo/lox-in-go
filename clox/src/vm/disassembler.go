package vm

import "fmt"

type Disassembler struct{}

func (disassembler *Disassembler) DisassembleChunk(chunk Chunk, name string) {
	fmt.Printf("=== %s ==\n", name)
	for offset := 0; offset < len(chunk.Code); {
		offset = disassembler.disassembleInstruction(chunk, offset)
	}
}

func (disassembler *Disassembler) disassembleInstruction(chunk Chunk, offset int) int {
	if offset > 0 && offset == chunk.Lines[offset-1] {
		fmt.Print(" | ")
	} else {
		fmt.Printf("%04d %d ", offset, chunk.Lines[offset])
	}
	instruction := chunk.Code[offset]
	switch instruction {
	case OP_CONSTANT:
		return disassembler.constantInstruction("OP_CONSTANT", chunk, offset)
	case OP_RETURN:
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

func (disassembler *Disassembler) constantInstruction(name string, chunk Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%s %04d '", name, constant)
	fmt.Printf("%v", chunk.Constants.Value[constant-1])
	fmt.Printf("'\n")
	return offset + 2
}
