package main

import (
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

// GetDefaultValue obtains the default value from `default` field.
func GetDefaultValue(_var *tfconfig.Variable, def bool, tp string) string {
	var df []byte

	if !def || _var.Default == nil {
		return GetDefaultValueByType(tp)
	}

	df, _ = json.Marshal(_var.Default)

	return string(df)
}

// GetDefaultValueByType returns the default value considering type.
func GetDefaultValueByType(_type string) string {
	r := regexp.MustCompile(`\w+`)
	tp := r.FindString(_type)

	switch tp {
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
}
