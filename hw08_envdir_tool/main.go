package main

import (
	"fmt"
	"os"
)

func printUsage(name string) {
	fmt.Printf("Usage:\n%s <env files directory> <file to run> [run params]\n", name)
}

func main() {
	if len(os.Args) < 3 {
		printUsage(os.Args[0])
		os.Exit(1)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("Error while reading env dir: %s\n", err)
	}
	os.Exit(RunCmd(os.Args[2:], env))
}
