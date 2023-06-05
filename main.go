// main package is the main of tfmodblock.
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version string

	_sort, def, v, desc, vscode bool
	tabsize                     int
)

func init() {
	flag.BoolVar(&_sort, "sort", true, "sort results")
	flag.BoolVar(&def, "default", true, "use default value if exists")
	flag.BoolVar(&v, "v", false, "tfmodblock version")
	flag.BoolVar(&desc, "description", true, "insert a variable description")
	flag.BoolVar(&vscode, "vscode", false, "VSCode extension mode")
	flag.IntVar(&tabsize, "tabsize", 4, "tab size for indent")
}

func main() {
	flag.Parse()

	if v {
		if version == "" {
			version = "v0.0.0"
		}
		fmt.Println(version)
		os.Exit(0)
	}

	if tabsize < 0 {
		fmt.Fprintln(os.Stderr, "tabsize must be >= 0")
		os.Exit(1)
	}

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	block, err := GenerateModuleBlockString(path, _sort, def, tabsize, desc, vscode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(block)
}
