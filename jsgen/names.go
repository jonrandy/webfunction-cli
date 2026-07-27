package jsgen

import "strings"

// camelCase folds dashes/underscores/case the same way webfunction-js
// itself normalizes endpoint names for method lookup, e.g. "list-items" and
// "list_items" both become "listItems".
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
// typedef names.
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
			// lower-to-upper boundary, e.g. "findUser" -> "find", "User"
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

// jsReserved is the set of identifiers that can't be used as method names on
// the wrapper object without being quoted, either because they're real JS
// reserved words or because the wrapper itself already uses them for the
// passthrough client surface.
var jsReserved = map[string]bool{
	// wrapper's own passthrough properties
	"package": true,
	"call":    true,
	// a conservative subset of JS reserved words worth avoiding as bare
	// method names
	"class": true, "function": true, "return": true, "delete": true,
	"new": true, "typeof": true, "in": true, "of": true, "this": true,
	"super": true, "import": true, "export": true, "default": true,
	"const": true, "let": true, "var": true,
}