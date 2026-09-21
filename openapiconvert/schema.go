package openapiconvert

import (
	"strconv"

	"github.com/webfunction-protocol/webfunction-go"
)

// refinementSchema returns the partial Schema fragment (format/minimum/
// description) for a recognized string/number refinement, or nil if the
// refinement isn't one this package knows how to render specially - the
// base type alone (already set by the caller) is all that's emitted then.
// Ported directly from the reference webfunction-to-openapi JS
// converter's HINT_SCHEMA table, adapted from a separate "hints" array to
// webfunction-go's dotted-refinement Type model.
func refinementSchema(refinement string) *Schema {
	switch refinement {
	case "u32":
		return &Schema{Format: "int32", Minimum: floatPtr(0)}
	case "u64":
		return &Schema{Format: "int64", Minimum: floatPtr(0)}
	case "i32":
		return &Schema{Format: "int32"}
	case "i64":
		return &Schema{Format: "int64"}
	case "f32":
		return &Schema{Format: "float"}
	case "f64":
		return &Schema{Format: "double"}
	case "timestamp":
		return &Schema{Format: "int64", Description: "Unix timestamp (seconds since the epoch)."}
	case "date":
		return &Schema{Format: "date"}
	case "time":
		return &Schema{Format: "time"}
	case "datetime":
		return &Schema{Format: "date-time"}
	case "uuid":
		return &Schema{Format: "uuid"}
	case "base64":
		return &Schema{Format: "byte"}
	case "email":
		return &Schema{Format: "email"}
	case "phone":
		return &Schema{Description: "Phone number in E.164 format."}
	case "url":
		return &Schema{Format: "uri"}
	case "uri":
		return &Schema{Format: "uri"}
	case "ipv4":
		return &Schema{Format: "ipv4"}
	case "ipv6":
		return &Schema{Format: "ipv6"}
	case "hostname":
		return &Schema{Format: "hostname"}
	default:
		return nil
	}
}

// openObjectSchema is the fallback shape for an object reference that
// doesn't resolve to a real object.<n> definition (or has no members in
// the requested context) - left open rather than forbidding properties.
func openObjectSchema() *Schema {
	return &Schema{Type: "object", AdditionalProperties: boolPtr(true)}
}

// wrapNullable folds "or null" onto an existing schema, preferring to
// merge into the type array where that's semantically sound (plain base
// types) and falling back to a oneOf wrapper where a bare type-array
// merge would be wrong (a $ref or an existing oneOf - merging "null" into
// a sibling "type" keyword next to $ref would, under JSON Schema 2020-12,
// require values to satisfy an intersection of the ref schema AND
// type:null, which no non-null value referencing an object schema ever
// does; a oneOf branch is what actually expresses "either the real thing
// or null").
func wrapNullable(s *Schema) *Schema {
	if s.Ref != "" {
		return &Schema{OneOf: []*Schema{s, {Type: "null"}}}
	}
	if s.OneOf != nil {
		s.OneOf = append(s.OneOf, &Schema{Type: "null"})
		return s
	}
	switch t := s.Type.(type) {
	case nil:
		s.Type = "null"
	case string:
		s.Type = []string{t, "null"}
	case []string:
		for _, existing := range t {
			if existing == "null" {
				return s
			}
		}
		s.Type = append(append([]string{}, t...), "null")
	}
	return s
}

// applyEnum attaches an argument's choices or attribute's values to a
// built schema, on the items schema for an array type and directly on
// the schema itself otherwise. Ported from the reference converter's
// fieldSchema; like the reference, a multi-alternative union (oneOf) with
// declared choices/values is a combination the wire format doesn't
// actually produce in practice and isn't specially handled.
func applyEnum(s *Schema, values []any) {
	if len(values) == 0 {
		return
	}
	if t, ok := s.Type.(string); ok && t == "array" && s.Items != nil {
		s.Items.Enum = values
		return
	}
	s.Enum = values
}

// altSchema builds the schema for a single non-null TypeAlt. ctx says
// which member set to resolve an object.<n> ref against; fallback, when
// non-empty, is used for a *bare* (unrefined) "object" or "array" base -
// this only ever applies at the top level of an endpoint's own Returns
// type, using that endpoint's own declared Attributes as the shape (real
// packages use Attributes this way even though the spec's letter only
// defines it for a bare "object" return - carried over from this
// project's jsgen, which found the same pattern needed for a bare
// "array" return too).
func altSchema(alt webfunction.TypeAlt, ctx webfunction.ObjectContext, fallback []webfunction.Attribute, resolver *schemaSet) *Schema {
	switch alt.Base {
	case "object":
		if alt.IsObjectRef() {
			return resolver.resolveRef(alt.Refinement, ctx)
		}
		if len(fallback) > 0 {
			return attributesToSchema(fallback, resolver)
		}
		return openObjectSchema()
	case "array":
		if alt.Of != nil {
			return &Schema{Type: "array", Items: buildSchema(*alt.Of, ctx, nil, resolver)}
		}
		if len(fallback) > 0 {
			return &Schema{Type: "array", Items: attributesToSchema(fallback, resolver)}
		}
		return &Schema{Type: "array", Items: &Schema{}}
	case "any":
		return &Schema{}
	case "null":
		// Only reached if a Type's *entire* union is just "null" (see
		// buildSchema, which otherwise strips "null" out before calling
		// altSchema and folds it back on via wrapNullable instead).
		return &Schema{Type: "null"}
	default: // "string", "number", "boolean"
		s := &Schema{Type: alt.Base}
		if frag := refinementSchema(alt.Refinement); frag != nil {
			s.Format = frag.Format
			s.Minimum = frag.Minimum
			s.Description = frag.Description
		}
		return s
	}
}

