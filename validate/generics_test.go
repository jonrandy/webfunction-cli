package validate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"
)

func runFixture(t *testing.T, name string) *Report {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "testdata", "validate", name))
	if err != nil {
		t.Fatal(err)
	}
	var pkg webfunction.Package
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return Run(&pkg, "file://"+name)
}

func countByCheck(r *Report) map[string]int {
	m := map[string]int{}
	for _, f := range r.Findings {
		m[f.Check]++
	}
	return m
}

func findFinding(t *testing.T, r *Report, check string) Finding {
	t.Helper()
	for _, f := range r.Findings {
		if f.Check == check {
			return f
		}
	}
	t.Fatalf("no %q finding in report", check)
	return Finding{}
}

func TestGenericsCleanFixtureHasNoFindings(t *testing.T) {
	r := runFixture(t, "generics-clean.json")
	for _, f := range r.Findings {
		t.Errorf("unexpected finding: %s %s: %s", f.Check, f.Location, f.Message)
	}
}

func TestGenericsFlawedFixture(t *testing.T) {
	r := runFixture(t, "generics-flawed.json")
	want := map[string]int{
		"generic-unapplied":                  1,
		"generic-argument-on-static-object":  1,
		"generic-application-key-count":      1,
		"generic-application-key-prefix":     1,
		"generic-application-argument":       1,
		"dangling-object-ref":                1, // the ghost inside an application's argument
		"generic-parameter-outside-object":   3, // t-return, t-argument, t-applied
		"generic-application-cycle":          2, // chain applies itself; ping and pong apply each other (reported once)
		"generic-name-collision":             1,
		"generic-argument-unused-in-context": 1,
	}
	got := countByCheck(r)
	for check, n := range want {
		if got[check] != n {
			t.Errorf("%s: got %d, want %d", check, got[check], n)
		}
	}
	for check := range got {
		if _, ok := want[check]; !ok {
			t.Errorf("unexpected check fired: %s (%d)", check, got[check])
		}
	}
	if r.ErrorCount != 12 || r.WarningCount != 1 || r.InfoCount != 0 {
		t.Errorf("counts: %d errors / %d warnings / %d info, want 12/1/0", r.ErrorCount, r.WarningCount, r.InfoCount)
	}

	// Findings point at the right places.
	loc := map[string]string{
		"generic-unapplied":                  `endpoint "unapplied" returns`,
		"generic-argument-on-static-object":  `endpoint "static-arg" returns`,
		"generic-application-key-count":      `endpoint "empty-app" returns`,
		"generic-application-key-prefix":     `endpoint "bad-key" returns`,
		"generic-application-argument":       `endpoint "bad-arg" returns`,
		"dangling-object-ref":                `endpoint "ghost-arg" returns`,
		"generic-parameter-outside-object":   `endpoint "t-return" returns`,
		"generic-application-cycle":          `object "chain"`,
		"generic-name-collision":             `object "ListOfUser"`,
		"generic-argument-unused-in-context": `endpoint "wrong-context" argument "m"`,
	}
	for check, want := range loc {
		if f := findFinding(t, r, check); f.Location != want {
			t.Errorf("%s: location %q, want %q", check, f.Location, want)
		}
	}
	if f := findFinding(t, r, "generic-argument-unused-in-context"); f.Severity != Warning {
		t.Errorf("unused-in-context should be a warning, got %s", f.Severity)
	}
	for _, check := range []string{"generic-parameter-outside-object", "generic-application-cycle"} {
		if f := findFinding(t, r, check); f.Severity != Error {
			t.Errorf("%s should be an error, got %s", check, f.Severity)
		}
	}
}

func runJSON(t *testing.T, raw string) *Report {
	t.Helper()
	var pkg webfunction.Package
	if err := json.Unmarshal([]byte(raw), &pkg); err != nil {
		t.Fatalf("%v\n%s", err, raw)
	}
	return Run(&pkg, "inline")
}

// packageJSON wraps objects and endpoints (comma separated JSON) into a minimal package.
func packageJSON(objects, endpoints string) string {
	return `{"base_url":"https://api.example.com","name":"Inline","objects":[` + objects + `],"endpoints":[` + endpoints + `]}`
}

