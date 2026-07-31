// Package webfunction provides a minimal Go representation of the Web
// Function package/endpoint schema (see https://webfunction.org/package),
// along with a function to fetch one over the wire.
//
// This is intentionally kept small and self-contained. If/when a proper
// "webfunction-go" client library exists, the intent is to swap this
// package out for that (or have this package wrap it) with minimal change
// to callers - so keep the exported surface here narrow and stable.
package webfunction

// Package is the top-level Web Function package definition.
// See https://webfunction.org/package#package-definition.
type Package struct {
	BaseURL        string     `json:"base_url"`
	EventSourceURL string     `json:"event_source_url,omitempty"`
	PipelineURL    string     `json:"pipeline_url,omitempty"`
	Name           string     `json:"name,omitempty"`
	Flags          []string   `json:"flags,omitempty"`
	Version        string     `json:"version,omitempty"`
	Versions       []string   `json:"versions,omitempty"`
	Docs           string     `json:"docs,omitempty"`
	Endpoints      []Endpoint `json:"endpoints"`
	Events         []Event    `json:"events,omitempty"`
	Errors         []ErrorDef `json:"errors,omitempty"`
	Objects        []Object   `json:"objects,omitempty"`
}

// Versioned reports whether the package declares the "versioned" flag.
func (p *Package) Versioned() bool {
	return hasFlag(p.Flags, "versioned")
}

// Object returns the named object definition, or nil if the package has
// none by that name.
func (p *Package) Object(name string) *Object {
	for i := range p.Objects {
		if p.Objects[i].Name == name {
			return &p.Objects[i]
		}
	}
	return nil
}

// Endpoint describes a single callable endpoint within a package.
// See https://webfunction.org/package#endpoint-definition.
type Endpoint struct {
	Name       string      `json:"name"`
	Returns    Type        `json:"returns"`
	Flags      []string    `json:"flags,omitempty"`
	Group      string      `json:"group,omitempty"`
	Docs       string      `json:"docs,omitempty"`
	Errors     []ErrorDef  `json:"errors,omitempty"`
	Arguments  []Argument  `json:"arguments"`
	Attributes []Attribute `json:"attributes,omitempty"`
}

// HasFlag reports whether the endpoint declares the given flag
// (e.g. "package", "paginated", "bearer_auth", "private", "error_triple",
// "capture_bearer", "event_source").
func (e *Endpoint) HasFlag(flag string) bool {
	return hasFlag(e.Flags, flag)
}

// Argument returns the named argument, or nil if the endpoint has none by
// that name.
func (e *Endpoint) Argument(name string) *Argument {
	for i := range e.Arguments {
		if e.Arguments[i].Name == name {
			return &e.Arguments[i]
		}
	}
	return nil
}

// Argument describes a single argument accepted by an endpoint.
// See https://webfunction.org/package#argument-definition.
type Argument struct {
	Name    string        `json:"name"`
	Type    Type          `json:"type"`
	Group   string        `json:"group,omitempty"`
	Choices []interface{} `json:"choices,omitempty"`
	Flags   []string      `json:"flags,omitempty"`
	Docs    string        `json:"docs,omitempty"`
}

// Required reports whether the argument declares the "required" flag.
func (a *Argument) Required() bool {
	return hasFlag(a.Flags, "required")
}

// Attribute describes a single attribute of an object returned by an
// endpoint (or carried by an event).
// See https://webfunction.org/package#attribute-definition.
type Attribute struct {
	Name   string        `json:"name"`
	Type   Type          `json:"type"`
	Values []interface{} `json:"values,omitempty"`
	Flags  []string      `json:"flags,omitempty"`
	Docs   string        `json:"docs,omitempty"`
}

// Nullable reports whether the attribute declares the "nullable" flag.
//
// Per spec this means more than just "the value may be null": the key MAY
// be absent from the object entirely, and when present its value MAY be
// null. Consumers SHOULD treat a missing key and a null value equivalently.
func (a *Attribute) Nullable() bool {
	return hasFlag(a.Flags, "nullable")
}

// Event describes a real-time event a package can emit.
// See https://webfunction.org/package#event-definition.
type Event struct {
	Name       string      `json:"name"`
	Group      string      `json:"group,omitempty"`
	Docs       string      `json:"docs,omitempty"`
	Attributes []Attribute `json:"attributes"`
}

// ErrorDef describes a single named ERROR_CODE a package or endpoint can
// return. Named ErrorDef (not Error) to avoid colliding with Go's built-in
// error interface.
// See https://webfunction.org/package#error-definition.
type ErrorDef struct {
	Code string `json:"code"`
	Docs string `json:"docs,omitempty"`
}

func hasFlag(flags []string, want string) bool {
	for _, f := range flags {
		if f == want {
			return true
		}
	}
	return false
}