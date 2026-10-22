package validate

import (
	"fmt"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"

	"wfn/generics"
)

// checkGenerics validates generic object usage. Generic objects are
// object definitions whose member types use the bare token "T". A type
// applies one to a type with { "object.<name>": <type> }, and
// { "array": <type> } applies an array to a type. An application may appear
// anywhere a type is expected, including inside an object definition. "T"
// is only in scope inside an object definition, and a generic object must
// not apply itself, directly or through other generic objects.
//
// These are authoring-time hard errors per the spec. Clients never fail
// on any of them at runtime - they fall back (an unapplied "T" becomes
// "any", an argument on a static object is ignored, a malformed
// application is "any") - so the messages say what the fallback is.
func (v *validator) checkGenerics() {
	for _, e := range v.pkg.Endpoints {
		v.genericsType(e.Returns, webfunction.AttributeContext, EndpointLoc(e.Name).Returns(), false)
		for _, a := range e.Arguments {
			v.genericsType(a.Type, webfunction.ArgumentContext, EndpointLoc(e.Name).Argument(a.Name), false)
		}
		for _, a := range e.Attributes {
			v.genericsType(a.Type, webfunction.AttributeContext, EndpointLoc(e.Name).Attribute(a.Name), false)
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			v.genericsType(a.Type, webfunction.ArgumentContext, ObjectLoc(o.Name).Argument(a.Name), true)
		}
		for _, a := range o.Attributes {
			v.genericsType(a.Type, webfunction.AttributeContext, ObjectLoc(o.Name).Attribute(a.Name), true)
		}
	}
	v.checkGenericCycles()
	for _, c := range generics.Collisions(v.pkg) {
		v.emit("generic-name-collision", Error, ObjectLoc(c.Name),
			fmt.Sprintf("the object generated for the generic application %s is named %q, which a user-defined object already uses; the generated object is renamed in output, so rename the user-defined object to avoid the clash", c.Application, c.Name))
	}
}

// genericsType checks a type written at loc. inObject is true while
// walking an object definition's member types, the only place "T" is in
// scope.
func (v *validator) genericsType(t webfunction.Type, ctx webfunction.ObjectContext, loc Loc, inObject bool) {
	if !inObject && typeUsesParam(t) {
		v.emit("generic-parameter-outside-object", Error, loc,
			`uses the type parameter "T" outside an object definition; "T" is only in scope inside the generic object that uses it, so at runtime it reads as "any"`)
	}
	v.genericsAlts(t, ctx, loc)
}

// genericsAlts checks every alternative in t, recursing into array
// elements and application arguments.
func (v *validator) genericsAlts(t webfunction.Type, ctx webfunction.ObjectContext, loc Loc) {
	ctxName := "attribute"
	if ctx == webfunction.ArgumentContext {
		ctxName = "argument"
	}
	for _, alt := range t.Union {
		switch {
		case alt.Malformed == webfunction.MalformedKeyCount:
			v.emit("generic-application-key-count", Error, loc,
				`an application must be an object with exactly one key ({"object.<name>": <type>} or {"array": <type>}); at runtime this type resolves to "any"`)
		case alt.Malformed == webfunction.MalformedKeyPrefix:
			v.emit("generic-application-key-prefix", Error, loc,
				`an application's key must be "object.<name>" (a generic object applied to a type) or "array" (an array of a type); at runtime this type resolves to "any"`)
		case alt.Malformed == webfunction.MalformedArgument:
			v.emit("generic-application-argument", Error, loc,
				`an application's value must be a valid type; at runtime this type resolves to "any"`)
		case alt.Arg != nil:
			// A missing object is already reported as dangling-object-ref.
			if obj := v.pkg.Object(alt.Refinement); obj != nil {
				switch {
				case !generics.IsGeneric(obj):
					v.emit("generic-argument-on-static-object", Error, loc,
						fmt.Sprintf(`applies a type argument to object.%s, which has no "T" type parameter in any member; at runtime the argument is ignored`, alt.Refinement))
				case v.pkg.ObjectInContext(alt.Refinement, ctx) != nil && !generics.UsesParam(obj, ctx):
					v.emit("generic-argument-unused-in-context", Warning, loc,
						fmt.Sprintf(`applies a type argument to object.%s in %s context, but its %s members don't use "T" (only its other context does), so the argument has no effect here`, alt.Refinement, ctxName, ctxName))
				}
			}
			v.genericsAlts(*alt.Arg, ctx, loc)
		case alt.Base == "array" && alt.Of != nil:
			v.genericsAlts(*alt.Of, ctx, loc)
		case alt.IsObjectRef():
			if generics.IsGeneric(v.pkg.Object(alt.Refinement)) {
				v.emit("generic-unapplied", Error, loc,
					fmt.Sprintf(`references generic object.%s without applying a type argument; write {"object.%s": <type>} (at runtime its "T" falls back to "any")`, alt.Refinement, alt.Refinement))
			}
		}
	}
}

