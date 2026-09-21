// Package openapiconvert converts a webfunction.Package into an OpenAPI
// 3.1 document.
//
// Written from scratch against webfunction-go's real parsed Package model
// (not a port of the reference JS converter, webfunction-to-openapi, which
// operates on an older/looser reading of the wire format - notably it
// never resolves object.<n> refs, wraps a paginated response's item shape
// as if it were the whole response, and reads event_source_url/events/
// pipeline_url fields webfunction-go's Package deliberately doesn't
// expose). The reference's per-refinement format/description mapping and
// operationId casing were still directly useful and are ported as-is; see
// refinementSchema (schema.go) and toOperationId (names.go).
//
// Design decisions (confirmed with Jon before writing):
//  1. object.<n> refs become real, named, reusable, cycle-safe
//     components.schemas entries with $ref - not inlined - matching every
//     other codegen target in this project (jsgen/phpgen/gogen).
//  2. A paginated endpoint's response schema is the real wire envelope
//     ({"page": [...items], "next": ..., "previous": ...}, confirmed
//     against webfunction-go's page.go), not the bare item shape the
//     reference emits.
//  3. Private endpoints are excluded from the output by default; the
//     convert command's --private flag (shared with postmanconvert) opts
//     them back in.
package openapiconvert

import (
	"encoding/json"
	"fmt"

	"github.com/webfunction-protocol/webfunction-go"
)

// Generate converts pkg into an OpenAPI 3.1 document and returns it as
// indented JSON.
func Generate(pkg *webfunction.Package, includePrivate bool) (string, error) {
	resolver := newSchemaSet(pkg)

	title := pkg.Name
	if title == "" {
		title = "Web Function API"
	}
	version := pkg.Version
	if version == "" {
		version = "1.0.0"
	}

	doc := &Document{
		OpenAPI: "3.1.0",
		Info: Info{
			Title:                title,
			Version:              version,
			Description:          pkg.Docs,
			XWebfunctionVersions: pkg.Versions,
		},
		Servers:               []Server{{URL: pkg.BaseURL}},
		Paths:                 map[string]*PathItem{},
		XWebfunctionVersioned: pkg.Versioned(),
	}
	if len(pkg.Errors) > 0 {
		doc.XWebfunctionPackageErrors = errorDocs(pkg.Errors)
	}

	needsBearerAuth := false
	needsErrorTriple := false

	for i := range pkg.Endpoints {
		endpoint := &pkg.Endpoints[i]
		if endpoint.Private() && !includePrivate {
			continue
		}
		op := buildOperation(pkg, endpoint, resolver)
		doc.Paths["/"+endpoint.Name] = &PathItem{Post: op}
		if endpoint.BearerAuth() {
			needsBearerAuth = true
		}
		if endpoint.HasFlag("error_triple") {
			needsErrorTriple = true
		}
	}

	if needsBearerAuth || needsErrorTriple || len(resolver.schemas) > 0 {
		doc.Components = &Components{}
		if needsBearerAuth {
			doc.Components.SecuritySchemes = map[string]SecurityScheme{
				"bearerAuth": {Type: "http", Scheme: "bearer"},
			}
		}
		if needsErrorTriple || len(resolver.schemas) > 0 {
			doc.Components.Schemas = map[string]any{}
			for name, schema := range resolver.schemas {
				doc.Components.Schemas[name] = schema
			}
			if needsErrorTriple {
				doc.Components.Schemas["ErrorTriple"] = errorTripleSchema()
			}
		}
	}

	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling OpenAPI document: %w", err)
	}
	return string(out), nil
}

func errorDocs(defs []webfunction.ErrorDef) []ErrorDoc {
	docs := make([]ErrorDoc, len(defs))
	for i, d := range defs {
		docs[i] = ErrorDoc{Code: d.Code, Docs: d.Docs}
	}
	return docs
}

// errorTripleSchema is the fixed [code, message, details] tuple shape a
// bad-request body takes when an endpoint declares "error_triple" - see
// https://webfunction.org and webfunction-go's error.go. Modeled as a raw
// map[string]any (rather than the general Schema struct) since it needs
// "items": false, a JSON Schema 2020-12 boolean schema, tuple-closing the
// prefixItems list - a shape the general Schema struct doesn't otherwise
// need to express.
func errorTripleSchema() map[string]any {
	return map[string]any{
		"type":        "array",
		"description": "[code, message, details] triple returned on a 400 response.",
		"prefixItems": []any{
			map[string]any{"type": "string", "description": "Machine-readable error code."},
			map[string]any{"type": "string", "description": "Human-readable error message."},
			map[string]any{"description": "Additional error-specific details."},
		},
		"minItems": 2,
		"items":    false,
	}
}

