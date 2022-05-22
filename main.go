/*
memo
- object の扱い
  - 展開するためには object を自前でパースする必要がある？
- unknown な型の扱い
*/

// main package is the main of tfvergen.
package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

const VERSION string = "0.0.1"

func getRequiredVersion() (string, error) {
	module, _ := tfconfig.LoadModule(".")
	/*
		if len(module.RequiredCore) < 1 {
			return "", fmt.Errorf("tfvergen %s\nThere is no required version.", VERSION)
		}
	*/
	variables := module.Variables
	for k, v := range variables {
		fmt.Printf("%s: %s\n", k, v.Type)
	}

	/*
		r := regexp.MustCompile(`\d+\.\d+\.\d+`)
		version := r.FindString(constraint)
		if version == "" {
			fmt.Fprintf(os.Stderr, "Fail to extract version from %s\n", constraint)
		}
	*/

	return "", nil
}

func main() {
	version, err := getRequiredVersion()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", version)
}
