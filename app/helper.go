package main

import "strings"

func strProcessor(line []byte) string {
	return strings.TrimSpace(string(line))
}

func getLongestCommonPrefix(strs []string) string {
	if len(strs) < 1 {
		return ""
	}
	res := ""
	for idx, r := range strs[0] {
		for _, str := range strs {
			if idx >= len(str) || r != rune(str[idx]) {
				return res
			}
		}
		res = res + string(r)
	}
	return res
}
