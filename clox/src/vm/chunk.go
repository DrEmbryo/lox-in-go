package vm

const (
	OP_CONSTANT = iota
	OP_RETURN
)

type Chunk struct {
	Code      []byte
	Constants ValuePool
}

func (chunk *Chunk) WriteChunk(b byte) {
	chunk.Code = append(chunk.Code, b)
}
