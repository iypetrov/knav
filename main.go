package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		return
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

	actionIndex := knavIndex + 1
	action := args[actionIndex]
	err = IsPermitted(trg, action)
	if err != nil {
		fmt.Println(err)
		os.Exit(StatusAbortedOperation)
	}

	kubectlArgs := args[knavIndex+1:]
	cmd := exec.Command("kubectl", kubectlArgs...)
	env := os.Environ()
	kPath, err := expandPath(trg.KubeconfigPath)
	if err != nil {
		panic(err)
	}
	env = append(env, "KUBECONFIG="+kPath)
	for _, e := range trg.Envs {
		if e.Name == "" {
			continue
		}
		env = append(env, e.Name+"="+e.Value)
	}
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(StatusKubectlCommandFailed)
	}
}

func expandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, p[1:])
	}
	return filepath.Abs(p)
}
