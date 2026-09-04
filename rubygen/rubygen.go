// Package rubygen generates an RBS type-signature file for a webfunction
// package, targeting the real reference gem (github.com/robinclart/
// web_function, "web_function" on RubyGems - the public mirror of the
// private github.com/webfunction-protocol/webfunction-ruby) and its
// already dynamic-dispatch-based WebFunction::Client.
//
// Design decisions (confirmed with Jon):
//  1. Pure `.rbs` signature file only - no generated wrapper class/
//     subclass, unlike every other target. The real Client uses
//     method_missing-based dynamic dispatch (confirmed directly from
//     client.rb: `client.list_items(a: "b")` already works with zero
//     generated code), and RBS lets you declare method signatures for a
//     class regardless of whether the runtime implementation uses `def`
//     or `method_missing` - Steep/rbs type-checks against the declared
//     signature, not the dispatch mechanism. So there's nothing for
//     rubygen to generate beyond the signatures themselves.
//  2. Endpoint method names are dash-to-underscore ONLY, matching
//     Client#initialize's own `endpoints.to_h { |e| [e.gsub("-", "_")
//     .to_sym, e] }` exactly (confirmed from client.rb) - no camelCase
//     folding, unlike jsgen. A mismatched name would type-check a call
//     the real dispatcher actually raises NoMethodError on.
//  3. Arguments get full per-field RBS record-type / keyword-argument
//     precision, since `client.list_items(a: "b")` really is Ruby
//     keyword-call syntax against a symbol-keyed hash (confirmed:
//     method_missing declares no keyword params of its own, so this
//     is genuine keyword sugar, not a positional-hash-shaped param).
//     RBS keyword-param grammar (`?name: Type`) is the natural match.
//  4. Structured RETURN values do NOT use RBS record types. Confirmed
//     from Request#execute: `JSON.parse(body)` with no
//     `symbolize_names: true`, so every returned object has STRING
//     keys - RBS record types are specifically symbol-keyed and would
//     be wrong here. Instead, a structured return gets a named
//     `interface _FooAttributes` with one overloaded `def []` per
//     field (real per-key precision on Hash#[] access, at the cost of
//     only [] itself being typed - #dig/#each/#keys/etc are out of
//     scope for v1, deliberately not modeled).
//  5. Numeric refinements: u32/u64/i32/i64/timestamp -> Integer,
//     f32/f64 -> Float, unrefined -> "Integer | Float" (RBS has no
//     single "number" type, same fundamental ambiguity phpgen hit).
//  6. Choices/values render as real RBS literal types (String/Integer/
//     true/false/nil) - RBS has no Float literal type, so a non-integer
//     numeric choice widens to plain "Float" rather than a false-
//     precision literal.
//  7. Pagination: the real gem always wraps a paginated response in a
//     WebFunction::Page instance (confirmed from page.rb - `page`,
//     `next?`, `previous?`, `next_page`, `previous_page`, Enumerable).
//     Rather than a per-endpoint wrapper class (impossible anyway under
//     decision #1 - no generated code at all), rubygen declares Page
//     itself as RBS-generic (`class Page[out Item]`) ONCE, globally, in
//     the same sig file - purely a type-level construct (Ruby has no
//     runtime generics, and RBS doesn't require one; this is the
//     standard RBS pattern for typing Ruby's own generic-ish
//     collections). Each paginated endpoint's return type is then
//     simply `Page[SpecificItemType]`.
//  8. No `--namespace` flag - doesn't apply to this target at all,
//     unlike php/go/java/csharp, since there's no generated class to
//     namespace.
//
// Scope note: this generates the package-specific dynamically-dispatched
// endpoint methods, plus the fixed instance-level surface every Client
// has regardless of package (call/package/bearer_auth=/version=/
// pipeline=/methods/nil?, all confirmed directly from client.rb) so the
// sig file is usable standalone. Class-level factory methods
// (from_package_endpoint/from_url/from_package) and the Pipeline class
// itself are NOT modeled (declared `untyped` where referenced) - these
// are gem-wide, not package-specific, and arguably belong in the gem's
// own sig/ directory (which currently only has a one-line VERSION stub)
// rather than being duplicated per generated package file.
package rubygen

import (
	"fmt"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
)

