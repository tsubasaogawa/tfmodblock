// main package is the main of tfmodblock.
package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

// ModuleBlock includes output values consisted of variables.
type ModuleBlock struct {
	Name      string
	Source    string
	Variables []tfconfig.Variable
}

var (
	version string
	//go:embed module_block.tmpl
	tmpl string
	//go:embed module_block_vscode.tmpl
	vsc_tmpl string
)

const INDENT = "  "

// getDefaultValue
func getDefaultValue(_var *tfconfig.Variable, def bool, tp string) interface{} {
	var df interface{}

	// TODO: enable default for map and object as well
	if def && _var.Default != nil && (tp == "string" || tp == "number" || tp == "bool" || strings.HasPrefix(tp, "list(")) {
		df = _var.Default
	}

	return df
}

// indentByReplacingWords replaces string with space.
func indentByReplacingWords(s string, size int, w string) string {
	re := regexp.MustCompile("^" + w)
	var result string
	for _, l := range strings.Split(s, "\n") {
		result = fmt.Sprintf("%s%s\n", result, re.ReplaceAllString(l, strings.Repeat(" ", size)))
	}
	return result
}

// constructModuleBlock constructs ModuleBlock from tfconfig.Variable.
func constructModuleBlock(mb *ModuleBlock, vars map[string]*tfconfig.Variable, _sort bool, def bool, tabSize int) {
	r := regexp.MustCompile(`\w+`)

	for k, v := range vars {
		tp := r.FindString(v.Type)
		if tp == "" {
			tp = v.Type
		}
		desc := v.Description
		df := getDefaultValue(v, def, tp)

		mb.Variables = append(
			mb.Variables,
			tfconfig.Variable{Type: tp, Name: k, Description: desc, Default: df},
		)
	}
	if !_sort {
		return
	}
	sort.Slice(mb.Variables, func(i, j int) bool { return mb.Variables[i].Name < mb.Variables[j].Name })
}

// generateFuncMap
func generateFuncMap() template.FuncMap {
	return template.FuncMap{
		"convertTypeToLiteral": func(_type string) string {
			switch _type {
			case "string":
				return "\"\""
			case "number":
				return "0"
			case "list", "set", "tuple":
				return "[]"
			case "bool":
				return "true/false"
			case "object", "map":
				return "{}"
			default:
				return "null"
			}
		},
	}
}

// generateModuleBlockString returns Terraform module block string.
func generateModuleBlockString(path string, _sort bool, def bool, tabSize int, vscode bool) (string, error) {
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
	constructModuleBlock(modBlock, module.Variables, _sort, def, tabSize)

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

	return indentByReplacingWords(buffer.String(), tabSize, INDENT), nil
}

func main() {
	var (
		_sort   = flag.Bool("sort", true, "sort results")
		def     = flag.Bool("default", true, "use default value if exists")
		tabSize = flag.Int("tabsize", 4, "tab size for indent")
		v       = flag.Bool("v", false, "tfmodblock version")
		vscode  = flag.Bool("vscode", false, "VSCode extension mode")
	)
	flag.Parse()

	if *v {
		if version == "" {
			version = "v0.0.0"
		}
		fmt.Println(version)
		os.Exit(0)
	}

	if *tabSize < 0 {
		fmt.Fprintln(os.Stderr, "tabsize must be >= 0")
		os.Exit(1)
	}

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	block, err := generateModuleBlockString(path, *_sort, *def, *tabSize, *vscode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(block)
}
