package module_block

import (
	"log"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestDo(t *testing.T) {
	tests := map[string]struct {
		vars    []byte
		sort    bool
		def     bool
		desc    bool
		tabsize int
		needs   []string
	}{
		"TypeString": {
			vars:    []byte(`variable "foo" { type = string }`),
			sort:    true,
			def:     false,
			desc:    true,
			tabsize: 4,
			needs:   []string{`foo = ""`},
		},
		"TypeStringWithDefaultValue": {
			vars: []byte(`
				variable "foo" {
					type = string
					default = "bar"
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{`foo = "bar"`},
		},
		"TypeNumberWithDefaultValue": {
			vars: []byte(`
			    variable "foo" {
					type = number
					default = 100
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{`foo = 100`},
		},
		"TypeListWithDefaultValue": {
			vars: []byte(`
				variable "foo" {
					type = list(string)
					default = ["bar", "baz"]
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{`foo = ["bar","baz"]`},
		},
		"TypeMapWithDefaultValue": {
			vars: []byte(`
			    variable "foo" {
					type = map(string)
					default = {
						foo = "bar"
						baz = "qux"
					}
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{`foo = {"baz":"qux","foo":"bar"}`},
		},
		"TypeObjectWithDefaultValue": {
			vars: []byte(`
				variable "foo" {
					type = object({ bar = string, baz = number })
					default = {
						bar = "bar"
						baz = 100
					}
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{`foo = {"bar":"bar","baz":100}`},
		},
		"Description": {
			vars: []byte(`
			    variable "foo" {
					type        = string
					description = "bar"
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{"// bar"},
		},
		"MultiLineDescription": {
			vars: []byte(`
			    variable "foo" {
					type        = string
					description = <<EOD
bar
baz
				    EOD
				}
			`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 4,
			needs:   []string{"// bar", "// baz"},
		},
		"No Description": {
			vars: []byte(`
			    variable "foo" {
					type        = string
					description = "bar"
				}
			`),
			sort:    true,
			def:     true,
			desc:    false,
			tabsize: 4,
			needs:   []string{""},
		},
		"TabSize Is 8": {
			vars:    []byte(`variable "foo" { type = string }`),
			sort:    true,
			def:     true,
			desc:    true,
			tabsize: 8,
			needs:   []string{strings.Repeat(" ", 8) + `foo = ""`},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dir, _ := createTfFile(tt.vars)
			defer os.RemoveAll(dir)
			mb, _ := NewModuleBlock(dir, tt.sort, tt.def, tt.tabsize, tt.desc, true).Do()
			for _, need := range tt.needs {
				if !strings.Contains(mb, need) {
					t.Errorf("modblock (the following) does not include `%s`:\n %s", need, mb)
				}
			}
		})
	}
}

func createTfFile(vars []byte) (string, string) {
	dir, err := os.MkdirTemp("/tmp", "tfmodblock-test-*")
	if err != nil {
		log.Fatal(err)
	}
	const FILE = "variables.tf"
	err = os.WriteFile(dir+"/"+FILE, vars, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return dir, FILE
}
