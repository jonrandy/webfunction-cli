// Package phpgen generates a typed PHP client for a webfunction package,
// targeting webfunction-php (github.com/webfunction-protocol/webfunction-php) and PHP
// 8.1+ (webfunction-php's own stated floor).
//
// Design decisions (confirmed with Jon before writing this):
//  1. Return/argument shapes use plain PHPDoc array-shape annotations
//     (`array{key: type, ...}`), not generated DTO/value classes.
//  2. Named, reusable PHPDoc type aliases (`@phpstan-type`/`@psalm-type`)
//     are used for object.<name> refs specifically - decision #1 doesn't
//     cover those, since a plain inline array-shape can't self-reference
//     or be reused, which named object refs sometimes need to (mirrors
//     jsgen's named JSDoc typedefs for the same case). Everything else
//     (each endpoint's own Args/Return shape) stays inline.
//  3. Target PHP 8.1+ (matches webfunction-php's own README-stated floor;
//     no reason to hedge for compatibility the library itself lacks).
//  4. Pagination needs no special shape-resolution handling beyond what
//     object-ref resolution already does (mirrors the JS resolution).
//  5. `@throws` per endpoint, mirroring jsgen's design and scoping
//     (endpoint's own Errors only, not merged with package-level Errors -
//     see phpDocType/jsgen's equivalent reasoning) - viable here because
//     webfunction-php's error model (confirmed from its README) matches
//     webfunction-js's closely: WebFunction\Error subclasses
//     (BadRequestError/UnexpectedStatusCodeError/JsonParseError/
//     UnresolvedPromiseError), same [code, message, details] triple.
//
// Known open items:
//   - UnresolvedPromiseError (pipelining) isn't referenced in generated
//     @throws - it's not tied to any specific endpoint or ordinary call
//     path, unlike the other three.
//
// Confirmed from the real Client.php source (github.com/webfunction-protocol/webfunction-php/blob/main/src/Client.php,
// not just its README) - two corrections this required:
//   - call()'s real signature is `call(string $endpointName, array $args = [])`
//     - non-nullable, empty-array default - NOT `?array $args = null` as
//     first assumed. Both the generated escape-hatch call() and every
//     per-endpoint method's optional-args case now default to `= []` to
//     match exactly (passing `null` through would have thrown a
//     TypeError at that boundary - a real bug this correction fixes).
//   - call()'s real return type is declared `mixed`, not the narrower
//     array|string|int|float|bool|null first derived from the README's
//     prose alone - because it can also return a real \WebFunction\Page
//     object when the raw endpoint name happens to be paginated, which
//     the generic passthrough (unlike a per-endpoint method) can't rule
//     out. Per-endpoint native return types are still narrowed
//     (nativeReturnType) since those DO know which endpoint they are.
//   - Client::fromPackageEndpoint(url, bearerAuth: ..., version: ...)
//     really does take named arguments (matches what was already used).
//   - A getPackage(): ?Package accessor exists - now exposed on the
//     generated class too, for parity with what jsgen exposes via .package.
//   - Numeric refinements (u32/u64/i32/i64/timestamp -> int, f32/f64 ->
//     float) are confirmed and narrow phpDocAlt's number handling.
package phpgen

import (
	"fmt"
	"regexp"
	"strings"

	"wfn/webfunction"
)

// Generate builds a typed PHP client for pkg, targeting webfunction-php.
// namespace is the PHP namespace for the generated class (validated: each
// backslash-separated segment must be a valid PHP identifier).
func Generate(pkg *webfunction.Package, sourceURL, namespace string) (string, error) {
	if err := validateNamespace(namespace); err != nil {
		return "", err
	}

	aliases := newAliasSet(pkg)
	endpoints := visibleEndpoints(pkg)

	// Pre-compute every endpoint's shapes first, so alias resolution
	// (which mutates aliases as it goes, same as jsgen's typedefSet) is
	// complete before the class docblock (which lists every alias) is
	// written.
	infos := make(map[string]endpointInfo, len(endpoints))
	pageClassNames := make(map[string]bool)
	for _, ep := range endpoints {
		infos[ep.Name] = buildEndpointInfo(aliases, ep, pageClassNames)
	}

	var b strings.Builder
	b.WriteString("<?php\n\ndeclare(strict_types=1);\n\n")
	b.WriteString("namespace " + namespace + ";\n\n")
	b.WriteString("use WebFunction\\Page;\n")
	b.WriteString("use WebFunction\\Package;\n")
	b.WriteString("use WebFunction\\Pipeline;\n")
	b.WriteString("use WebFunction\\BadRequestError;\n")
	b.WriteString("use WebFunction\\UnexpectedStatusCodeError;\n")
	b.WriteString("use WebFunction\\JsonParseError;\n\n")
	b.WriteString("// IMPORTANT - composer require webfunction-protocol/webfunction\n\n")

	writeClientClass(&b, pkg, sourceURL, endpoints, infos, aliases)

	for _, ep := range endpoints {
		if info := infos[ep.Name]; info.pageWrapper != nil {
			writePageWrapper(&b, info.pageWrapper, aliases)
		}
	}

	return b.String(), nil
}

