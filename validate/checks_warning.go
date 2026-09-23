package validate

import (
	"encoding/json"
	"fmt"

	"github.com/webfunction-protocol/webfunction-go"
)

// checkBareReturns flags an endpoint whose Returns includes a bare
// (unrefined) "object" or a bare "array" alternative with no endpoint-
// level Attributes to fall back on. Technically valid per the wire
// grammar, but genuinely unhelpful to a consumer - an unconstrained
// shape with nothing to type or document against.
func (v *validator) checkBareReturns() {
	for _, e := range v.pkg.Endpoints {
		if len(e.Attributes) > 0 {
			continue
		}
		for _, alt := range e.Returns.Union {
			switch {
			case alt.Base == "object" && alt.Refinement == "":
				v.emit("bare-object-return", Warning, EndpointLoc(e.Name).Returns(),
					`return type includes a bare (unrefined) "object" with no endpoint-level "attributes" to describe its shape`)
			case alt.Base == "array" && alt.Of == nil:
				v.emit("bare-array-return", Warning, EndpointLoc(e.Name).Returns(),
					`return type includes a bare "array" (no element type) with no endpoint-level "attributes" to describe its item shape`)
			}
		}
	}
}

// checkChoicesAndValues flags, for every Argument.Choices and
// Attribute.Values list: an entry that doesn't validate against the
// field's declared type, and duplicate entries within the same list.
//
// When the declared type is (or includes) an array, choices/values
// entries are individual allowed *elements*, not whole allowed arrays -
// e.g. an argument typed array<string> with choices ["red","green"]
// means each element of the submitted array must be "red" or "green",
// not that the whole array must equal ["red"] or ["green"]. So entries
// are checked against the array's element type, not the array type
// itself (t.Valid(val) directly would ask "is val itself an array", which
// is the wrong question here).
func (v *validator) checkChoicesAndValues() {
	for _, e := range v.pkg.Endpoints {
		for _, a := range e.Arguments {
			v.checkChoiceList(a.Type, a.Choices, EndpointLoc(e.Name).Argument(a.Name), "choices")
		}
		for _, a := range e.Attributes {
			v.checkChoiceList(a.Type, a.Values, EndpointLoc(e.Name).Attribute(a.Name), "values")
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			v.checkChoiceList(a.Type, a.Choices, ObjectLoc(o.Name).Argument(a.Name), "choices")
		}
		for _, a := range o.Attributes {
			v.checkChoiceList(a.Type, a.Values, ObjectLoc(o.Name).Attribute(a.Name), "values")
		}
	}
}

func (v *validator) checkChoiceList(t webfunction.Type, values []any, loc Loc, field string) {
	if len(values) == 0 {
		return
	}

	elementType, skipTypeCheck, describedAs := elementValidationType(t)

	seen := map[string]bool{}
	for _, val := range values {
		if !skipTypeCheck && !elementType.Valid(val) {
			v.emit("choice-value-type-mismatch", Warning, loc,
				fmt.Sprintf("%s entry %s doesn't validate against %s", field, jsonRepr(val), describedAs))
		}
		key := jsonRepr(val)
		if seen[key] {
			v.emit("duplicate-choice-value", Warning, loc,
				fmt.Sprintf("%s contains a duplicate entry %s", field, key))
		}
		seen[key] = true
	}
}

// elementValidationType derives the type that individual choices/values
// entries should be checked against, given the field's own declared
// type t:
//
//   - A non-array alternative contributes itself directly (an entry can
//     simply be a valid value of that alternative).
//   - An array alternative with an element type (array<T>) contributes
//     T's own alternatives - entries are elements of the array, not
//     whole arrays.
//   - A bare "array" alternative (no element type, i.e. array of any)
//     means there's no meaningful per-element constraint at all, so
//     skipTypeCheck is returned true - any entry is fine.
//
// describedAs is a human-readable phrase for the finding message,
// naming the array's element type explicitly when that's what's
// actually being checked against, so the message doesn't just repeat
// the field's own (possibly array) type and confuse the two.
func elementValidationType(t webfunction.Type) (elementType webfunction.Type, skipTypeCheck bool, describedAs string) {
	var alts []webfunction.TypeAlt
	arrayDerived := false
	for _, alt := range t.Union {
		if alt.Base == "array" {
			if alt.Of == nil {
				// Array of any - no per-element constraint to check.
				return webfunction.Type{}, true, ""
			}
			alts = append(alts, alt.Of.Union...)
			arrayDerived = true
			continue
		}
		alts = append(alts, alt)
	}

	elementType = webfunction.Type{Union: alts}
	if arrayDerived {
		describedAs = fmt.Sprintf("the array's own element type %q", elementType.String())
	} else {
		describedAs = fmt.Sprintf("the field's own declared type %q", elementType.String())
	}
	return elementType, false, describedAs
}

// jsonRepr renders a decoded-JSON value (as found in Choices/Values) in
// a short, canonical, comparable form for messages and dedup keys.
func jsonRepr(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}