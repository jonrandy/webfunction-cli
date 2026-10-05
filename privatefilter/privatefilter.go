// Package privatefilter strips members flagged "private" out of a
// package, so the converters and codegen targets - which otherwise
// walk a package's arguments and attributes directly in many places -
// can be handed a package that simply no longer contains them.
//
// Per https://webfunction.org/package, "private" (internal use only)
// is valid on endpoints, arguments, and attributes. Private
// *endpoints* are excluded by each consumer's own existing check; this
// package handles the member level: arguments and attributes of every
// endpoint, and of every object definition.
package privatefilter

import "github.com/webfunction-protocol/webfunction-go"

// Apply returns a copy of pkg with every private argument and
// attribute removed, from endpoints and from object definitions alike.
// pkg itself is never mutated. Private endpoints are left in place -
// callers already skip them via Endpoint.Private().
//
// Note a side effect worth knowing: an object whose attributes (or
// arguments) are *all* private ends up with none, so it resolves as an
// object with no members in that context, exactly like an object that
// never defined any.
func Apply(pkg *webfunction.Package) *webfunction.Package {
	out := *pkg

	out.Endpoints = make([]webfunction.Endpoint, len(pkg.Endpoints))
	for i, e := range pkg.Endpoints {
		e.Arguments = publicArguments(e.Arguments)
		e.Attributes = publicAttributes(e.Attributes)
		out.Endpoints[i] = e
	}

	if pkg.Objects != nil {
		out.Objects = make([]webfunction.Object, len(pkg.Objects))
		for i, o := range pkg.Objects {
			o.Arguments = publicArguments(o.Arguments)
			o.Attributes = publicAttributes(o.Attributes)
			out.Objects[i] = o
		}
	}
	return &out
}

func publicArguments(args []webfunction.Argument) []webfunction.Argument {
	if len(args) == 0 {
		return args
	}
	out := make([]webfunction.Argument, 0, len(args))
	for _, a := range args {
		if !a.Private() {
			out = append(out, a)
		}
	}
	return out
}

func publicAttributes(attrs []webfunction.Attribute) []webfunction.Attribute {
	if len(attrs) == 0 {
		return attrs
	}
	out := make([]webfunction.Attribute, 0, len(attrs))
	for _, a := range attrs {
		if !a.Private() {
			out = append(out, a)
		}
	}
	return out
}