package template

import (
	"fmt"
	"html/template"
	"strings"
)

// StringFuncMap return string func map.
func StringFuncMap() template.FuncMap {
	return template.FuncMap{
		"to_upper":   strings.ToUpper,
		"to_lower":   strings.ToLower,
		"to_string":  toString,
		"trim":       trim,
		"trim_left":  trimLeft,
		"trim_right": trimRight,
		"has_prefix": strings.HasPrefix,
	}
}

func toString(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func trim(cutset string, s string) string {
	return strings.Trim(s, cutset)
}

func trimLeft(cutset string, s string) string {
	return strings.TrimLeft(s, cutset)
}

func trimRight(cutset string, s string) string {
	return strings.TrimRight(s, cutset)
}

func hasPrefix(prefix string, s string) bool {
	return strings.HasPrefix(s, prefix)
}
