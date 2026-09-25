package validate

import "fmt"

// flagLevel maps every flag documented at
// https://webfunction.org/package#available-flags to the single level
// spec assigns it to: "package", "endpoint", "argument", or
// "attribute". A flag not present here at all is genuinely unrecognized
// (checkFlags treats that as info - could be a typo, or a legitimate
// newer flag this validator doesn't know about yet). A flag that IS
// present here but at a level OTHER than where it was found is a
// different, more serious case: spec states "A flag MUST only be used
// at its designated level ... and MUST NOT appear at any other level" -
// that's a known, unambiguous violation, not a guess, so checkFlags
// reports it as an error instead of lumping it in with genuinely
// unknown flags.
var flagLevel = map[string]string{
	"versioned":      "package",
	"package":        "endpoint",
	"event_source":   "endpoint",
	"error_triple":   "endpoint",
	"bearer_auth":    "endpoint",
	"capture_bearer": "endpoint",
	"paginated":      "endpoint",
	"private":        "endpoint",
	"required":       "argument",
	"nullable":       "attribute",
}

// checkFlags checks every flag string, at every level (package,
// endpoint, argument, attribute), against flagLevel.
func (v *validator) checkFlags() {
	v.checkFlagList(v.pkg.Flags, "package", PackageLoc())
	for _, e := range v.pkg.Endpoints {
		v.checkFlagList(e.Flags, "endpoint", EndpointLoc(e.Name))
		for _, a := range e.Arguments {
			v.checkFlagList(a.Flags, "argument", EndpointLoc(e.Name).Argument(a.Name))
		}
		for _, a := range e.Attributes {
			v.checkFlagList(a.Flags, "attribute", EndpointLoc(e.Name).Attribute(a.Name))
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			v.checkFlagList(a.Flags, "argument", ObjectLoc(o.Name).Argument(a.Name))
		}
		for _, a := range o.Attributes {
			v.checkFlagList(a.Flags, "attribute", ObjectLoc(o.Name).Attribute(a.Name))
		}
	}
}

// checkFlagList checks each flag in flags, found at the given level and
// location, against flagLevel.
func (v *validator) checkFlagList(flags []string, level string, loc Loc) {
	for _, f := range flags {
		properLevel, known := flagLevel[f]
		switch {
		case !known:
			v.emit("unrecognized-flag", Info, loc,
				fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
		case properLevel != level:
			v.emit("flag-wrong-level", Error, loc,
				fmt.Sprintf("flag %q is a valid flag, but only at %s level per spec - it MUST NOT appear at %s level", f, properLevel, level))
		}
	}
}