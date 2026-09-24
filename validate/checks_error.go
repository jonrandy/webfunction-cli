package validate

import (
	"fmt"

	"github.com/webfunction-protocol/webfunction-go"
)

// checkEventSourceReturnType flags an endpoint with the "event_source"
// flag whose Returns isn't exactly ["string"]. Per
// https://webfunction.org/package#available-flags: "An endpoint with
// this flag MUST declare a returns of [\"string\"]" - a real MUST, not
// a style preference, since event_source is meant to return a single
// event source URL string.
func (v *validator) checkEventSourceReturnType() {
	for _, e := range v.pkg.Endpoints {
		if !e.HasFlag("event_source") {
			continue
		}
		if !isBareStringReturn(e.Returns) {
			v.emit("event-source-invalid-return-type", Error, EndpointLoc(e.Name).Returns(),
				`endpoint declares the "event_source" flag, which per spec requires returns to be exactly ["string"], but its declared returns type doesn't match`)
		}
	}
}

// isBareStringReturn reports whether t is exactly the single-alternative
// union {"string"} - no other alternatives, no refinement.
func isBareStringReturn(t webfunction.Type) bool {
	return len(t.Union) == 1 && t.Union[0].Base == "string" && t.Union[0].Refinement == ""
}

// checkDuplicateEndpointNames flags an endpoint name appearing more than
// once. webfunction.Package.Endpoint(name) returns only the first match,
// so a duplicate silently masks the later definition(s) - a real
// authoring mistake, not just untidiness.
func (v *validator) checkDuplicateEndpointNames() {
	seen := map[string]int{}
	for _, e := range v.pkg.Endpoints {
		seen[e.Name]++
	}
	for name, count := range seen {
		if count > 1 {
			v.emit("duplicate-endpoint-name", Error, EndpointLoc(name),
				fmt.Sprintf("endpoint name %q appears %d times; only the first definition is reachable via Package.Endpoint(), the rest are silently masked", name, count))
		}
	}
}

// checkDuplicateObjectNames flags an object name appearing more than once
// in pkg.Objects. webfunction.Package.Object(name) (and therefore
// ObjectInContext) returns only the first match regardless of which
// context is being resolved - so if a later entry with the same name
// defines a *different* context's members (e.g. the first entry only has
// Arguments, a later one only has Attributes), that second entry's data
// is completely unreachable, not just redundant.
func (v *validator) checkDuplicateObjectNames() {
	seen := map[string]int{}
	for _, o := range v.pkg.Objects {
		seen[o.Name]++
	}
	for name, count := range seen {
		if count > 1 {
			v.emit("duplicate-object-name", Error, ObjectLoc(name),
				fmt.Sprintf("object name %q appears %d times; only the first entry is ever resolved by Package.ObjectInContext(), so a later entry's arguments/attributes may be silently unreachable", name, count))
		}
	}
}

// checkVersioning flags an internally inconsistent "versioned" flag: the
// flag set with no (or an empty) Versions list, or a Version value that
// isn't actually present in the package's own Versions list.
func (v *validator) checkVersioning() {
	if !v.pkg.Versioned() {
		return
	}
	if len(v.pkg.Versions) == 0 {
		v.emit("versioned-empty-versions", Error, PackageLoc(),
			`package declares the "versioned" flag but its "versions" list is empty or missing`)
		return
	}
	if v.pkg.Version != "" {
		found := false
		for _, ver := range v.pkg.Versions {
			if ver == v.pkg.Version {
				found = true
				break
			}
		}
		if !found {
			v.emit("version-not-in-versions", Error, PackageLoc(),
				fmt.Sprintf("package's declared version %q is not present in its own \"versions\" list %v", v.pkg.Version, v.pkg.Versions))
		}
	}
}

// checkObjectRefs walks every Type in the package - endpoint returns,
// argument types, attribute types, and the types nested inside every
// object.<n> definition itself - and flags any object.<n> reference that
// doesn't resolve via ObjectInContext (either the name doesn't exist at
// all, or it exists but defines no members for the context it was
// referenced from).
//
// Iterating pkg.Objects as a flat list (rather than recursively
// resolving each ref into the referenced object's own members) already
// covers every reachable reference exactly once, so no cycle-safety
// bookkeeping is needed here - unlike openapiconvert's schemaSet, which
// needs it because it's *building* named output, not just checking
// existence.
func (v *validator) checkObjectRefs() {
	for _, e := range v.pkg.Endpoints {
		v.walkType(e.Returns, webfunction.AttributeContext, EndpointLoc(e.Name).Returns())
		for _, a := range e.Arguments {
			v.walkType(a.Type, webfunction.ArgumentContext, EndpointLoc(e.Name).Argument(a.Name))
		}
		for _, a := range e.Attributes {
			v.walkType(a.Type, webfunction.AttributeContext, EndpointLoc(e.Name).Attribute(a.Name))
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			v.walkType(a.Type, webfunction.ArgumentContext, ObjectLoc(o.Name).Argument(a.Name))
		}
		for _, a := range o.Attributes {
			v.walkType(a.Type, webfunction.AttributeContext, ObjectLoc(o.Name).Attribute(a.Name))
		}
	}
}

// walkType checks every alternative in t (recursing into array element
// types) for a dangling object.<n> reference.
func (v *validator) walkType(t webfunction.Type, ctx webfunction.ObjectContext, loc Loc) {
	for _, alt := range t.Union {
		if alt.IsObjectRef() {
			if v.pkg.ObjectInContext(alt.Refinement, ctx) == nil {
				ctxName := "attribute"
				if ctx == webfunction.ArgumentContext {
					ctxName = "argument"
				}
				v.emit("dangling-object-ref", Error, loc,
					fmt.Sprintf("references object.%s in %s context, which doesn't exist or defines no %s members", alt.Refinement, ctxName, ctxName))
			}
		}
		if alt.Of != nil {
			v.walkType(*alt.Of, ctx, loc)
		}
	}
}