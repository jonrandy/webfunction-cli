package generics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/webfunction-protocol/webfunction-go"
)

// load builds a package from an endpoints-and-objects JSON snippet.
func load(t *testing.T, returns string, objects string) *webfunction.Package {
	t.Helper()
	raw := `{"base_url":"https://api.example.com","endpoints":[{"name":"e","returns":` + returns + `,"arguments":[]}],"objects":[` + objects + `]}`
	return loadRaw(t, raw)
}

func loadRaw(t *testing.T, raw string) *webfunction.Package {
	t.Helper()
	var pkg webfunction.Package
	if err := json.Unmarshal([]byte(raw), &pkg); err != nil {
		t.Fatalf("unmarshal: %v\n%s", err, raw)
	}
	return &pkg
}

const (
	userObj = `{"name":"user","attributes":[{"name":"id","type":"string"}]}`
	listObj = `{"name":"list","attributes":[{"name":"total","type":"number"},{"name":"items","type":[["T"]]}]}`
	pairObj = `{"name":"pair","attributes":[{"name":"first","type":"T"},{"name":"second","type":"T"}]}`
	pageObj = `{"name":"page","attributes":[{"name":"results","type":[["T"]]}]}`
)

func objectNames(pkg *webfunction.Package) []string {
	var names []string
	for _, o := range pkg.Objects {
		names = append(names, o.Name)
	}
	return names
}

func attrType(t *testing.T, pkg *webfunction.Package, object, attr string) string {
	t.Helper()
	o := pkg.Object(object)
	if o == nil {
		t.Fatalf("object %q missing; have %v", object, objectNames(pkg))
	}
	for _, a := range o.Attributes {
		if a.Name == attr {
			return a.Type.String()
		}
	}
	t.Fatalf("object %q has no attribute %q", object, attr)
	return ""
}

// Non-generic packages must come out completely unchanged.
func TestNonGenericPackagesUnchanged(t *testing.T) {
	files, _ := filepath.Glob("../testdata/*/*.json")
	if len(files) == 0 {
		t.Fatal("no fixtures found")
	}
	checked := 0
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		var pkg webfunction.Package
		if json.Unmarshal(raw, &pkg) != nil || len(pkg.Endpoints) == 0 {
			continue
		}
		got := Expand(&pkg)
		if !reflect.DeepEqual(got, &pkg) {
			t.Errorf("%s: Expand changed a non-generic package", f)
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no package fixtures were checked")
	}
}

func TestExpandDoesNotMutateInput(t *testing.T) {
	pkg := load(t, `{"object.list":"object.user"}`, userObj+","+listObj)
	before, _ := json.Marshal(pkg)
	Expand(pkg)
	after, _ := json.Marshal(pkg)
	if string(before) != string(after) {
		t.Errorf("input mutated:\nbefore %s\nafter  %s", before, after)
	}
	if got := pkg.Endpoints[0].Returns.String(); got != "object.list<object.user>" {
		t.Errorf("input returns changed: %q", got)
	}
}

func TestSingleApplication(t *testing.T) {
	out := Expand(load(t, `{"object.list":"object.user"}`, userObj+","+listObj))
	if got := out.Endpoints[0].Returns.String(); got != "object.ListOfUser" {
		t.Fatalf("returns: got %q", got)
	}
	if got := attrType(t, out, "ListOfUser", "items"); got != "array<object.user>" {
		t.Errorf("items: got %q", got)
	}
	if got := attrType(t, out, "ListOfUser", "total"); got != "number" {
		t.Errorf("total: got %q", got)
	}
	// The generic template itself is not part of the expanded package.
	if out.Object("list") != nil {
		t.Error("generic template should be removed from the expanded package")
	}
	if out.Object("user") == nil {
		t.Error("ordinary objects must be kept")
	}
}

func TestTypeParameterUsedTwice(t *testing.T) {
	out := Expand(load(t, `{"object.pair":"string"}`, pairObj))
	if attrType(t, out, "PairOfString", "first") != "string" || attrType(t, out, "PairOfString", "second") != "string" {
		t.Error("every T must be replaced by the same argument")
	}
}

