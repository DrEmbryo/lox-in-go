package vm

type Value float32

type ValuePool struct {
	Value []Value
}

func (pool *ValuePool) AddConstant(b Value) int {
	pool.Value = append(pool.Value, b)
	return len(pool.Value)
}
