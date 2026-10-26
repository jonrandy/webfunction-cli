package openapiconvert

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/webfunction-protocol/webfunction-go"
)

var update = flag.Bool("update", false, "rewrite the golden files from the current converter output")

// goldenInputs lists the fixture packages whose OpenAPI output is pinned
// byte-for-byte: every validate and convert fixture, generic or not, plus
// the generic packages made for codegen.
func goldenInputs(t *testing.T) []string {
	t.Helper()
	var files []string
	for _, pattern := range []string{"validate/*.json", "convert/*.json", "codegen/generics/generic.json", "codegen/generics/collision.json"} {
		matches, err := filepath.Glob(filepath.Join("..", "testdata", pattern))
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, matches...)
	}
	if len(files) == 0 {
		t.Fatal("no fixtures found")
	}
	return files
}

func loadPackage(t *testing.T, path string) *webfunction.Package {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var pkg webfunction.Package
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return &pkg
}

// TestGoldenOutput pins the converter's output for every fixture,
// generic or not, with and without --private, so later changes can prove they
// leave existing packages byte-identical. Run with -update to rewrite the
// golden files deliberately.
func TestGoldenOutput(t *testing.T) {
	for _, path := range goldenInputs(t) {
		dir := filepath.Base(filepath.Dir(path))
		base := strings.TrimSuffix(filepath.Base(path), ".json")
		for _, includePrivate := range []bool{false, true} {
			mode := "public"
			if includePrivate {
				mode = "private"
			}
			name := dir + "-" + base + "." + mode
			t.Run(name, func(t *testing.T) {
				got, err := Generate(loadPackage(t, path), includePrivate)
				if err != nil {
					t.Fatal(err)
				}
				// Generate must be deterministic, or a golden comparison is meaningless.
				again, err := Generate(loadPackage(t, path), includePrivate)
				if err != nil {
					t.Fatal(err)
				}
				if got != again {
					t.Fatal("Generate is not deterministic")
				}

				golden := filepath.Join("..", "testdata", "convert", "golden", "openapi", name+".json")
				if *update {
					if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(golden, []byte(got), 0o644); err != nil {
						t.Fatal(err)
					}
					return
				}
				want, err := os.ReadFile(golden)
				if err != nil {
					t.Fatalf("%v (run go test -update to create it)", err)
				}
				if got != string(want) {
					t.Errorf("output differs from %s", golden)
				}
			})
		}
	}
}