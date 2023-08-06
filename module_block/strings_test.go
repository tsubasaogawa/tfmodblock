package module_block

import (
	"strings"
	"testing"
)

func TestIndentByReplacingWords(t *testing.T) {
	tests := map[string]struct {
		s    string
		size int
		need string
	}{
		"OneLine": {
			s:    "  lorem",
			size: 4,
			need: strings.Repeat(" ", 4) + "lorem",
		},
		"MultiLine": {
			s:    "  lorem\n  ipsum",
			size: 4,
			need: strings.Repeat(" ", 4) + "lorem" + "\n" + strings.Repeat(" ", 4) + "ipsum",
		},
		"DoNotIndentWhenNoSpaceIsGiven": {
			s:    "lorem",
			size: 4,
			need: "lorem",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := IndentByReplacingWords(tt.s, tt.size)
			if result != tt.need {
				t.Errorf("Need: %s, Actual: %s\n", tt.need, result)
			}
		})
	}
}
