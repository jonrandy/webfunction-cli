package jsgen

import (
	"fmt"
	"strings"

	"wfn/webfunction"
)

// typedef is one generated JSDoc @typedef, shared by every endpoint whose
// object return shape matches it exactly.
type typedef struct {
	name       string
	attributes []webfunction.Attribute
}

// typedefSet builds the set of typedefs needed for a package's endpoints,
// and a lookup from endpoint name to the typedef name that describes its
// return shape (only present for endpoints that return an object with a
// known attribute list).
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

// forEndpoint returns the typedef name for the endpoint's return shape,
// creating a new typedef if no existing one has an identical shape, or ""
// if the endpoint doesn't return a documented object shape.
func (s *typedefSet) forEndpoint(ep webfunction.Endpoint) string {
	if !contains(ep.Returns, "object") || len(ep.Attributes) == 0 {
		return ""
	}

	sig := attributeSignature(ep.Attributes)
	if existing, ok := s.bySig[sig]; ok {
		return existing.name
	}

	name := s.uniqueName(pascalCase(ep.Name) + "Result")
	td := &typedef{name: name, attributes: ep.Attributes}
	s.bySig[sig] = td
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

// render writes out every collected typedef as a JSDoc block.
func (s *typedefSet) render() string {
	if len(s.ordered) == 0 {
		return ""
	}

	var b strings.Builder
	for _, td := range s.ordered {
		b.WriteString("/**\n")
		b.WriteString(" * @typedef {Object} " + td.name + "\n")
		for _, attr := range td.attributes {
			typ := jsdocType(attr.Type, "", attr.Nullable())
			line := " * @property {" + typ + "} " + attr.Name
			if attr.Docs != "" {
				line += " - " + strings.ReplaceAll(strings.TrimSpace(attr.Docs), "\n", " ")
			}
			b.WriteString(line + "\n")
		}
		b.WriteString(" */\n\n")
	}
	return b.String()
}