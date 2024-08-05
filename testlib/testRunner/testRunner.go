package testRunner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DrEmbryo/lox/src/grammar"
	"github.com/DrEmbryo/lox/src/lexer"
	"github.com/DrEmbryo/lox/src/parser"
)

type LoxTestRunnerConfig struct {
	RootDir string
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

	// add env dump for the eval stage
	// hooks can be usefull

}

func (runner *LoxTestRunner) checkSnapshot(path string, postfix string, payload string) {
	snapshotPath := strings.Join([]string{path, postfix}, "")
	snapshot, err := os.ReadFile(snapshotPath)
	if err != nil {
		os.WriteFile(snapshotPath, []byte(payload), os.ModeAppend)
	} else {
		fmt.Println(strings.Compare(string(snapshot), payload))
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