// escapeHatchReturnType is the native PHP return type for the generic
// call() escape-hatch passthrough specifically. Confirmed exactly from
// the real Client.php source (github.com/webfunction-protocol/webfunction-php/blob/main/src/Client.php):
// its own call() is declared `: mixed` - not the narrower
// array|string|int|float|bool|null this was first assumed to be from the
// README's prose alone, because it can also return a real \WebFunction\Page
// object when the raw endpoint name happens to be a paginated one (this
// passthrough doesn't know which endpoint is being called at generation
// time, so it can't rule that out the way a per-endpoint method can).
const escapeHatchReturnType = "mixed"

// nativeReturnType maps a specific endpoint's Returns type to the actual
// native PHP return-type declaration for its generated method - narrower
// than nonPaginatedReturnType whenever the wire type says something more
// specific (e.g. a bare "boolean" return narrows to `bool`, not the full
// union), since PHP enforces this at runtime and a caller benefits from
// knowing precisely what a given endpoint always returns. Structured
// shapes (bare object/array, or a named object.<name> ref) still
// collapse to the native "array" - PHP arrays don't carry shape at the
// language level, so PHPDoc's @return is what carries that precision,
// not the native declaration. "any" (or anything unrecognized) collapses
// the whole type to `mixed` alone - PHP doesn't allow combining `mixed`
// with other types in a union, unlike `null`, which can be combined.
func nativeReturnType(t webfunction.Type) string {
	var parts []string
	sawAny := false
	for _, alt := range t.Union {
		switch alt.Base {
		case "string":
			parts = append(parts, "string")
		case "boolean":
			parts = append(parts, "bool")
		case "number":
			switch alt.Refinement {
			case "u32", "u64", "i32", "i64", "timestamp":
				parts = append(parts, "int")
			case "f32", "f64":
				parts = append(parts, "float")
			default:
				parts = append(parts, "int", "float")
			}
		case "null":
			parts = append(parts, "null")
		case "array", "object":
			parts = append(parts, "array")
		default: // "any", or anything not in the confirmed base-type list
			sawAny = true
		}
	}
	if sawAny || len(parts) == 0 {
		return "mixed"
	}
	return strings.Join(dedupe(parts), "|")
}

// endpointInfo holds everything about one endpoint needed to write both
// its method and (for paginated endpoints) its companion wrapper class.
type endpointInfo struct {
	argsShape    string // inline "array{...}" PHPDoc shape, or "" if no args
	argsOptional bool
	returnShape  string       // PHPDoc @return shape, ignored if pageWrapper != nil
	pageWrapper  *pageWrapper // non-nil for a paginated endpoint
}

type pageWrapper struct {
	className string
	itemShape string
}

func buildEndpointInfo(aliases *aliasSet, ep webfunction.Endpoint, pageClassNames map[string]bool) endpointInfo {
	argFields := argumentFields(ep.Arguments)
	var argsShape string
	if len(argFields) > 0 {
		argsShape = aliases.renderShape(argFields, "argument")
	}

	if ep.HasFlag("paginated") {
		itemShape := "mixed"
		resolve := func(refName string) string { return aliases.resolve(refName, "attribute") }
		for _, a := range ep.Attributes {
			if a.Name != "page" {
				continue
			}
			for _, alt := range a.Type.Union {
				if alt.Base == "array" && alt.Of != nil {
					itemShape = phpDocType(*alt.Of, localShapes{}, false, resolve, nil)
				}
			}
		}
		className := uniqueClassName(pageClassNames, pascalCase(ep.Name)+"Page")
		return endpointInfo{
			argsShape: argsShape, argsOptional: allOptional(argFields),
			pageWrapper: &pageWrapper{className: className, itemShape: itemShape},
		}
	}

	local := forEndpointReturn(aliases, ep)
	resolve := func(refName string) string { return aliases.resolve(refName, "attribute") }
	returnShape := phpDocType(ep.Returns, local, false, resolve, nil)

	return endpointInfo{
		argsShape: argsShape, argsOptional: allOptional(argFields),
		returnShape: returnShape,
	}
}

