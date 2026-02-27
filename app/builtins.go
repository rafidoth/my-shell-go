package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtins = [...]string{"type", "exit", "echo", "pwd", "cd"}

func Type(arg string) {
	matched := false
	for _, val := range builtins {
		if val == arg {
			fmt.Println(arg, "is a shell builtin")
			matched = true
			return
		}
	}

	if !matched {
		path, err := exec.LookPath(arg)
		if err == nil {
			fmt.Println(arg, "is", path)
			return
		}
		fmt.Printf("%v: not found\n", arg)
	}
}

func Echo(args ...string) {
	redirectIndex := -1
	isStderrRedirect := false
	isStdoutRedirect := false
	isStdoutAppend := false
	isStderrAppend := false

	if idx := slices.Index(args, ">>"); idx != -1 {
		redirectIndex = idx
		isStdoutAppend = true
	} else if idx := slices.Index(args, "1>>"); idx != -1 {
		redirectIndex = idx
		isStdoutAppend = true
	} else if idx := slices.Index(args, "2>>"); idx != -1 {
		redirectIndex = idx
		isStderrAppend = true
	} else if idx := slices.Index(args, "2>"); idx != -1 {
		redirectIndex = idx
		isStderrRedirect = true
	} else if idx := slices.Index(args, ">"); idx != -1 {
		redirectIndex = idx
		isStdoutRedirect = true
	} else if idx := slices.Index(args, "1>"); idx != -1 {
		redirectIndex = idx
		isStdoutRedirect = true
	}

	if redirectIndex != -1 && redirectIndex+1 < len(args) {
		filename := args[redirectIndex+1]
		output := strings.Join(args[:redirectIndex], " ") + "\n"

		var outFile *os.File
		var err error

		if isStdoutRedirect || isStderrRedirect {
			outFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		} else if isStdoutAppend || isStderrAppend {
			outFile, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		}

		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer outFile.Close()

		if isStdoutRedirect || isStdoutAppend {
			outFile.WriteString(output)
		} else if isStderrRedirect || isStderrAppend {
			fmt.Print(output)
		}
	} else {
		output := strings.Join(args, " ")
		fmt.Println(output)
	}
}

func Pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(dir)
	}
}

func Cd(arg string) {
	dir := arg
	if dir == "~" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("error :", err)
		}
		dir = userHome
	}
	err := os.Chdir(dir)
	if err != nil {
		fmt.Printf("cd: %v: No such file or directory\n", dir)
	}
}
