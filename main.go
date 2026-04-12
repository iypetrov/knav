package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	StatusWrongNumberOfArguments = 1
	StatusKnavCommandNoFound     = 2
	StatusAbortedOperation       = 3
	StatusKubectlCommandFailed   = 4
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: knav <kubectl args>")
		os.Exit(StatusWrongNumberOfArguments)
	}

	args := os.Args[0:]

	knavIndex := -1
	for i, arg := range args {
		if arg == "knav" {
			knavIndex = i
			break
		}
	}

	if knavIndex == -1 {
		fmt.Println("Error: 'knav' not found in arguments")
		os.Exit(StatusKnavCommandNoFound)
	}

	err := confirm("local")
	if err != nil {
		fmt.Println(err)
		os.Exit(StatusAbortedOperation)
	}

	kubectlArgs := args[knavIndex+1:]

	cmd := exec.Command("kubectl", kubectlArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(StatusKubectlCommandFailed)
	}
}
