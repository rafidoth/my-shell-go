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
	redirectIndex := -1
	if idx := slices.Index(args, ">"); idx != -1 {
		redirectIndex = idx
	} else if idx := slices.Index(args, "1>"); idx != -1 {
		redirectIndex = idx
	}

	var outFile *os.File
	var err error

	if redirectIndex != -1 && redirectIndex+1 < len(args) {
		filename := args[redirectIndex+1]

		outFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer outFile.Close()

		args = args[:redirectIndex]
	}

	cmd := exec.Command(primary, args...)

	if outFile != nil {
		cmd.Stdout = outFile
	} else {
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr

	if runErr := cmd.Run(); runErr != nil {
		if _, pathErr := exec.LookPath(primary); pathErr != nil {
			fmt.Printf("%v: command not found\n", primary)
		}
		return
	}
}
