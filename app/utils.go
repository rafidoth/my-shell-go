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
	// d := doubleQuoteHandler(wholeCommand)
	// fmt.Println(d)
	splitted := parseCommand(wholeCommand)
	return splitted
}
func parseCommand(input string) []string {
	var parsed []string
	var builder strings.Builder

	singleQuote := false
	doubleQuote := false
	escaped := false

	input = strings.TrimSpace(input)

	for _, r := range input {
		if escaped {
			builder.WriteRune(r)
			escaped = false
			continue
		}

		switch {
		case r == '\\' && !singleQuote:
			escaped = true

		case r == '"' && !singleQuote:
			doubleQuote = !doubleQuote

		case r == '\'' && !doubleQuote:
			singleQuote = !singleQuote

		case unicode.IsSpace(r) && !singleQuote && !doubleQuote:
			if builder.Len() > 0 {
				parsed = append(parsed, builder.String())
				builder.Reset()
			}

		default:
			builder.WriteRune(r)
		}
	}

	if escaped {
		builder.WriteRune('\\')
	}
	if builder.Len() > 0 {
		parsed = append(parsed, builder.String())
	}

	return parsed
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
	singleQuote := false
	doubleQuote := false
	lastRuneWasSpace := false
	escaped := false

	input := strings.TrimSpace(wholeCommand)
	for _, r := range input {
		if escaped {
			builder.WriteRune(r)
			escaped = false
			lastRuneWasSpace = false
			continue
		}

		switch {
		case r == '\\' && !singleQuote:
			escaped = true

		case r == '"' && !singleQuote:
			doubleQuote = !doubleQuote
			builder.WriteRune(r)
			lastRuneWasSpace = false

		case r == '\'' && !doubleQuote:
			singleQuote = !singleQuote
			builder.WriteRune(r)
			lastRuneWasSpace = false

		case unicode.IsSpace(r):
			if singleQuote || doubleQuote {
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

	if escaped {
		builder.WriteRune('\\')
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
