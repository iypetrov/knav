package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func confirm(target string) error {
	fmt.Printf("You are about to run a kubectl command against %s\n", target)
	fmt.Print("Do you want to continue? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		return errors.New("Command was aborted")
	}
	return nil
}
