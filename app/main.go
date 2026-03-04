package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func strProcessor(line []byte) string {
	return strings.TrimSpace(string(line))
}

func readRawInput(fd int) string {
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error :  setting terminal : %v \n", err)
	}
	defer term.Restore(fd, oldState)

	var line []byte
	buf := make([]byte, 1)

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
			partial := strProcessor(line)
			completed := autocomplete(partial)
			if len(partial) < len(completed) {
				fmt.Print(completed[len(partial):])
				line = []byte(completed)
			}
			continue
		}

		if char == 127 || char == 8 {
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Print("\b \b")
			}
			continue
		}

		if char == 3 {
			line = nil
			fmt.Print("\r\n^C")
			break
		}

		line = append(line, char)
		fmt.Print(string(char))
	}

	cmdString := string(line)
	command := strings.TrimSpace(cmdString)
	return command
}

func autocomplete(partial string) string {
	builtins := []string{"echo", "exit"}
	for _, str := range builtins {
		if strings.HasPrefix(str, partial) {
			return str
		}
	}
	return ""
}

func main() {
	fd := int(os.Stdin.Fd())
	for {
		fmt.Print("$ ")
		command := readRawInput(fd)
		if command == "exit" {
			break
		}

		if len(command) == 0 {
			break
		}

		command = strings.TrimSpace(command)
		commandWithExtractedArgs := extractArguments(command)
		primary := commandWithExtractedArgs[0]
		args := commandWithExtractedArgs[1:]

		if command == "exit" {
			break
		} else if primary == "echo" {
			Echo(args...)
		} else if primary == "type" {
			if len(args) == 0 {
				fmt.Println("no arguments")
				continue
			}
			Type(args[0])
		} else if primary == "pwd" {
			Pwd()
		} else if primary == "cd" {
			Cd(args[0])
		} else {
			execute(primary, args)
		}
		// fmt.Println(command)
	}
}
