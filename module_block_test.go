package main

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

func TestGenerateModuleBlockString(t *testing.T) {
	tests := map[string]struct {
		vars    []byte
		sort    bool
		def     bool
		tabsize int
		needs   []string
	}{
		"TypeString": {
			vars:    []byte(`variable "foo" { type = string }`),
			sort:    true,
			def:     false,
			tabsize: 4,
			needs:   []string{"foo = \"\""},
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
			tabsize: 4,
			needs:   []string{`foo = "bar"`},
		},
		"Description": {
			vars:    []byte(`variable "foo" { description = "bar" }`),
			sort:    true,
			def:     true,
			tabsize: 4,
			needs:   []string{"// bar"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dir, _ := createTfFile(tt.vars)
			defer os.RemoveAll(dir)
			modblock, _ := GenerateModuleBlockString(dir, tt.sort, tt.def, tt.tabsize, false)
			for _, need := range tt.needs {
				if !strings.Contains(modblock, need) {
					t.Errorf("modblock (the following) does not include `%s`:\n %s", need, modblock)
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
