package jsgen

import (
	"fmt"
	"sort"
	"strings"

	"wfn/webfunction"
)

// ImportSpecifier is the module specifier the generated file imports
// { Client } from - the npm package name, once webfunction-js is
// published.
const ImportSpecifier = "webfunction"

// endpointTypedefs holds the typedef(s) (if any) describing an endpoint's
// argument shape and return shape.
type endpointTypedefs struct {
	args    string
	returns localTypedefs
}

// Generate turns a fetched Package into the source of a JS module. It
// exports an async createClient(options) factory: options are passed
// straight through to Client.fromPackageEndpoint (e.g. { bearerAuth,
// version }), and the resolved client is wrapped with a named method per
// endpoint, plus the underlying client's own package/call surface.
//
// Every method's argument and return shape gets its own JSDoc @typedef
// (deduped across endpoints with identical shapes), and createClient's own
// @returns points at a composite typedef listing every method - so editors
// get full intellisense on both the call args and the result, not just a
// generic Object.
func Generate(pkg *webfunction.Package, sourceURL string) (string, error) {
	var b strings.Builder

	writeHeader(&b, pkg, sourceURL)

	typedefs := newTypedefSet(pkg)
	endpoints := visibleEndpoints(pkg)

	// Pre-compute every endpoint's arg/return typedefs first, so the
	// typedef block above the factory is complete before anything
	// references one - including the composite client typedef below.
	perEndpoint := make(map[string]endpointTypedefs, len(endpoints))
	for _, ep := range endpoints {
		perEndpoint[ep.Name] = endpointTypedefs{
			args:    typedefs.forFields(pascalCase(ep.Name)+"Args", argumentFields(ep.Arguments), "argument"),
			returns: forEndpointReturn(typedefs, ep),
		}
	}

	clientTypedef := buildClientTypedef(typedefs, pkg, endpoints, perEndpoint)

	b.WriteString(typedefs.render())

	writeFactory(&b, pkg, sourceURL, endpoints, perEndpoint, clientTypedef, typedefs)

	return b.String(), nil
}

// forEndpointReturn builds typedefs describing the endpoint's own inline
// return shape from its attributes - not a named object.<name> ref, which
// resolves separately in returnType. Per spec, attributes describes the
// object itself when returns is (or includes) bare "object". In practice,
// this covers two different real conventions seen for `paginated`
// endpoints specifically:
//
//   - Some document the {previous, page, next} envelope directly as their
//     own attributes (a literal "page" field among them) - nothing special
//     to do, attributes already describe the whole return shape.
//   - Others document only the per-item fields (no "page" attribute) and
//     leave the standard paginated envelope unstated - the real return is
//     still {previous, page: Array<Item>, next}, so that's synthesized
//     here via the generic Page<T> typedef.
//
// Outside the paginated case, a bare, untyped "array" return with
// attributes is - in every real-world package seen so far - describing
// each item's shape too, so that's honored as well even though the spec's
// letter only defines attributes for the bare "object" case. Returns the
// zero value if there's nothing to build.
func forEndpointReturn(typedefs *typedefSet, ep webfunction.Endpoint) localTypedefs {
	if len(ep.Attributes) == 0 {
		return localTypedefs{}
	}

	if ep.HasFlag("paginated") && ep.Returns.HasBase("object") && !hasAttributeNamed(ep.Attributes, "page") {
		itemName := typedefs.forFields(pascalCase(ep.Name)+"Result", attributeFields(ep.Attributes), "attribute")
		if itemName == "" {
			return localTypedefs{}
		}
		pageTypedef := typedefs.ensurePageTypedef()
		return localTypedefs{object: pageTypedef + "<" + itemName + ">"}
	}

	name := typedefs.forFields(pascalCase(ep.Name)+"Result", attributeFields(ep.Attributes), "attribute")
	if name == "" {
		return localTypedefs{}
	}

	switch {
	case ep.Returns.HasBase("object"):
		return localTypedefs{object: name}
	case ep.Returns.HasBareArray():
		return localTypedefs{arrayOfItem: name}
	default:
		return localTypedefs{}
	}
}

func hasAttributeNamed(attrs []webfunction.Attribute, name string) bool {
	for _, a := range attrs {
		if a.Name == name {
			return true
		}
	}
	return false
}

// buildClientTypedef adds the composite typedef describing everything
// createClient's returned object has: the passthrough package/call
// properties, plus one function-typed property per endpoint method.
func buildClientTypedef(typedefs *typedefSet, pkg *webfunction.Package, endpoints []webfunction.Endpoint, perEndpoint map[string]endpointTypedefs) string {
	methodNames := reservedNames()

	lines := []string{
		"@property {Object} package",
		"@property {(name: string, args?: any) => Promise<any>} call",
	}

	for _, ep := range endpoints {
		methodName := uniqueMethodName(methodNames, camelCase(ep.Name))
		td := perEndpoint[ep.Name]

		retType := returnType(typedefs, ep, td.returns)

		var argSig string
		if td.args != "" {
			argSig = "args: " + td.args
		} else {
			argSig = "args?: any"
		}

		lines = append(lines, fmt.Sprintf("@property {(%s) => Promise<%s>} %s", argSig, retType, methodName))
	}

	base := "Client"
	if pkg.Name != "" {
		base = pascalCase(pkg.Name) + "Client"
	}
	return typedefs.addComposite(base, lines)
}