func TestUnionArgumentAndSplice(t *testing.T) {
	out := Expand(load(t, `{"object.list":["object.user","null"]}`, userObj+","+listObj))
	if got := out.Endpoints[0].Returns.String(); got != "object.ListOfUserOrNull" {
		t.Fatalf("returns: got %q", got)
	}
	// T inside an array element list is spliced as union members.
	if got := attrType(t, out, "ListOfUserOrNull", "items"); got != "array<object.user|null>" {
		t.Errorf("items: got %q", got)
	}

	// T inside a union is spliced too, and duplicates collapse.
	opt := `{"name":"opt","attributes":[{"name":"v","type":["T","null"]}]}`
	out = Expand(load(t, `{"object.opt":["string","null"]}`, opt))
	if got := attrType(t, out, "OptOfStringOrNull", "v"); got != "string|null" {
		t.Errorf("spliced union: got %q, want string|null", got)
	}

	// An array-typed argument spliced into a union stays an array member.
	out = Expand(load(t, `{"object.opt":[["string"]]}`, opt))
	if got := attrType(t, out, "OptOfArrayOfString", "v"); got != "array<string>|null" {
		t.Errorf("array arg in union: got %q, want array<string>|null", got)
	}
	// ...and an array-typed argument inside an array nests.
	out = Expand(load(t, `{"object.list":[["string"]]}`, listObj))
	if got := attrType(t, out, "ListOfArrayOfString", "items"); got != "array<array<string>>" {
		t.Errorf("nested array: got %q", got)
	}
}

func TestNestedApplicationInnermostFirst(t *testing.T) {
	out := Expand(load(t, `{"object.page":{"object.page":"object.user"}}`, userObj+","+pageObj))
	if got := out.Endpoints[0].Returns.String(); got != "object.PageOfPageOfUser" {
		t.Fatalf("returns: got %q", got)
	}
	if got := attrType(t, out, "PageOfPageOfUser", "results"); got != "array<object.PageOfUser>" {
		t.Errorf("outer results: got %q", got)
	}
	if got := attrType(t, out, "PageOfUser", "results"); got != "array<object.user>" {
		t.Errorf("inner results: got %q", got)
	}
	want := []string{"user", "PageOfUser", "PageOfPageOfUser"}
	if got := objectNames(out); !reflect.DeepEqual(got, want) {
		t.Errorf("object order: got %v, want %v", got, want)
	}
}

func TestInstantiationsAreReused(t *testing.T) {
	raw := `{"base_url":"https://x","endpoints":[
	  {"name":"a","returns":{"object.list":"object.user"},"arguments":[]},
	  {"name":"b","returns":{"object.list":"object.user"},"arguments":[]},
	  {"name":"c","returns":{"object.list":"string"},"arguments":[]}],
	  "objects":[` + userObj + `,` + listObj + `]}`
	out := Expand(loadRaw(t, raw))
	want := []string{"user", "ListOfUser", "ListOfString"}
	if got := objectNames(out); !reflect.DeepEqual(got, want) {
		t.Errorf("objects: got %v, want %v", got, want)
	}
}

func TestApplicationsInArgumentsAndAttributes(t *testing.T) {
	raw := `{"base_url":"https://x","endpoints":[
	  {"name":"a","returns":"object","arguments":[{"name":"f","type":[{"object.pair":"string"},"null"]}],
	   "attributes":[{"name":"p","type":{"object.pair":"number"}}]}],
	  "objects":[` + pairObj + `]}`
	out := Expand(loadRaw(t, raw))
	e := out.Endpoints[0]
	if got := e.Arguments[0].Type.String(); got != "object.PairOfString|null" {
		t.Errorf("argument type: got %q", got)
	}
	if got := e.Attributes[0].Type.String(); got != "object.PairOfNumber" {
		t.Errorf("endpoint attribute type: got %q", got)
	}
}

