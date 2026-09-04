package rubygen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
)

// field is a name/type/optional/docs tuple built from either an
// endpoint's Arguments or Attributes, or an object.<n>'s own
// Arguments/Attributes when resolving a named ref - mirrors every other
// target's field model exactly.
//
// optional means the key may be absent; nullable means the value itself
// may additionally be nil. An Attribute's `nullable` flag means both at
// once per spec (mirrors jsgen/phpgen).
type field struct {
	name     string
	jsonType webfunction.Type
	optional bool
	nullable bool
	docs     string
	choices  []interface{}
}

func attributeFields(attrs []webfunction.Attribute) []field {
	fields := make([]field, len(attrs))
	for i, a := range attrs {
		fields[i] = field{
			name: a.Name, jsonType: a.Type,
			optional: a.Nullable(), nullable: a.Nullable(),
			docs: a.Docs, choices: a.Values,
		}
	}
	return fields
}

func argumentFields(args []webfunction.Argument) []field {
	fields := make([]field, len(args))
	for i, a := range args {
		fields[i] = field{
			name: a.Name, jsonType: a.Type,
			optional: !a.Required(),
			docs:     a.Docs, choices: a.Choices,
		}
	}
	return fields
}

// ---- argument side: real Ruby keyword args, symbol-keyed record types ----
//
// Confirmed from the real gem: client.list_items(a: "b") is genuinely
// Ruby keyword-argument syntax (method_missing declares no keyword
// params of its own, so Ruby packs it into a positional Hash - but the
// CALL SITE syntax is keyword-style either way, and RBS's keyword-param
// grammar (?name: Type) is the natural, precise match for that - not a
// single positional Hash-shaped param). Nested argument objects (e.g.
// filters: { first_name: "Joe" }) are themselves symbol-keyed Ruby hash
// literals, so RBS record types (also symbol-keyed by definition) are
// exactly right for those.

// argLocalShapes holds inline (anonymous, unnamed) record-type strings
// for an endpoint's own bare object/array argument - built from that
// argument's own nested Attributes-equivalent... endpoints don't carry
// per-argument attributes in this spec, so this only ever applies via
// object.<n> refs in practice; kept for symmetry with the return side.
type argLocalShapes struct {
	object      string
	arrayOfItem string
}

type argResolver func(refName string) string

// rbsArgType renders the RBS type for an argument-side Type: full
// per-field record-type precision, since arguments really are
// symbol-keyed Ruby hashes.
func rbsArgType(t webfunction.Type, local argLocalShapes, forceNullable bool, resolve argResolver, choices []interface{}) string {
	var parts []string
	if lits := literalUnion(choices); len(lits) > 0 {
		parts = lits
	} else {
		parts = make([]string, 0, len(t.Union))
		for _, alt := range t.Union {
			parts = append(parts, rbsArgAlt(alt, local, resolve))
		}
	}

	if forceNullable || t.HasBase("null") {
		parts = append(parts, "nil")
	}
	parts = dedupe(parts)
	return strings.Join(parts, " | ")
}

func rbsArgAlt(alt webfunction.TypeAlt, local argLocalShapes, resolve argResolver) string {
	switch alt.Base {
	case "string":
		return "String"
	case "number":
		switch alt.Refinement {
		case "u32", "u64", "i32", "i64", "timestamp":
			return "Integer"
		case "f32", "f64":
			return "Float"
		default:
			return "Integer | Float"
		}
	case "boolean":
		return "bool"
	case "null":
		return "nil"
	case "any":
		return "untyped"
	case "array":
		if alt.Of != nil {
			return "Array[" + rbsArgType(*alt.Of, argLocalShapes{}, false, resolve, nil) + "]"
		}
		if local.arrayOfItem != "" {
			return "Array[" + local.arrayOfItem + "]"
		}
		return "Array[untyped]"
	case "object":
		if alt.Refinement != "" {
			return resolve(alt.Refinement)
		}
		if local.object != "" {
			return local.object
		}
		return "Hash[Symbol, untyped]"
	default:
		return "untyped"
	}
}

// argAliasSet builds named `type foo_args = { ... }` aliases for
// object.<n> refs referenced from the argument side, memoized per
// name so a self-referential (or mutually referential) object resolves
// cleanly - registered before its own shape is built, mirroring every
// other target's cycle-safe alias resolution.
type argAliasSet struct {
	pkg     *webfunction.Package
	names   map[string]bool
	byKey   map[string]string
	ordered []aliasDef
}

