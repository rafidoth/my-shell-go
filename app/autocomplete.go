package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func searchExecutables(partialEx string) (string, error) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", fmt.Errorf("PATH environment variable is empty.\n")
	}

	directories := filepath.SplitList(pathEnv)

	for _, dir := range directories {
		dir = filepath.Clean(dir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			// ignoring non-existed directory
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				// not going to sub directories for now
				continue
			}
			info, err := entry.Info()
			if err != nil {
				continue
			}

			if info.Mode().IsRegular() && info.Mode()&0111 != 0 {
				if strings.HasPrefix(entry.Name(), partialEx) {
					return entry.Name(), nil
				}
			}
		}
	}
	return "", nil
}

func autocomplete(partial string) (string, error) {
	complete, err := searchExecutables(partial)
	if err != nil {
		return "", err
	}

	if complete != "" {
		return complete, nil
	}

	builtins := []string{"echo", "exit"}
	for _, str := range builtins {
		if strings.HasPrefix(str, partial) {
			return str, nil
		}
	}
	return "", nil
}