func buildOperation(pkg *webfunction.Package, endpoint *webfunction.Endpoint, resolver *schemaSet) *Operation {
	op := &Operation{
		OperationID: toOperationId(endpoint.Name),
		Summary:     endpoint.Name,
		Description: endpoint.Docs,
		RequestBody: &RequestBody{
			Required: true,
			Content: map[string]MediaType{
				"application/json": {Schema: argumentsToSchema(endpoint.Arguments, resolver)},
			},
		},
		Responses: map[string]*Response{
			"200": successResponse(endpoint, resolver),
		},
		XWebfunctionPaginated:      endpoint.Paginated(),
		XWebfunctionReturnsPackage: endpoint.HasFlag("package"),
		XWebfunctionEventSource:    endpoint.HasFlag("event_source"),
		XWebfunctionCaptureBearer:  endpoint.HasFlag("capture_bearer"),
	}
	if endpoint.Group != "" {
		op.Tags = []string{endpoint.Group}
	}
	if pkg.Versioned() {
		op.Parameters = []Parameter{{
			Name:        "Api-Version",
			In:          "header",
			Description: "The package version to call against.",
			Schema:      &Schema{Type: "string", Enum: stringsToAny(pkg.Versions)},
		}}
	}
	if codes := errorCodes(endpoint); len(codes) > 0 {
		op.XWebfunctionErrorCodes = codes
		op.Responses["400"] = errorResponse(endpoint)
	}
	if endpoint.BearerAuth() {
		op.Security = []SecurityRequirement{{"bearerAuth": {}}}
		op.Responses["401"] = &Response{Description: "Missing or invalid bearer token."}
	}
	return op
}

func successResponse(endpoint *webfunction.Endpoint, resolver *schemaSet) *Response {
	var schema *Schema
	if endpoint.Paginated() {
		schema = pageEnvelopeSchema(endpoint, resolver)
	} else {
		schema = buildSchema(endpoint.Returns, webfunction.AttributeContext, endpoint.Attributes, resolver)
	}
	return &Response{
		Description: "Successful response.",
		Content: map[string]MediaType{
			"application/json": {Schema: schema},
		},
	}
}

// pageEnvelopeSchema builds the real wire envelope for a paginated
// endpoint's response - {"page": [...items], "next": ..., "previous":
// ...} (see webfunction-go's page.go) - using the endpoint's declared
// Returns type as the shape of one item, not the whole response.
func pageEnvelopeSchema(endpoint *webfunction.Endpoint, resolver *schemaSet) *Schema {
	itemSchema := buildSchema(endpoint.Returns, webfunction.AttributeContext, endpoint.Attributes, resolver)
	openOrNull := &Schema{Type: []string{"object", "null"}, Description: "Opaque pagination cursor, or null if there is no such page."}
	return &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"page":     {Type: "array", Items: itemSchema},
			"next":     openOrNull,
			"previous": openOrNull,
		},
		Required: []string{"page"},
	}
}

func errorResponse(endpoint *webfunction.Endpoint) *Response {
	var schema *Schema
	if endpoint.HasFlag("error_triple") {
		schema = &Schema{Ref: "#/components/schemas/ErrorTriple"}
	} else {
		schema = &Schema{Description: "Error response body (shape not further specified by the package)."}
	}
	return &Response{
		Description: "Bad request.",
		Content: map[string]MediaType{
			"application/json": {Schema: schema},
		},
	}
}

// errorCodes collects the endpoint's own declared error codes only - not
// merged with package-level errors, matching this project's established
// caution elsewhere (jsgen/phpgen/gogen's @throws scoping) that whether
// package-level errors implicitly apply to every endpoint isn't confirmed
// by anything read from the spec.
func errorCodes(endpoint *webfunction.Endpoint) []string {
	if len(endpoint.Errors) == 0 {
		return nil
	}
	codes := make([]string, len(endpoint.Errors))
	for i, e := range endpoint.Errors {
		codes[i] = e.Code
	}
	return codes
}

func stringsToAny(ss []string) []any {
	if len(ss) == 0 {
		return nil
	}
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}