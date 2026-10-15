package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"
)

var warningLine = regexp.MustCompile(`^warning: generics: .+ \[(generic-[a-z-]+|dangling-object-ref)\]$`)

func fixturePath(name string) string {
	return filepath.Join("..", "testdata", "validate", name)
}

func loadFixture(t *testing.T, name string) *webfunction.Package {
	t.Helper()
	return parsePackageJSON(t, readFile(t, fixturePath(name)))
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func parsePackageJSON(t *testing.T, raw string) *webfunction.Package {
	t.Helper()
	var pkg webfunction.Package
	if err := json.Unmarshal([]byte(raw), &pkg); err != nil {
		t.Fatal(err)
	}
	return &pkg
}

// warnings returns what warnGenerics writes for pkg, split into lines.
func warnings(t *testing.T, pkg *webfunction.Package) []string {
	t.Helper()
	var buf bytes.Buffer
	warnGenerics(&buf, pkg, "https://api.example.com/pkg")
	if buf.Len() == 0 {
		return nil
	}
	return strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
}

func TestWarnGenericsFlawedPackage(t *testing.T) {
	lines := warnings(t, loadFixture(t, "generics-flawed.json"))
	if len(lines) == 0 {
		t.Fatal("expected warnings for a package with generics problems")
	}

	// Every line but the last is one finding; the last is a single note.
	note := lines[len(lines)-1]
	if !strings.HasPrefix(note, "note: ") || !strings.Contains(note, "wfn validate") {
		t.Errorf("last line should be a note pointing at wfn validate, got %q", note)
	}
	findings := lines[:len(lines)-1]
	for _, l := range findings {
		if !warningLine.MatchString(l) {
			t.Errorf("malformed warning line: %q", l)
		}
	}

	all := strings.Join(findings, "\n")
	for _, check := range []string{
		"generic-unapplied",
		"generic-argument-on-static-object",
		"generic-application-key-count",
		"generic-application-key-prefix",
		"generic-application-argument",
		"generic-name-collision",
		"dangling-object-ref", // the ghost inside an application's argument
		"generic-argument-unused-in-context",
	} {
		if !strings.Contains(all, "["+check+"]") {
			t.Errorf("no warning for %s in:\n%s", check, all)
		}
	}
	if !strings.Contains(all, `endpoint "unapplied" returns`) {
		t.Errorf("warnings should say where the problem is, got:\n%s", all)
	}
	if !strings.Contains(all, `endpoint "ghost-arg" returns`) {
		t.Errorf("the dangling ref inside an application should be reported at its endpoint, got:\n%s", all)
	}
}

func TestWarnGenericsErrorsComeBeforeWarnings(t *testing.T) {
	lines := warnings(t, loadFixture(t, "generics-flawed.json"))
	seenWarningLevel := false
	for _, l := range lines {
		if strings.Contains(l, "[generic-argument-unused-in-context]") {
			seenWarningLevel = true
			continue
		}
		if seenWarningLevel && strings.HasPrefix(l, "warning:") {
			t.Errorf("error-level finding listed after a warning-level one: %q", l)
		}
	}
	if !seenWarningLevel {
		t.Fatal("fixture should produce a generic-argument-unused-in-context finding")
	}
}

func TestWarnGenericsCleanPackageIsSilent(t *testing.T) {
	if lines := warnings(t, loadFixture(t, "generics-clean.json")); lines != nil {
		t.Errorf("a clean generic package must produce no output, got:\n%s", strings.Join(lines, "\n"))
	}
}

// Problems the converters have never warned about must stay silent for
// packages that do not use generics, so existing users see no new noise.
func TestWarnGenericsIgnoresNonGenericProblems(t *testing.T) {
	pkg := parsePackageJSON(t, `{"base_url":"https://api.example.com","name":"Plain","objects":[],"endpoints":[
		{"name":"a","returns":"object.ghost","arguments":[]},
		{"name":"a","returns":"boolean","arguments":[]}]}`)
	if lines := warnings(t, pkg); lines != nil {
		t.Errorf("non-generic problems must not be reported, got:\n%s", strings.Join(lines, "\n"))
	}
}

// Once a package defines a generic, a dangling reference anywhere may be
// the product of an application, so it is reported.
func TestWarnGenericsReportsDanglingRefsInGenericPackages(t *testing.T) {
	pkg := parsePackageJSON(t, `{"base_url":"https://api.example.com","name":"G","objects":[
		{"name":"list","attributes":[{"name":"items","type":[["T"]]}]}],"endpoints":[
		{"name":"get","returns":{"object.list":"object.ghost"},"arguments":[]}]}`)
	lines := warnings(t, pkg)
	if len(lines) != 2 || !strings.Contains(lines[0], "[dangling-object-ref]") {
		t.Errorf("want one dangling-object-ref warning and a note, got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestWarnGenericsDoesNotMutatePackage(t *testing.T) {
	pkg := loadFixture(t, "generics-flawed.json")
	before, _ := json.Marshal(pkg)
	warnings(t, pkg)
	after, _ := json.Marshal(pkg)
	if string(before) != string(after) {
		t.Error("warnGenerics mutated the package")
	}
}

// serve answers every request with the fixture, standing in for a Web
// Function package endpoint.
func serve(t *testing.T, fixture string) string {
	t.Helper()
	body := readFile(t, fixturePath(fixture))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func runConvert(t *testing.T, fixture string) (string, string) {
	t.Helper()
	var warned bytes.Buffer
	old := warnWriter
	warnWriter = &warned
	t.Cleanup(func() { warnWriter = old })

	out := filepath.Join(t.TempDir(), "out.json")
	if err := (&ConvertCommand{}).Run([]string{"--target", "openapi", "--url", serve(t, fixture), "-o", out}); err != nil {
		t.Fatalf("convert failed: %v", err)
	}
	return warned.String(), readFile(t, out)
}

// The warnings are advisory: a package with generics problems still
// converts, exits cleanly, and the problems are reported on the side.
func TestConvertWarnsButStillConverts(t *testing.T) {
	warned, doc := runConvert(t, "generics-flawed.json")
	if !strings.Contains(warned, "[generic-unapplied]") {
		t.Errorf("convert should warn about the unapplied generic, got %q", warned)
	}
	if !json.Valid([]byte(doc)) || !strings.Contains(doc, `"openapi"`) {
		t.Errorf("convert should still write an OpenAPI document, got %.200q", doc)
	}
}

func TestConvertIsSilentForCleanGenericPackage(t *testing.T) {
	warned, doc := runConvert(t, "generics-clean.json")
	if warned != "" {
		t.Errorf("no warnings expected, got %q", warned)
	}
	if !strings.Contains(doc, "ListOfUser") {
		t.Error("the clean generic package should convert to a document containing ListOfUser")
	}
}