package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/runtime"
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
	loxTokens, errs := grammar.Lexer{}.Tokenize(source);
	if len(errs) > 0 {
		for _, e := range errs {
			grammar.LoxError.Print(e)
		}
	}
	fmt.Println("tokens generated from source:")
	fmt.Println(loxTokens)
	ast, err := grammar.Parser{}.Parse(loxTokens)
	if err != nil {
		grammar.LoxError.Print(err)
	}
	fmt.Println("ast generated from tokens:")
	fmt.Println(ast)
	value, err := runtime.Interpriter{}.Interpret(ast)
	if err != nil {
		grammar.LoxError.Print(err)
		return 
	}
	fmt.Printf("%v", value)
}
