package jsgen

import (
	"strings"

	"wfn/webfunction"
)

// objectResolver looks up (creating if needed) the typedef name for a
// referenced object.<name>, given the name that followed "object.".
type objectResolver func(refName string) string

// jsdocAlt maps a single type alternative to its JSDoc equivalent.
// localObjectTypedef is used for a bare "object" alternative that isn't a
// named reference (i.e. the endpoint's own inline return shape); resolve
// is used for an actual object.<name> reference.
func jsdocAlt(alt webfunction.TypeAlt, localObjectTypedef string, resolve objectResolver) string {
	switch alt.Base {
	case "object":
		if alt.IsObjectRef() && resolve != nil {
			return resolve(alt.Refinement)
		}
		if localObjectTypedef != "" {
			return localObjectTypedef
		}
		return "Object"
	case "array":
		if alt.Of != nil {
			return "Array<" + jsdocType(*alt.Of, "", false, resolve) + ">"
		}
		return "Array<any>"
	case "string", "number", "boolean", "null":
		return alt.Base
	default: // "any", or anything unrecognized
		return "any"
	}
}

// jsdocType builds the JSDoc type string for an Argument or Attribute's
// Type, given the typedef name to use for a bare "object" entry (or "" if
// there's no known shape for it), whether the value may additionally be
// null on top of whatever the type itself says, and a resolver for any
// object.<name> references encountered.
func jsdocType(t webfunction.Type, localObjectTypedef string, forceNullable bool, resolve objectResolver) string {
	parts := make([]string, 0, len(t.Union)+1)
	for _, alt := range t.Union {
		parts = append(parts, jsdocAlt(alt, localObjectTypedef, resolve))
	}

	if forceNullable && !t.HasBase("null") {
		parts = append(parts, "null")
	}

	parts = dedupe(parts)

	if len(parts) == 1 {
		return parts[0]
	}
	return "(" + strings.Join(parts, "|") + ")"
}

// refinementNote turns any dotted refinements in a type (e.g. "email" in
// "string.email") into a short parenthetical note to append to a docs
// description, e.g. "(email format)". Object refinements (object.<name>,
// a type reference rather than a format constraint) are excluded. Returns
// "" if there's nothing worth noting.
func refinementNote(t webfunction.Type) string {
	var notes []string
	for _, alt := range t.Union {
		if alt.Refinement != "" && alt.Base != "object" {
			notes = append(notes, alt.Refinement)
		}
	}
	if len(notes) == 0 {
		return ""
	}
	return "(" + strings.Join(notes, ", ") + " format)"
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

func boolStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}