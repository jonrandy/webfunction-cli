package validate

import (
	"fmt"
	"strings"
)

// flagLevel maps every flag documented at
// https://webfunction.org/package#available-flags to the level(s) spec
// allows it at: "package", "endpoint", "argument", or "attribute". Most
// flags have exactly one; "private" is allowed at endpoint, argument,
// and attribute level. A flag not present here at all is genuinely unrecognized
// (checkFlags treats that as info - could be a typo, or a legitimate
// newer flag this validator doesn't know about yet). A flag that IS
// present here but at a level OTHER than where it was found is a
// different, more serious case: spec states "A flag MUST only be used
// at its designated level(s) ... and MUST NOT appear at any other level" -
// that's a known, unambiguous violation, not a guess, so checkFlags
// reports it as an error instead of lumping it in with genuinely
// unknown flags.
var flagLevel = map[string][]string{
	"versioned":      {"package"},
	"package":        {"endpoint"},
	"event_source":   {"endpoint"},
	"error_triple":   {"endpoint"},
	"bearer_auth":    {"endpoint"},
	"capture_bearer": {"endpoint"},
	"paginated":      {"endpoint"},
	"private":        {"endpoint", "argument", "attribute"},
	"required":       {"argument"},
	"nullable":       {"attribute"},
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
		properLevels, known := flagLevel[f]
		switch {
		case !known:
			v.emit("unrecognized-flag", Info, loc,
				fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
		case !levelAllowed(properLevels, level):
			v.emit("flag-wrong-level", Error, loc,
				fmt.Sprintf("flag %q is a valid flag, but only at %s level per spec - it MUST NOT appear at %s level", f, describeLevels(properLevels), level))
		}
	}
}

// levelAllowed reports whether level is one of levels.
func levelAllowed(levels []string, level string) bool {
	for _, l := range levels {
		if l == level {
			return true
		}
	}
	return false
}

// describeLevels renders levels for a human: "endpoint", "endpoint or
// argument", "endpoint, argument, or attribute".
func describeLevels(levels []string) string {
	switch len(levels) {
	case 1:
		return levels[0]
	case 2:
		return levels[0] + " or " + levels[1]
	}
	return strings.Join(levels[:len(levels)-1], ", ") + ", or " + levels[len(levels)-1]
}

// hasFlag reports whether flags contains flag.
func hasFlag(flags []string, flag string) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}