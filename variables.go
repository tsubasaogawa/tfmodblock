package main

import (
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

// GetDefaultValue obtains the default value from `default` field.
func GetDefaultValue(_var *tfconfig.Variable, def bool, tp string) interface{} {
	var df interface{}

	// TODO: enable default for map and object as well
	if def && _var.Default != nil && (tp == "string" || tp == "number" || tp == "bool" || strings.HasPrefix(tp, "list(")) {
		df = _var.Default
	}

	return df
}
