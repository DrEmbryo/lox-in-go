package main

import (
	"fmt"
	"os"
	"path/filepath"

	TestRunner "github.com/DrEmbryo/lox/testlib/testRunner"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	runner := &TestRunner.LoxTestRunner{Config: TestRunner.LoxTestRunnerConfig{RootDir: filepath.Join(path, "tests/")}}
	runner.Run()
	fmt.Println(runner)
}
