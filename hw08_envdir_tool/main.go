package main

import (
	"flag"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		return // missing arguments
	}
	env, err := ReadDir(args[0])
	if err != nil {
		return // failed to read env directory
	}

	RunCmd(args[1:], env)
}
