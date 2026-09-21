// Package postmanconvert converts a webfunction.Package into a Postman
// Collection Format v2.1.0 document.
//
// Written from scratch - there's no reference converter project for this
// target the way webfunction-to-openapi existed for openapiconvert.
// Confirmed with Jon before writing:
//  1. Endpoints are grouped into folders by Endpoint.Group when a package
//     uses groups; an endpoint with no group is a flat top-level item.
//     Folders and flat items interleave in the endpoints' declared order
//     (a folder is created at the position of its first member).
//  2. A package needing bearer auth at all gets collection-level bearer
//     auth too (via a {{bearerToken}} variable, useful in a full Postman
//     client), but every item's own auth is always set explicitly -
//     bearer where the endpoint requires it, {"type": "noauth"}
//     otherwise - rather than relying on Postman's real "inherit auth
//     from parent" resolution. Inheritance is legitimate, standard
//     Postman behavior, but real-world testing found that at least one
//     popular lightweight/third-party collection *viewer* only reads an
//     item's own auth field and never walks up to the collection level,
//     showing no auth requirement on any endpoint as a result - explicit
//     per-item auth is what actually renders correctly everywhere.
//  3. Only {{baseUrl}} (and {{apiVersion}} for a versioned package,
//     {{bearerToken}} when needed) are parameterized as collection
//     variables - nothing else.
//  4. Items carry only a request - no canned example "saved" responses.
//  5. Private endpoints are excluded unless includePrivate is true (the
//     convert command's --private flag, shared with openapiconvert).
package postmanconvert

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
)

const schemaURL = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

// Generate converts pkg into a Postman collection and returns it as
// indented JSON.
func Generate(pkg *webfunction.Package, includePrivate bool) (string, error) {
	title := pkg.Name
	if title == "" {
		title = "Web Function API"
	}

	var included []*webfunction.Endpoint
	needsBearer := false
	for i := range pkg.Endpoints {
		endpoint := &pkg.Endpoints[i]
		if endpoint.Private() && !includePrivate {
			continue
		}
		included = append(included, endpoint)
		if endpoint.BearerAuth() {
			needsBearer = true
		}
	}

	items := buildItems(pkg, included, needsBearer)

	doc := &Collection{
		Info: Info{
			Name:        title,
			Description: pkg.Docs,
			Schema:      schemaURL,
		},
		Item:     items,
		Variable: buildVariables(pkg, needsBearer),
	}
	if needsBearer {
		doc.Auth = &Auth{Type: "bearer", Bearer: []AuthAttr{{Key: "token", Value: "{{bearerToken}}", Type: "string"}}}
	}

	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling Postman collection: %w", err)
	}
	return string(out), nil
}

func trimTrailingSlashes(url string) string {
	return strings.TrimRight(url, "/")
}

func buildVariables(pkg *webfunction.Package, needsBearer bool) []Variable {
	vars := []Variable{{Key: "baseUrl", Value: trimTrailingSlashes(pkg.BaseURL)}}
	if needsBearer {
		vars = append(vars, Variable{Key: "bearerToken", Value: "your-bearer-token-here"})
	}
	if pkg.Versioned() {
		version := pkg.Version
		if version == "" && len(pkg.Versions) > 0 {
			version = pkg.Versions[len(pkg.Versions)-1]
		}
		vars = append(vars, Variable{Key: "apiVersion", Value: version})
	}
	return vars
}

// buildItems lays endpoints out in their declared order, opening a folder
// the first time a given Group is seen and appending later same-group
// endpoints into that existing folder wherever it landed.
func buildItems(pkg *webfunction.Package, endpoints []*webfunction.Endpoint, needsBearer bool) []*Item {
	folderAt := map[string]int{}
	var items []*Item
	for _, endpoint := range endpoints {
		item := buildItem(pkg, endpoint, needsBearer)
		if endpoint.Group == "" {
			items = append(items, item)
			continue
		}
		if idx, ok := folderAt[endpoint.Group]; ok {
			items[idx].Item = append(items[idx].Item, item)
			continue
		}
		folder := &Item{Name: endpoint.Group, Item: []*Item{item}}
		items = append(items, folder)
		folderAt[endpoint.Group] = len(items) - 1
	}
	return items
}

func buildItem(pkg *webfunction.Package, endpoint *webfunction.Endpoint, needsBearer bool) *Item {
	body := requestBodyExample(endpoint.Arguments, pkg)
	raw, _ := json.MarshalIndent(body, "", "  ") // orderedMap.MarshalJSON never errors on the plain values this package produces

	header := []Header{{Key: "Content-Type", Value: "application/json"}}
	if pkg.Versioned() {
		header = append(header, Header{Key: "Api-Version", Value: "{{apiVersion}}"})
	}

	req := &Request{
		Method: "POST",
		Header: header,
		Body: &Body{
			Mode:    "raw",
			Raw:     string(raw),
			Options: &BodyOptions{Raw: &RawBodyOptions{Language: "json"}},
		},
		URL: URL{
			Raw:  "{{baseUrl}}/" + endpoint.Name,
			Host: []string{"{{baseUrl}}"},
			Path: []string{endpoint.Name},
		},
	}
	if needsBearer {
		if endpoint.BearerAuth() {
			// Set explicitly per-item rather than relying on Postman's
			// real "inherit auth from parent" resolution (walking up to
			// the collection-level auth set below) - that's legitimate,
			// standard Postman behavior, but plenty of lightweight/
			// third-party collection *viewers* (as opposed to full
			// Postman-compatible clients) only ever read an item's own
			// auth field and never implement the inheritance walk, so
			// they show no auth requirement at all on an item that
			// relies on it.
			req.Auth = &Auth{Type: "bearer", Bearer: []AuthAttr{{Key: "token", Value: "{{bearerToken}}", Type: "string"}}}
		} else {
			req.Auth = &Auth{Type: "noauth"}
		}
	}

	return &Item{
		Name:        endpoint.Name,
		Description: itemDescription(endpoint),
		Request:     req,
	}
}

// itemDescription appends brief, informational notes about pagination and
// declared errors onto the endpoint's own docs - Postman has no
// structured per-response schema/vendor-extension mechanism the way
// OpenAPI does, so this is folded into plain description text instead.
func itemDescription(endpoint *webfunction.Endpoint) string {
	desc := endpoint.Docs
	if endpoint.Paginated() {
		desc = appendNote(desc, "**Paginated.** Response is `{ page: [...], next, previous }`.")
	}
	if len(endpoint.Errors) > 0 {
		note := "**Errors:**"
		for _, e := range endpoint.Errors {
			note += fmt.Sprintf("\n- `%s`", e.Code)
			if e.Docs != "" {
				note += " - " + e.Docs
			}
		}
		if endpoint.HasFlag("error_triple") {
			note += "\n\nReturned as a `[code, message, details]` triple on a 400 response."
		}
		desc = appendNote(desc, note)
	}
	return desc
}

func appendNote(desc, note string) string {
	if desc == "" {
		return note
	}
	return desc + "\n\n" + note
}