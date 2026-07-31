package jsgen

import (
	"fmt"
	"sort"
	"strings"

	"wfn/webfunction"
)

// field is a generic name/type/docs tuple, used to build a typedef from
// either an endpoint's Attributes (return shape) or its Arguments (call
// shape). optional and nullable are deliberately separate concepts:
// optional means the property may be absent entirely (JSDoc "[name]"
// bracket syntax); nullable means its value may additionally be null (part
// of the type union, e.g. "string|null"). Per the spec, an Attribute's
// "nullable" flag actually means both at once - the key MAY be absent AND,
// when present, its value MAY be null - so attributeFields sets both.
// Argument's "required" flag only ever controls optional (arguments have
// no separate null concept in the spec).
type field struct {
	name     string
	jsonType webfunction.Type
	optional bool
	nullable bool
	docs     string
}

func attributeFields(attrs []webfunction.Attribute) []field {
	fields := make([]field, len(attrs))
	for i, a := range attrs {
		docs := withNotes(a.Docs, refinementNote(a.Type), choicesNote(a.Values))
		fields[i] = field{name: a.Name, jsonType: a.Type, optional: a.Nullable(), nullable: a.Nullable(), docs: docs}
	}
	return fields
}

func argumentFields(args []webfunction.Argument) []field {
	fields := make([]field, len(args))
	for i, a := range args {
		docs := withNotes(a.Docs, refinementNote(a.Type), choicesNote(a.Choices))
		fields[i] = field{name: a.Name, jsonType: a.Type, optional: !a.Required(), docs: docs}
	}
	return fields
}

// withNotes appends any non-empty parenthetical notes to a docs string.
func withNotes(docs string, notes ...string) string {
	docs = strings.TrimSpace(docs)
	for _, n := range notes {
		if n == "" {
			continue
		}
		if docs == "" {
			docs = n
		} else {
			docs = docs + " " + n
		}
	}
	return docs
}

func fieldSignature(fields []field) string {
	parts := make([]string, len(fields))
	for i, f := range fields {
		parts[i] = f.name + ":" + f.jsonType.String() + ":" + boolStr(f.optional) + ":" + boolStr(f.nullable)
	}
	sort.Strings(parts)
	return strings.Join(parts, "|")
}

// typedef is one generated JSDoc @typedef. lines holds each already-
// rendered "@property ..." (or similar) line, without the leading " * ".
type typedef struct {
	name  string
	lines []string
}

// typedefSet builds and dedupes the set of typedefs needed for a package:
// endpoint return shapes, endpoint argument shapes, resolved named object
// definitions (object.<name> refs), and a composite client typedef
// describing every method.
type typedefSet struct {
	pkg     *webfunction.Package
	ordered []*typedef
	bySig   map[string]*typedef
	names   map[string]bool
	// objects memoizes resolved object.<name> typedefs, keyed by
	// "<context>:<name>" - the same object name can resolve to two
	// different typedefs depending on whether it's referenced in an
	// argument or attribute context (per spec, each context uses a
	// different member set). This also guards against infinite recursion
	// on self-referential objects.
	objects map[string]string
}

func newTypedefSet(pkg *webfunction.Package) *typedefSet {
	return &typedefSet{
		pkg:     pkg,
		bySig:   make(map[string]*typedef),
		names:   make(map[string]bool),
		objects: make(map[string]string),
	}
}

// forFields returns the typedef name describing this exact set of fields,
// creating a new typedef if no existing one has an identical shape.
// Returns "" if fields is empty (nothing to describe). context ("argument"
// or "attribute") determines which member set is used to resolve any
// object.<name> references found within the fields' types.
func (s *typedefSet) forFields(baseName string, fields []field, context string) string {
	if len(fields) == 0 {
		return ""
	}

	sig := fieldSignature(fields)
	if existing, ok := s.bySig[sig]; ok {
		return existing.name
	}

	lines := s.renderFieldLines(fields, context)
	name := s.uniqueName(baseName)
	td := &typedef{name: name, lines: lines}
	s.bySig[sig] = td
	s.ordered = append(s.ordered, td)
	return name
}

// resolveObjectTypedef returns the typedef name for a named object
// definition (an object.<name> reference), building it (recursively - the
// context propagates into any nested object refs within it, per spec) if
// it hasn't been resolved yet in this context. Falls back to a generic
// "Object" if the referenced object doesn't exist, or has no members
// defined for this context.
func (s *typedefSet) resolveObjectTypedef(name, context string) string {
	key := context + ":" + name
	if existing, ok := s.objects[key]; ok {
		return existing
	}

	obj := s.pkg.Object(name)

	var fields []field
	suffix := "Attributes"
	if context == "argument" {
		suffix = "Args"
		if obj != nil {
			fields = argumentFields(obj.Arguments)
		}
	} else if obj != nil {
		fields = attributeFields(obj.Attributes)
	}

	if len(fields) == 0 {
		s.objects[key] = "Object"
		return "Object"
	}

	typedefName := s.uniqueName(pascalCase(name) + suffix)
	// Register before rendering lines, so a self-referential (or mutually
	// referential) object resolves to this name instead of recursing
	// forever.
	s.objects[key] = typedefName

	lines := s.renderFieldLines(fields, context)
	s.ordered = append(s.ordered, &typedef{name: typedefName, lines: lines})
	return typedefName
}

// renderFieldLines turns fields into fully-rendered "@property ..." lines,
// resolving any object.<name> references encountered via context.
func (s *typedefSet) renderFieldLines(fields []field, context string) []string {
	resolve := func(refName string) string { return s.resolveObjectTypedef(refName, context) }

	lines := make([]string, len(fields))
	for i, f := range fields {
		typ := jsdocType(f.jsonType, "", f.nullable, resolve)
		name := f.name
		if f.optional {
			name = "[" + name + "]"
		}
		line := "@property {" + typ + "} " + name
		if f.docs != "" {
			line += " - " + strings.ReplaceAll(strings.TrimSpace(f.docs), "\n", " ")
		}
		lines[i] = line
	}
	return lines
}

// addComposite adds a typedef built from already-rendered lines (e.g. the
// client-shape typedef, whose function-type properties don't fit the
// simple field model above). Unlike forFields, this is never deduped
// against an existing typedef - it always creates a new one.
func (s *typedefSet) addComposite(baseName string, lines []string) string {
	name := s.uniqueName(baseName)
	td := &typedef{name: name, lines: lines}
	s.ordered = append(s.ordered, td)
	return name
}

func (s *typedefSet) uniqueName(base string) string {
	name := base
	for i := 2; s.names[name]; i++ {
		name = fmt.Sprintf("%s%d", base, i)
	}
	s.names[name] = true
	return name
}

// render writes out every collected typedef as a JSDoc block, in the order
// they were created.
func (s *typedefSet) render() string {
	if len(s.ordered) == 0 {
		return ""
	}

	var b strings.Builder
	for _, td := range s.ordered {
		b.WriteString("/**\n")
		b.WriteString(" * @typedef {Object} " + td.name + "\n")
		for _, line := range td.lines {
			b.WriteString(" * " + line + "\n")
		}
		b.WriteString(" */\n\n")
	}
	return b.String()
}