// Generate builds an RBS signature file for pkg, targeting the real
// web_function gem. sourceURL is recorded in the header comment only
// (there's no generated constructor to bake it into, per decision #1).
func Generate(pkg *webfunction.Package, sourceURL string) (string, error) {
	argAliases := newArgAliasSet(pkg)
	returnIfaces := newIfaceSet(pkg)
	endpoints := visibleEndpoints(pkg)

	infos := make(map[string]endpointInfo, len(endpoints))
	paginated := false
	for _, ep := range endpoints {
		info := buildEndpointInfo(argAliases, returnIfaces, ep)
		infos[ep.Name] = info
		if info.paginated {
			paginated = true
		}
	}

	var b strings.Builder
	name := pkg.Name
	if name == "" {
		name = "(unnamed package)"
	}
	b.WriteString("# Generated RBS signatures for the " + name + " webfunction package.\n")
	b.WriteString("# Source: " + sourceURL + "\n")
	b.WriteString("#\n")
	b.WriteString("# Targets the real web_function gem (github.com/robinclart/web_function).\n")
	b.WriteString("# Load alongside the gem via `rbs` or Steep - no generated Ruby code, this\n")
	b.WriteString("# file declares signatures for WebFunction::Client's existing dynamically-\n")
	b.WriteString("# dispatched methods (method_missing-based - the real gem already answers\n")
	b.WriteString("# to these calls, this file only adds the types).\n\n")

	b.WriteString("module WebFunction\n")

	writeFixedSurface(&b)

	if paginated {
		writePageClass(&b)
	}

	writeClient(&b, endpoints, infos)

	for _, a := range argAliases.ordered {
		b.WriteString("  type " + a.name + " = " + a.rbsType + "\n")
	}
	for _, iface := range returnIfaces.ordered {
		writeInterface(&b, iface)
	}

	b.WriteString("end\n")

	return b.String(), nil
}

// writeFixedSurface declares the gem-wide Error hierarchy - confirmed
// directly from lib/web_function.rb: `Error < StandardError` with
// `code`/`details` readers, and four subclasses (UnresolvedPromiseError,
// UnexpectedStatusCodeError, JsonParseError, BadRequestError), same
// [code, message, details] triple every other target's error model
// matches.
func writeFixedSurface(b *strings.Builder) {
	b.WriteString("  class Error < StandardError\n")
	b.WriteString("    def code: () -> String\n")
	b.WriteString("    def details: () -> untyped\n")
	b.WriteString("  end\n\n")
	for _, sub := range []string{"UnresolvedPromiseError", "UnexpectedStatusCodeError", "JsonParseError", "BadRequestError"} {
		b.WriteString("  class " + sub + " < Error\n  end\n\n")
	}
}

// writePageClass declares WebFunction::Page as RBS-generic - see
// decision #7. Only emitted when the package has at least one paginated
// endpoint.
func writePageClass(b *strings.Builder) {
	b.WriteString("  class Page[out Item]\n")
	b.WriteString("    include Enumerable[Item]\n\n")
	b.WriteString("    def page: () -> Array[Item]\n")
	b.WriteString("    def next?: () -> bool\n")
	b.WriteString("    def previous?: () -> bool\n")
	b.WriteString("    def next_page: () -> Page[Item]?\n")
	b.WriteString("    def previous_page: () -> Page[Item]?\n")
	b.WriteString("  end\n\n")
}

func writeClient(b *strings.Builder, endpoints []webfunction.Endpoint, infos map[string]endpointInfo) {
	b.WriteString("  class Client < BasicObject\n")
	b.WriteString("    # Escape hatch for calling an endpoint this file doesn't declare a\n")
	b.WriteString("    # typed method for.\n")
	b.WriteString("    def call: (String endpoint_name, ?Hash[Symbol, untyped] args) -> untyped\n\n")
	b.WriteString("    def package: () -> untyped\n")
	b.WriteString("    def bearer_auth=: (String?) -> void\n")
	b.WriteString("    def version=: (String?) -> void\n")
	b.WriteString("    def pipeline=: (untyped) -> void\n")
	b.WriteString("    def methods: () -> ::Array[Symbol]\n")
	b.WriteString("    def nil?: () -> bool\n\n")

	used := map[string]bool{}
	for _, ep := range endpoints {
		methodName := uniqueMethodName(used, rubyMethodName(ep.Name))
		writeMethod(b, ep, methodName, infos[ep.Name])
	}

	b.WriteString("  end\n\n")
}

type endpointInfo struct {
	params     []paramInfo
	returnType string
	paginated  bool
	errorCodes []webfunction.ErrorDef
}

type paramInfo struct {
	name     string
	rbsType  string
	optional bool
}

