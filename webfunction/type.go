package webfunction

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Type represents a Web Function type specification - what appears in an
// endpoint's `returns`, or an argument's or attribute's `type`.
// See https://webfunction.org/package#types.
//
// The wire grammar is recursive: a type is a single base/refined type, or a
// union of types written as a JSON array, where an array *entry* that is
// itself an array denotes an array type (whose own entries are, in turn, a
// union of the types its items may take). Type models this directly rather
// than flattening it to a list of strings, so array element types and
// object.<name> references survive parsing.
type Type struct {
	// Union holds every alternative this type may be. A plain (non-union)
	// type has exactly one entry.
	Union []TypeAlt
}

// TypeAlt is one alternative within a Type's union.
type TypeAlt struct {
	// Base is one of: "object", "array", "string", "number", "boolean",
	// "null", or "any".
	Base string

	// Refinement is the dotted suffix narrowing Base, e.g. "email" for
	// "string.email", or the referenced object's name for "object.<name>".
	// Empty when there's no refinement. The "any" base MUST NOT carry one.
	Refinement string

	// Of is the element type, present only when Base == "array" and the
	// wire form was a nested array (e.g. [["string"]] means "array of
	// string"). Nil for a bare "array" entry (array of any).
	Of *Type
}

// IsObjectRef reports whether this alternative is a reference to a named
// object definition (an "object.<name>" refinement).
func (a TypeAlt) IsObjectRef() bool {
	return a.Base == "object" && a.Refinement != ""
}

func (a TypeAlt) String() string {
	s := a.Base
	if a.Refinement != "" {
		s += "." + a.Refinement
	}
	if a.Base == "array" && a.Of != nil {
		s += "<" + a.Of.String() + ">"
	}
	return s
}

// HasBase reports whether any alternative in the union has the given base
// type, e.g. t.HasBase("object").
func (t Type) HasBase(base string) bool {
	for _, alt := range t.Union {
		if alt.Base == base {
			return true
		}
	}
	return false
}

// HasBareArray reports whether the union includes an "array" alternative
// with no nested element type (i.e. array of any, at the wire level) -
// distinct from HasBase("array"), which also matches a typed array.
func (t Type) HasBareArray() bool {
	for _, alt := range t.Union {
		if alt.Base == "array" && alt.Of == nil {
			return true
		}
	}
	return false
}

// ObjectNames returns the names of every object definition referenced
// anywhere within this type, including inside array element types -
// mirroring webfunction-php's Type::objects().
func (t Type) ObjectNames() []string {
	var names []string
	for _, alt := range t.Union {
		if alt.IsObjectRef() {
			names = append(names, alt.Refinement)
		}
		if alt.Of != nil {
			names = append(names, alt.Of.ObjectNames()...)
		}
	}
	return names
}

func (t Type) String() string {
	parts := make([]string, len(t.Union))
	for i, alt := range t.Union {
		parts[i] = alt.String()
	}
	return strings.Join(parts, "|")
}

// UnmarshalJSON parses the recursive type grammar: a JSON null (treated as
// "any" - only actually expected for choices/values elements, not type
// fields themselves, but handled defensively), a bare string, or a JSON
// array whose entries are either strings or nested arrays.
func (t *Type) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	parsed, err := parseTypeValue(raw)
	if err != nil {
		return err
	}
	*t = parsed
	return nil
}

func parseTypeValue(raw interface{}) (Type, error) {
	switch v := raw.(type) {
	case nil:
		return Type{Union: []TypeAlt{{Base: "any"}}}, nil
	case string:
		return Type{Union: []TypeAlt{parseTypeString(v)}}, nil
	case []interface{}:
		if len(v) == 0 {
			// Not permitted per spec ("Arrays MUST NOT be empty at any
			// depth"), but fail soft to "any" rather than erroring
			// outright on a non-conformant package.
			return Type{Union: []TypeAlt{{Base: "any"}}}, nil
		}
		alts := make([]TypeAlt, 0, len(v))
		for _, elem := range v {
			switch e := elem.(type) {
			case string:
				alts = append(alts, parseTypeString(e))
			case []interface{}:
				inner, err := parseTypeValue(e)
				if err != nil {
					return Type{}, err
				}
				alts = append(alts, TypeAlt{Base: "array", Of: &inner})
			default:
				return Type{}, fmt.Errorf("unexpected type entry %T", elem)
			}
		}
		return Type{Union: alts}, nil
	default:
		return Type{}, fmt.Errorf("unexpected type value %T", raw)
	}
}

func parseTypeString(s string) TypeAlt {
	if s == "array" {
		return TypeAlt{Base: "array"}
	}
	if i := strings.IndexByte(s, '.'); i >= 0 {
		return TypeAlt{Base: s[:i], Refinement: s[i+1:]}
	}
	return TypeAlt{Base: s}
}