// main package is the main of tfmodblock.
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version string
)

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

	block, err := GenerateModuleBlockString(path, *_sort, *def, *tabSize, *vscode)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(block)
}
