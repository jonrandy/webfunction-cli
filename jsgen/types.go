package jsgen

import (
	"fmt"
	"math"
	"strconv"
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
//
// paginated/pageItemType are a separate case: webfunction-js (confirmed
// from its real source at github.com/jonrandy/webfunction-js) auto-wraps
// a canonical {previous, page, next} response in its own `Page` class
// rather than handing back the raw envelope - matching what the
// PHP/Ruby clients already do. So a paginated endpoint's return type is
// `Page` itself, not a typedef of the envelope shape; pageItemType is
// only the JSDoc type string describing what `.page` holds (e.g.
// "Array<PersonAttributes>"), for the @returns description - not
// something the type system itself narrows, since nothing confirms
// whether the real Page class is declared generic (e.g. `Page<T>`). Not
// inventing that pending being able to check the actual class source.
type localTypedefs struct {
	object       string // shape of a bare "object" return
	arrayOfItem  string // shape of each item in a bare "array" return
	paginated    bool   // true if this endpoint's return is Page instead
	pageItemType string // JSDoc type of what `.page` holds, for docs only
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
			// Array item types have no choices concept of their own -
			// choices/values apply to the field as a whole, not to each
			// array element individually.
			return "Array<" + jsdocType(*alt.Of, localTypedefs{}, false, resolve, nil) + ">"
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
// be null on top of whatever the type itself says, a resolver for any
// object.<name> references encountered, and any closed set of legal
// values (an Attribute's "values" or Argument's "choices"; pass nil when
// there are none, e.g. for an endpoint's own Returns type, which has no
// choices concept).
//
// When choices is non-empty, it takes over the type entirely: the result
// is a JSDoc literal union of the exact values (e.g. ("active"|"pending"))
// rather than the bare base type (e.g. "string") - since the base type
// alone doesn't add anything a literal union doesn't already convey, and
// the literal union is what actually lets an editor or tsc catch a typo
// or an out-of-set value at the call site.
func jsdocType(t webfunction.Type, local localTypedefs, forceNullable bool, resolve objectResolver, choices []interface{}) string {
	var parts []string
	if lits := literalUnion(choices); len(lits) > 0 {
		parts = lits
	} else {
		parts = make([]string, 0, len(t.Union))
		for _, alt := range t.Union {
			parts = append(parts, jsdocAlt(alt, local, resolve))
		}
	}

	if forceNullable || t.HasBase("null") {
		parts = append(parts, "null")
	}

	parts = dedupe(parts)

	if len(parts) == 1 {
		return parts[0]
	}
	return "(" + strings.Join(parts, "|") + ")"
}

// literalUnion renders a closed set of legal values (an Attribute's
// "values" or an Argument's "choices", per spec) as JSDoc literal type
// parts: a string becomes a quoted literal ("active"), a bool or number
// becomes its own bare literal (true, 3), and an explicit null entry
// becomes the "null" keyword. Returns nil if choices is empty.
func literalUnion(choices []interface{}) []string {
	if len(choices) == 0 {
		return nil
	}
	parts := make([]string, 0, len(choices))
	for _, c := range choices {
		switch v := c.(type) {
		case nil:
			parts = append(parts, "null")
		case string:
			parts = append(parts, strconv.Quote(v))
		case bool:
			parts = append(parts, strconv.FormatBool(v))
		case float64: // encoding/json decodes all JSON numbers as float64
			parts = append(parts, formatNumberLiteral(v))
		default:
			// Shouldn't happen for values decoded from JSON, but don't
			// drop the entry silently.
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return parts
}

// formatNumberLiteral renders a float64 as a JS numeric literal without a
// spurious ".0" on whole numbers (choices are very often small integers,
// e.g. status codes or ratings).
func formatNumberLiteral(f float64) string {
	if !math.IsInf(f, 0) && f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
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