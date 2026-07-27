package jsgen

import (
	"sort"
	"strings"

	"wfn/webfunction"
)

// jsdocBaseType maps a single wire base type to its JSDoc equivalent.
// "object" and "array" are handled by the caller, since "object" may need
// to resolve to a typedef name and "array" has no documented element type
// in the wire spec to draw on.
func jsdocBaseType(t string) string {
	switch t {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	case "null":
		return "null"
	case "object":
		return "Object"
	case "array":
		return "Array"
	default:
		return "*"
	}
}

// jsdocType builds the JSDoc type string for an Argument or Attribute's
// JSONType, given the typedef name to use for an "object" entry (or "" if
// there's no known shape for it) and whether the value may be null
// (attributes only - the "nullable" flag).
func jsdocType(types webfunction.JSONType, objectTypedef string, nullable bool) string {
	parts := make([]string, 0, len(types)+1)
	for _, t := range types {
		if t == "object" && objectTypedef != "" {
			parts = append(parts, objectTypedef)
			continue
		}
		parts = append(parts, jsdocBaseType(t))
	}

	if nullable && !contains(types, "null") {
		parts = append(parts, "null")
	}

	parts = dedupe(parts)

	if len(parts) == 1 {
		return parts[0]
	}
	return "(" + strings.Join(parts, "|") + ")"
}

// hintNote turns a hint list into a short parenthetical note to append to a
// docs description, e.g. "(email format)". Returns "" if there are no
// hints worth noting.
func hintNote(hints []string) string {
	if len(hints) == 0 {
		return ""
	}
	return "(" + strings.Join(hints, ", ") + " format)"
}

func contains(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}

func dedupe(ss []string) []string {
	seen := make(map[string]bool, len(ss))
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// attributeSignature builds a canonical string describing the shape of an
// attribute list, order-independent, so identical shapes appearing on
// different endpoints can share one typedef.
func attributeSignature(attrs []webfunction.Attribute) string {
	parts := make([]string, len(attrs))
	for i, a := range attrs {
		parts[i] = a.Name + ":" + a.Type.String() + ":" + boolStr(a.Nullable())
	}
	sort.Strings(parts)
	return strings.Join(parts, "|")
}

func boolStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}