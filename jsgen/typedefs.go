package jsgen

import (
	"fmt"
	"sort"
	"strings"

	"wfn/webfunction"
)

// field is a generic name/type/docs tuple, used to build a typedef from
// either an endpoint's Attributes (return shape) or its Arguments (call
// shape). optional and nullable are deliberately separate: optional means
// the property may be absent entirely (JSDoc "[name]" bracket syntax) -
// true for a non-required Argument; nullable means the property is always
// present but its value may be null (part of the type union, e.g.
// "string|null") - true for an Attribute with the "nullable" flag.
// Conflating the two would make a nullable-but-always-present return
// field look absent to a type checker, which it isn't.
type field struct {
	name     string
	jsonType webfunction.JSONType
	optional bool
	nullable bool
	docs     string
}

func attributeFields(attrs []webfunction.Attribute) []field {
	fields := make([]field, len(attrs))
	for i, a := range attrs {
		docs := withNotes(a.Docs, hintNote(a.Hints), choicesNote(a.Values))
		fields[i] = field{name: a.Name, jsonType: a.Type, nullable: a.Nullable(), docs: docs}
	}
	return fields
}

func argumentFields(args []webfunction.Argument) []field {
	fields := make([]field, len(args))
	for i, a := range args {
		docs := withNotes(a.Docs, hintNote(a.Hints), choicesNote(a.Choices))
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
// endpoint return shapes, endpoint argument shapes, and (once those are
// known) a composite client typedef describing every method.
type typedefSet struct {
	ordered []*typedef
	bySig   map[string]*typedef
	names   map[string]bool
}

func newTypedefSet() *typedefSet {
	return &typedefSet{
		bySig: make(map[string]*typedef),
		names: make(map[string]bool),
	}
}

// forFields returns the typedef name describing this exact set of fields,
// creating a new typedef if no existing one has an identical shape.
// Returns "" if fields is empty (nothing to describe).
func (s *typedefSet) forFields(baseName string, fields []field) string {
	if len(fields) == 0 {
		return ""
	}

	sig := fieldSignature(fields)
	if existing, ok := s.bySig[sig]; ok {
		return existing.name
	}

	lines := make([]string, len(fields))
	for i, f := range fields {
		typ := jsdocType(f.jsonType, "", f.nullable)
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

	name := s.uniqueName(baseName)
	td := &typedef{name: name, lines: lines}
	s.bySig[sig] = td
	s.ordered = append(s.ordered, td)
	return name
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