const (
	userObject = `{"name":"user","attributes":[{"name":"id","type":"string"}]}`
	listObject = `{"name":"list","attributes":[{"name":"items","type":[["T"]]}]}`
)

// Spec (Generics, Application): an application may appear anywhere a type is expected,
// including inside an object definition.
func TestApplicationsInsideObjectDefinitionsAreAllowed(t *testing.T) {
	cases := map[string]string{
		"generic applied in an attribute": `{"name":"holder","attributes":[{"name":"inner","type":{"object.list":"object.user"}}]}`,
		"generic applied in an argument":  `{"name":"holder","arguments":[{"name":"inner","type":{"object.list":"object.user"}}]}`,
		"nested applications":             `{"name":"holder","attributes":[{"name":"deep","type":{"object.list":{"object.list":"object.user"}}}]}`,
		"in a union":                      `{"name":"holder","attributes":[{"name":"maybe","type":[{"object.list":"object.user"},"null"]}]}`,
		"in an array":                     `{"name":"holder","attributes":[{"name":"many","type":[[{"object.list":"object.user"}]]}]}`,
		"array application":               `{"name":"holder","attributes":[{"name":"many","type":{"array":"object.user"}}]}`,
	}
	for name, holder := range cases {
		r := runJSON(t, packageJSON(userObject+","+listObject+","+holder, `{"name":"get","returns":"object.holder","arguments":[]}`))
		for _, f := range r.Findings {
			if strings.HasPrefix(f.Check, "generic-") {
				t.Errorf("%s: unexpected %s at %s: %s", name, f.Check, f.Location, f.Message)
			}
		}
	}
}

func TestArrayApplicationChecks(t *testing.T) {
	cases := []struct{ name, returns, want string }{
		{"valid element", `{"array":"object.user"}`, ""},
		{"valid union element", `{"array":["object.user","null"]}`, ""},
		{"valid nested", `{"array":{"array":"string"}}`, ""},
		{"valid applied element", `{"object.list":{"array":"object.user"}}`, ""},
		{"element must be a type", `{"array":5}`, "generic-application-argument"},
		{"exactly one key", `{"array":"string","object.user":"string"}`, "generic-application-key-count"},
		{"only the exact key array", `{"arrays":"string"}`, "generic-application-key-prefix"},
		{"other keys are still wrong", `{"string":"object.user"}`, "generic-application-key-prefix"},
	}
	for _, c := range cases {
		r := runJSON(t, packageJSON(userObject+","+listObject, `{"name":"get","returns":`+c.returns+`,"arguments":[]}`))
		var got []string
		for _, f := range r.Findings {
			got = append(got, f.Check)
		}
		if c.want == "" && len(got) != 0 {
			t.Errorf("%s: unexpected findings %v", c.name, got)
		}
		if c.want != "" && (len(got) != 1 || got[0] != c.want) {
			t.Errorf("%s: got %v, want exactly [%s]", c.name, got, c.want)
		}
	}
}

