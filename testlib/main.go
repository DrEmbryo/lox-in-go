package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	TestRunner "github.com/DrEmbryo/lox/testlib/testRunner"
)

func main() {
	options := flag.NewFlagSet("options", flag.ContinueOnError)
	options.Bool("update", false, "Update snapshots without compare")
	options.Parse(os.Args[1:])
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	updateSnapshots, parseErr := strconv.ParseBool(options.Lookup("update").Value.String())
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	runner := &TestRunner.LoxTestRunner{Config: TestRunner.LoxTestRunnerConfig{RootDir: filepath.Join(path, "tests/"), UpdateSnapshots: updateSnapshots}}
	runner.Run()
	fmt.Println(runner)
}