// buildSchema builds the schema for a full (possibly union) Type. See
// altSchema for ctx/fallback.
func buildSchema(t webfunction.Type, ctx webfunction.ObjectContext, fallback []webfunction.Attribute, resolver *schemaSet) *Schema {
	var nonNull []webfunction.TypeAlt
	nullable := false
	for _, alt := range t.Union {
		if alt.Base == "null" {
			nullable = true
			continue
		}
		nonNull = append(nonNull, alt)
	}

	var schema *Schema
	switch len(nonNull) {
	case 0:
		schema = &Schema{}
	case 1:
		schema = altSchema(nonNull[0], ctx, fallback, resolver)
	default:
		oneOf := make([]*Schema, len(nonNull))
		for i, alt := range nonNull {
			oneOf[i] = altSchema(alt, ctx, fallback, resolver)
		}
		schema = &Schema{OneOf: oneOf}
	}

	if nullable {
		schema = wrapNullable(schema)
	}
	return schema
}

// argumentsToSchema builds the object schema (properties + required) for
// an arguments list - used both for an endpoint's request body and for an
// object.<n> ref resolved in ArgumentContext.
func argumentsToSchema(args []webfunction.Argument, resolver *schemaSet) *Schema {
	properties := map[string]*Schema{}
	var required []string
	for _, arg := range args {
		s := buildSchema(arg.Type, webfunction.ArgumentContext, nil, resolver)
		applyEnum(s, arg.Choices)
		if arg.Docs != "" {
			s.Description = arg.Docs
		}
		properties[arg.Name] = s
		if arg.Required() {
			required = append(required, arg.Name)
		}
	}
	schema := &Schema{Type: "object", Properties: properties, AdditionalProperties: boolPtr(false)}
	if len(required) > 0 {
		schema.Required = required
	}
	return schema
}

// attributesToSchema builds the object schema (properties only - webfunction
// has no "required" concept for returned attributes, only "nullable") for
// an attributes list - used both for an endpoint's response body and for
// an object.<n> ref resolved in AttributeContext.
func attributesToSchema(attrs []webfunction.Attribute, resolver *schemaSet) *Schema {
	properties := map[string]*Schema{}
	for _, attr := range attrs {
		s := buildSchema(attr.Type, webfunction.AttributeContext, nil, resolver)
		applyEnum(s, attr.Values)
		if attr.Nullable() {
			s = wrapNullable(s)
		}
		if attr.Docs != "" {
			s.Description = attr.Docs
		}
		properties[attr.Name] = s
	}
	schema := &Schema{Type: "object", Properties: properties}
	if len(properties) == 0 {
		schema.AdditionalProperties = boolPtr(true)
	}
	return schema
}

// schemaSet resolves object.<n> refs into named, reusable, cycle-safe
// entries under components.schemas - the same pattern this project's
// codegen targets use for typedefs/structs/aliases (jsgen's typedefSet,
// phpgen's aliasSet, gogen's structSet), adapted to OpenAPI schemas.
//
// The same object name can carry different members depending on context
// (Arguments vs Attributes - see webfunction.Package.ObjectInContext), so
// a name used in both contexts gets two distinct schema entries; the
// ArgumentContext one is suffixed "Input" to keep the two apart and
// readable (an object used only as a return shape keeps the plain name).
type schemaSet struct {
	pkg *webfunction.Package

	// keyFor maps a "ctx:name" resolution key to the component schema
	// name already assigned to it - present (even with an empty/
	// placeholder body in s.schemas) as soon as resolution starts, which
	// is what makes a self-referential object safe to resolve.
	keyFor map[string]string
	used   map[string]bool
	// schemas holds the actual schema body per component name, filled in
	// once resolution of that object completes.
	schemas map[string]*Schema
}

func newSchemaSet(pkg *webfunction.Package) *schemaSet {
	return &schemaSet{
		pkg:     pkg,
		keyFor:  map[string]string{},
		used:    map[string]bool{},
		schemas: map[string]*Schema{},
	}
}

func ctxKey(ctx webfunction.ObjectContext) string {
	if ctx == webfunction.ArgumentContext {
		return "arg"
	}
	return "attr"
}

// resolveRef returns a fresh $ref-only Schema pointing at the given
// object's component schema, registering (and, on first resolution,
// building) that component schema as a side effect. Always returns a new
// *Schema value even for a repeat resolution of the same object/context,
// since callers may attach per-use fields (description, nullability) to
// what they get back and must never mutate a shared instance.
func (s *schemaSet) resolveRef(name string, ctx webfunction.ObjectContext) *Schema {
	key := ctxKey(ctx) + ":" + name
	if compName, ok := s.keyFor[key]; ok {
		return &Schema{Ref: "#/components/schemas/" + compName}
	}

	compName := s.assignName(name, ctx)
	s.keyFor[key] = compName
	s.schemas[compName] = &Schema{} // placeholder: makes a cycle back to this key resolve immediately above, instead of recursing forever

	obj := s.pkg.ObjectInContext(name, ctx)
	var built *Schema
	switch {
	case obj == nil:
		built = openObjectSchema()
	case ctx == webfunction.ArgumentContext:
		built = argumentsToSchema(obj.Arguments, s)
	default:
		built = attributesToSchema(obj.Attributes, s)
	}
	s.schemas[compName] = built

	return &Schema{Ref: "#/components/schemas/" + compName}
}

func (s *schemaSet) assignName(name string, ctx webfunction.ObjectContext) string {
	base := schemaKeyName(name)
	if ctx == webfunction.ArgumentContext {
		base += "Input"
	}
	candidate := base
	for i := 2; s.used[candidate]; i++ {
		candidate = base + strconv.Itoa(i)
	}
	s.used[candidate] = true
	return candidate
}