// Spec (Generics, Where T may appear): T is only in scope inside the object definition that
// uses it; anything outside an object definition must not use it.
func TestParameterOutsideObjectDefinition(t *testing.T) {
	cases := []struct{ name, endpoint, location string }{
		{"bare return", `{"name":"e","returns":"T","arguments":[]}`, `endpoint "e" returns`},
		{"argument", `{"name":"e","returns":"string","arguments":[{"name":"a","type":"T"}]}`, `endpoint "e" argument "a"`},
		{"in a union", `{"name":"e","returns":["T","null"],"arguments":[]}`, `endpoint "e" returns`},
		{"in a nested array", `{"name":"e","returns":[["T"]],"arguments":[]}`, `endpoint "e" returns`},
		{"in an array application", `{"name":"e","returns":{"array":"T"},"arguments":[]}`, `endpoint "e" returns`},
		{"as an application value", `{"name":"e","returns":{"object.list":"T"},"arguments":[]}`, `endpoint "e" returns`},
		{"deep in an application value", `{"name":"e","returns":{"object.list":{"object.list":[["T"]]}},"arguments":[]}`, `endpoint "e" returns`},
		{"endpoint attribute", `{"name":"e","returns":"string","arguments":[],"attributes":[{"name":"x","type":"T"}]}`, `endpoint "e" attribute "x"`},
		{"reported once per type", `{"name":"e","returns":["T",{"array":"T"}],"arguments":[]}`, `endpoint "e" returns`},
	}
	for _, c := range cases {
		r := runJSON(t, packageJSON(userObject+","+listObject, c.endpoint))
		var found []Finding
		for _, f := range r.Findings {
			if f.Check == "generic-parameter-outside-object" {
				found = append(found, f)
			}
		}
		if len(found) != 1 {
			t.Errorf("%s: got %d generic-parameter-outside-object findings, want 1 (%v)", c.name, len(found), countByCheck(r))
			continue
		}
		if found[0].Severity != Error || found[0].Location != c.location {
			t.Errorf("%s: got %s at %q, want error at %q", c.name, found[0].Severity, found[0].Location, c.location)
		}
	}
}

func TestParameterInsideObjectDefinitionsIsFine(t *testing.T) {
	objects := userObject + "," + listObject + "," +
		`{"name":"pair","attributes":[{"name":"a","type":"T"},{"name":"b","type":["T","null"]}]},` +
		`{"name":"wrapper","attributes":[{"name":"items","type":{"array":"T"}},{"name":"pair","type":{"object.pair":"T"}}],"arguments":[{"name":"v","type":"T"}]}`
	r := runJSON(t, packageJSON(objects, `{"name":"get","returns":{"object.wrapper":"object.user"},"arguments":[]}`))
	if got := countByCheck(r)["generic-parameter-outside-object"]; got != 0 {
		t.Errorf("T inside object definitions must not be flagged, got %d", got)
	}
}

