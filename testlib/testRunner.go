package testlib

import (
	"fmt"
	"os"
)

type LoxTestRunnerConfig struct {
	rootDir string
}

type LoxTestRunner struct {
	Config    LoxTestRunnerConfig
	testPaths []string
}

func (runner *LoxTestRunner) Run(path string) {
	sourceRaw, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		panic("Unable to read source")
	}
	source := string(sourceRaw)
	fmt.Println(source)
}

func main() {
	path, err := os.Executable()
	if err != nil {
		panic("Unable to find current working dir")
	}

	runner := &LoxTestRunner{Config: LoxTestRunnerConfig{rootDir: path}}
	fmt.Println(runner)
}
