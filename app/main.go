package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func extractArguments(wholeCommand string) []string {
	return strings.Split(wholeCommand, " ")
}

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
		primary := commandWithExtractedArgs[0]
		args := commandWithExtractedArgs[1:]

		if command == "exit" {
			break
		} else if primary == "echo" {
			fmt.Println(strings.Join(args, " "))
		} else if primary == "type" {
			if len(args) == 0 {
				fmt.Println("no arguments")
				continue
			}
			builtins := [...]string{"type", "exit", "echo"}
			matched := false
			for _, val := range builtins {
				if val == args[0] {
					fmt.Println(args[0], "is a shell builtin")
					matched = true
					break
				}
			}

			if !matched {
				path, err := exec.LookPath(args[0])
				if err == nil {
					fmt.Println(args[0], "is", path)
					continue
				}
				fmt.Printf("%v: not found\n", args[0])
			}
		} else {
			fmt.Printf("%v: command not found\n", command)
		}
	}
}