type aliasDef struct {
	name    string
	rbsType string
}

func newArgAliasSet(pkg *webfunction.Package) *argAliasSet {
	return &argAliasSet{pkg: pkg, names: map[string]bool{}, byKey: map[string]string{}}
}

func (s *argAliasSet) resolve(name string) string {
	key := name
	if existing, ok := s.byKey[key]; ok {
		return existing
	}

	obj := s.pkg.Object(name)
	var fields []field
	if obj != nil {
		fields = argumentFields(obj.Arguments)
	}

	if len(fields) == 0 {
		s.byKey[key] = "Hash[Symbol, untyped]"
		return "Hash[Symbol, untyped]"
	}

	aliasName := uniqueName(s.names, snakeCase(name)+"_args")
	s.byKey[key] = aliasName

	shape := s.renderRecord(fields)
	s.ordered = append(s.ordered, aliasDef{name: aliasName, rbsType: shape})
	return aliasName
}

func (s *argAliasSet) renderRecord(fields []field) string {
	if len(fields) == 0 {
		return "{ }"
	}
	resolve := func(refName string) string { return s.resolve(refName) }
	parts := make([]string, len(fields))
	for i, f := range fields {
		// NOTE: unlike a method type's keyword-parameter list (where
		// `?name: Type` for "may be absent" works fine, confirmed via
		// rbs validate), a record-type LITERAL's `?key:` optional-field
		// marker is a newer RBS syntax addition not supported by this
		// RBS version (2.8.2, confirmed via rbs validate - "unexpected
		// record key token, token=`?`"). Rather than gamble on which
		// RBS version Jon's own environment has, optional nested-object
		// fields fold into nilable here instead (same fallback every
		// RBS-record-type generator needs pre-3.x) - loses the "absent
		// vs present-but-nil" distinction only for nested object.<n>
		// argument fields; top-level endpoint arguments are unaffected,
		// since those render as real keyword params, not record keys.
		typ := rbsArgType(f.jsonType, argLocalShapes{}, f.nullable || f.optional, resolve, f.choices)
		parts[i] = recordKey(f.name) + " " + typ
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

// ---- return side: real JSON.parse output, string-keyed, no RBS record ----
//
// Confirmed from the real gem's Request#execute: JSON.parse(body) with no
// symbolize_names, so returned objects have STRING keys - RBS record
// types (symbol-keyed by definition) would be flatly wrong here. Instead,
// a structured return gets a named `interface _FooAttributes` with
// overloaded `[]` signatures per key (confirmed viable with Jon) - real
// per-field precision on Hash#[] access, at the cost of only [] being
// typed (not #dig/#each/#keys/etc, deliberately out of scope for v1).

type returnLocalShapes struct {
	object      string // interface name, if this endpoint's own bare "object" return has attributes
	arrayOfItem string // interface name (or scalar type), if this endpoint's own bare "array" return has attributes
}

type returnResolver func(refName string) string

func rbsReturnType(t webfunction.Type, local returnLocalShapes, forceNullable bool, resolve returnResolver, choices []interface{}) string {
	var parts []string
	if lits := literalUnion(choices); len(lits) > 0 {
		parts = lits
	} else {
		parts = make([]string, 0, len(t.Union))
		for _, alt := range t.Union {
			parts = append(parts, rbsReturnAlt(alt, local, resolve))
		}
	}

	if forceNullable || t.HasBase("null") {
		parts = append(parts, "nil")
	}
	parts = dedupe(parts)
	return strings.Join(parts, " | ")
}

func rbsReturnAlt(alt webfunction.TypeAlt, local returnLocalShapes, resolve returnResolver) string {
	switch alt.Base {
	case "string":
		return "String"
	case "number":
		switch alt.Refinement {
		case "u32", "u64", "i32", "i64", "timestamp":
			return "Integer"
		case "f32", "f64":
			return "Float"
		default:
			return "Integer | Float"
		}
	case "boolean":
		return "bool"
	case "null":
		return "nil"
	case "any":
		return "untyped"
	case "array":
		if alt.Of != nil {
			return "Array[" + rbsReturnType(*alt.Of, returnLocalShapes{}, false, resolve, nil) + "]"
		}
		if local.arrayOfItem != "" {
			return "Array[" + local.arrayOfItem + "]"
		}
		return "Array[untyped]"
	case "object":
		if alt.Refinement != "" {
			return resolve(alt.Refinement)
		}
		if local.object != "" {
			return local.object
		}
		return "Hash[String, untyped]"
	default:
		return "untyped"
	}
}

// ifaceDef is one generated `interface _Name ... end` block: one `def []`
// overload per field, each returning that field's own rbsReturnType
// (nil-unioned already baked in per-field, since Hash#[] can't
// distinguish "absent" from "present-but-nil" either way).
type ifaceDef struct {
	name   string
	fields []ifaceField
}

type ifaceField struct {
	name    string
	rbsType string
}

// ifaceSet builds named return-side interfaces: one per object.<n> ref
// actually referenced from a return/attribute position (memoized by
// name, cycle-safe the same way argAliasSet is), plus one per endpoint
// whose own bare object/array return has inline Attributes (synthesized
// name, resolveLocal).
type ifaceSet struct {
	pkg     *webfunction.Package
	names   map[string]bool
	byKey   map[string]string
	ordered []ifaceDef
}

func newIfaceSet(pkg *webfunction.Package) *ifaceSet {
	return &ifaceSet{pkg: pkg, names: map[string]bool{}, byKey: map[string]string{}}
}

func (s *ifaceSet) resolve(name string) string {
	key := name
	if existing, ok := s.byKey[key]; ok {
		return existing
	}

	obj := s.pkg.Object(name)
	var fields []field
	if obj != nil {
		fields = attributeFields(obj.Attributes)
	}

	if len(fields) == 0 {
		s.byKey[key] = "Hash[String, untyped]"
		return "Hash[String, untyped]"
	}

	ifaceName := uniqueName(s.names, "_"+pascalCase(name)+"Attributes")
	s.byKey[key] = ifaceName
	s.registerFields(ifaceName, fields)
	return ifaceName
}

// resolveLocal synthesizes a fresh, uniquely-named interface for an
// endpoint's own inline (non-object.<n>) return attributes - e.g.
// _ListItemsResult - not memoized by any reusable key, since each
// endpoint's own bare-return shape is inherently one-off.
func (s *ifaceSet) resolveLocal(endpointName string, fields []field) string {
	ifaceName := uniqueName(s.names, "_"+pascalCase(endpointName)+"Result")
	s.registerFields(ifaceName, fields)
	return ifaceName
}

func (s *ifaceSet) registerFields(ifaceName string, fields []field) {
	resolve := func(refName string) string { return s.resolve(refName) }
	ifields := make([]ifaceField, len(fields))
	for i, f := range fields {
		typ := rbsReturnType(f.jsonType, returnLocalShapes{}, f.nullable, resolve, f.choices)
		ifields[i] = ifaceField{name: f.name, rbsType: typ}
	}
	s.ordered = append(s.ordered, ifaceDef{name: ifaceName, fields: ifields})
}

func (s *ifaceSet) hasInterfaces() bool {
	return len(s.ordered) > 0
}

// ---- shared helpers ----

// literalUnion renders a closed set of legal values (Attribute#values or
// Argument#choices) as RBS literal-type parts. RBS supports literal
// types for String/Symbol/Integer/true/false, but NOT Float - a non-
// integer numeric choice widens to "Float" rather than a false-precision
// literal.
func literalUnion(choices []interface{}) []string {
	if len(choices) == 0 {
		return nil
	}
	parts := make([]string, 0, len(choices))
	for _, c := range choices {
		switch v := c.(type) {
		case nil:
			parts = append(parts, "nil")
		case string:
			parts = append(parts, rbsStringLiteral(v))
		case bool:
			parts = append(parts, strconv.FormatBool(v))
		case float64: // encoding/json decodes all JSON numbers as float64
			if v == float64(int64(v)) {
				parts = append(parts, strconv.FormatInt(int64(v), 10))
			} else {
				parts = append(parts, "Float")
			}
		default:
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return parts
}

func rbsStringLiteral(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
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

func pascalCase(name string) string {
	words := splitWords(name)
	var b strings.Builder
	for _, w := range words {
		if w == "" {
			continue
		}
		b.WriteString(strings.ToUpper(w[:1]))
		if len(w) > 1 {
			b.WriteString(strings.ToLower(w[1:]))
		}
	}
	return b.String()
}