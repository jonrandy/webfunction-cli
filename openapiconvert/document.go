package openapiconvert

// This file defines the (small) subset of the OpenAPI 3.1 object model
// this package emits. Struct field order is preserved by encoding/json,
// so field declaration order below is the order keys appear in the
// output document; map keys (Paths, Properties, Schemas, ...) are sorted
// alphabetically by encoding/json itself, which is fine - OpenAPI doesn't
// care about property order.

// Document is the top-level OpenAPI 3.1 document.
type Document struct {
	OpenAPI    string               `json:"openapi"`
	Info       Info                 `json:"info"`
	Servers    []Server             `json:"servers"`
	Paths      map[string]*PathItem `json:"paths"`
	Components *Components          `json:"components,omitempty"`

	// x-webfunction-* vendor extensions, package-level. Endpoint-level
	// ones live on Operation instead.
	XWebfunctionVersioned     bool       `json:"x-webfunction-versioned,omitempty"`
	XWebfunctionPackageErrors []ErrorDoc `json:"x-webfunction-package-errors,omitempty"`
}

// Info is the OpenAPI Info Object.
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
	// x-webfunction-versions lives here rather than at the document root,
	// mirroring the reference JS converter's placement.
	XWebfunctionVersions []string `json:"x-webfunction-versions,omitempty"`
}

// Server is the OpenAPI Server Object.
type Server struct {
	URL string `json:"url"`
}

// PathItem is the OpenAPI Path Item Object. Web Function endpoints are
// always POST-only, so Post is the only operation ever populated.
type PathItem struct {
	Post *Operation `json:"post,omitempty"`
}

// Operation is the OpenAPI Operation Object for one endpoint.
type Operation struct {
	OperationID string                `json:"operationId"`
	Summary     string                `json:"summary"`
	Tags        []string              `json:"tags,omitempty"`
	Description string                `json:"description,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody"`
	Responses   map[string]*Response  `json:"responses"`
	Security    []SecurityRequirement `json:"security,omitempty"`

	// x-webfunction-* vendor extensions, per endpoint.
	XWebfunctionErrorCodes     []string `json:"x-webfunction-error-codes,omitempty"`
	XWebfunctionPaginated      bool     `json:"x-webfunction-paginated,omitempty"`
	XWebfunctionReturnsPackage bool     `json:"x-webfunction-returns-package,omitempty"`
	XWebfunctionEventSource    bool     `json:"x-webfunction-event-source,omitempty"`
	XWebfunctionCaptureBearer  bool     `json:"x-webfunction-capture-bearer,omitempty"`
}

// Parameter is the OpenAPI Parameter Object - used here only for the
// "Api-Version" header a versioned package's endpoints all require.
type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Required    bool    `json:"required,omitempty"`
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

// RequestBody is the OpenAPI Request Body Object.
type RequestBody struct {
	Required bool                 `json:"required"`
	Content  map[string]MediaType `json:"content"`
}

// Response is the OpenAPI Response Object.
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// MediaType is the OpenAPI Media Type Object.
type MediaType struct {
	Schema *Schema `json:"schema"`
}

// SecurityRequirement is the OpenAPI Security Requirement Object, e.g.
// {"bearerAuth": []}.
type SecurityRequirement map[string][]string

// Components is the OpenAPI Components Object.
type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	// Schemas holds both real object.<n>-derived *Schema entries and the
	// one hand-built ErrorTriple schema (a map[string]any literal, since
	// it needs "items": false - a JSON Schema 2020-12 boolean schema -
	// which the general Schema struct doesn't model).
	Schemas map[string]any `json:"schemas,omitempty"`
}

// SecurityScheme is the OpenAPI Security Scheme Object - only the bearer
// HTTP scheme is ever emitted here.
type SecurityScheme struct {
	Type   string `json:"type"`
	Scheme string `json:"scheme"`
}

// ErrorDoc mirrors webfunction.ErrorDef's JSON shape, reused directly for
// the x-webfunction-package-errors extension.
type ErrorDoc struct {
	Code string `json:"code"`
	Docs string `json:"docs,omitempty"`
}

// Schema is the (recursive) subset of JSON Schema 2020-12 / OpenAPI 3.1
// Schema Object this package needs.
type Schema struct {
	Ref                  string             `json:"$ref,omitempty"`
	Type                 any                `json:"type,omitempty"` // string or []string
	Format               string             `json:"format,omitempty"`
	Description          string             `json:"description,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty"`
	Required             []string           `json:"required,omitempty"`
	AdditionalProperties *bool              `json:"additionalProperties,omitempty"`
	Items                *Schema            `json:"items,omitempty"`
	Enum                 []any              `json:"enum,omitempty"`
	OneOf                []*Schema          `json:"oneOf,omitempty"`
	Minimum              *float64           `json:"minimum,omitempty"`
}

func boolPtr(b bool) *bool { return &b }

func floatPtr(f float64) *float64 { return &f }