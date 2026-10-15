package openapiconvert

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"
)

const refPrefix = "#/components/schemas/"

// pkgJSON wraps object and endpoint JSON fragments in a minimal package.
func pkgJSON(objects, endpoints string) string {
	return `{"base_url":"https://api.example.com","name":"Generics Convert Test","objects":[` + objects + `],"endpoints":[` + endpoints + `]}`
}

func parsePackage(t *testing.T, raw string) *webfunction.Package {
	t.Helper()
	var pkg webfunction.Package
	if err := json.Unmarshal([]byte(raw), &pkg); err != nil {
		t.Fatal(err)
	}
	return &pkg
}

// generate converts a package and decodes the OpenAPI document.
func generate(t *testing.T, pkg *webfunction.Package, includePrivate bool) map[string]any {
	t.Helper()
	out, err := Generate(pkg, includePrivate)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("output is not JSON: %v", err)
	}
	return doc
}

func schemaKeys(doc map[string]any) []string {
	comps, _ := doc["components"].(map[string]any)
	schemas, _ := comps["schemas"].(map[string]any)
	keys := make([]string, 0, len(schemas))
	for k := range schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func schemaOf(doc map[string]any, key string) map[string]any {
	comps, _ := doc["components"].(map[string]any)
	schemas, _ := comps["schemas"].(map[string]any)
	s, _ := schemas[key].(map[string]any)
	return s
}

// refsIn collects every $ref string found anywhere inside v.
func refsIn(v any) []string {
	var refs []string
	switch x := v.(type) {
	case map[string]any:
		for k, val := range x {
			if s, ok := val.(string); ok && k == "$ref" {
				refs = append(refs, s)
				continue
			}
			refs = append(refs, refsIn(val)...)
		}
	case []any:
		for _, val := range x {
			refs = append(refs, refsIn(val)...)
		}
	}
	sort.Strings(refs)
	return refs
}

func contains(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

// assertRefsResolve fails if any $ref in the document points at a schema
// that is not defined, which is what an unexpanded generic would cause.
func assertRefsResolve(t *testing.T, doc map[string]any) {
	t.Helper()
	keys := schemaKeys(doc)
	for _, ref := range refsIn(doc) {
		if !strings.HasPrefix(ref, refPrefix) || !contains(keys, strings.TrimPrefix(ref, refPrefix)) {
			t.Errorf("dangling $ref %q (schemas: %v)", ref, keys)
		}
	}
}

// TestGenericSchemaNames checks the components generated for the fixture:
// one concrete schema per application, named in Pascal case, with the
// generic templates themselves gone.
func TestGenericSchemaNames(t *testing.T) {
	doc := generate(t, loadPackage(t, "../testdata/convert/generics-types.json"), false)
	want := []string{"FilterOfStringInput", "ListOfUser", "ListOfUserOrNull", "PageOfPageOfUser", "PageOfUser", "PairOfNumber", "User"}
	if got := schemaKeys(doc); !reflect.DeepEqual(got, want) {
		t.Errorf("schema keys:\n got %v\nwant %v", got, want)
	}
	assertRefsResolve(t, doc)
}

func TestGenericApplicationSchema(t *testing.T) {
	doc := generate(t, loadPackage(t, "../testdata/convert/generics-types.json"), false)
	list := schemaOf(doc, "ListOfUser")
	if list == nil {
		t.Fatal("no ListOfUser schema")
	}
	props, _ := list["properties"].(map[string]any)
	if _, ok := props["total"]; !ok {
		t.Error("ListOfUser lost its non-generic total property")
	}
	if !contains(refsIn(props["items"]), refPrefix+"User") {
		t.Errorf("ListOfUser.items should reference User, got %v", props["items"])
	}
	if !contains(refsIn(doc["paths"]), refPrefix+"ListOfUser") {
		t.Error("an endpoint response should reference ListOfUser")
	}
	filter := schemaOf(doc, "FilterOfStringInput")
	if req, _ := filter["required"].([]any); len(req) != 1 || req[0] != "value" {
		t.Errorf("FilterOfStringInput should require value, got %v", filter["required"])
	}
}

// TestGenericEqualsMonomorphic is the core check: a package using
// generics must convert to exactly the document its hand-written
// monomorphic equivalent converts to.
func TestGenericEqualsMonomorphic(t *testing.T) {
	user := `{"name":"user","attributes":[{"name":"id","type":"string"}]}`
	list := `{"name":"list","attributes":[{"name":"total","type":"number"},{"name":"items","type":[["T"]]}]}`
	page := `{"name":"page","attributes":[{"name":"results","type":[["T"]]}]}`
	pair := `{"name":"pair","attributes":[{"name":"first","type":"T"},{"name":"second","type":"T"}]}`
	box := `{"name":"box","arguments":[{"name":"value","type":"T"}],"attributes":[{"name":"value","type":"T"}]}`
	filter := `{"name":"filter","arguments":[{"name":"value","type":"T","flags":["required"]}]}`
	ret := func(returns string) string {
		return `{"name":"get","returns":` + returns + `,"arguments":[]}`
	}

	cases := []struct {
		name    string
		generic string
		mono    string
	}{
		{
			"single application",
			pkgJSON(user+","+list, ret(`{"object.list":"object.user"}`)),
			pkgJSON(user+`,{"name":"ListOfUser","attributes":[{"name":"total","type":"number"},{"name":"items","type":[["object.user"]]}]}`, ret(`"object.ListOfUser"`)),
		},
		{
			"union argument is spliced, not nested",
			pkgJSON(user+","+list, ret(`{"object.list":["object.user","null"]}`)),
			pkgJSON(user+`,{"name":"ListOfUserOrNull","attributes":[{"name":"total","type":"number"},{"name":"items","type":[["object.user","null"]]}]}`, ret(`"object.ListOfUserOrNull"`)),
		},
		{
			"nested application",
			pkgJSON(user+","+page, ret(`{"object.page":{"object.page":"object.user"}}`)),
			pkgJSON(user+`,{"name":"PageOfUser","attributes":[{"name":"results","type":[["object.user"]]}]},{"name":"PageOfPageOfUser","attributes":[{"name":"results","type":[["object.PageOfUser"]]}]}`, ret(`"object.PageOfPageOfUser"`)),
		},
		{
			"application inside an array",
			pkgJSON(pair, ret(`[[{"object.pair":"number"}]]`)),
			pkgJSON(`{"name":"PairOfNumber","attributes":[{"name":"first","type":"number"},{"name":"second","type":"number"}]}`, ret(`[["object.PairOfNumber"]]`)),
		},
		{
			"argument-only generic stays argument-only",
			pkgJSON(filter, `{"name":"search","returns":"boolean","arguments":[{"name":"filter","type":{"object.filter":"string"}}]}`),
			pkgJSON(`{"name":"FilterOfString","arguments":[{"name":"value","type":"string","flags":["required"]}]}`, `{"name":"search","returns":"boolean","arguments":[{"name":"filter","type":"object.FilterOfString"}]}`),
		},
		{
			"generic used in both contexts",
			pkgJSON(box, `{"name":"swap","returns":{"object.box":"string"},"arguments":[{"name":"box","type":{"object.box":"string"}}]}`),
			pkgJSON(`{"name":"BoxOfString","arguments":[{"name":"value","type":"string"}],"attributes":[{"name":"value","type":"string"}]}`, `{"name":"swap","returns":"object.BoxOfString","arguments":[{"name":"box","type":"object.BoxOfString"}]}`),
		},
		{
			"argument applied to a static object is ignored",
			pkgJSON(user, ret(`{"object.user":"string"}`)),
			pkgJSON(user, ret(`"object.user"`)),
		},
		{
			"malformed application with two keys resolves to any",
			pkgJSON(user+","+list, ret(`{"object.list":"object.user","object.page":"string"}`)),
			pkgJSON(user+","+list, ret(`"any"`)),
		},
		{
			"malformed application with a non-object key resolves to any",
			pkgJSON(user, ret(`{"user":"string"}`)),
			pkgJSON(user, ret(`"any"`)),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, includePrivate := range []bool{false, true} {
				got := generate(t, parsePackage(t, c.generic), includePrivate)
				want := generate(t, parsePackage(t, c.mono), includePrivate)
				if !reflect.DeepEqual(got, want) {
					g, _ := json.MarshalIndent(got, "", "  ")
					w, _ := json.MarshalIndent(want, "", "  ")
					t.Errorf("includePrivate=%v: generic output differs from its monomorphic equivalent\n--- got\n%s\n--- want\n%s", includePrivate, g, w)
				}
				assertRefsResolve(t, got)
			}
		})
	}
}

// An unapplied generic resolves T to any rather than failing or leaving a
// dangling reference (the validator, not the converter, reports it).
func TestUnappliedGenericConverts(t *testing.T) {
	list := `{"name":"list","attributes":[{"name":"total","type":"number"},{"name":"items","type":[["T"]]}]}`
	doc := generate(t, parsePackage(t, pkgJSON(list, `{"name":"get","returns":"object.list","arguments":[]}`)), false)
	assertRefsResolve(t, doc)
	if len(refsIn(doc["paths"])) == 0 {
		t.Error("the endpoint should still reference a schema for the list")
	}
	for _, k := range schemaKeys(doc) {
		if k == "List" {
			t.Error("the generic template itself must not become a schema")
		}
	}
}

// Private members of a generic template must carry over to every
// instantiation and be filtered like any other private member, and a
// private endpoint must take its instantiations with it.
func TestGenericPrivateFlags(t *testing.T) {
	box := `{"name":"box","attributes":[{"name":"value","type":"T"},{"name":"secret","type":"T","flags":["private"]}]}`
	list := `{"name":"list","attributes":[{"name":"items","type":[["T"]]}]}`
	user := `{"name":"user","attributes":[{"name":"id","type":"string"}]}`
	raw := pkgJSON(box+","+list+","+user,
		`{"name":"get-box","returns":{"object.box":"string"},"arguments":[]},`+
			`{"name":"internal-list","returns":{"object.list":"object.user"},"flags":["private"],"arguments":[]}`)

	pub := generate(t, parsePackage(t, raw), false)
	props, _ := schemaOf(pub, "BoxOfString")["properties"].(map[string]any)
	if _, ok := props["value"]; !ok {
		t.Errorf("public BoxOfString should keep value, got %v", props)
	}
	if _, ok := props["secret"]; ok {
		t.Error("public BoxOfString must not contain the private member")
	}
	if schemaOf(pub, "ListOfUser") != nil {
		t.Error("a private endpoint's instantiation must not appear in public output")
	}
	assertRefsResolve(t, pub)

	priv := generate(t, parsePackage(t, raw), true)
	props, _ = schemaOf(priv, "BoxOfString")["properties"].(map[string]any)
	if _, ok := props["secret"]; !ok {
		t.Error("private output should keep the private member")
	}
	if schemaOf(priv, "ListOfUser") == nil {
		t.Error("private output should include the private endpoint's instantiation")
	}
	assertRefsResolve(t, priv)
}

// Generate must neither mutate its input nor produce unstable output.
func TestGenerateWithGenericsIsPureAndDeterministic(t *testing.T) {
	pkg := loadPackage(t, "../testdata/convert/generics-types.json")
	before := pkg.Endpoints[0].Returns.String()
	first, err := Generate(pkg, false)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Generate(pkg, false)
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Error("Generate is not deterministic for a generic package")
	}
	if after := pkg.Endpoints[0].Returns.String(); after != before || len(pkg.Objects) != 5 {
		t.Errorf("Generate mutated its input package (returns %q -> %q, %d objects)", before, after, len(pkg.Objects))
	}
}