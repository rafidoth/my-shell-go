package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Print("$ ")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if len(command) == 0 {
			break
		}

		command = strings.TrimSpace(command)
		commandWithExtractedArgs := extractArguments(command)
		// fmt.Println(commandWithExtractedArgs)
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
