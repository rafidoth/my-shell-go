package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/term"
)

func handleTabPress(line []byte, doubleMatch bool) ([]byte, bool) {
	partial := strProcessor(line)
	completed, err := autocomplete(partial)
	if err != nil {
		fmt.Print(err)
		return line, doubleMatch
	}

	// single completion exists
	if len(completed) == 1 && len(partial) < len(completed[0]) {
		final := completed[0] + " "
		fmt.Print(final[len(partial):])
		return []byte(final), doubleMatch
	}

	// multiple completions exists
	// and no longest prefix
	if len(completed) > 1 {
		fmt.Print("\x07")
		lcp := getLongestCommonPrefix(completed)
		if lcp != "" && len(lcp) > len(partial) {
			fmt.Print(lcp[len(partial):])
			return []byte(lcp), true
		}
		fmt.Print("\r\n")
		sort.Strings(completed)
		for _, sug := range completed {
			fmt.Printf("%v  ", sug)
		}
		fmt.Print("\r\n$ ", string(line))
		return line, true
	}

	if len(completed) == 0 {
		fmt.Print("\x07")
	}
	return line, doubleMatch
}

func handleBackspace(line []byte) []byte {
	if len(line) > 0 {
		line = line[:len(line)-1]
		fmt.Print("\b \b")
	}
	return line
}

func handleCtrlC() {
	fmt.Print("\r\n^C")
}

func readRawInput(fd int) string {
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error :  setting terminal : %v \n", err)
	}
	defer term.Restore(fd, oldState)

	var line []byte
	buf := make([]byte, 1)
	doubleMatchBell := false

	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			break
		}
		char := buf[0]

		if char == '\n' || char == '\r' {
			fmt.Print("\r\n")
			return strProcessor(line)
		}

		if char == '\t' {
			line, doubleMatchBell = handleTabPress(line, doubleMatchBell)
			continue
		}

		if char == 127 || char == 8 {
			line = handleBackspace(line)
			continue
		}

		if char == 3 {
			handleCtrlC()
			break
		}

		line = append(line, char)
		fmt.Print(string(char))
	}

	cmdString := string(line)
	command := strings.TrimSpace(cmdString)
	return command
}
