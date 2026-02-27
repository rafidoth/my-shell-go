package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func extractArguments(wholeCommand string) []string {
	// s := singleQuoteHandler(wholeCommand)
	// fmt.Println(s)
	d := doubleQuoteHandler(wholeCommand)
	// fmt.Println(d)
	splitted := commandParser(d)
	return splitted
}

func commandParser(quoteHandled string) []string {
	// fmt.Println(singleQuoteHandledCmd)
	var parsed []string

	var builder strings.Builder
	insideQuotes := false
	doubleQuote := false
	singleQuote := false
	for _, r := range quoteHandled {
		switch {
		case r == '"' && !singleQuote:
			insideQuotes = !insideQuotes
			doubleQuote = !doubleQuote
		case r == '\'' && !doubleQuote:
			insideQuotes = !insideQuotes
			singleQuote = !singleQuote
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

func doubleQuoteHandler(wholeCommand string) string {
	var builder strings.Builder
	insideQuotes := false
	singleQuote := false
	dobleQuote := false
	lastRuneWasSpace := false
	input := strings.TrimSpace(wholeCommand)
	for _, r := range input {
		switch {
		case r == '"' && !singleQuote:
			insideQuotes = !insideQuotes
			dobleQuote = !dobleQuote
			builder.WriteRune(r)
			lastRuneWasSpace = false
		case r == '\'' && !dobleQuote:
			insideQuotes = !insideQuotes
			singleQuote = !singleQuote
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