func TestContextIsMaintained(t *testing.T) {
	filter := `{"name":"filter","arguments":[{"name":"value","type":"T"}]}`
	both := `{"name":"both","arguments":[{"name":"in","type":"T"}],"attributes":[{"name":"out","type":[["T"]]}]}`
	raw := `{"base_url":"https://x","endpoints":[
	  {"name":"a","returns":{"object.both":"string"},"arguments":[{"name":"f","type":{"object.filter":"number"}}]}],
	  "objects":[` + filter + `,` + both + `]}`
	out := Expand(loadRaw(t, raw))

	if out.ObjectInContext("FilterOfNumber", webfunction.ArgumentContext) == nil {
		t.Error("an argument-only generic must resolve in the argument context")
	}
	if out.ObjectInContext("FilterOfNumber", webfunction.AttributeContext) != nil {
		t.Error("an argument-only generic must not gain attributes")
	}
	b := out.Object("BothOfString")
	if b == nil || b.Arguments[0].Type.String() != "string" || b.Attributes[0].Type.String() != "array<string>" {
		t.Errorf("both contexts must be substituted independently: %+v", b)
	}
}

func TestUnappliedGenericResolvesToAny(t *testing.T) {
	out := Expand(load(t, `"object.list"`, listObj))
	if got := out.Endpoints[0].Returns.String(); got != "object.ListOfAny" {
		t.Fatalf("returns: got %q", got)
	}
	if got := attrType(t, out, "ListOfAny", "items"); got != "array<any>" {
		t.Errorf("items: got %q", got)
	}
}

func TestArgumentOnStaticObjectIsIgnored(t *testing.T) {
	out := Expand(load(t, `{"object.user":"string"}`, userObj))
	if got := out.Endpoints[0].Returns.String(); got != "object.user" {
		t.Errorf("got %q, want plain object.user", got)
	}
	if len(out.Objects) != 1 {
		t.Errorf("no instantiation expected: %v", objectNames(out))
	}
	out = Expand(load(t, `{"object.nope":"string"}`, userObj))
	if got := out.Endpoints[0].Returns.String(); got != "object.nope" {
		t.Errorf("dangling application: got %q, want object.nope", got)
	}
}

func TestMalformedApplicationStaysAny(t *testing.T) {
	out := Expand(load(t, `{}`, userObj))
	if got := out.Endpoints[0].Returns.String(); got != "any" {
		t.Errorf("got %q, want any", got)
	}
}

func TestGeneratedNameCollisions(t *testing.T) {
	// A user-defined object already owns the natural name.
	taken := `{"name":"ListOfUser","attributes":[{"name":"x","type":"string"}]}`
	out := Expand(load(t, `{"object.list":"object.user"}`, userObj+","+listObj+","+taken))
	if got := out.Endpoints[0].Returns.String(); got == "object.ListOfUser" {
		t.Fatalf("generated object must not take a user-defined name")
	}
	if attrType(t, out, "ListOfUser", "x") != "string" {
		t.Error("the user-defined object must be untouched")
	}

	// Two different arguments that flatten to the same name stay distinct.
	raw := `{"base_url":"https://x","endpoints":[
	  {"name":"a","returns":{"object.list":[["string"],"null"]},"arguments":[]},
	  {"name":"b","returns":{"object.list":[["string","null"]]},"arguments":[]}],
	  "objects":[` + listObj + `]}`
	out = Expand(loadRaw(t, raw))
	ra, rb := out.Endpoints[0].Returns.String(), out.Endpoints[1].Returns.String()
	if ra == rb {
		t.Fatalf("distinct instantiations collapsed into %q", ra)
	}
	ia := attrType(t, out, strings.TrimPrefix(ra, "object."), "items")
	ib := attrType(t, out, strings.TrimPrefix(rb, "object."), "items")
	if ia != "array<array<string>|null>" || ib != "array<array<string|null>>" {
		t.Errorf("contents wrong: %q / %q", ia, ib)
	}
}

// Invalid polymorphic recursion must not hang the expander.
func TestPolymorphicRecursionTerminates(t *testing.T) {
	loop := `{"name":"a","attributes":[{"name":"next","type":{"object.a":[["T"]]}}]}`
	pkg := load(t, `{"object.a":"string"}`, loop)
	done := make(chan struct{})
	go func() { Expand(pkg); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Expand did not terminate on polymorphic recursion")
	}
}

func TestIsGeneric(t *testing.T) {
	pkg := load(t, `"object.user"`, userObj+","+listObj+","+pairObj)
	if IsGeneric(pkg.Object("user")) || !IsGeneric(pkg.Object("list")) || !IsGeneric(pkg.Object("pair")) {
		t.Error("IsGeneric wrong")
	}
}