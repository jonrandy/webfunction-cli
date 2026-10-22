package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// codegenTargets is every target the codegen command knows.
var codegenTargets = []string{"js", "go", "php", "python", "java", "csharp", "ruby"}

// runCodegenWarned runs codegen and returns what it sent to the generics
// warning writer, what it printed to stdout, and the path of the main output
// file.
func runCodegenWarned(t *testing.T, target string, args ...string) (warned, stdout, out string) {
	t.Helper()
	var buf bytes.Buffer
	old := warnWriter
	warnWriter = &buf
	t.Cleanup(func() { warnWriter = old })

	out = filepath.Join(t.TempDir(), "client"+codegenExts[target][0])
	all := append([]string{"--target", target, "-o", out}, args...)
	stdout, err := runCommand(t, codegenRun, all...)
	if err != nil {
		t.Fatalf("codegen --target %s failed: %v", target, err)
	}
	return buf.String(), stdout, out
}

// Whatever convert says about a package's generics problems, codegen says
// too, in the same words.
func TestCodegenWarnsAboutGenericsProblemsLikeConvert(t *testing.T) {
	want := strings.Join(warnings(t, loadFixture(t, "generics-flawed.json")), "\n") + "\n"
	for _, target := range codegenTargets {
		t.Run(target, func(t *testing.T) {
			warned, _, _ := runCodegenWarned(t, target, "--url", goldenSourceURL, "--file", fixturePath("generics-flawed.json"))
			if warned != want {
				t.Errorf("codegen --target %s should print the same generics warnings as convert.\n got: %q\nwant: %q", target, warned, want)
			}
		})
	}
}

func TestCodegenWarningsNameTheNewChecks(t *testing.T) {
	warned, _, _ := runCodegenWarned(t, "js", "--url", goldenSourceURL, "--file", fixturePath("generics-flawed.json"))
	for _, check := range []string{"generic-unapplied", "generic-parameter-outside-object", "generic-application-cycle"} {
		if !strings.Contains(warned, "["+check+"]") {
			t.Errorf("no warning for %s in:\n%s", check, warned)
		}
	}
	lines := strings.Split(strings.TrimRight(warned, "\n"), "\n")
	if note := lines[len(lines)-1]; !strings.HasPrefix(note, "note: ") || !strings.Contains(note, "wfn validate") {
		t.Errorf("the last line should be the note pointing at wfn validate, got %q", note)
	}
}

// The warnings are advisory: a package with generics problems still
// generates, and the warnings stay off stdout.
func TestCodegenWarnsButStillGenerates(t *testing.T) {
	for _, target := range codegenTargets {
		t.Run(target, func(t *testing.T) {
			warned, stdout, out := runCodegenWarned(t, target, "--url", goldenSourceURL, "--file", fixturePath("generics-flawed.json"))
			if warned == "" {
				t.Fatal("expected generics warnings")
			}
			raw, err := os.ReadFile(out)
			if err != nil || len(raw) == 0 {
				t.Errorf("the generated client should still be written, got %d bytes, err %v", len(raw), err)
			}
			if strings.Contains(stdout, "warning: generics:") {
				t.Errorf("warnings must not go to stdout, got:\n%s", stdout)
			}
		})
	}
}

func TestCodegenWarnsWhenThePackageIsFetched(t *testing.T) {
	warned, stdout, _ := runCodegenWarned(t, "js", "--url", serveFile(t, fixturePath("generics-flawed.json")))
	if !strings.Contains(warned, "[generic-unapplied]") {
		t.Errorf("expected generics warnings for a fetched package, got %q", warned)
	}
	if !strings.Contains(stdout, "Fetched ") {
		t.Errorf("a fetched package should still report Fetched, got %q", stdout)
	}
}

func TestCodegenWarnsWithPrivateToo(t *testing.T) {
	want := strings.Join(warnings(t, loadFixture(t, "generics-flawed.json")), "\n") + "\n"
	warned, _, _ := runCodegenWarned(t, "go", "--url", goldenSourceURL, "--file", fixturePath("generics-flawed.json"), "--private")
	if warned != want {
		t.Errorf("--private should not change the warnings.\n got: %q\nwant: %q", warned, want)
	}
}

// Invalid recursion expands to a pile of ever-longer generated types, so
// this is where the warning matters most.
func TestCodegenWarnsAboutGenericRecursion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "poly.json")
	poly := `{"base_url":"https://api.example.com","name":"Poly","objects":[
		{"name":"chain","attributes":[{"name":"value","type":"T"},{"name":"next","type":{"object.chain":[["T"]]}}]}
	],"endpoints":[{"name":"get-chain","returns":{"object.chain":"string"},"arguments":[]}]}`
	if err := os.WriteFile(path, []byte(poly), 0o644); err != nil {
		t.Fatal(err)
	}
	warned, _, _ := runCodegenWarned(t, "js", "--url", goldenSourceURL, "--file", path)
	if !strings.Contains(warned, `object "chain"`) || !strings.Contains(warned, "[generic-application-cycle]") {
		t.Errorf("expected a cycle warning at object chain, got:\n%s", warned)
	}
}

func TestCodegenIsSilentForCleanGenericPackage(t *testing.T) {
	for _, target := range codegenTargets {
		t.Run(target, func(t *testing.T) {
			if warned, _, _ := runCodegenWarned(t, target, "--url", goldenSourceURL, "--file", fixturePath("generics-clean.json")); warned != "" {
				t.Errorf("no warnings expected for a clean generic package, got %q", warned)
			}
		})
	}
}

// Existing users must see nothing new: packages without generics stay silent.
func TestCodegenIsSilentForPackagesWithoutGenerics(t *testing.T) {
	for _, fixture := range codegenFixtures(t) {
		t.Run(filepath.Base(fixture), func(t *testing.T) {
			if warned, _, _ := runCodegenWarned(t, "js", "--url", goldenSourceURL, "--file", fixture); warned != "" {
				t.Errorf("no warnings expected without generics, got %q", warned)
			}
		})
	}
}