package phpgen

import (
	"fmt"
	"strings"
)

// camelCase folds dashes/underscores/case the same way jsgen's does, for
// consistency across targets - used for generated method names.
func camelCase(name string) string {
	words := splitWords(name)
	if len(words) == 0 {
		return ""
	}

	var b strings.Builder
	for i, w := range words {
		if i == 0 {
			b.WriteString(strings.ToLower(w))
			continue
		}
		b.WriteString(strings.ToUpper(w[:1]))
		if len(w) > 1 {
			b.WriteString(strings.ToLower(w[1:]))
		}
	}
	return b.String()
}

// pascalCase is camelCase with the first letter upper-cased, used for
// generated class names (Page wrapper classes, type alias names).
func pascalCase(name string) string {
	c := camelCase(name)
	if c == "" {
		return c
	}
	return strings.ToUpper(c[:1]) + c[1:]
}

// splitWords splits on dashes, underscores, and existing camelCase
// boundaries.
func splitWords(name string) []string {
	var words []string
	var cur strings.Builder

	flush := func() {
		if cur.Len() > 0 {
			words = append(words, cur.String())
			cur.Reset()
		}
	}

	runes := []rune(name)
	for i, r := range runes {
		switch {
		case r == '-' || r == '_' || r == ' ':
			flush()
		case i > 0 && isUpper(r) && !isUpper(runes[i-1]):
			flush()
			cur.WriteRune(r)
		default:
			cur.WriteRune(r)
		}
	}
	flush()

	return words
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// phpReserved is the set of identifiers that can't be used as a method
// name on the generated client class - either because PHP itself
// reserves them (can't be a function/method name at all) or because the
// wrapper class already uses them for its own passthrough surface.
var phpReserved = map[string]bool{
	// wrapper's own surface
	"call": true, "package": true, "client": true, "__construct": true,
	// PHP keywords that can't be used as a method name
	"list": true, "array": true, "echo": true, "print": true, "class": true,
	"function": true, "return": true, "new": true, "clone": true,
	"instanceof": true, "namespace": true, "use": true, "static": true,
	"const": true, "default": true, "switch": true, "case": true,
	"global": true, "isset": true, "unset": true, "empty": true,
	"require": true, "include": true, "throw": true, "try": true,
	"catch": true, "finally": true, "match": true, "fn": true,
}

// uniqueMethodName suffixes a name with an incrementing number if it
// collides with a reserved name or one already used by this client -
// mirrors jsgen's call -> call2 handling.
func uniqueMethodName(used map[string]bool, name string) string {
	candidate := name
	if phpReserved[candidate] {
		candidate += "2"
	}
	for i := 2; used[candidate]; i++ {
		candidate = fmt.Sprintf("%s%d", name, i)
	}
	used[candidate] = true
	return candidate
}