package main

import (
	_ "embed"

	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

var (
	//go:embed module_block.tmpl
	tmpl string
	//go:embed module_block_vscode.tmpl
	vsc_tmpl string
)

// ModuleBlock includes output values consisted of variables.
type ModuleBlock struct {
	Name      string
	Source    string
	Variables []tfconfig.Variable
}

// GenerateModuleBlockString returns Terraform module block string.
func GenerateModuleBlockString(path string, _sort bool, def bool, tabSize int, desc bool, vscode bool) (string, error) {
	if !tfconfig.IsModuleDir(path) {
		return "", fmt.Errorf("given path does not contain tf files")
	}
	// Pass tf file path to tfconfig
	module, _ := tfconfig.LoadModule(path)

	modBlock := new(ModuleBlock)
	fullpath, _ := filepath.Abs(path)
	modBlock.Name = filepath.Base(fullpath)
	cwd, _ := os.Getwd()
	modBlock.Source, _ = filepath.Rel(cwd, fullpath)
	// The result from tfconfig is used to construct modBlock
	constructModuleBlock(modBlock, module.Variables, _sort, def, tabSize, desc)

	_template := tmpl
	if vscode {
		_template = vsc_tmpl
	}

	block, err := template.New("block").Funcs(generateFuncMap()).Parse(_template)
	if err != nil {
		return "", err
	}
	buffer := &bytes.Buffer{}
	// Apply to template
	block.Execute(buffer, modBlock)

	return IndentByReplacingWords(buffer.String(), tabSize), nil
}

// constructModuleBlock constructs ModuleBlock from tfconfig.Variable.
func constructModuleBlock(mb *ModuleBlock, vars map[string]*tfconfig.Variable, _sort bool, def bool, tabSize int, useDesc bool) {
	maxLen := getLongestKeySize(vars)

	for k, v := range vars {
		nm := k + strings.Repeat(" ", maxLen-len(k))
		tp := v.Type
		desc := v.Description
		if !useDesc {
			desc = ""
		}
		df := GetDefaultValue(v, def, tp)

		mb.Variables = append(
			mb.Variables,
			tfconfig.Variable{Type: tp, Name: nm, Description: desc, Default: df},
		)
	}

	if !_sort {
		return
	}

	sort.Slice(mb.Variables, func(i, j int) bool { return mb.Variables[i].Name < mb.Variables[j].Name })
}

// getLongestKeySize returns the longest key size in the map.
func getLongestKeySize(vars map[string]*tfconfig.Variable) int {
	max := 0
	for key := range vars {
		_len := len(key)
		if _len > max {
			max = _len
		}
	}
	return max
}

// generateFuncMap returns FuncMap used in template.
func generateFuncMap() template.FuncMap {
	return template.FuncMap{
		"getDefaultValueByType": GetDefaultValueByType,
	}
}
