package validate

import (
	"fmt"

	"github.com/webfunction-protocol/webfunction-go"

	"wfn/generics"
)

// checkGenerics validates generic object usage. Generic objects are
// object definitions whose member types use the bare token "T"; they are
// applied to a type with { "object.<name>": <type> } wherever a type is
// written outside an object definition (endpoint returns, argument types
// and attribute types).
//
// These are authoring-time hard errors per the spec. Clients never fail
// on any of them at runtime - they fall back (an unapplied "T" becomes
// "any", an argument on a static object is ignored, a malformed
// application is "any") - so the messages say what the fallback is.
func (v *validator) checkGenerics() {
	for _, e := range v.pkg.Endpoints {
		v.genericsType(e.Returns, webfunction.AttributeContext, EndpointLoc(e.Name).Returns(), false, false)
		for _, a := range e.Arguments {
			v.genericsType(a.Type, webfunction.ArgumentContext, EndpointLoc(e.Name).Argument(a.Name), false, false)
		}
		for _, a := range e.Attributes {
			v.genericsType(a.Type, webfunction.AttributeContext, EndpointLoc(e.Name).Attribute(a.Name), false, false)
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			v.genericsType(a.Type, webfunction.ArgumentContext, ObjectLoc(o.Name).Argument(a.Name), true, false)
		}
		for _, a := range o.Attributes {
			v.genericsType(a.Type, webfunction.AttributeContext, ObjectLoc(o.Name).Attribute(a.Name), true, false)
		}
	}
	for _, c := range generics.Collisions(v.pkg) {
		v.emit("generic-name-collision", Error, ObjectLoc(c.Name),
			fmt.Sprintf("the object generated for the generic application %s is named %q, which a user-defined object already uses; the generated object is renamed in output, so rename the user-defined object to avoid the clash", c.Application, c.Name))
	}
}

// genericsType checks every alternative in t. inObject is true while
// walking an object definition's member types (where applications are not
// allowed); nested is true once inside an application's own argument, so
// a nested application inside an object is reported once, not per level.
func (v *validator) genericsType(t webfunction.Type, ctx webfunction.ObjectContext, loc Loc, inObject, nested bool) {
	ctxName := "attribute"
	if ctx == webfunction.ArgumentContext {
		ctxName = "argument"
	}
	for _, alt := range t.Union {
		switch {
		case alt.Malformed == webfunction.MalformedKeyCount:
			v.emit("generic-application-key-count", Error, loc,
				`a generic application must be an object with exactly one key ({"object.<name>": <type>}); at runtime this type resolves to "any"`)
		case alt.Malformed == webfunction.MalformedKeyPrefix:
			v.emit("generic-application-key-prefix", Error, loc,
				`a generic application's key must be "object." followed by an object name; at runtime this type resolves to "any"`)
		case alt.Malformed == webfunction.MalformedArgument:
			v.emit("generic-application-argument", Error, loc,
				`a generic application's value must be a valid type; at runtime this type resolves to "any"`)
		case alt.Arg != nil:
			if inObject && !nested {
				v.emit("generic-application-in-object", Error, loc,
					fmt.Sprintf("applies object.%s to a type inside an object definition, which is not allowed; type application is only valid in endpoint returns, argument types and attribute types", alt.Refinement))
			}
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
			v.genericsType(*alt.Arg, ctx, loc, inObject, true)
		case alt.Base == "array" && alt.Of != nil:
			v.genericsType(*alt.Of, ctx, loc, inObject, nested)
		case alt.IsObjectRef():
			if generics.IsGeneric(v.pkg.Object(alt.Refinement)) {
				v.emit("generic-unapplied", Error, loc,
					fmt.Sprintf(`references generic object.%s without applying a type argument; write {"object.%s": <type>} (at runtime its "T" falls back to "any")`, alt.Refinement, alt.Refinement))
			}
		}
	}
}