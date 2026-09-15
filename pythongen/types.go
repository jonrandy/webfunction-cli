package pythongen

import (
	"strconv"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
)

// typeResolver resolves a named object.<n> reference (in a given context,
// already baked in by the caller) to the name of a generated Python
// dataclass. Mirrors every other target's identical typeResolver.
type typeResolver func(refName string) string

// localTypes holds type names for an endpoint's own inline bare
// object/array attributes (as opposed to a named object.<n> reference,
// which resolves separately via typeResolver). Mirrors every other
// target's identical localTypes/localStructs.
type localTypes struct {
	object      string // dataclass name for a bare "object" return
	arrayOfItem string // item type name for a bare "array" return
}

// pythonType returns the concrete Python type-hint expression for t, with
// Optional[...] wrapping applied uniformly for both an Attribute's
// "nullable" flag and an Argument's "not required" flag - the same single-
// mechanism collapsing gogen/javagen/csharpgen do for their own languages'
// equivalent (Go pointers, C#'s "?"), since typing.Optional is Python's
// one real mechanism covering both "key may be absent" and "value may be
// null" at the type level.
func pythonType(t webfunction.Type, local localTypes, forceNullable bool, resolve typeResolver) string {
	base := pythonBaseType(t, local, resolve)
	nullable := forceNullable || t.HasBase("null")
	if nullable {
		return "Optional[" + base + "]"
	}
	return base
}

// pythonBaseType renders a Type's union as a Python type-hint expression.
//
// Unlike Go/Java/C# (none of which have a real type-union mechanism, so a
// multi-alternative union collapses to any/Object/object), Python's
// typing.Union is a genuine, natively-expressible union type - a real
// advantage this target gets to use rather than a gap to work around.
func pythonBaseType(t webfunction.Type, local localTypes, resolve typeResolver) string {
	var nonNull []webfunction.TypeAlt
	for _, alt := range t.Union {
		if alt.Base != "null" {
			nonNull = append(nonNull, alt)
		}
	}
	if len(nonNull) == 0 {
		return "Any"
	}
	if len(nonNull) == 1 {
		return pythonAltType(nonNull[0], local, resolve)
	}
	parts := make([]string, len(nonNull))
	seen := map[string]bool{}
	var unique []string
	for i, alt := range nonNull {
		parts[i] = pythonAltType(alt, local, resolve)
		if !seen[parts[i]] {
			seen[parts[i]] = true
			unique = append(unique, parts[i])
		}
	}
	if len(unique) == 1 {
		return unique[0]
	}
	return "Union[" + strings.Join(unique, ", ") + "]"
}

func pythonAltType(alt webfunction.TypeAlt, local localTypes, resolve typeResolver) string {
	switch alt.Base {
	case "string":
		return "str"
	case "boolean":
		return "bool"
	case "number":
		// Python has no real sized-integer or single-precision-float
		// types at all - ints are arbitrary precision and every float is
		// the equivalent of a 64-bit double - so u32/u64/i32/i64/timestamp
		// collapse to plain int (Python's int already covers every one of
		// these ranges exactly, unlike Java's u32/u64->long compromise
		// which can lose precision) and f32/f64/no-refinement collapse to
		// plain float. A real, flagged expressive-power gap in the
		// opposite direction from every other target: Python can't
		// distinguish int from float at the type-hint level any more
		// precisely than this, but never loses precision doing so.
		switch alt.Refinement {
		case "u32", "u64", "i32", "i64", "timestamp":
			return "int"
		default:
			return "float"
		}
	case "array":
		if alt.Of != nil {
			item := pythonType(*alt.Of, localTypes{}, false, resolve)
			return "List[" + item + "]"
		}
		if local.arrayOfItem != "" {
			return "List[" + local.arrayOfItem + "]"
		}
		return "List[Any]"
	case "object":
		if alt.IsObjectRef() && resolve != nil {
			return resolve(alt.Refinement)
		}
		if local.object != "" {
			return local.object
		}
		return "Dict[str, Any]"
	default: // "any", or anything unrecognized
		return "Any"
	}
}

// choiceLiteralKind picks which Python literal kind a generated
// typing.Literal[...] should use, based on the field's own declared type -
// not by inspecting the choices/values list itself (keeps a choices
// field's wire representation consistent with what its type would have
// been without choices at all). Mirrors javagen's/csharpgen's
// choiceWireType exactly, with one Python-specific addition: "float" is
// called out as its own kind (see buildLiteralField in dataclasses.go for
// why it matters - PEP 586 does not allow float values in a
// typing.Literal).
func choiceLiteralKind(t webfunction.Type) string {
	for _, alt := range t.Union {
		switch alt.Base {
		case "string":
			return "string"
		case "boolean":
			return "bool"
		case "number":
			switch alt.Refinement {
			case "u32", "u64", "i32", "i64", "timestamp":
				return "int"
			default:
				return "float"
			}
		}
	}
	return "string"
}

// pythonLiteral renders one raw choice/value as a Python literal
// expression matching kind (see choiceLiteralKind).
func pythonLiteral(v interface{}, kind string) string {
	switch val := v.(type) {
	case string:
		return strconv.Quote(val)
	case bool:
		// Python capitalizes its boolean literals, unlike JSON/Go/C#.
		if val {
			return "True"
		}
		return "False"
	case float64: // encoding/json decodes all JSON numbers as float64
		if kind == "int" {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'g', -1, 64)
	default:
		return "None"
	}
}

// refinementNote turns any dotted refinements in a type (e.g. "email" in
// "string.email") into a short parenthetical note to append to a
// docstring line, e.g. "(email format)". Mirrors every other target's
// identical refinementNote.
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

// withNotes appends any non-empty parenthetical notes to a docs string.
func withNotes(docs string, notes ...string) string {
	docs = strings.TrimSpace(docs)
	for _, n := range notes {
		if n == "" {
			continue
		}
		if docs == "" {
			docs = n
		} else {
			docs = docs + " " + n
		}
	}
	return docs
}