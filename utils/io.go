package utils

import (
	"fmt"
	"strings"
)

func Confirm(question string, def bool) bool {
	var s string

ask:

	fmt.Printf("%s ", question)
	_, err := fmt.Scan(&s)
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
