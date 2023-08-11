package module_block

import (
	"regexp"
	"strings"
)

const INDENT = "  "

var RE = regexp.MustCompile("^" + INDENT)

// IndentByReplacingWords replaces string with space.
func IndentByReplacingWords(s string, size int) string {
	var result string

	lfCnt := strings.Count(s, "\n")
	for i, l := range strings.Split(s, "\n") {
		result = result + RE.ReplaceAllString(l, strings.Repeat(" ", size))
		// do not insert lf when `l` is eol
		if i != lfCnt {
			result = result + "\n"
		}
	}
	return result
}
