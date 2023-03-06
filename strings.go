package main

import (
	"fmt"
	"regexp"
	"strings"
)

// IndentByReplacingWords replaces string with space.
func IndentByReplacingWords(s string, size int, w string) string {
	re := regexp.MustCompile("^" + w)
	var result string
	for _, l := range strings.Split(s, "\n") {
		result = fmt.Sprintf("%s%s\n", result, re.ReplaceAllString(l, strings.Repeat(" ", size)))
	}
	return result
}
