package validate

import "fmt"

// knownPackageFlags/knownEndpointFlags/knownArgumentFlags are the flags
// referenced anywhere in this project so far. A flag outside these sets
// could be a typo, or could be a legitimate newer flag this project just
// doesn't know about yet - so this is an info-level note either way, not
// a hard failure. "package" and "event_source" are themselves
// unconfirmed (see checkSpecAmbiguities) but are included here as
// "known-but-ambiguous" rather than "unknown", since flagging them AGAIN
// as merely-unrecognized would be redundant with the dedicated note.
var (
	knownPackageFlags  = map[string]bool{"versioned": true}
	knownEndpointFlags = map[string]bool{
		"paginated": true, "bearer_auth": true, "private": true,
		"error_triple": true, "capture_bearer": true, "versioned": true,
		"package": true, "event_source": true,
	}
	knownArgumentFlags = map[string]bool{"required": true}
)

// checkFlags flags any flag string, at package/endpoint/argument level,
// outside the known sets above.
func (v *validator) checkFlags() {
	for _, f := range v.pkg.Flags {
		if !knownPackageFlags[f] {
			v.emit("unrecognized-flag", Info, "package",
				fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
		}
	}
	for _, e := range v.pkg.Endpoints {
		for _, f := range e.Flags {
			if !knownEndpointFlags[f] {
				v.emit("unrecognized-flag", Info, fmt.Sprintf("endpoint %q", e.Name),
					fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
			}
		}
		for _, a := range e.Arguments {
			for _, f := range a.Flags {
				if !knownArgumentFlags[f] {
					v.emit("unrecognized-flag", Info, fmt.Sprintf("endpoint %q argument %q", e.Name, a.Name),
						fmt.Sprintf("unrecognized flag %q - could be a typo, or a newer flag this validator doesn't know about yet", f))
				}
			}
		}
	}
}

// checkSpecAmbiguities emits the two open spec questions this project
// has never resolved, as info-level notes rather than silently assuming
// an answer either way - deliberately data-driven (only raised when the
// package actually exercises the ambiguous area), not blanket
// disclaimers on every report.
func (v *validator) checkSpecAmbiguities() {
	if len(v.pkg.Errors) > 0 {
		v.emit("ambiguous-package-error-inheritance", Info, "package",
			`package declares package-level "errors"; whether these implicitly apply to every endpoint or only an endpoint's own declared "errors" is unconfirmed by anything read from the spec so far - this validator does not check endpoint @throws-style completeness against package-level errors because of that ambiguity`)
	}

	eventSource := false
	for _, e := range v.pkg.Endpoints {
		if e.HasFlag("event_source") {
			eventSource = true
			break
		}
	}
	if eventSource || v.pkg.PipelineURL != "" {
		v.emit("ambiguous-events-concept", Info, "package",
			`package references "event_source" and/or a pipeline URL; whether "events"/event_source_url/pipeline-based delivery are real, still-unimplemented spec concepts or an early over-read of the spec site is unconfirmed - webfunction-go's Package type deliberately has no structural fields for "events" because of this`)
	}
}