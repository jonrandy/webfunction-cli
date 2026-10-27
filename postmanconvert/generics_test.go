package postmanconvert

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"
)

// pkgJSON wraps object and endpoint JSON fragments in a minimal package.
func pkgJSON(objects, endpoints string) string {
	return `{"base_url":"https://api.example.com","name":"Generics Postman Test","objects":[` + objects + `],"endpoints":[` + endpoints + `]}`
}

func parsePackage(t *testing.T, raw string) *webfunction.Package {
	t.Helper()
	var pkg webfunction.Package
	if err := json.Unmarshal([]byte(raw), &pkg); err != nil {
		t.Fatal(err)
	}
	return &pkg
}

// findItem returns the request item named name, searching folders.
func findItem(items []any, name string) map[string]any {
	for _, it := range items {
		m, _ := it.(map[string]any)
		if m["name"] == name {
			if _, ok := m["request"]; ok {
				return m
			}
		}
		if sub, ok := m["item"].([]any); ok {
			if found := findItem(sub, name); found != nil {
				return found
			}
		}
	}
	return nil
}

// requestBody converts pkg and returns the decoded JSON request body of
// the endpoint called name, or nil when the endpoint is absent.
func requestBody(t *testing.T, pkg *webfunction.Package, includePrivate bool, name string) map[string]any {
	t.Helper()
	out, err := Generate(pkg, includePrivate)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("output is not JSON: %v", err)
	}
	items, _ := doc["item"].([]any)
	item := findItem(items, name)
	if item == nil {
		return nil
	}
	req, _ := item["request"].(map[string]any)
	body, _ := req["body"].(map[string]any)
	raw, _ := body["raw"].(string)
	var decoded map[string]any
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
			t.Fatalf("request body is not JSON: %v\n%s", err, raw)
		}
	}
	return decoded
}

const (
	wrapper = `{"name":"wrapper","arguments":[{"name":"value","type":"T"},{"name":"items","type":[["T"]]}]}`
	address = `{"name":"address","arguments":[{"name":"street","type":"string"}]}`
)

// run is an endpoint taking one argument, arg, of the given type.
func run(argType string) string {
	return `{"name":"run","returns":"boolean","arguments":[{"name":"arg","type":` + argType + `}]}`
}

