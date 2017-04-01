package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm shows a prompt on the screen and waits for a user response
func Confirm(question string, def bool) bool {
	var s string

ask:

	fmt.Printf("%s ", question)
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" {
		return true
	} else if s == "n" {
		return false
	} else if s == "" {
		return def
	}
	fmt.Println("\nHm. Please enter y or n.")
	goto ask
}
