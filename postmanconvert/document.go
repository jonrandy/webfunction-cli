package postmanconvert

// This file defines the (small) subset of the real Postman Collection
// Format v2.1.0 object model this package emits, checked directly against
// the raw schema (http://schema.getpostman.com/collection/json/v2.1.0/
// draft-04/collection.json) rather than assumed from memory - there is no
// reference converter project for this target the way webfunction-to-
// openapi existed for openapiconvert.

// Collection is the top-level Postman collection document.
type Collection struct {
	Info     Info       `json:"info"`
	Item     []*Item    `json:"item"`
	Auth     *Auth      `json:"auth,omitempty"`
	Variable []Variable `json:"variable,omitempty"`
}

// Info is the collection's required info block.
type Info struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

// Item is either a folder (Item populated, Request nil) or a leaf request
// item (Request populated, Item nil) - the real schema's "item-group" and
// "item" definitions respectively. This package never populates both on
// the same value.
type Item struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Item        []*Item  `json:"item,omitempty"`
	Request     *Request `json:"request,omitempty"`
	Auth        *Auth    `json:"auth,omitempty"`
}

// Request is the real schema's Request object, restricted to what a
// webfunction endpoint ever needs: POST, JSON body, no query/path params.
type Request struct {
	Method string   `json:"method"`
	Header []Header `json:"header,omitempty"`
	Body   *Body    `json:"body,omitempty"`
	URL    URL      `json:"url"`
	Auth   *Auth    `json:"auth,omitempty"`
}

// Header is one HTTP header.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Body is a raw-mode request body carrying a JSON example.
type Body struct {
	Mode    string       `json:"mode"`
	Raw     string       `json:"raw"`
	Options *BodyOptions `json:"options,omitempty"`
}

type BodyOptions struct {
	Raw *RawBodyOptions `json:"raw,omitempty"`
}

type RawBodyOptions struct {
	Language string `json:"language"`
}

// URL is always built in object form here (raw + host + path), matching
// how every real generator that uses a {{baseUrl}}-style variable renders
// it - a single host segment that IS the variable, not a literal domain.
type URL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host"`
	Path []string `json:"path"`
}

// Auth is the real schema's Auth object - only "bearer" and "noauth" are
// ever emitted here.
type Auth struct {
	Type   string     `json:"type"`
	Bearer []AuthAttr `json:"bearer,omitempty"`
}

// AuthAttr is one auth-attribute entry (key/value/type), per the schema's
// auth-attribute definition.
type AuthAttr struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
	Type  string `json:"type,omitempty"`
}

// Variable is one collection-level variable.
type Variable struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}