func decode(t *testing.T, raw string) any {
	t.Helper()
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

// TestGenericArgumentExamples pins the example value generated for an
// argument whose type is a generic application: the type argument must
// reach every T, so examples are typed rather than null.
func TestGenericArgumentExamples(t *testing.T) {
	cases := []struct {
		name    string
		argType string
		want    string
	}{
		{"number argument", `{"object.wrapper":"number"}`, `{"value":0,"items":[0]}`},
		{"union argument is spliced, first non-null wins", `{"object.wrapper":["string.email","null"]}`, `{"value":"user@example.com","items":["user@example.com"]}`},
		{"object argument", `{"object.wrapper":"object.address"}`, `{"value":{"street":""},"items":[{"street":""}]}`},
		{"dangling object reference inside an application renders empty", `{"object.wrapper":"object.ghost"}`, `{"value":{},"items":[{}]}`},
		{"nested application", `{"object.wrapper":{"object.wrapper":"boolean"}}`, `{"value":{"value":false,"items":[false]},"items":[{"value":false,"items":[false]}]}`},
		{"application inside an array", `[[{"object.wrapper":"number"}]]`, `[{"value":0,"items":[0]}]`},
		{"unapplied generic resolves T to any", `"object.wrapper"`, `{"value":null,"items":[null]}`},
		{"malformed application resolves to any", `{"object.wrapper":"number","object.address":"string"}`, `null`},
		{"argument on a static object is ignored", `{"object.address":"string"}`, `{"street":""}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pkg := parsePackage(t, pkgJSON(wrapper+","+address, run(c.argType)))
			got := requestBody(t, pkg, false, "run")["arg"]
			if want := decode(t, c.want); !reflect.DeepEqual(got, want) {
				g, _ := json.Marshal(got)
				t.Errorf("arg example:\n got %s\nwant %s", g, c.want)
			}
		})
	}
}

func TestGenericFixtureRequestBody(t *testing.T) {
	body := requestBody(t, loadPackage(t, "../testdata/convert/generics-types.json"), false, "list-users")
	want := decode(t, `{"filter":{"value":""}}`)
	if !reflect.DeepEqual(any(body), want) {
		g, _ := json.Marshal(body)
		t.Errorf("list-users body: got %s, want %s", g, `{"filter":{"value":""}}`)
	}
}

// TestGenericEqualsMonomorphic: a package using generics must convert to
// exactly the collection its hand-written monomorphic equivalent does.
// Postman inlines examples, so the equivalent's object name is irrelevant.
func TestGenericEqualsMonomorphic(t *testing.T) {
	box := `{"name":"box","arguments":[{"name":"value","type":"T"}],"attributes":[{"name":"value","type":"T"}]}`
	cases := []struct{ name, generic, mono string }{
		{
			"single application",
			pkgJSON(wrapper, run(`{"object.wrapper":"number"}`)),
			pkgJSON(`{"name":"m","arguments":[{"name":"value","type":"number"},{"name":"items","type":[["number"]]}]}`, run(`"object.m"`)),
		},
		{
			"union argument",
			pkgJSON(wrapper, run(`{"object.wrapper":["string.email","null"]}`)),
			pkgJSON(`{"name":"m","arguments":[{"name":"value","type":["string.email","null"]},{"name":"items","type":[["string.email","null"]]}]}`, run(`"object.m"`)),
		},
		{
			"nested application",
			pkgJSON(wrapper, run(`{"object.wrapper":{"object.wrapper":"boolean"}}`)),
			pkgJSON(`{"name":"inner","arguments":[{"name":"value","type":"boolean"},{"name":"items","type":[["boolean"]]}]},`+
				`{"name":"m","arguments":[{"name":"value","type":"object.inner"},{"name":"items","type":[["object.inner"]]}]}`, run(`"object.m"`)),
		},
		{
			"generic with both contexts, used as an argument",
			pkgJSON(box, run(`{"object.box":"string"}`)),
			pkgJSON(`{"name":"m","arguments":[{"name":"value","type":"string"}],"attributes":[{"name":"value","type":"string"}]}`, run(`"object.m"`)),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, includePrivate := range []bool{false, true} {
				got, err := Generate(parsePackage(t, c.generic), includePrivate)
				if err != nil {
					t.Fatal(err)
				}
				want, err := Generate(parsePackage(t, c.mono), includePrivate)
				if err != nil {
					t.Fatal(err)
				}
				if got != want {
					t.Errorf("includePrivate=%v: generic output differs from its monomorphic equivalent\n--- got\n%s\n--- want\n%s", includePrivate, got, want)
				}
			}
		})
	}
}

// Postman output has no responses, so what an endpoint returns, generic or
// not, must not change the collection.
func TestGenericReturnsDoNotAffectOutput(t *testing.T) {
	list := `{"name":"list","attributes":[{"name":"items","type":[["T"]]}]}`
	withGeneric := pkgJSON(list, `{"name":"get","returns":{"object.list":"string"},"arguments":[{"name":"id","type":"string"}]}`)
	plain := pkgJSON(list, `{"name":"get","returns":"boolean","arguments":[{"name":"id","type":"string"}]}`)
	got, err := Generate(parsePackage(t, withGeneric), false)
	if err != nil {
		t.Fatal(err)
	}
	want, err := Generate(parsePackage(t, plain), false)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Error("a generic return type changed the collection")
	}
}

// Private members of a generic must carry over to the instantiation and be
// filtered like any other; a private endpoint or argument takes its
// application with it.
func TestGenericPrivateFlags(t *testing.T) {
	box := `{"name":"box","arguments":[{"name":"value","type":"T"},{"name":"secret","type":"T","flags":["private"]}]}`
	endpoints := `{"name":"run","returns":"boolean","arguments":[` +
		`{"name":"arg","type":{"object.box":"string"}},` +
		`{"name":"debug","type":{"object.box":"number"},"flags":["private"]}]},` +
		`{"name":"internal","returns":"boolean","flags":["private"],"arguments":[{"name":"arg","type":{"object.box":"string"}}]}`
	raw := pkgJSON(box, endpoints)

	pub := requestBody(t, parsePackage(t, raw), false, "run")
	if want := decode(t, `{"arg":{"value":""}}`); !reflect.DeepEqual(any(pub), want) {
		g, _ := json.Marshal(pub)
		t.Errorf("public body: got %s, want %s", g, `{"arg":{"value":""}}`)
	}
	if requestBody(t, parsePackage(t, raw), false, "internal") != nil {
		t.Error("a private endpoint must not appear in public output")
	}

	priv := requestBody(t, parsePackage(t, raw), true, "run")
	want := decode(t, `{"arg":{"value":"","secret":""},"debug":{"value":0,"secret":0}}`)
	if !reflect.DeepEqual(any(priv), want) {
		g, _ := json.Marshal(priv)
		t.Errorf("private body: got %s, want %s", g, `{"arg":{"value":"","secret":""},"debug":{"value":0,"secret":0}}`)
	}
	if requestBody(t, parsePackage(t, raw), true, "internal") == nil {
		t.Error("a private endpoint should appear in private output")
	}
}

// Generate must neither mutate its input nor produce unstable output.
func TestGenerateWithGenericsIsPureAndDeterministic(t *testing.T) {
	pkg := loadPackage(t, "../testdata/convert/generics-types.json")
	before, _ := json.Marshal(pkg)
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
	if after, _ := json.Marshal(pkg); string(after) != string(before) {
		t.Error("Generate mutated its input package")
	}
}

// ---- Array applications, applications inside objects, and recursion -------

const (
	nodeObject   = `{"name":"node","arguments":[{"name":"label","type":"string"},{"name":"children","type":{"array":"object.node"}}]}`
	listObject   = `{"name":"list","arguments":[{"name":"items","type":[["T"]]}]}`
	threadObject = `{"name":"thread","arguments":[{"name":"comments","type":{"object.list":"object.thread"}}]}`
	holderObject = `{"name":"holder","arguments":[{"name":"inner","type":{"object.wrapper":"number"}},{"name":"many","type":{"array":{"object.wrapper":"string"}}}]}`
	chainObject  = `{"name":"chain","arguments":[{"name":"value","type":"T"},{"name":"next","type":{"object.chain":[["T"]]}}]}`
)

// exampleOf converts a package whose one endpoint takes an argument of
// argType and returns that argument's example value.
func exampleOf(t *testing.T, objects, argType string) any {
	t.Helper()
	return requestBody(t, parsePackage(t, pkgJSON(objects, run(argType))), false, "run")["arg"]
}

// Spec (Generics): { "array": X } is an array of X, the same as [[X]].
func TestArrayApplicationExamples(t *testing.T) {
	cases := []struct{ name, argType, legacy, want string }{
		{"array of a number", `{"array":"number"}`, `[["number"]]`, `[0]`},
		{"array of arrays", `{"array":{"array":"string"}}`, `[[["string"]]]`, `[[""]]`},
		{"application of an array application", `{"object.wrapper":{"array":"number"}}`, `{"object.wrapper":[["number"]]}`, `{"value":[0],"items":[[0]]}`},
		{"array of applications", `{"array":{"object.wrapper":"number"}}`, `[[{"object.wrapper":"number"}]]`, `[{"value":0,"items":[0]}]`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := exampleOf(t, wrapper, c.argType)
			if want := decode(t, c.want); !reflect.DeepEqual(got, want) {
				g, _ := json.Marshal(got)
				t.Errorf("example:\n got %s\nwant %s", g, c.want)
			}
			if legacy := exampleOf(t, wrapper, c.legacy); !reflect.DeepEqual(got, legacy) {
				t.Errorf("%s should give the same example as %s", c.argType, c.legacy)
			}
		})
	}
}

func TestApplicationInsideObjectExample(t *testing.T) {
	got := exampleOf(t, wrapper+","+holderObject, `"object.holder"`)
	want := decode(t, `{"inner":{"value":0,"items":[0]},"many":[{"value":"","items":[""]}]}`)
	if !reflect.DeepEqual(got, want) {
		g, _ := json.Marshal(got)
		t.Errorf("holder example: got %s", g)
	}
}

// A self-referential argument stops where it would repeat, rendering the
// repeat as an empty object.
func TestRecursiveArgumentExamplesTerminate(t *testing.T) {
	cases := []struct{ name, objects, argType, want string }{
		{"object with an array of itself", nodeObject, `"object.node"`, `{"label":"","children":[{}]}`},
		{"object applying a generic to itself", listObject + "," + threadObject, `"object.thread"`, `{"comments":{"items":[{}]}}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := exampleOf(t, c.objects, c.argType)
			if want := decode(t, c.want); !reflect.DeepEqual(got, want) {
				g, _ := json.Marshal(got)
				t.Errorf("example:\n got %s\nwant %s", g, c.want)
			}
		})
	}
}

// Invalid polymorphic recursion is cut off by Expand's depth limit, so the
// example nests deeply but finitely.
func TestPolymorphicRecursionExampleIsFinite(t *testing.T) {
	depth := 0
	for v, _ := exampleOf(t, chainObject, `{"object.chain":"string"}`).(map[string]any); v != nil; depth++ {
		v, _ = v["next"].(map[string]any)
	}
	if depth < 2 || depth > 40 {
		t.Errorf("expected a finite nesting, got depth %d", depth)
	}
}