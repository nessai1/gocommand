package gocommand

import (
	"fmt"
)

func ExampleReadCommand() {
	ListenAndServe(func(cmd *Command) error {
		if cmd.Name == "exit" {
			return ErrGracefulExit
		}

		fmt.Printf("Command: %s\n", cmd.Name)
		fmt.Printf("Arguments: %v\n", cmd.Args)

		return nil
	})
}
