package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var builtins = [...]string{"type", "exit", "echo", "pwd"}

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
	fmt.Println(strings.Join(args, " "))
}

func Pwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(dir)
	}
}
