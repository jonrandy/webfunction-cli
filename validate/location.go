package validate

import "fmt"

// Loc identifies where in a package a Finding was found, structured
// enough to group findings by subject (e.g. one endpoint) in the
// Markdown/HTML reports, while still rendering to a single human-
// readable string for Finding.Location and JSON consumers who don't
// want to deal with the structured form at all.
type Loc struct {
	// Scope is "package", "endpoint", or "object".
	Scope string
	// Subject is the endpoint or object name. Empty for package scope.
	Subject string
	// Part is "", "returns", "argument", or "attribute" - which part of
	// Subject this finding is about. Empty means the subject as a whole
	// (e.g. a duplicate name, a flag on the endpoint/package itself).
	Part string
	// Field is the argument/attribute name. Empty unless Part is
	// "argument" or "attribute".
	Field string
}

// PackageLoc is the location for a package-level finding.
func PackageLoc() Loc { return Loc{Scope: "package"} }

// EndpointLoc is the location for a finding about an endpoint as a
// whole (e.g. its name, its own flags).
func EndpointLoc(name string) Loc { return Loc{Scope: "endpoint", Subject: name} }

// ObjectLoc is the location for a finding about an object definition as
// a whole (e.g. its name).
func ObjectLoc(name string) Loc { return Loc{Scope: "object", Subject: name} }

// Returns narrows an endpoint location to its return type.
func (l Loc) Returns() Loc { l.Part = "returns"; return l }

// Argument narrows a location to one of the subject's arguments.
func (l Loc) Argument(name string) Loc { l.Part = "argument"; l.Field = name; return l }

// Attribute narrows a location to one of the subject's attributes.
func (l Loc) Attribute(name string) Loc { l.Part = "attribute"; l.Field = name; return l }

// GroupKey is what report rendering groups findings by: everything
// about the same endpoint or object lands in the same group regardless
// of which specific Part/Field triggered each individual finding.
func (l Loc) GroupKey() string { return l.Scope + ":" + l.Subject }

// String renders a short human-readable location, e.g.
// `endpoint "list-people" argument "status"`. This is what populates
// Finding.Location.
func (l Loc) String() string {
	switch l.Scope {
	case "package":
		return "package"
	case "endpoint", "object":
		base := fmt.Sprintf("%s %q", l.Scope, l.Subject)
		switch l.Part {
		case "returns":
			return base + " returns"
		case "argument":
			return base + fmt.Sprintf(" argument %q", l.Field)
		case "attribute":
			return base + fmt.Sprintf(" attribute %q", l.Field)
		default:
			return base
		}
	default:
		return ""
	}
}