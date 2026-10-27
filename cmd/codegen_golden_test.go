package cmd

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite the golden files from the current generator output")

// codegenExts lists the files each target writes, as extensions of the
// output path. Ruby writes a companion .rbs signature file.
var codegenExts = map[string][]string{
	"java":   {".java"},
	"go":     {".go"},
	"php":    {".php"},
	"js":     {".js"},
	"csharp": {".cs"},
	"ruby":   {".rb", ".rbs"},
	"python": {".py"},
}

const goldenSourceURL = "https://example.test/package"

// codegenFixtures lists the fixture packages without generics whose generated
// code is pinned byte-for-byte. Fixtures using generics (generics-*.json) are
// kept apart, in genericFixtures, because codegen warns about generics
// problems and the non-generic fixtures must stay silent.
func codegenFixtures(t *testing.T) []string {
	t.Helper()
	var files []string
	for _, path := range fixtureFiles(t) {
		if strings.HasPrefix(filepath.Base(path), "generics-") {
			continue
		}
		files = append(files, path)
	}
	return files
}

// genericFixtures lists the fixture packages that use generics, whose
// generated code is pinned too: the generics-*.json validate and convert
// fixtures, plus the codegen-specific generic packages.
func genericFixtures(t *testing.T) []string {
	t.Helper()
	var files []string
	for _, path := range fixtureFiles(t) {
		if strings.HasPrefix(filepath.Base(path), "generics-") {
			files = append(files, path)
		}
	}
	return append(files, genericsFixture("generic.json"), genericsFixture("collision.json"))
}

// generateCode runs the codegen command on a fixture and returns each file
// it wrote, keyed by extension. Generics warnings are discarded: the
// flawed fixtures are meant to trigger them.
func generateCode(t *testing.T, target, fixture string, private bool) map[string]string {
	t.Helper()
	oldWarn := warnWriter
	warnWriter = io.Discard
	defer func() { warnWriter = oldWarn }()
	exts := codegenExts[target]
	out := filepath.Join(t.TempDir(), "client"+exts[0])
	args := []string{"--target", target, "--url", goldenSourceURL, "--file", fixture, "-o", out}
	if private {
		args = append(args, "--private")
	}
	if _, err := runCommand(t, codegenRun, args...); err != nil {
		t.Fatalf("codegen --target %s failed on %s: %v", target, fixture, err)
	}
	files := map[string]string{}
	for _, ext := range exts {
		files[ext] = readFile(t, strings.TrimSuffix(out, exts[0])+ext)
	}
	return files
}

// TestCodegenGolden pins the code every target generates for every
// fixture, generic or not, so later changes can prove they leave existing
// packages byte-identical. The --private output is stored only where it
// differs from the public output; everywhere else it must equal it. Run
// with -update to rewrite the golden files deliberately.
func TestCodegenGolden(t *testing.T) {
	for _, target := range validTargets {
		for _, fixture := range append(codegenFixtures(t), genericFixtures(t)...) {
			name := filepath.Base(filepath.Dir(fixture)) + "-" + strings.TrimSuffix(filepath.Base(fixture), ".json")
			t.Run(target+"/"+name, func(t *testing.T) {
				pub := generateCode(t, target, fixture, false)
				// Generation must be deterministic, or a golden comparison is meaningless.
				if again := generateCode(t, target, fixture, false); !reflect.DeepEqual(pub, again) {
					t.Fatal("code generation is not deterministic")
				}
				priv := generateCode(t, target, fixture, true)
				samePrivate := reflect.DeepEqual(pub, priv)

				for ext, content := range pub {
					if strings.Contains(content, "file://") || strings.Contains(content, os.TempDir()) {
						t.Errorf("%s output mentions where the package file was read from; it must not", ext)
					}
				}

				dir := filepath.Join("..", "testdata", "codegen", "golden", target)
				goldenPath := func(mode, ext string) string { return filepath.Join(dir, name+"."+mode+ext) }

				if *update {
					if err := os.MkdirAll(dir, 0o755); err != nil {
						t.Fatal(err)
					}
					for ext, content := range pub {
						if err := os.WriteFile(goldenPath("public", ext), []byte(content), 0o644); err != nil {
							t.Fatal(err)
						}
						if samePrivate {
							_ = os.Remove(goldenPath("private", ext))
						} else if err := os.WriteFile(goldenPath("private", ext), []byte(priv[ext]), 0o644); err != nil {
							t.Fatal(err)
						}
					}
					return
				}

				for ext, content := range pub {
					want, err := os.ReadFile(goldenPath("public", ext))
					if err != nil {
						t.Fatalf("%v (run go test ./cmd -run TestCodegenGolden -update to create it)", err)
					}
					if content != string(want) {
						t.Errorf("public output differs from %s", goldenPath("public", ext))
					}
					privGolden := goldenPath("private", ext)
					if samePrivate {
						if _, err := os.Stat(privGolden); err == nil {
							t.Errorf("--private output now equals the public output, so %s is stale", privGolden)
						}
						continue
					}
					wantPriv, err := os.ReadFile(privGolden)
					if err != nil {
						t.Fatalf("--private output differs from public but %v", err)
					}
					if priv[ext] != string(wantPriv) {
						t.Errorf("--private output differs from %s", privGolden)
					}
				}
			})
		}
	}
}