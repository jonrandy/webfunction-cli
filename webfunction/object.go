package webfunction

// Object is a named, reusable object type definition, referenced elsewhere
// in the package as a refined "object.<name>" type.
// See https://webfunction.org/package#object-definition.
//
// Per the spec, the same object name can carry two different member lists
// depending on the context it's referenced from: Arguments when the object
// is referenced as an argument's type, Attributes when referenced as an
// endpoint's return type or an attribute's type. An object MAY define
// both, if it's referenced in both contexts somewhere in the package.
type Object struct {
	Name       string      `json:"name"`
	Arguments  []Argument  `json:"arguments,omitempty"`
	Attributes []Attribute `json:"attributes,omitempty"`
}