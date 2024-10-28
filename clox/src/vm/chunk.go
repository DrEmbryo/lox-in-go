package vm

const (
	OP_RETURN = iota
)

type Chunk struct {
	Code []byte
}

func (chunk *Chunk) WriteChunk(b byte) {
	chunk.Code = append(chunk.Code, b)
}
