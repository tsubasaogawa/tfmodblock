package main

import (
	"testing"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

func TestGetDefaultValue(t *testing.T) {
	tests := map[string]struct {
		_var *tfconfig.Variable
		def  bool
		tp   string
		need string
	}{
		"DefaultIsNull": {
			_var: &tfconfig.Variable{},
			def:  true,
			tp:   "string",
			need: `""`,
		},
		"DefaultIsFoo": {
			_var: &tfconfig.Variable{Default: "foo"},
			def:  true,
			tp:   "string",
			need: `"foo"`,
		},
		"DefaultIs123": {
			_var: &tfconfig.Variable{Default: 123},
			def:  true,
			tp:   "number",
			need: "123",
		},
		"DefaultIsList": {
			_var: &tfconfig.Variable{Default: []string{"a", "b", "c"}},
			def:  true,
			tp:   "list",
			need: `["a","b","c"]`,
		},
		"UseType": {
			_var: &tfconfig.Variable{},
			def:  false,
			tp:   "string",
			need: `""`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := GetDefaultValue(tt._var, tt.def, tt.tp)
			if result != tt.need {
				t.Errorf("Need: %s, Actual: %s\n", tt.need, result)
			}
		})
	}
}
