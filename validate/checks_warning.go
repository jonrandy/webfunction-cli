package validate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/webfunction-protocol/webfunction-go"
)

// checkBaseURL flags a base_url ending in "/". Every codegen target and
// both convert targets already defend against this downstream (their own
// "/" + endpoint-name join would otherwise double up), but the defect is
// upstream in the package definition itself - worth flagging at the
// source instead of relying on every consumer to keep working around it.
func (v *validator) checkBaseURL() {
	if strings.HasSuffix(v.pkg.BaseURL, "/") {
		v.emit("base-url-trailing-slash", Warning, "package",
			fmt.Sprintf("base_url %q ends with a trailing slash; joining it with an endpoint name (base_url + \"/\" + name) produces a double slash", v.pkg.BaseURL))
	}
}

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
				v.emit("bare-object-return", Warning, fmt.Sprintf("endpoint %q returns", e.Name),
					`return type includes a bare (unrefined) "object" with no endpoint-level "attributes" to describe its shape`)
			case alt.Base == "array" && alt.Of == nil:
				v.emit("bare-array-return", Warning, fmt.Sprintf("endpoint %q returns", e.Name),
					`return type includes a bare "array" (no element type) with no endpoint-level "attributes" to describe its item shape`)
			}
		}
	}
}

// checkChoicesAndValues flags, for every Argument.Choices and
// Attribute.Values list: a value that doesn't validate against the
// field's own declared Type (using Type.Valid, the same validity check
// webfunction-go itself would apply to a real call), and duplicate
// values within the same list.
func (v *validator) checkChoicesAndValues() {
	for _, e := range v.pkg.Endpoints {
		for _, a := range e.Arguments {
			v.checkChoiceList(a.Type, a.Choices, fmt.Sprintf("endpoint %q argument %q", e.Name, a.Name), "choices")
		}
		for _, a := range e.Attributes {
			v.checkChoiceList(a.Type, a.Values, fmt.Sprintf("endpoint %q attribute %q", e.Name, a.Name), "values")
		}
	}
	for _, o := range v.pkg.Objects {
		for _, a := range o.Arguments {
			v.checkChoiceList(a.Type, a.Choices, fmt.Sprintf("object %q argument %q", o.Name, a.Name), "choices")
		}
		for _, a := range o.Attributes {
			v.checkChoiceList(a.Type, a.Values, fmt.Sprintf("object %q attribute %q", o.Name, a.Name), "values")
		}
	}
}

func (v *validator) checkChoiceList(t webfunction.Type, values []any, location, field string) {
	if len(values) == 0 {
		return
	}
	seen := map[string]bool{}
	for _, val := range values {
		if !t.Valid(val) {
			v.emit("choice-value-type-mismatch", Warning, location,
				fmt.Sprintf("%s entry %s doesn't validate against the field's own declared type %q", field, jsonRepr(val), t.String()))
		}
		key := jsonRepr(val)
		if seen[key] {
			v.emit("duplicate-choice-value", Warning, location,
				fmt.Sprintf("%s contains a duplicate entry %s", field, key))
		}
		seen[key] = true
	}
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