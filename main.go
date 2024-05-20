package main

import (
	"bufio"
	"fmt"
	"os"

	lexer "github.com/DrEmbryo/lox/frontend"
)


func main() {
	var source string
	if len(os.Args) < 2 {
		fmt.Println("Lox REPL 0.0.1: ")
		for {			
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			source, _ = reader.ReadString('\n')
			eval(source)
		}
		} else {
		sourceRaw, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			panic("Unable to read source")
		}
		source = string(sourceRaw)
	}
	eval(source)
}

func eval (source string) {
	tokens, err := lexer.Tokenize(source)
	if err != nil {
		panic(err)
	}

	fmt.Println(tokens)
}
