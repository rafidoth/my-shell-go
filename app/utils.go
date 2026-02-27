package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func extractArguments(wholeCommand string) []string {
	return strings.Split(wholeCommand, " ")
}

func execute(primary string, args []string) {
	_, err := exec.LookPath(primary)
	if err == nil {
		cmd := exec.Command(primary, args...)
		out, cmdErr := cmd.Output()
		if cmdErr != nil {
			fmt.Println("error :", cmdErr)
			return
		}
		fmt.Print(string(out))
		return
	}

	fmt.Printf("%v: command not found\n", primary)
}