// typeUsesParam reports whether the type parameter "T" appears anywhere
// in t: bare, in a union, as an array element, or inside an application's
// argument.
func typeUsesParam(t webfunction.Type) bool {
	for _, alt := range t.Union {
		if alt.Param {
			return true
		}
		if alt.Of != nil && typeUsesParam(*alt.Of) {
			return true
		}
		if alt.Arg != nil && typeUsesParam(*alt.Arg) {
			return true
		}
	}
	return false
}

// checkGenericCycles reports generic objects that apply themselves,
// directly or through other generic objects. Each cycle is reported once,
// at the first of its objects in definition order, naming the path.
func (v *validator) checkGenericCycles() {
	var order []string
	edges := map[string][]string{}
	for i := range v.pkg.Objects {
		o := &v.pkg.Objects[i]
		if !generics.IsGeneric(o) {
			continue
		}
		order = append(order, o.Name)
		seen := map[string]bool{}
		add := func(target string) {
			if seen[target] {
				return
			}
			seen[target] = true
			if t := v.pkg.Object(target); t != nil && generics.IsGeneric(t) {
				edges[o.Name] = append(edges[o.Name], target)
			}
		}
		for _, a := range o.Arguments {
			applicationTargets(a.Type, add)
		}
		for _, a := range o.Attributes {
			applicationTargets(a.Type, add)
		}
	}
	reported := map[string]bool{}
	for _, name := range order {
		if reported[name] {
			continue
		}
		path := findCycle(name, edges)
		if path == nil {
			continue
		}
		for _, n := range path {
			reported[n] = true
		}
		v.emit("generic-application-cycle", Error, ObjectLoc(name),
			fmt.Sprintf("generic object.%s applies itself (%s); a generic object must not apply itself, directly or through other generic objects", name, strings.Join(path, " -> ")))
	}
}

// applicationTargets calls add with the name of every object applied to a
// type anywhere inside t, including inside array elements and the
// arguments of other applications.
func applicationTargets(t webfunction.Type, add func(string)) {
	for _, alt := range t.Union {
		if alt.Arg != nil {
			add(alt.Refinement)
			applicationTargets(*alt.Arg, add)
		}
		if alt.Of != nil {
			applicationTargets(*alt.Of, add)
		}
	}
}

// findCycle returns a path start -> ... -> start through edges, or nil if
// start cannot reach itself.
func findCycle(start string, edges map[string][]string) []string {
	visited := map[string]bool{start: true}
	var path []string
	var walk func(n string) bool
	walk = func(n string) bool {
		path = append(path, n)
		for _, m := range edges[n] {
			if m == start {
				path = append(path, m)
				return true
			}
			if !visited[m] {
				visited[m] = true
				if walk(m) {
					return true
				}
			}
		}
		path = path[:len(path)-1]
		return false
	}
	if walk(start) {
		return path
	}
	return nil
}