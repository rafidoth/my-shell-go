package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"
)

func extractArguments(wholeCommand string) []string {
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

func execute(primary string, args []string) {
	gT := slices.Contains(args, ">")
	gT1 := slices.Contains(args, "1>")
	if gT || gT1 {
		var filenameIndex int
		if gT {
			filenameIndex = slices.Index(args, ">") + 1
		} else if gT1 {
			filenameIndex = slices.Index(args, "1>") + 1
		}
		if filenameIndex < len(args) {
			file, err := os.OpenFile(args[filenameIndex], os.O_RDONLY|os.O_CREATE, 0644)
			if err != nil {
			}
			file.Close()
		}
	}

	cmd := exec.Command(primary, args...)
	cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		_, pathErr := exec.LookPath(primary)
		if pathErr != nil {
			fmt.Printf("%v: command not found\n", primary)
		}
		return
	}

}
