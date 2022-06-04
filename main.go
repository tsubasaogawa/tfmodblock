// main package is the main of tfmodblock.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"text/template"

	"github.com/gobuffalo/packr/v2"
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

const VERSION string = "0.0.0"
const TMPL_FILE string = "module_block.tmpl"

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
func generateModuleBlock(path string) (string, error) {
	module, _ := tfconfig.LoadModule(path)

	modBlock := new(ModuleBlock)
	fullpath, _ := filepath.Abs(path)
	modBlock.Name = filepath.Base(fullpath)
	applyModuleBlock(modBlock, module.Variables)

	box := packr.New("root", ".")
	s, err := box.FindString(TMPL_FILE)
	if err != nil {
		return "", err
	}

	block, err := template.New("tmpl").Funcs(generateFuncMap()).Parse(s)
	if err != nil {
		return "", err
	}
	buffer := &bytes.Buffer{}
	block.Execute(buffer, modBlock)

	return buffer.String(), nil
}

func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	block, err := generateModuleBlock(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(block)
}
