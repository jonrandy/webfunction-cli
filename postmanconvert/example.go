package postmanconvert

import (
	"github.com/webfunction-protocol/webfunction-go"
)

// exampleString returns a plausible placeholder for a recognized string
// refinement, or "" for an unrefined/unrecognized one.
func exampleString(refinement string) string {
	switch refinement {
	case "email":
		return "user@example.com"
	case "uuid":
		return "00000000-0000-0000-0000-000000000000"
	case "date":
		return "2026-01-01"
	case "time":
		return "00:00:00"
	case "datetime":
		return "2026-01-01T00:00:00Z"
	case "url", "uri":
		return "https://example.com"
	case "ipv4":
		return "127.0.0.1"
	case "ipv6":
		return "::1"
	case "hostname":
		return "example.com"
	case "phone":
		return "+15555550100"
	default:
		return ""
	}
}

// exampleValue builds a concrete placeholder value for a (possibly union)
// Type - the first non-null alternative is used, since a raw JSON example
// body can only ever show one concrete shape at a time; a Type that's
// only ever "null" (rare in practice) example-values as nil. visiting
// guards against infinite recursion on a self-referential object.<n> ref
// - unlike schema.go's schemaSet, this can't break the cycle with a
// $ref, so it breaks it by rendering an empty object the second time
// around instead.
func exampleValue(t webfunction.Type, ctx webfunction.ObjectContext, pkg *webfunction.Package, visiting map[string]bool) any {
	for _, alt := range t.Union {
		if alt.Base == "null" {
			continue
		}
		return exampleAlt(alt, ctx, pkg, visiting)
	}
	return nil
}

func exampleAlt(alt webfunction.TypeAlt, ctx webfunction.ObjectContext, pkg *webfunction.Package, visiting map[string]bool) any {
	switch alt.Base {
	case "object":
		if alt.IsObjectRef() {
			return exampleObjectRef(alt.Refinement, ctx, pkg, visiting)
		}
		return newOrderedMap()
	case "array":
		if alt.Of != nil {
			return []any{exampleValue(*alt.Of, ctx, pkg, visiting)}
		}
		return []any{}
	case "string":
		return exampleString(alt.Refinement)
	case "number":
		return 0
	case "boolean":
		return false
	default: // "any", "null" (shouldn't reach "null" here - see exampleValue)
		return nil
	}
}

// exampleObjectRef resolves an object.<n> ref into a concrete example
// object, using the ref'd object's own arguments/attributes (per ctx) as
// its fields - an argument's choices or attribute's values, when
// declared, take priority over the field's own type-derived example
// (mirrors the ref'd object's real constraints better than a generic
// placeholder would).
func exampleObjectRef(name string, ctx webfunction.ObjectContext, pkg *webfunction.Package, visiting map[string]bool) any {
	key := ctxKey(ctx) + ":" + name
	if visiting[key] {
		return newOrderedMap() // cycle guard: stop recursing, render empty
	}
	visiting[key] = true
	defer delete(visiting, key)

	obj := pkg.ObjectInContext(name, ctx)
	if obj == nil {
		return newOrderedMap()
	}

	result := newOrderedMap()
	if ctx == webfunction.ArgumentContext {
		for _, arg := range obj.Arguments {
			result.set(arg.Name, argumentExample(arg, pkg, visiting))
		}
	} else {
		for _, attr := range obj.Attributes {
			result.set(attr.Name, attributeExample(attr, pkg, visiting))
		}
	}
	return result
}

func argumentExample(arg webfunction.Argument, pkg *webfunction.Package, visiting map[string]bool) any {
	if len(arg.Choices) > 0 {
		return arg.Choices[0]
	}
	return exampleValue(arg.Type, webfunction.ArgumentContext, pkg, visiting)
}

func attributeExample(attr webfunction.Attribute, pkg *webfunction.Package, visiting map[string]bool) any {
	if len(attr.Values) > 0 {
		return attr.Values[0]
	}
	return exampleValue(attr.Type, webfunction.AttributeContext, pkg, visiting)
}

// requestBodyExample builds the top-level request body example (always
// ArgumentContext) for an endpoint's arguments list.
func requestBodyExample(args []webfunction.Argument, pkg *webfunction.Package) *orderedMap {
	visiting := map[string]bool{}
	body := newOrderedMap()
	for _, arg := range args {
		body.set(arg.Name, argumentExample(arg, pkg, visiting))
	}
	return body
}

func ctxKey(ctx webfunction.ObjectContext) string {
	if ctx == webfunction.ArgumentContext {
		return "arg"
	}
	return "attr"
}