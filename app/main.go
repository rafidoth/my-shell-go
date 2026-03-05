package main

import (
	"fmt"
	"os"
	"strings"
)

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
	}
}