func buildEndpointInfo(argAliases *argAliasSet, returnIfaces *ifaceSet, ep webfunction.Endpoint) endpointInfo {
	argFields := argumentFields(ep.Arguments)
	resolveArg := func(refName string) string { return argAliases.resolve(refName) }
	params := make([]paramInfo, len(argFields))
	for i, f := range argFields {
		params[i] = paramInfo{
			name:     f.name,
			rbsType:  rbsArgType(f.jsonType, argLocalShapes{}, f.nullable, resolveArg, f.choices),
			optional: f.optional,
		}
	}

	if ep.HasFlag("paginated") {
		itemType := "untyped"
		resolveReturn := func(refName string) string { return returnIfaces.resolve(refName) }
		for _, a := range ep.Attributes {
			if a.Name != "page" {
				continue
			}
			for _, alt := range a.Type.Union {
				if alt.Base == "array" && alt.Of != nil {
					itemType = rbsReturnType(*alt.Of, returnLocalShapes{}, false, resolveReturn, nil)
				}
			}
		}
		return endpointInfo{
			params: params, paginated: true,
			returnType: "Page[" + itemType + "]",
			errorCodes: ep.Errors,
		}
	}

	local := forEndpointReturn(returnIfaces, ep)
	resolveReturn := func(refName string) string { return returnIfaces.resolve(refName) }
	returnType := rbsReturnType(ep.Returns, local, false, resolveReturn, nil)

	return endpointInfo{params: params, returnType: returnType, errorCodes: ep.Errors}
}

// forEndpointReturn synthesizes a named interface for an endpoint's own
// bare object/array return described by its own Attributes (not a named
// object.<n> ref, which resolves separately via ifaceSet.resolve) -
// mirrors phpgen's forEndpointReturn's spec-reading exactly (a bare,
// untyped "array" return with attributes describes each item's shape).
func forEndpointReturn(returnIfaces *ifaceSet, ep webfunction.Endpoint) returnLocalShapes {
	if len(ep.Attributes) == 0 {
		return returnLocalShapes{}
	}
	fields := attributeFields(ep.Attributes)

	switch {
	case ep.Returns.HasBase("object"):
		return returnLocalShapes{object: returnIfaces.resolveLocal(ep.Name, fields)}
	case ep.Returns.HasBareArray():
		return returnLocalShapes{arrayOfItem: returnIfaces.resolveLocal(ep.Name, fields)}
	default:
		return returnLocalShapes{}
	}
}

func writeMethod(b *strings.Builder, ep webfunction.Endpoint, methodName string, info endpointInfo) {
	writeDocComment(b, "    ", ep.Docs)
	for _, e := range info.errorCodes {
		code := strings.TrimSpace(e.Code)
		if code == "" {
			continue
		}
		line := "    # May raise BadRequestError with code " + code
		if dl := docLines(e.Docs); len(dl) > 0 {
			line += " - " + dl[0]
		}
		b.WriteString(line + "\n")
	}

	var kwParts []string
	hasUnnamed := false
	for _, p := range info.params {
		// RBS keyword-parameter names, like Ruby's own formal keyword
		// parameters, must be plain identifiers - unlike a record-type
		// hash key (recordKey), which can be a quoted string. A wire
		// argument name that isn't identifier-safe (e.g. digit-leading)
		// genuinely can't be declared as a named keyword param in a
		// method TYPE, even though the real call site can still pass it
		// via quoted-keyword syntax (client.foo("2fa-enabled": true)).
		// Rather than silently mistyping or dropping it, such params
		// fall through to a single trailing **untyped catch-all below,
		// so the real call still type-checks (just without per-field
		// precision for that one field).
		if !rubyIdentifier(p.name) {
			hasUnnamed = true
			continue
		}
		prefix := ""
		if p.optional {
			prefix = "?"
		}
		kwParts = append(kwParts, prefix+p.name+": "+p.rbsType)
	}
	if hasUnnamed {
		kwParts = append(kwParts, "**untyped")
	}
	sig := "()"
	if len(kwParts) > 0 {
		sig = "(" + strings.Join(kwParts, ", ") + ")"
	}
	returnType := info.returnType
	if strings.Contains(returnType, " | ") {
		returnType = "(" + returnType + ")"
	}
	b.WriteString(fmt.Sprintf("    def %s: %s -> %s\n\n", methodName, sig, returnType))
}

func writeInterface(b *strings.Builder, iface ifaceDef) {
	b.WriteString("  interface " + iface.name + "\n")
	overloads := make([]string, len(iface.fields))
	for i, f := range iface.fields {
		rt := f.rbsType
		if strings.Contains(rt, " | ") {
			rt = "(" + rt + ")"
		}
		overloads[i] = "(" + rbsStringLiteral(f.name) + ") -> " + rt
	}
	b.WriteString("    def []: " + strings.Join(overloads, "\n          | ") + "\n")
	b.WriteString("  end\n\n")
}

func visibleEndpoints(pkg *webfunction.Package) []webfunction.Endpoint {
	out := make([]webfunction.Endpoint, 0, len(pkg.Endpoints))
	for _, ep := range pkg.Endpoints {
		if ep.HasFlag("private") {
			continue
		}
		out = append(out, ep)
	}
	return out
}