// Spec (Generics, Recursion): a generic object must not apply itself, directly or through
// other generic objects.
func TestGenericApplicationCycles(t *testing.T) {
	generic := func(name string, members ...string) string {
		return `{"name":"` + name + `","attributes":[{"name":"v","type":"T"}` + strings.Join(members, "") + `]}`
	}
	member := func(name, typ string) string { return `,{"name":"` + name + `","type":` + typ + `}` }

	cases := []struct {
		name     string
		objects  string
		want     int
		location string
		path     string
	}{
		{"direct", generic("chain", member("next", `{"object.chain":"T"}`)), 1, `object "chain"`, "chain -> chain"},
		{"direct with a concrete type", generic("chain", member("next", `{"object.chain":"string"}`)), 1, `object "chain"`, "chain -> chain"},
		{"direct inside the argument of another application", generic("chain", member("next", `{"object.pair":{"object.chain":"T"}}`)) + "," + generic("pair"), 1, `object "chain"`, "chain -> chain"},
		{"mutual", generic("a", member("x", `{"object.b":"T"}`)) + "," + generic("b", member("y", `{"object.a":"T"}`)), 1, `object "a"`, "a -> b -> a"},
		{"three objects", generic("a", member("x", `{"object.b":"T"}`)) + "," + generic("b", member("y", `{"object.c":"T"}`)) + "," + generic("c", member("z", `{"object.a":"T"}`)), 1, `object "a"`, "a -> b -> c -> a"},
		{"through an array element", generic("a", member("x", `{"object.b":[["T"]]}`)) + "," + generic("b", member("y", `{"array":{"object.a":"T"}}`)), 1, `object "a"`, "a -> b -> a"},
		{"through a union", generic("a", member("x", `[{"object.b":"T"},"null"]`)) + "," + generic("b", member("y", `[[{"object.a":"T"}]]`)), 1, `object "a"`, "a -> b -> a"},
		{"through an argument member", `{"name":"a","arguments":[{"name":"x","type":{"object.b":"T"}}],"attributes":[{"name":"v","type":"T"}]},{"name":"b","arguments":[{"name":"p","type":"T"}],"attributes":[{"name":"v","type":"T"},{"name":"y","type":{"object.a":"T"}}]}`, 1, `object "a"`, "a -> b -> a"},
		{"reaching a cycle without being in it", generic("a", member("x", `{"object.b":"T"}`)) + "," + generic("b", member("y", `{"object.b":"T"}`)), 1, `object "b"`, "b -> b"},
		{"a search that backtracks", generic("a", member("x", `{"object.b":"T"}`), member("y", `{"object.c":"T"}`)) + "," + generic("b") + "," + generic("c", member("z", `{"object.a":"T"}`)), 1, `object "a"`, "a -> c -> a"},
		{"two separate cycles", generic("a", member("x", `{"object.a":"T"}`)) + "," + generic("b", member("y", `{"object.b":"T"}`)), 2, "", ""},
		// Not cycles.
		{"a generic applying another generic", generic("a", member("x", `{"object.b":"T"}`)) + "," + generic("b"), 0, "", ""},
		{"applying the same generic twice", generic("a", member("x", `{"object.b":"T"}`), member("y", `{"object.b":{"array":"T"}}`)) + "," + generic("b"), 0, "", ""},
		{"a chain that ends", generic("a", member("x", `{"object.b":"T"}`)) + "," + generic("b", member("y", `{"object.c":"T"}`)) + "," + generic("c"), 0, "", ""},
		{"non-generic recursion through a generic application", generic("list") + `,{"name":"thread","attributes":[{"name":"comments","type":{"object.list":"object.thread"}}]}`, 0, "", ""},
		{"non-generic recursion with an array application", `{"name":"comment","attributes":[{"name":"replies","type":{"array":"object.comment"}}]}`, 0, "", ""},
		{"a bare reference is not an application", generic("a", member("x", `"object.a"`)), 0, "", ""},
	}
	for _, c := range cases {
		r := runJSON(t, packageJSON(c.objects, ""))
		var found []Finding
		for _, f := range r.Findings {
			if f.Check == "generic-application-cycle" {
				found = append(found, f)
			}
		}
		if len(found) != c.want {
			t.Errorf("%s: got %d generic-application-cycle findings, want %d (%v)", c.name, len(found), c.want, countByCheck(r))
			continue
		}
		if c.want == 1 {
			f := found[0]
			if f.Severity != Error || f.Location != c.location || !strings.Contains(f.Message, c.path) {
				t.Errorf("%s: got %s at %q: %s\nwant error at %q mentioning %q", c.name, f.Severity, f.Location, f.Message, c.location, c.path)
			}
		}
	}
}

// Applying a generic at the use site is nesting, not a cycle.
func TestNestedApplicationAtTheUseSiteIsNotACycle(t *testing.T) {
	page := `{"name":"page","attributes":[{"name":"results","type":[["T"]]}]}`
	r := runJSON(t, packageJSON(userObject+","+page, `{"name":"get","returns":{"object.page":{"object.page":"object.user"}},"arguments":[]}`))
	if got := countByCheck(r)["generic-application-cycle"]; got != 0 {
		t.Errorf("got %d cycle findings, want 0", got)
	}
}

// Packages that don't use generics must never trigger a generics check.
func TestNoGenericFindingsOnOtherFixtures(t *testing.T) {
	files, _ := filepath.Glob(filepath.Join("..", "testdata", "validate", "*.json"))
	checked := 0
	for _, f := range files {
		name := filepath.Base(f)
		if strings.HasPrefix(name, "generics-") {
			continue
		}
		raw, _ := os.ReadFile(f)
		var pkg webfunction.Package
		if json.Unmarshal(raw, &pkg) != nil {
			continue
		}
		for _, fd := range Run(&pkg, name).Findings {
			if strings.HasPrefix(fd.Check, "generic-") {
				t.Errorf("%s: unexpected %s", name, fd.Check)
			}
		}
		checked++
	}
	if checked == 0 {
		t.Fatal("no other fixtures were checked")
	}
	if r := runFixture(t, "clean.json"); len(r.Findings) != 0 {
		t.Errorf("clean.json should still have zero findings, got %d", len(r.Findings))
	}
}