// forEndpointReturn builds localShapes describing the endpoint's own
// inline return shape from its attributes - not a named object.<name>
// ref, which resolves separately via aliasSet. Mirrors jsgen's
// forEndpointReturn (same spec-reading: attributes describes the object
// itself for bare "object"; a bare, untyped "array" return with
// attributes describes each item's shape, matching real-world packages
// even though the spec's letter only defines attributes for "object").
func forEndpointReturn(aliases *aliasSet, ep webfunction.Endpoint) localShapes {
	if len(ep.Attributes) == 0 {
		return localShapes{}
	}
	fields := attributeFields(ep.Attributes)
	shape := aliases.renderShape(fields, "attribute")

	switch {
	case ep.Returns.HasBase("object"):
		return localShapes{object: shape}
	case ep.Returns.HasBareArray():
		return localShapes{arrayOfItem: shape}
	default:
		return localShapes{}
	}
}

func writeClientClass(b *strings.Builder, pkg *webfunction.Package, sourceURL string, endpoints []webfunction.Endpoint, infos map[string]endpointInfo, aliases *aliasSet) {
	b.WriteString("/**\n")
	name := pkg.Name
	if name == "" {
		name = "(unnamed package)"
	}
	b.WriteString(" * Generated client for the " + name + " webfunction package.\n")
	b.WriteString(" * Source: " + sourceURL + "\n")
	if aliases.hasAliases() {
		b.WriteString(" *\n")
		for _, a := range aliases.ordered {
			b.WriteString(" * @phpstan-type " + a.name + " = " + a.phpType + "\n")
			b.WriteString(" * @psalm-type " + a.name + " = " + a.phpType + "\n")
		}
	}
	b.WriteString(" */\n")
	b.WriteString("final class Client\n{\n")
	b.WriteString("    private \\WebFunction\\Client $client;\n\n")
	b.WriteString("    public function __construct(?string $bearerAuth = null, ?string $version = null)\n")
	b.WriteString("    {\n")
	b.WriteString("        $this->client = \\WebFunction\\Client::fromPackageEndpoint('" + phpStringEscape(sourceURL) + "', bearerAuth: $bearerAuth, version: $version);\n")
	b.WriteString("    }\n\n")

	b.WriteString("    /**\n     * Escape hatch for calling an endpoint this client doesn't have a typed method for.\n     */\n")
	b.WriteString("    public function call(string $name, array $args = []): " + escapeHatchReturnType + "\n    {\n        return $this->client->call($name, $args);\n    }\n\n")

	b.WriteString("    public function getPackage(): ?Package\n    {\n        return $this->client->getPackage();\n    }\n\n")

	b.WriteString("    public function setBearerAuth(?string $bearerAuth): void\n    {\n        $this->client->setBearerAuth($bearerAuth);\n    }\n\n")
	b.WriteString("    public function setVersion(?string $version): void\n    {\n        $this->client->setVersion($version);\n    }\n\n")
	b.WriteString("    public function setPipeline(?Pipeline $pipeline): void\n    {\n        $this->client->setPipeline($pipeline);\n    }\n\n")

	used := map[string]bool{}
	for _, ep := range endpoints {
		methodName := uniqueMethodName(used, camelCase(ep.Name))
		writeMethod(b, ep, methodName, infos[ep.Name])
	}

	b.WriteString("}\n\n")
}

func writeMethod(b *strings.Builder, ep webfunction.Endpoint, methodName string, info endpointInfo) {
	b.WriteString("    /**\n")
	for _, l := range docLines(ep.Docs) {
		if l == "" {
			b.WriteString("     *\n")
		} else {
			b.WriteString("     * " + l + "\n")
		}
	}
	if info.argsShape != "" {
		b.WriteString("     * @param " + info.argsShape + " $args\n")
	}
	if info.pageWrapper != nil {
		b.WriteString("     * @return " + info.pageWrapper.className + "\n")
	} else {
		b.WriteString("     * @return " + info.returnShape + "\n")
	}

	for _, e := range ep.Errors {
		code := strings.TrimSpace(e.Code)
		if code == "" {
			continue
		}
		line := "     * @throws BadRequestError " + code
		if dl := docLines(e.Docs); len(dl) > 0 {
			line += " - " + dl[0]
			b.WriteString(line + "\n")
			for _, extra := range dl[1:] {
				b.WriteString("     * " + extra + "\n")
			}
		} else {
			b.WriteString(line + "\n")
		}
	}
	b.WriteString("     * @throws UnexpectedStatusCodeError Non-200, non-400 response.\n")
	b.WriteString("     * @throws JsonParseError Response body wasn't valid JSON.\n")
	b.WriteString("     */\n")

	var params string
	if info.argsShape != "" {
		if info.argsOptional {
			// Not `?array $args = null` - the real Client::call() takes
			// a non-nullable `array $args = []` (confirmed from its real
			// source), so `null` passed straight through would throw a
			// TypeError at that boundary. `[]` matches exactly and needs
			// no `?? []` coalesce at the call site below.
			params = "array $args = []"
		} else {
			params = "array $args"
		}
	}
	nativeReturn := nativeReturnType(ep.Returns)
	if info.pageWrapper != nil {
		nativeReturn = info.pageWrapper.className
	}
	b.WriteString(fmt.Sprintf("    public function %s(%s): %s\n", methodName, params, nativeReturn))
	b.WriteString("    {\n")

	callArgs := "'" + phpStringEscape(ep.Name) + "'"
	if info.argsShape != "" {
		callArgs += ", $args"
	}
	if info.pageWrapper != nil {
		b.WriteString("        $raw = $this->client->call(" + callArgs + ");\n")
		b.WriteString("        /** @var Page $raw */\n")
		b.WriteString("        return new " + info.pageWrapper.className + "($raw);\n")
	} else {
		b.WriteString("        return $this->client->call(" + callArgs + ");\n")
	}
	b.WriteString("    }\n\n")
}

