# Lox interptiter in Go

## About this project

Recreation of the (j/c)Lox programming language from the well-regarded book "Crafting Interpreters" by Robert Nystrom with usage of Go progtamming language. This project involves both learning and implementing the concepts and techniques required to build a complete interpreter from scratch.

Features implemented in jLox:

- tokens and lexing
- abstract syntax trees
- recursive descent parsing
- prefix and infix expressions
- runtime representation of objects
- interpreting code using the Visitor pattern
- lexical scope
- environment chains for storing variables
- control flow
- functions with parameters
- closures
- static variable resolution and error detection
- classes
- constructors
- fields
- methods
- inheritance

### Installation option:

- Nix flake:
  1. Allow [nix flakes](https://nixos.wiki/wiki/Flakes) in your configuration
  2. Open dev shell in your terminal by running `nix develop`
  3. Execute `go run main.go` script
- Without usage of nix
  1. Install golang (current version is 1.22.3)
  2. Execute `go run main.go` script

### Running the app:

- Execute the `go run main.go` script in src/repl directory
- The option list can be accessed by providing the `-h` flag.
- For the compilation of the app please see official [golang docs](https://pkg.go.dev/cmd/go@go1.23.2)

#### Example of running lox script

- `go run main.go <file>.lox` will run the script from the file with the provided path
- `go run main.go` will run lox in REPL mode
