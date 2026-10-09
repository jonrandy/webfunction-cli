// Package generics expands Web Function generic object applications into
// ordinary, concrete objects, so the converters and codegen targets - which
// walk a package's types directly in many places - can be handed a package
// that contains no generics at all.
//
// A generic object is an object definition whose member types use the bare
// token "T" (in its arguments, its attributes, or both). A generic is
// applied to a type wherever a type is written outside an object
// definition: { "object.list": "object.user" }. Expand replaces every
// application with a reference to a generated object (here "ListOfUser")
// whose members have "T" substituted, and drops the generic templates
// themselves from the result.
//
// Expand never fails: it applies the spec's runtime fallbacks, so an
// unapplied generic resolves "T" to "any", an argument applied to an
// object without "T" is ignored, and a malformed application stays "any".
// Reporting those as errors is the validator's job, not this package's.
package generics

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/webfunction-protocol/webfunction-go"
)

// maxDepth bounds how deeply instantiations may nest. Valid packages never
// get near it; it exists so invalid polymorphic recursion (a generic that
// applies itself to something larger than T) terminates, with the
// over-deep application resolving to "any".
const maxDepth = 32

// IsGeneric reports whether the object uses the type parameter "T" in any
// member type of either context.
func IsGeneric(o *webfunction.Object) bool {
	return UsesParam(o, webfunction.ArgumentContext) || UsesParam(o, webfunction.AttributeContext)
}

// UsesParam reports whether the object uses the type parameter "T" in the
// member types of one specific context.
func UsesParam(o *webfunction.Object, ctx webfunction.ObjectContext) bool {
	if o == nil {
		return false
	}
	if ctx == webfunction.ArgumentContext {
		for _, a := range o.Arguments {
			if hasParam(a.Type) {
				return true
			}
		}
		return false
	}
	for _, a := range o.Attributes {
		if hasParam(a.Type) {
			return true
		}
	}
	return false
}

// Collision records a generated object whose natural name was already
// taken by a user-defined object. Expand renames the generated object
// (with a numeric suffix) so output stays valid; Collisions lets a
// validator report the clash.
type Collision struct {
	// Name is the clashing object name.
	Name string
	// Application is the canonical form of the application that generated
	// it, e.g. "list<object.user>".
	Application string
}

// Collisions returns every generated-object name that clashes with a
// user-defined object in pkg, in the order they were encountered.
func Collisions(pkg *webfunction.Package) []Collision {
	_, e := expand(pkg)
	return e.collisions
}

func isParam(a webfunction.TypeAlt) bool {
	return a.Base == "T" && a.Refinement == "" && a.Of == nil && a.Arg == nil
}

func hasParam(t webfunction.Type) bool {
	for _, a := range t.Union {
		if isParam(a) {
			return true
		}
		if a.Of != nil && hasParam(*a.Of) {
			return true
		}
		if a.Arg != nil && hasParam(*a.Arg) {
			return true
		}
	}
	return false
}

// Expand returns a copy of pkg in which every generic application has been
// replaced by a reference to a concrete generated object, and the generic
// templates removed. pkg itself is never mutated. A package with no
// generics comes back unchanged.
//
// Generated objects are appended after the package's own objects, in the
// order they were first needed (an application's argument is expanded
// before the application itself, so inner instantiations come first).
func Expand(pkg *webfunction.Package) *webfunction.Package {
	out, _ := expand(pkg)
	return out
}

func expand(pkg *webfunction.Package) (*webfunction.Package, *expander) {
	e := &expander{
		templates: map[string]*webfunction.Object{},
		names:     map[string]bool{},
		userNames: map[string]bool{},
		instances: map[string]string{},
	}
	for i := range pkg.Objects {
		o := &pkg.Objects[i]
		e.names[o.Name] = true
		e.userNames[o.Name] = true
		if IsGeneric(o) {
			e.templates[o.Name] = o
		}
	}

	out := *pkg

	if pkg.Endpoints != nil {
		out.Endpoints = make([]webfunction.Endpoint, len(pkg.Endpoints))
		for i, ep := range pkg.Endpoints {
			ep.Returns = e.expandType(ep.Returns, 0)
			ep.Arguments = e.arguments(ep.Arguments, 0)
			ep.Attributes = e.attributes(ep.Attributes, 0)
			out.Endpoints[i] = ep
		}
	}

	objs := make([]webfunction.Object, 0, len(pkg.Objects)+len(e.out))
	for _, o := range pkg.Objects {
		if e.templates[o.Name] != nil {
			continue
		}
		o.Arguments = e.arguments(o.Arguments, 0)
		o.Attributes = e.attributes(o.Attributes, 0)
		objs = append(objs, o)
	}
	objs = append(objs, e.out...)
	if pkg.Objects == nil && len(objs) == 0 {
		out.Objects = nil
	} else {
		out.Objects = objs
	}
	return &out, e
}

type expander struct {
	templates  map[string]*webfunction.Object // generic objects by name
	names      map[string]bool                // every object name in use
	userNames  map[string]bool                // names of the package's own objects
	instances  map[string]string              // canonical application -> generated name
	out        []webfunction.Object           // generated objects, in completion order
	collisions []Collision                    // generated names that clashed with userNames
}

func (e *expander) arguments(in []webfunction.Argument, depth int) []webfunction.Argument {
	if in == nil {
		return nil
	}
	out := make([]webfunction.Argument, len(in))
	for i, a := range in {
		a.Type = e.expandType(a.Type, depth)
		out[i] = a
	}
	return out
}

func (e *expander) attributes(in []webfunction.Attribute, depth int) []webfunction.Attribute {
	if in == nil {
		return nil
	}
	out := make([]webfunction.Attribute, len(in))
	for i, a := range in {
		a.Type = e.expandType(a.Type, depth)
		out[i] = a
	}
	return out
}

