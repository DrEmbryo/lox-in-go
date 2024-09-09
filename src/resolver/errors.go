package resolver

import (
	"fmt"

	"github.com/DrEmbryo/lox/src/grammar"
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
