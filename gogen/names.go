package gogen

import (
	"fmt"
	"strings"
	"unicode"
)

// camelCase/pascalCase/splitWords mirror jsgen's/phpgen's exactly
// (duplicated rather than imported - matches the existing convention:
// phpgen already duplicates these rather than importing jsgen's).
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

// pascalCase is camelCase with the first letter upper-cased - used for
// every generated Go identifier here (struct names, exported field
// names, method names), since Go requires exported identifiers to start
// with an upper-case letter.
func pascalCase(name string) string {
	c := camelCase(name)
	if c == "" {
		return c
	}
	return strings.ToUpper(c[:1]) + c[1:]
}

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

// exportedName renders name as a valid, exported Go identifier:
// PascalCase, prefixed with "X" if the result would otherwise start with
// a digit (Go identifiers can't start with a digit - a real-world field
// name like "2fa_enabled" would otherwise produce the invalid "2faEnabled").
func exportedName(name string) string {
	p := pascalCase(name)
	if p == "" {
		return "Field"
	}
	if unicode.IsDigit(rune(p[0])) {
		return "X" + p
	}
	return p
}

// goReserved is the set of identifiers that can't be used as a method
// name on the generated Client - either because the wrapper's own
// passthrough surface already uses them, or because they're Go keywords.
// Method names are always PascalCase (exported), so an exact lower-case
// keyword collision can't happen there after casing - but packageName
// (see below) also consults this list, where a collision genuinely can
// happen since it's derived from arbitrary user input.
var goReserved = map[string]bool{
	// wrapper's own passthrough surface (lower-cased for comparison)
	"call": true, "package": true, "setbearerauth": true,
	"setversion": true, "setpipeline": true, "newclient": true, "client": true,
	// Go keywords
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true,
	"for": true, "func": true, "go": true, "goto": true, "if": true,
	"import": true, "interface": true, "map": true, "range": true,
	"return": true, "select": true, "struct": true, "switch": true,
	"type": true, "var": true,
}

// uniqueMethodName suffixes a name with an incrementing number if it
// collides with a reserved name or one already used - mirrors jsgen's
// call -> call2 handling.
func uniqueMethodName(used map[string]bool, name string) string {
	candidate := name
	if goReserved[strings.ToLower(candidate)] {
		candidate += "2"
	}
	for i := 2; used[candidate]; i++ {
		candidate = fmt.Sprintf("%s%d", name, i)
	}
	used[candidate] = true
	return candidate
}

// packageName derives a valid, idiomatic Go package name from the shared
// --namespace flag (see cmd/codegen.go). Go package names are
// conventionally short, lower-case, and contain no underscores or mixed
// case - which the shared default ("WebFunctionClient", styled for PHP/
// C#/Java-style class namespacing) doesn't match. Unlike phpgen (which
// uses --namespace largely verbatim as a PHP namespace), gogen sanitizes
// it into Go's own convention rather than requiring a separate flag just
// for this one target.
func packageName(namespace string) string {
	var b strings.Builder
	for _, r := range namespace {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
		// anything else (spaces, backslashes, underscores, punctuation)
		// is dropped entirely, matching Go's convention of avoiding
		// underscores in package names.
	}
	name := b.String()
	if name == "" || unicode.IsDigit(rune(name[0])) || goReserved[name] {
		name = "webfunctionclient"
	}
	return name
}
