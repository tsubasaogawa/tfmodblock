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

// constructModuleBlock constructs ModuleBlock from tfconfig.Variable.
func constructModuleBlock(mb *ModuleBlock, vars map[string]*tfconfig.Variable, _sort bool, def bool) {
	for k, v := range vars {
		r := regexp.MustCompile(`\w+`)

		tp := r.FindString(v.Type)
		if tp == "" {
			tp = v.Type
		}
		desc := v.Description

		var df interface{}
		// TODO: enable default for map and object as well
		if def && v.Default != nil && (tp == "string" || tp == "number" || tp == "bool" || strings.HasPrefix(tp, "list(")) {
			df = v.Default
		}
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

// printModuleBlock outputs Terraform module block.
func printModuleBlock(path string, _sort bool, def bool, vscode bool) (string, error) {
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
	constructModuleBlock(modBlock, module.Variables, _sort, def)

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

	return buffer.String(), nil
}

func main() {
	var (
		_sort  = flag.Bool("sort", true, "sort results")
		v      = flag.Bool("v", false, "tfmodblock version")
		def    = flag.Bool("default", true, "use default value if exists")
		vscode = flag.Bool("vscode", false, "VSCode extension mode")
	)
	flag.Parse()

	if *v {
		if version == "" {
			version = "v0.0.0"
		}
		fmt.Println(version)
		os.Exit(0)
	}

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	block, err := printModuleBlock(path, *_sort, *def, *vscode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(block)
}
