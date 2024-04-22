package gocommand

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
	"syscall"
)

// Command info about entered command
type Command struct {
	// Name first word of given command line
	Name string
	// Args arguments of written command, separated by spaces
	Args []string
}

var ErrGracefulExit = fmt.Errorf("graceful exit")

// ListenAndServe reads commands from input and calls handler for each command.
// If handler returns ErrGracefulExit, ListenAndServe stops without error message; If handler returns err != nil, ListenAndServe stops and print error message; If handler returns nil, ListenAndServe continues to read commands.
func ListenAndServe(handler func(*Command) error) {
	for {
		cmd, err := ReadCommand()
		if err != nil {
			fmt.Printf("\033[31mCannot read command: %s\033[0m\n", err)

			return
		}

		err = handler(cmd)
		if err != nil {
			if errors.Is(err, ErrGracefulExit) {
				fmt.Println("Bye!")
			} else {
				fmt.Printf("\033[31mError: %s\033[0m\n", err)
			}

			return
		}
	}
}

// ReadCommand prompts user to enter a command to input, writes command anchor to output
func ReadCommand() (*Command, error) {
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("cannot read command: %w", err)
	}

	strs := strings.Split(text, " ")
	strs[len(strs)-1] = strings.Trim(strs[len(strs)-1], "\n")

	return &Command{
		Name: strings.TrimSpace(strs[0]),
		Args: strs,
	}, nil
}

// AskText prompts user to enter a text, writes prompt to output as an anchor and separates it with a colon
func AskText(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	r := bufio.NewReader(os.Stdin)
	val, err := r.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("cannot ask text: %w", err)
	}

	return strings.Trim(val, "\n"), nil
}

// AskSecret prompts user to enter a secret, writes prompt to output as an anchor and separates it with a colon
func AskSecret(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	secret, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", fmt.Errorf("cannot read secret: %w", err)
	}
	fmt.Printf("\n")

	return string(secret), err
}
