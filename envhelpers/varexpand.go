package envhelpers

import (
	"regexp"
	"strings"
)

var (
	varUnquouteRe = regexp.MustCompile(`[\"']\$\{(\w+)@U\}[\"']`)
	varRe         = regexp.MustCompile(`\$\{(\w+)\}`)
)

// Expand replaces ${var} (and only that) in the string based on the mapping function.
// This signature is identical to the one in the Go standard library,
// but has a more restrictive scope with less ambiguity.
// The value inside ${} must match \w+ (see regexp.MatchString).
//
// A special form of environment variable syntax is supported to allow for removing surrounding
// quoutes (both single and double), useful to setting numeric values in TOML/JSON config files etc.
// Use the "@U" suffix to signal unquoting, e.g. `"${myvar@U}"` or `'${myvar@U}'`. Note that when using the "@U" suffix,
// the variable expression needs to be surrounded by quoutes.
func Expand(s string, mapping func(string) string) string {
	if !strings.Contains(s, "${") {
		return s
	}

	firstPass := varUnquouteRe.ReplaceAllStringFunc(s, func(match string) string {
		// Remove quoutes and "@U" suffix before passing to mapping function.
		return mapping(match[3 : len(match)-4])
	})

	return varRe.ReplaceAllStringFunc(firstPass, func(varName string) string {
		return mapping(varName[2 : len(varName)-1])
	})
}
