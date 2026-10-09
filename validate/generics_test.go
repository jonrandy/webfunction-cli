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
		"generic-application-in-object":      2, // holder, plus nested-holder reported once (not per level)
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
	if r.ErrorCount != 9 || r.WarningCount != 1 || r.InfoCount != 0 {
		t.Errorf("counts: %d errors / %d warnings / %d info, want 9/1/0", r.ErrorCount, r.WarningCount, r.InfoCount)
	}

	// Findings point at the right places.
	loc := map[string]string{
		"generic-unapplied":                  `endpoint "unapplied" returns`,
		"generic-argument-on-static-object":  `endpoint "static-arg" returns`,
		"generic-application-key-count":      `endpoint "empty-app" returns`,
		"generic-application-key-prefix":     `endpoint "bad-key" returns`,
		"generic-application-argument":       `endpoint "bad-arg" returns`,
		"dangling-object-ref":                `endpoint "ghost-arg" returns`,
		"generic-application-in-object":      `object "holder" attribute "inner"`,
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