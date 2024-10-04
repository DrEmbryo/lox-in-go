package runtime

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
)

type LoxClass struct {
	Name grammar.Token
}

func (class *LoxClass) ToString() string {
	return fmt.Sprintf("<class  %v>", class.Name.Lexeme)
}
