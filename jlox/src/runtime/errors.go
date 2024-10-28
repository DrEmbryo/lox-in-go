package runtime

import (
	"fmt"

	"github.com/DrEmbryo/jlox/src/grammar"
)

type RuntimeError struct {
	Token   grammar.Token
	Message string
}

func (e RuntimeError) Print() {
	fmt.Printf("[%v]: Runtime error: %s", e.Token, e.Message)
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[%v]: Runtime error: %s", e.Token, e.Message)
}
