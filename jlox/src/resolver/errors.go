package resolver

import (
	"fmt"

	"github.com/DrEmbryo/jlox/src/grammar"
)

type ResolverError struct {
	Token   grammar.Token
	Message string
}

func (e ResolverError) Print() {
	fmt.Printf("[%v]: Resolver error: %s", e.Token, e.Message)
}

func (e ResolverError) Error() string {
	return fmt.Sprintf("[%v]: Resolver error: %s", e.Token, e.Message)
}
