package module_block

import (
	_ "embed"

	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

	REGEX_HAT = regexp.MustCompile("^")
)

// ModuleBlock includes output values consisted of variables.
type ModuleBlock struct {
	path, Name, Source         string
	sort, def, useDesc, vscode bool
	tabSize                    int
	Variables                  []tfconfig.Variable
}

//
func NewModuleBlock(path string, _sort bool, def bool, tabSize int, useDesc bool, vscode bool) *ModuleBlock {
	// Pass tf file path to tfconfig
	module, _ := tfconfig.LoadModule(path)
	fullpath, _ := filepath.Abs(path)
	cwd, _ := os.Getwd()
	source, _ := filepath.Rel(cwd, fullpath)

	mb := ModuleBlock{
		path:    path,
		Name:    filepath.Base(fullpath),
		Source:  source,
		sort:    _sort,
		def:     def,
		useDesc: useDesc,
		vscode:  vscode,
		tabSize: tabSize,
	}

	mb.buildVariables(module.Variables)

	return &mb
}

func (mb *ModuleBlock) Do() (string, error) {
	if !tfconfig.IsModuleDir(mb.path) {
		return "", fmt.Errorf("given path does not contain tf files")
	}

	// The result from tfconfig is used to construct modBlock
	_template := tmpl
	if mb.vscode {
		_template = vsc_tmpl
	}

	block, err := template.New("block").Funcs(generateFuncMap()).Parse(_template)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}
	// Apply to template
	block.Execute(buffer, mb)

	return IndentByReplacingWords(buffer.String(), mb.tabSize), nil
}

//
func (mb *ModuleBlock) buildVariables(vars map[string]*tfconfig.Variable) {
	maxLen := getLongestKeySize(vars)

	for k, v := range vars {
		nm := k + strings.Repeat(" ", maxLen-len(k))
		tp := v.Type
		desc := v.Description
		if !mb.useDesc {
			desc = ""
		}
		df := GetDefaultValue(v, mb.def, tp)

		mb.Variables = append(
			mb.Variables,
			tfconfig.Variable{Type: tp, Name: nm, Description: desc, Default: df},
		)
	}

	if !mb.sort {
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
		"desc2comment":          desc2comment,
	}
}

// desc2comment adds comment `//` at the start of each lines
func desc2comment(desc string) string {
	var result string
	for _, l := range strings.Split(desc, "\n") {
		result = result + "\n" + REGEX_HAT.ReplaceAllString(l, INDENT+"// ")
	}

	return strings.TrimPrefix(result, "\n")
}
