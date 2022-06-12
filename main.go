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
	"text/template"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

// Variable is Terraform variable object.
type Variable struct {
	Type string
	Name string
}

// ModuleBlock is an output text consisted of variables.
type ModuleBlock struct {
	Name      string
	Variables []Variable
}

const VERSION string = "0.0.4"

var (
	//go:embed module_block.tmpl
	tmpl string
	//go:embed module_block_vscode.tmpl
	vsc_tmpl string
)

// applyModuleBlock
func applyModuleBlock(mb *ModuleBlock, vars map[string]*tfconfig.Variable) {
	for k, v := range vars {
		r := regexp.MustCompile(`\w+`)
		tp := r.FindString(v.Type)
		if tp == "" {
			tp = v.Type
		}
		mb.Variables = append(mb.Variables, Variable{Type: tp, Name: k})
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

// generateModuleBlock generates Terraform module block.
func generateModuleBlock(path string, vscode bool) (string, error) {
	module, _ := tfconfig.LoadModule(path)

	modBlock := new(ModuleBlock)
	fullpath, _ := filepath.Abs(path)
	modBlock.Name = filepath.Base(fullpath)
	applyModuleBlock(modBlock, module.Variables)

	_template := tmpl
	if vscode {
		_template = vsc_tmpl
	}

	block, err := template.New("block").Funcs(generateFuncMap()).Parse(_template)
	if err != nil {
		return "", err
	}
	buffer := &bytes.Buffer{}
	block.Execute(buffer, modBlock)

	return buffer.String(), nil
}

func main() {
	var (
		v      = flag.Bool("v", false, "tfmodblock version")
		vscode = flag.Bool("vscode", false, "VSCode extension mode")
	)
	flag.Parse()

	if *v {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	block, err := generateModuleBlock(path, *vscode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(block)
}