// visibleEndpoints returns the package's endpoints in their original
// order, excluding any flagged "private" (per the spec, tooling SHOULD
// omit these from generated/published output).
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

// reservedNames seeds a method-name set with the wrapper's own passthrough
// properties and JS reserved words, so endpoint methods can't collide with
// them.
func reservedNames() map[string]bool {
	names := make(map[string]bool, len(jsReserved))
	for k := range jsReserved {
		names[k] = true
	}
	return names
}

func uniqueMethodName(used map[string]bool, base string) string {
	name := base
	for i := 2; used[name]; i++ {
		name = fmt.Sprintf("%s%d", base, i)
	}
	used[name] = true
	return name
}

func writeHeader(b *strings.Builder, pkg *webfunction.Package, sourceURL string) {
	b.WriteString("// AUTO-GENERATED by wfn codegen --target js. Do not edit by hand -\n")
	b.WriteString("// re-run codegen instead.\n")
	if pkg.Name != "" {
		b.WriteString("// Package: " + pkg.Name + "\n")
	}
	b.WriteString("// Source:  " + sourceURL + "\n\n")
	b.WriteString("// IMPORTANT - Remember to import '" + ImportSpecifier + "' into your project - npm i " + ImportSpecifier + "\n\n")
	b.WriteString(fmt.Sprintf("import { Client } from %s;\n\n", jsStringLiteral(ImportSpecifier)))
}

// writeFactory writes the exported async createClient(options) function:
// it builds the underlying dynamic client (forwarding options, e.g.
// bearerAuth/version, straight to Client.fromPackageEndpoint) and returns
// the wrapper object, typed as clientTypedef.
func writeFactory(b *strings.Builder, pkg *webfunction.Package, sourceURL string, endpoints []webfunction.Endpoint, perEndpoint map[string]endpointTypedefs, clientTypedef string, typedefs *typedefSet) {
	b.WriteString("/**\n")
	b.WriteString(" * Builds a wrapped client for this package.\n")
	b.WriteString(" *\n")
	b.WriteString(" * @param {Object} [options]\n")
	b.WriteString(" * @param {string} [options.bearerAuth] - Sent as `Authorization: Bearer <token>` on every call.\n")
	b.WriteString(" * @param {string} [options.version] - Sent as the `Api-Version` header.\n")
	b.WriteString(" * @returns {Promise<" + clientTypedef + ">}\n")
	b.WriteString(" */\n")
	b.WriteString("export async function createClient(options = {}) {\n")
	b.WriteString(fmt.Sprintf("  const rawClient = await Client.fromPackageEndpoint(%s, options);\n\n", jsStringLiteral(sourceURL)))

	methodNames := reservedNames()

	b.WriteString("  return {\n")
	b.WriteString("    package: rawClient.package,\n")
	b.WriteString("    call: (name, args) => rawClient.call(name, args),\n")

	for _, ep := range endpoints {
		methodName := uniqueMethodName(methodNames, camelCase(ep.Name))
		b.WriteString("\n")
		writeMethod(b, ep, methodName, perEndpoint[ep.Name], typedefs)
	}

	b.WriteString("  };\n")
	b.WriteString("}\n")
}

func writeMethod(b *strings.Builder, ep webfunction.Endpoint, methodName string, td endpointTypedefs, typedefs *typedefSet) {
	b.WriteString("    /**\n")

	for _, line := range docLines(ep.Docs) {
		if line == "" {
			b.WriteString("     *\n")
			continue
		}
		b.WriteString("     * " + line + "\n")
	}
	if ep.Docs != "" {
		b.WriteString("     *\n")
	}

	if td.args != "" {
		b.WriteString("     * @param {" + td.args + "} args\n")
	} else {
		b.WriteString("     * @param {any} [args]\n")
	}

	b.WriteString("     * @returns {Promise<" + returnType(typedefs, ep, td.returns) + ">}\n")
	b.WriteString("     */\n")
	b.WriteString(fmt.Sprintf("    %s(args) {\n", methodName))
	b.WriteString(fmt.Sprintf("      return rawClient.call(%s, args);\n", jsStringLiteral(ep.Name)))
	b.WriteString("    },\n")
}

func choicesNote(choices []interface{}) string {
	if len(choices) == 0 {
		return ""
	}
	strs := make([]string, len(choices))
	for i, c := range choices {
		strs[i] = fmt.Sprintf("%v", c)
	}
	sort.Strings(strs)
	return "(one of: " + strings.Join(strs, ", ") + ")"
}

func returnType(typedefs *typedefSet, ep webfunction.Endpoint, local localTypedefs) string {
	resolve := func(refName string) string { return typedefs.resolveObjectTypedef(refName, "attribute") }
	return jsdocType(ep.Returns, local, false, resolve)
}

// jsStringLiteral renders a Go string as a single-quoted JS string literal.
func jsStringLiteral(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	return "'" + s + "'"
}