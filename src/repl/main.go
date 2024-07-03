package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/lexer"
	"github.com/DrEmbryo/lox/src/parser"
	"github.com/DrEmbryo/lox/src/runtime"
)

func main() {
	options := flag.NewFlagSet("options", flag.ContinueOnError)
	options.Bool("debug", false, "Run REPL in debug mode")

	var source string
	if len(os.Args) < 2 || strings.Contains(os.Args[1], "-") {
		options.Parse(os.Args[1:])
		fmt.Println("Lox REPL 0.2: ")
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")
			source, _ = reader.ReadString('\n')
			eval(source, options)
		}

	} else {
		options.Parse(os.Args[2:])
		sourceRaw, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			panic("Unable to read source")
		}
		source = string(sourceRaw)
		eval(source, options)
	}
}

func eval(source string, options *flag.FlagSet) {
	debugOption, parseErr := strconv.ParseBool(options.Lookup("debug").Value.String())
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	lexer := &lexer.Lexer{Source: []rune(source)}
	loxTokens, lexErrs := lexer.Tokenize()
	if len(lexErrs) > 0 {
		for _, e := range lexErrs {
			grammar.LoxError.Print(e)
		}
	}
	if debugOption {
		fmt.Println("tokens generated from source:")
		fmt.Println(loxTokens)
	}
	parser := parser.Parser{Tokens: loxTokens}
	stmts, err := parser.Parse()
	if err != nil {
		grammar.LoxError.Print(err)
	}
	if debugOption {
		fmt.Println("ast generated from tokens:")
		fmt.Println(stmts)
	}
	Interpreter := runtime.Interpreter{Env: runtime.Environment{Values: make(map[string]any), Parent: nil}}
	errs := Interpreter.Interpret(stmts)
	if len(errs) > 0 {
		for _, e := range errs {
			grammar.LoxError.Print(e)
		}
	}
}
