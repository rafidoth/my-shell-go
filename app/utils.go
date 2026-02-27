package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"unicode"
)

type RedirectTarget int

const (
	RedirectNone RedirectTarget = iota
	RedirectStdout
	RedirectStderr
)

type Redirect struct {
	Target   RedirectTarget
	Append   bool
	Filename string
}

func extractArguments(wholeCommand string) []string {
	return parseCommand(wholeCommand)
}

func parseRedirect(args []string) (*Redirect, []string) {
	redirectOps := []string{">>", "1>>", "2>>", "2>", ">", "1>"}

	for _, op := range redirectOps {
		if idx := slices.Index(args, op); idx != -1 && idx+1 < len(args) {
			target := RedirectStdout
			append := false

			switch op {
			case ">>", "1>>":
				append = true
			case "2>>":
				target = RedirectStderr
				append = true
			case "2>":
				target = RedirectStderr
			}

			return &Redirect{
				Target:   target,
				Append:   append,
				Filename: args[idx+1],
			}, slices.Delete(slices.Delete(slices.Clone(args), idx, idx+2), idx, idx+2)
		}
	}

	return nil, args
}

func openRedirect(r *Redirect) (*os.File, error) {
	if r == nil {
		return nil, nil
	}

	flags := os.O_WRONLY | os.O_CREATE
	if r.Append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	return os.OpenFile(r.Filename, flags, 0644)
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
	redir, args := parseRedirect(args)

	outFile, err := openRedirect(redir)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	if outFile != nil {
		defer outFile.Close()
	}

	cmd := exec.Command(primary, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if outFile != nil {
		if redir.Target == RedirectStdout {
			cmd.Stdout = outFile
		} else if redir.Target == RedirectStderr {
			cmd.Stderr = outFile
		}
	}

	if runErr := cmd.Run(); runErr != nil {
		if _, pathErr := exec.LookPath(primary); pathErr != nil {
			fmt.Printf("%v: command not found\n", primary)
		}
	}
}
