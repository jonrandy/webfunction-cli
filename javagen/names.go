package javagen

import (
	"fmt"
	"strings"
	"unicode"
)

// camelCase/pascalCase/splitWords mirror gogen's exactly (duplicated
// rather than imported - matches the existing convention: phpgen/gogen
// already duplicate these rather than importing jsgen's).
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
// every generated Java type name (records, enums, the Client class
// itself), matching Java's convention that type names start upper-case.
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

// javaKeywords is the set of reserved words that can't be used as a Java
// identifier at all (keywords, plus the three reserved literals true/
// false/null, which are lexically reserved even though they aren't
// technically keywords).
var javaKeywords = map[string]bool{
	"abstract": true, "assert": true, "boolean": true, "break": true, "byte": true,
	"case": true, "catch": true, "char": true, "class": true, "const": true,
	"continue": true, "default": true, "do": true, "double": true, "else": true,
	"enum": true, "extends": true, "final": true, "finally": true, "float": true,
	"for": true, "goto": true, "if": true, "implements": true, "import": true,
	"instanceof": true, "int": true, "interface": true, "long": true, "native": true,
	"new": true, "package": true, "private": true, "protected": true, "public": true,
	"return": true, "short": true, "static": true, "strictfp": true, "super": true,
	"switch": true, "synchronized": true, "this": true, "throw": true, "throws": true,
	"transient": true, "try": true, "void": true, "volatile": true, "while": true,
	"var": true, "yield": true, "record": true, "sealed": true, "permits": true,
	"true": true, "false": true, "null": true,
}

// clientReserved is what a generated per-endpoint method name can't
// collide with: the Client class's own escape-hatch method names, plus
// every Java keyword. Mirrors gogen's goReserved.
var clientReserved = func() map[string]bool {
	m := map[string]bool{
		"call": true, "getpackage": true, "setbearerauth": true,
		"setversion": true, "setpipeline": true, "newclient": true,
	}
	for k := range javaKeywords {
		m[k] = true
	}
	return m
}()

// exportedFieldName renders name as a valid Java field/record-component
// identifier: camelCase, prefixed with "_" if the result would otherwise
// start with a digit (Java identifiers can't start with a digit - e.g.
// "2fa_enabled" would otherwise produce the invalid "2faEnabled") or
// collide with a reserved word (a record component literally named
// "class" is invalid, since Records implicitly can't shadow certain
// things and reserved words are outright illegal as identifiers
// regardless).
func exportedFieldName(name string) string {
	c := camelCase(name)
	if c == "" {
		return "field"
	}
	if unicode.IsDigit(rune(c[0])) {
		c = "_" + c
	}
	if javaKeywords[c] {
		c += "_"
	}
	return c
}

// exportedTypeName is exportedFieldName's PascalCase counterpart, used
// for generated record/enum type names.
func exportedTypeName(name string) string {
	p := pascalCase(name)
	if p == "" {
		return "Type"
	}
	if unicode.IsDigit(rune(p[0])) {
		p = "_" + p
	}
	return p
}

// uniqueMethodName suffixes a name with an incrementing number if it
// collides with a reserved name or one already used on the Client class -
// mirrors gogen's uniqueMethodName / jsgen's call -> call2 handling.
func uniqueMethodName(used map[string]bool, name string) string {
	candidate := name
	if clientReserved[strings.ToLower(candidate)] {
		candidate += "2"
	}
	for i := 2; used[candidate]; i++ {
		candidate = fmt.Sprintf("%s%d", name, i)
	}
	used[candidate] = true
	return candidate
}

// javaPackageName derives a valid, idiomatic Java package name from the
// shared --namespace flag (see cmd/codegen.go). Unlike phpgen (which
// uses --namespace largely verbatim as a PHP namespace) or gogen (which
// collapses it to one lower-case identifier), Java package names are
// naturally dot-separated - so a namespace like "com.example.api" maps
// directly, while the PHP-style default "WebFunctionClient" (no dots)
// becomes the single lower-cased segment "webfunctionclient". Each
// dot-separated segment is sanitized independently: non-letter/digit
// characters dropped, lower-cased (Java package segments are
// conventionally all lower-case), and re-checked against Java's
// identifier rules (can't start with a digit, can't be a keyword).
// Falls back to "webfunctionclient" for an empty or fully-invalid
// namespace.
func javaPackageName(namespace string) string {
	rawSegments := strings.Split(namespace, ".")
	var segments []string
	for _, seg := range rawSegments {
		var b strings.Builder
		for _, r := range seg {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				b.WriteRune(unicode.ToLower(r))
			}
		}
		s := b.String()
		if s == "" {
			continue
		}
		if unicode.IsDigit(rune(s[0])) {
			s = "_" + s
		}
		if javaKeywords[s] {
			s += "_"
		}
		segments = append(segments, s)
	}
	if len(segments) == 0 {
		return "webfunctionclient"
	}
	return strings.Join(segments, ".")
}
