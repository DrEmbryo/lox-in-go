package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/lexer"
	"github.com/DrEmbryo/lox/src/parser"
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
	lexer := &lexer.Lexer{Source: []rune(source)}
	loxTokens, lexErrs := lexer.Tokenize();
	if len(lexErrs) > 0 {
		for _, e := range lexErrs {
			grammar.LoxError.Print(e)
		}
	}
	fmt.Println("tokens generated from source:")
	fmt.Println(loxTokens)
	parser := parser.Parser{Tokens: loxTokens}
	stmts, err := parser.Parse()
	if err != nil {
		grammar.LoxError.Print(err)
	}
	fmt.Println("ast generated from tokens:")
	fmt.Println(stmts)
	// parsErr := runtime.Interpriter{}.Interpret(stmts)
	// if len(parsErr) > 0 {
	// 	for _, e := range parsErr {
	// 		grammar.LoxError.Print(e)
	// 	}
	// }
}
