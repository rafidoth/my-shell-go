package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func searchExecutables(partialEx string) ([]string, error) {
	var matched []string
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return nil, fmt.Errorf("PATH environment variable is empty.\n")
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
					matched = append(matched, entry.Name())
				}
			}
		}
	}
	return matched, nil
}

func autocomplete(partial string) ([]string, error) {
	matches, err := searchExecutables(partial)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		builtins := []string{"echo", "exit"}
		for _, str := range builtins {
			if strings.HasPrefix(str, partial) {
				matches = append(matches, str)
			}
		}
	}
	return matches, nil
}

