package main

import (
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

func extractArguments(wholeCommand string) []string {
	splitted := strings.Split(singleQuoteHandler(wholeCommand), " ")
	for i, _ := range splitted {
		splitted[i] = strings.ReplaceAll(splitted[i], "'", "")
	}

	return splitted
}

func singleQuoteHandler(wholeCommand string) string {
	var builder strings.Builder
	insideQuotes := false
	lastRuneWasSpace := false

	input := strings.TrimSpace(wholeCommand)
	// echo 'hello    world'
	// echo     'hello world'
	// echo     'hello       world'
	//
	// echo 'hello
	// echo hello world
	for _, r := range input {
		switch {
		case r == '\'':
			insideQuotes = !insideQuotes
			builder.WriteRune(r)
			lastRuneWasSpace = false
		case unicode.IsSpace(r):
			if insideQuotes {
				builder.WriteRune(r)
			} else {
				if lastRuneWasSpace {
					continue
				} else {
					builder.WriteRune(r)
					lastRuneWasSpace = true
				}
			}
		default:
			builder.WriteRune(r)
			lastRuneWasSpace = false
		}
	}
	return builder.String()
}

func execute(primary string, args []string) {
	// _, err := exec.LookPath(primary)
	// if err == nil {
	cmd := exec.Command(primary, args...)
	out, cmdErr := cmd.Output()
	if cmdErr != nil {
		fmt.Println("error :", cmdErr)
		return
	}
	fmt.Print(string(out))

	// fmt.Printf("%v: command not found\n", primary)
}
