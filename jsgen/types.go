package jsgen

import (
	"strings"

	"wfn/webfunction"
)

// objectResolver looks up (creating if needed) the typedef name for a
// referenced object.<name>, given the name that followed "object.".
type objectResolver func(refName string) string

// localTypedefs holds typedefs derived from an endpoint's own inline
// attributes (as opposed to a named object.<name> reference, which
// resolves separately via objectResolver). Per spec, an endpoint's
// `attributes` describes the shape of its `returns` when that's the bare
// `object` type. In practice, packages also use it to describe each
// item's shape when `returns` is a bare, untyped `array` - the spec's
// letter doesn't cover that case, but it's what the attributes are
// clearly there for, so it's honored here too.
type localTypedefs struct {
	object      string // shape of a bare "object" return
	arrayOfItem string // shape of each item in a bare "array" return
}

// jsdocAlt maps a single type alternative to its JSDoc equivalent.
// local supplies typedefs for the endpoint's own inline shapes (see
// localTypedefs); resolve is used for an actual object.<name> reference.
func jsdocAlt(alt webfunction.TypeAlt, local localTypedefs, resolve objectResolver) string {
	switch alt.Base {
	case "object":
		if alt.IsObjectRef() && resolve != nil {
			return resolve(alt.Refinement)
		}
		if local.object != "" {
			return local.object
		}
		return "Object"
	case "array":
		if alt.Of != nil {
			return "Array<" + jsdocType(*alt.Of, localTypedefs{}, false, resolve) + ">"
		}
		if local.arrayOfItem != "" {
			return "Array<" + local.arrayOfItem + ">"
		}
		return "Array<any>"
	case "string", "number", "boolean", "null":
		return alt.Base
	default: // "any", or anything unrecognized
		return "any"
	}
}

// jsdocType builds the JSDoc type string for an Argument or Attribute's
// Type, given typedefs for the endpoint's own inline shapes (see
// localTypedefs; pass the zero value when there are none, e.g. when
// rendering a typedef's own fields), whether the value may additionally
// be null on top of whatever the type itself says, and a resolver for any
// object.<name> references encountered.
func jsdocType(t webfunction.Type, local localTypedefs, forceNullable bool, resolve objectResolver) string {
	parts := make([]string, 0, len(t.Union)+1)
	for _, alt := range t.Union {
		parts = append(parts, jsdocAlt(alt, local, resolve))
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