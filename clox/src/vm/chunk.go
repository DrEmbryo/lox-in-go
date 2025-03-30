package vm

const (
	OP_CONSTANT = iota
	OP_RETURN
)

type Chunk struct {
	Code      []byte
	Lines     []int
	Constants ValuePool
}

func (chunk *Chunk) WriteChunk(b byte, line int) {
	chunk.Code = append(chunk.Code, b)
	chunk.Lines = append(chunk.Lines, line)
}
