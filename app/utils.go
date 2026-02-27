package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func extractArguments(wholeCommand string) []string {
	splitted := commandParser(singleQuoteHandler(wholeCommand))
	return splitted
}

func commandParser(singleQuoteHandledCmd string) []string {
	// fmt.Println(singleQuoteHandledCmd)
	var parsed []string

	var builder strings.Builder
	insideQuotes := false
	for _, r := range singleQuoteHandledCmd {
		switch {
		case r == '\'':
			insideQuotes = !insideQuotes
		case unicode.IsSpace(r) && !insideQuotes:
			parsed = append(parsed, builder.String())
			builder.Reset()
		default:
			builder.WriteRune(r)
		}
	}
	if builder.Len() > 0 {
		parsed = append(parsed, builder.String())
	}
	return parsed
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
	cmd := exec.Command(primary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		_, pathErr := exec.LookPath(primary)
		if pathErr != nil {
			fmt.Printf("%v: command not found\n", primary)
		}
		return
	}

}
