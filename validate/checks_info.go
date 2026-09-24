package validate

import "fmt"

// knownPackageFlags/knownEndpointFlags/knownArgumentFlags/
// knownAttributeFlags are the flags documented at
// https://webfunction.org/package#available-flags, scoped exactly to
// the level the spec assigns them. A flag outside its level's known set
// could be a typo, or a legitimate newer flag this validator doesn't
// know about yet - so this is an info-level note either way, not a hard
// failure.
var (
	knownPackageFlags  = map[string]bool{"versioned": true}
	knownEndpointFlags = map[string]bool{
		"paginated": true, "bearer_auth": true, "private": true,
		"error_triple": true, "capture_bearer": true,
		"package": true, "event_source": true,
	}
	knownArgumentFlags  = map[string]bool{"required": true}
	knownAttributeFlags = map[string]bool{"nullable": true}
)

// checkFlags flags any flag string, at package/endpoint/argument/
// attribute level, outside the known set for that level.
func (v *validator) checkFlags() {
	for _, f := range v.pkg.Flags {
		if !knownPackageFlags[f] {
			v.emit("unrecognized-flag", Info, PackageLoc(),
				fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
		}
	}
	for _, e := range v.pkg.Endpoints {
		for _, f := range e.Flags {
			if !knownEndpointFlags[f] {
				v.emit("unrecognized-flag", Info, EndpointLoc(e.Name),
					fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
			}
		}
		for _, a := range e.Arguments {
			for _, f := range a.Flags {
				if !knownArgumentFlags[f] {
					v.emit("unrecognized-flag", Info, EndpointLoc(e.Name).Argument(a.Name),
						fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
				}
			}
		}
		for _, a := range e.Attributes {
			for _, f := range a.Flags {
				if !knownAttributeFlags[f] {
					v.emit("unrecognized-flag", Info, EndpointLoc(e.Name).Attribute(a.Name),
						fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
				}
			}
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			for _, f := range a.Flags {
				if !knownArgumentFlags[f] {
					v.emit("unrecognized-flag", Info, ObjectLoc(o.Name).Argument(a.Name),
						fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
				}
			}
		}
		for _, a := range o.Attributes {
			for _, f := range a.Flags {
				if !knownAttributeFlags[f] {
					v.emit("unrecognized-flag", Info, ObjectLoc(o.Name).Attribute(a.Name),
						fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
				}
			}
		}
	}
}