package envhelpers

import (
	"regexp"
	"strings"
)

var varRe = regexp.MustCompile(`\$\{(\w+)\}`)

// Expand replaces ${var} (and only that) in the string based on the mapping function.
// This signature is identical to the one in the Go standard library,
// but has a more restrictive scope with less ambiguity.
// The value inside ${} must match \w+ (see regexp.MatchString).
func Expand(s string, mapping func(string) string) string {
	if !strings.Contains(s, "${") {
		return s
	}
	return varRe.ReplaceAllStringFunc(s, func(varName string) string {
		return mapping(varName[2 : len(varName)-1])
	})
}