// expandType replaces every application inside t (and every unapplied
// reference to a generic object) with a reference to a generated object.
func (e *expander) expandType(t webfunction.Type, depth int) webfunction.Type {
	if t.Union == nil {
		return t
	}
	alts := make([]webfunction.TypeAlt, 0, len(t.Union))
	for _, a := range t.Union {
		switch {
		case a.Arg != nil:
			alts = append(alts, e.apply(a, depth))
		case a.Base == "array" && a.Of != nil:
			of := e.expandType(*a.Of, depth)
			a.Of = &of
			alts = append(alts, a)
		case a.IsObjectRef() && e.templates[a.Refinement] != nil:
			// Unapplied generic: "T" falls back to "any".
			anyType := webfunction.Type{Union: []webfunction.TypeAlt{{Base: "any"}}}
			alts = append(alts, e.instantiate(e.templates[a.Refinement], anyType, depth))
		default:
			alts = append(alts, a)
		}
	}
	out := t
	out.Union = alts
	return out
}

// apply resolves one application alternative.
func (e *expander) apply(a webfunction.TypeAlt, depth int) webfunction.TypeAlt {
	tmpl := e.templates[a.Refinement]
	if tmpl == nil {
		// Not generic (or not defined): the argument is ignored.
		return webfunction.TypeAlt{Base: "object", Refinement: a.Refinement}
	}
	arg := e.expandType(*a.Arg, depth+1) // innermost application first
	return e.instantiate(tmpl, arg, depth)
}

// instantiate returns a reference to the generated object for tmpl applied
// to arg (already expanded), creating it on first use.
func (e *expander) instantiate(tmpl *webfunction.Object, arg webfunction.Type, depth int) webfunction.TypeAlt {
	if depth >= maxDepth {
		return webfunction.TypeAlt{Base: "any"}
	}
	key := tmpl.Name + "<" + arg.String() + ">"
	if name, ok := e.instances[key]; ok {
		return ref(name)
	}
	name := e.uniqueName(pascal(tmpl.Name)+"Of"+typeName(arg), key)
	e.instances[key] = name

	obj := webfunction.Object{Name: name}
	if tmpl.Arguments != nil {
		args := make([]webfunction.Argument, len(tmpl.Arguments))
		for i, a := range tmpl.Arguments {
			a.Type = subst(a.Type, arg)
			args[i] = a
		}
		obj.Arguments = e.arguments(args, depth+1)
	}
	if tmpl.Attributes != nil {
		attrs := make([]webfunction.Attribute, len(tmpl.Attributes))
		for i, a := range tmpl.Attributes {
			a.Type = subst(a.Type, arg)
			attrs[i] = a
		}
		obj.Attributes = e.attributes(attrs, depth+1)
	}
	e.out = append(e.out, obj)
	return ref(name)
}

func ref(name string) webfunction.TypeAlt {
	return webfunction.TypeAlt{Base: "object", Refinement: name}
}

// uniqueName reserves and returns base, or base with a numeric suffix if
// the name is already taken (by a user-defined object, or by another
// instantiation that flattens to the same name).
func (e *expander) uniqueName(base, application string) string {
	if e.userNames[base] {
		e.collisions = append(e.collisions, Collision{Name: base, Application: application})
	}
	name := base
	for i := 2; e.names[name]; i++ {
		name = base + strconv.Itoa(i)
	}
	e.names[name] = true
	return name
}

// subst replaces the type parameter in t with arg. A bare "T" becomes the
// whole argument; "T" inside a union or an array's element list is spliced
// in as union members, so an argument that is itself a union never ends up
// nested inside another one. Duplicate members are dropped.
func subst(t webfunction.Type, arg webfunction.Type) webfunction.Type {
	if !hasParam(t) {
		return t
	}
	var alts []webfunction.TypeAlt
	seen := map[string]bool{}
	add := func(a webfunction.TypeAlt) {
		if k := a.String(); !seen[k] {
			seen[k] = true
			alts = append(alts, a)
		}
	}
	for _, a := range t.Union {
		switch {
		case isParam(a):
			for _, x := range arg.Union {
				add(x)
			}
		case a.Of != nil:
			of := subst(*a.Of, arg)
			a.Of = &of
			add(a)
		case a.Arg != nil:
			x := subst(*a.Arg, arg)
			a.Arg = &x
			add(a)
		default:
			add(a)
		}
	}
	out := t
	out.Union = alts
	return out
}

// typeName flattens a type into the name fragment used for generated
// objects: alternatives joined by "Or", arrays as "ArrayOf<element>".
func typeName(t webfunction.Type) string {
	parts := make([]string, 0, len(t.Union))
	for _, a := range t.Union {
		parts = append(parts, altName(a))
	}
	return strings.Join(parts, "Or")
}

func altName(a webfunction.TypeAlt) string {
	switch {
	case a.Base == "array":
		if a.Of != nil {
			return "ArrayOf" + typeName(*a.Of)
		}
		return "Array"
	case a.IsObjectRef():
		return pascal(a.Refinement)
	default:
		n := pascal(a.Base)
		if a.Refinement != "" {
			n += pascal(a.Refinement)
		}
		return n
	}
}

// pascal upper-cases the first letter of each alphanumeric run and drops
// everything else ("user_profile" -> "UserProfile"). Letters after the
// first in a run are left alone, so "ListOfUser" is already Pascal case.
func pascal(s string) string {
	var b strings.Builder
	up := true
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			up = true
			continue
		}
		if up {
			r = unicode.ToUpper(r)
			up = false
		}
		b.WriteRune(r)
	}
	return b.String()
}