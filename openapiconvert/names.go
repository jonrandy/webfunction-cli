package openapiconvert

import (
	"strings"
	"unicode"
)

// toOperationId camelCases an endpoint's kebab/underscore/space-separated
// name for use as an OpenAPI operationId, e.g. "find-user" -> "findUser".
// Ported directly from the reference webfunction-to-openapi JS converter.
func toOperationId(name string) string {
	parts := splitWords(name)
	var b strings.Builder
	for i, part := range parts {
		if part == "" {
			continue
		}
		lower := strings.ToLower(part)
		if i == 0 {
			b.WriteString(lower)
			continue
		}
		b.WriteString(strings.ToUpper(lower[:1]))
		b.WriteString(lower[1:])
	}
	return b.String()
}

// schemaKeyName PascalCases an arbitrary webfunction object name into a
// valid OpenAPI component schema key (OpenAPI restricts these to
// [a-zA-Z0-9._-], but PascalCase identifiers are the readable convention
// used everywhere else in this project's generated output).
func schemaKeyName(raw string) string {
	parts := splitWords(raw)
	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		if keepsCasing(part) {
			b.WriteString(part)
			continue
		}
		lower := strings.ToLower(part)
		b.WriteString(strings.ToUpper(lower[:1]))
		b.WriteString(lower[1:])
	}
	name := b.String()
	if name == "" {
		name = "Object"
	}
	if unicode.IsDigit(rune(name[0])) {
		name = "_" + name
	}
	return name
}

// keepsCasing reports whether a word is already mixed case with a leading
// capital (ListOfUser, Page), as the names generated for generic
// instantiations are. Such words are kept as written; all other words are
// normalised to a leading capital followed by lower case.
func keepsCasing(word string) bool {
	runes := []rune(word)
	if !unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes[1:] {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// splitWords splits on any run of characters that aren't letters or
// digits (hyphens, underscores, spaces, dots, ...).
func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}