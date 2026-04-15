package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

func IsPermitted(trg Target, action string) error {
	if trg.Restricted {
		isAllowed := slices.Contains(trg.AllowedActions, action)
		if !isAllowed {
			err := confirm(trg, action)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func confirm(target Target, action string) error {
	fmt.Printf("You are about to run a \"kubectl %s\" against \"%s\" enviroment.\n", action, target.Name)
	fmt.Print("Do you want to continue? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" && input != "Y" {
		return errors.New("Command was aborted")
	}
	return nil
}
