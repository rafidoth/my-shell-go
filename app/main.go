package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

		if command == "exit" {
			break
		} else if commandWithExtractedArgs[0] == "echo" {
			fmt.Println(strings.Join(commandWithExtractedArgs[1:], " "))
		} else {
			fmt.Printf("%v: command not found\n", command)
		}

	}
}
