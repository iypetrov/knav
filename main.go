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
	var trg Target

	err := CreateDefaultConfigIfEmpty()
	if err != nil {
		panic(err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		cfg = DefaultConfig()
	}

	if len(os.Args) == 1 {
		trg, err = cfg.PickTarget()
		if err != nil {
			panic(err)
		}
	} else {
		trg, err = cfg.CurrentTarget()
		if err != nil {
			panic(err)
		}
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

	err = confirm(trg.Name)
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