// writePageWrapper generates a small real PHP class per paginated
// endpoint, wrapping the real \WebFunction\Page instance with item-typed
// access. Necessary (not just nice-to-have) because Page.php's own
// getItems() only promises `array<int, mixed>` in its PHPDoc - there's no
// generics/@template on the real class for jsgen-style per-caller
// item-typing the way a docs-only annotation could paper over in JS;
// PHP's native return-type declarations are runtime-enforced, so a
// generated method can't just *claim* a different return shape than what
// it actually returns (a real \WebFunction\Page object) without an actual
// object of that claimed type - hence a real wrapper class, not just a
// fancier docblock.
func writePageWrapper(b *strings.Builder, w *pageWrapper, aliases *aliasSet) {
	b.WriteString("/**\n")
	b.WriteString(" * A page of results, wrapping the real \\WebFunction\\Page instance with\n")
	b.WriteString(" * item-level typing on getItems() that the real class's own PHPDoc\n")
	b.WriteString(" * (`@return array<int, mixed>`) can't provide by itself.\n")
	if imports := aliasesReferencedIn(w.itemShape, aliases); len(imports) > 0 {
		b.WriteString(" *\n")
		for _, name := range imports {
			b.WriteString(" * @phpstan-import-type " + name + " from Client\n")
			b.WriteString(" * @psalm-import-type " + name + " from Client\n")
		}
	}
	b.WriteString(" */\n")
	b.WriteString("final class " + w.className + " implements \\IteratorAggregate, \\Countable\n{\n")
	b.WriteString("    private Page $page;\n\n")
	b.WriteString("    public function __construct(Page $page)\n    {\n        $this->page = $page;\n    }\n\n")
	b.WriteString("    /**\n     * @return list<" + w.itemShape + ">\n     */\n")
	b.WriteString("    public function getItems(): array\n    {\n        return $this->page->getItems();\n    }\n\n")
	b.WriteString("    public function hasNext(): bool\n    {\n        return $this->page->hasNext();\n    }\n\n")
	b.WriteString("    public function hasPrevious(): bool\n    {\n        return $this->page->hasPrevious();\n    }\n\n")
	b.WriteString("    public function nextPage(): ?self\n    {\n        $next = $this->page->nextPage();\n        return $next === null ? null : new self($next);\n    }\n\n")
	b.WriteString("    public function previousPage(): ?self\n    {\n        $prev = $this->page->previousPage();\n        return $prev === null ? null : new self($prev);\n    }\n\n")
	b.WriteString("    public function count(): int\n    {\n        return $this->page->count();\n    }\n\n")
	b.WriteString("    public function getIterator(): \\ArrayIterator\n    {\n        return $this->page->getIterator();\n    }\n")
	b.WriteString("}\n\n")
}

// aliasesReferencedIn returns the names (in registration order) of every
// known alias that appears as a whole identifier within shape - used to
// decide which @…-import-type lines a Page wrapper class needs, since it
// lives outside the main Client class's own docblock where aliases are
// defined.
func aliasesReferencedIn(shape string, aliases *aliasSet) []string {
	var found []string
	for _, a := range aliases.ordered {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(a.name) + `\b`)
		if re.MatchString(shape) {
			found = append(found, a.name)
		}
	}
	return found
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

func uniqueClassName(used map[string]bool, base string) string {
	name := base
	for i := 2; used[name]; i++ {
		name = fmt.Sprintf("%s%d", base, i)
	}
	used[name] = true
	return name
}

var namespaceSegmentRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func validateNamespace(namespace string) error {
	namespace = strings.Trim(namespace, `\`)
	if namespace == "" {
		return fmt.Errorf("namespace must not be empty")
	}
	for _, seg := range strings.Split(namespace, `\`) {
		if !namespaceSegmentRe.MatchString(seg) {
			return fmt.Errorf("invalid namespace segment %q in %q: must match %s", seg, namespace, namespaceSegmentRe.String())
		}
	}
	return nil
}

func phpStringEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return s
}