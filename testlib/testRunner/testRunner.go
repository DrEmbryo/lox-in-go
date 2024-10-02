package testRunner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/lexer"
	"github.com/DrEmbryo/lox/src/parser"
	"github.com/DrEmbryo/lox/src/runtime"
)

type LoxTestRunnerConfig struct {
	RootDir         string
	UpdateSnapshots bool
}

type LoxTestRunner struct {
	Config    LoxTestRunnerConfig
	testPaths []string
}

func (runner *LoxTestRunner) Run() {
	runner.collectPaths()
	for _, path := range runner.testPaths {
		runner.runTest(path)
	}
}

func (runner *LoxTestRunner) collectPaths() {
	err := filepath.Walk(runner.Config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() && !strings.Contains(path, ".snap") {
			runner.testPaths = append(runner.testPaths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (runner *LoxTestRunner) runTest(path string) {
	sourceRaw, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		panic("Unable to read source")
	}
	source := string(sourceRaw)

	loxTokens := runner.lex(source)
	runner.checkSnapshot(path, ".token.snap", fmt.Sprintf("%v", loxTokens))

	ast := runner.parse(loxTokens)
	runner.checkSnapshot(path, ".ast.snap", fmt.Sprintf("%v", ast))

	env := runner.interpret(ast)
	runner.checkSnapshot(path, ".env.snap", fmt.Sprintf("%v", env))

}

func (runner *LoxTestRunner) checkSnapshot(path string, postfix string, payload string) {
	snapshotPath := strings.Join([]string{path, postfix}, "")
	snapshot, err := os.ReadFile(snapshotPath)
	if err != nil {
		os.WriteFile(snapshotPath, []byte(payload), os.ModeAppend)
	} else {
		if runner.Config.UpdateSnapshots {
			os.WriteFile(snapshotPath, []byte(payload), fs.FileMode(os.O_TRUNC))
			fmt.Printf("Update snap at %s\n", snapshotPath)
			return
		}
		snapshotMatch := strings.Compare(string(snapshot), payload)
		if snapshotMatch != 0 {
			fmt.Printf("Failed snap at %s;\n expected %v;\n got %v;\n", snapshotPath, string(snapshot), payload)
		}
	}
}

func (runner *LoxTestRunner) lex(source string) []grammar.Token {
	lexer := &lexer.Lexer{Source: []rune(source)}
	loxTokens, lexErrs := lexer.Tokenize()
	if len(lexErrs) > 0 {
		for _, e := range lexErrs {
			grammar.LoxError.Print(e)
		}
	}
	return loxTokens
}

func (runner *LoxTestRunner) parse(loxTokens []grammar.Token) []grammar.Statement {
	parser := parser.Parser{Tokens: loxTokens}
	stmts, err := parser.Parse()
	if err != nil {
		grammar.LoxError.Print(err)
	}
	return stmts
}
func (runner *LoxTestRunner) interpret(stmts []grammar.Statement) runtime.Environment {
	env := runtime.Environment{Values: make(map[string]any), Parent: nil}
	interpreter := runtime.Interpreter{Env: env}
	interpreter.Interpret(stmts)
	return env
}
